//	@title			QuietStore API
//	@version		0.1
//	@description	QuietStore file storage API: auth, users, files.
//	@BasePath		/api/v1

//	@contact.name	QuietStore Devs
//	@contact.url	https://github.com/alexfisher03/quietstore-service
//	@contact.email	you@example.com

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Paste as: Bearer <ACCESS_JWT>

package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/middleware/helmet"

	"github.com/alexfisher03/quietstore-service/QuietStore/internal/config"
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/db"
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/handlers"
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/repo"
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"

	v1 "github.com/alexfisher03/quietstore-service/QuietStore/api/v1"
	fiberSwagger "github.com/gofiber/swagger"
)

//go:embed openapi.json
var openapiSpec []byte

func mustConnectDB(dsn string) *pgxpool.Pool {
	deadline := time.Now().Add(30 * time.Second)
	var pool *pgxpool.Pool
	var err error
	for {
		pool, err = db.Connect(context.Background(), dsn)
		if err == nil {
			return pool
		}
		if time.Now().After(deadline) {
			log.Fatalf("DB connect failed (timeout): %v", err)
		}
		log.Printf("DB not ready yet: %v; retrying...", err)
		time.Sleep(2 * time.Second)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dsn := os.Getenv("DB_DSN")
	pool := mustConnectDB(dsn)
	defer pool.Close()
	if err := db.Migrate(context.Background(), pool); err != nil {
		log.Fatalf("DB migrate failed: %v", err)
	}

	usersRepo := repo.NewUsersPGX(pool)
	filesRepo := repo.NewFilesPGX(pool)
	refreshRepo := repo.NewRefreshPGX(pool)

	// minio
	endpoint := os.Getenv("MINIO_ENDPOINT")
	ak := os.Getenv("MINIO_ACCESS_KEY")
	sk := os.Getenv("MINIO_SECRET_KEY")
	bucket := os.Getenv("MINIO_BUCKET")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	s3c := newMinIOS3Client(endpoint, ak, sk, useSSL)
	ensureBucket(context.Background(), s3c, bucket)
	storage := service.NewMinIOStorageService(s3c, bucket, filesRepo)

	requireHTTPS := os.Getenv("REQUIRE_HTTPS") == "true"

	app := fiber.New(fiber.Config{
		ServerHeader:            "QuietStore/1.0",
		ReadTimeout:             cfg.Server.ReadTimeout,
		WriteTimeout:            cfg.Server.WriteTimeout,
		BodyLimit:               cfg.Server.BodyLimit,
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"0.0.0.0/0"},
	})

	app.Use(helmet.New())

	if requireHTTPS {
		app.Use(func(c *fiber.Ctx) error {
			if c.Path() == "/api/v1/health" || strings.HasPrefix(c.Path(), "/docs/") || c.Path() == "/openapi.yaml" {
				return c.Next()
			}
			if c.Protocol() == "https" || strings.EqualFold(c.Get("X-Forwarded-Proto"), "https") {
				return c.Next()
			}
			return fiber.NewError(fiber.StatusUpgradeRequired, "TLS required")
		})
	}

	app.Get("/api/v1/health", handlers.HealthCheck)
	app.Get("/api/v1/ready", handlers.ReadyCheck(pool, s3c, bucket))
	app.Use(logger.New())
	app.Use(recover.New())

	v1.RegisterRoutes(app, cfg.App, storage, usersRepo, refreshRepo)

	app.Get("/openapi.json", func(c *fiber.Ctx) error {
		c.Type("json")
		return c.Send(openapiSpec)
	})
	if os.Getenv("ENABLE_SWAGGER") == "true" {
		app.Get("/docs/*", fiberSwagger.New(fiberSwagger.Config{
			URL: "/openapi.json",
		}))
	}

	go func() {
		ticker := time.NewTicker(6 * time.Hour)
		defer ticker.Stop()

		const revokedRetention = 30 * time.Hour

		for t := range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			now := t.UTC()

			var aboutToExpire, aboutToDrop int64
			_ = pool.QueryRow(ctx, `SELECT count(*) FROM refresh_tokens WHERE expires_at < $1`, now).Scan(&aboutToExpire)
			_ = pool.QueryRow(ctx, `SELECT count(*) FROM refresh_tokens WHERE revoked_at IS NOT NULL AND revoked_at < $1`, now.Add(-revokedRetention)).Scan(&aboutToDrop)

			deleted, err := refreshRepo.Purge(ctx, now, now.Add(-revokedRetention))
			cancel()

			if err != nil {
				log.Printf("[refresh-purge] ran at %s UTC, would-expire=%d, would-revoke=%d, ERROR: %v",
					now.Format(time.RFC3339), aboutToExpire, aboutToDrop, err)
				continue
			}
			log.Printf("[refresh-purge] ran at %s UTC, would-expire=%d, would-revoke=%d, deleted=%d",
				now.Format(time.RFC3339), aboutToExpire, aboutToDrop, deleted)
		}
	}()

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Fatal(app.Listen(addr))
}
