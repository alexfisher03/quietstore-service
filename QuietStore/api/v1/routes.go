package v1

import (
	"os"
	"strconv"
	"time"

	"github.com/alexfisher03/quietstore-service/QuietStore/internal/config"
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/handlers"
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/repo"
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func parseIntEnv(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil || i <= 0 {
		return def
	}
	return i
}

func RegisterRoutes(app *fiber.App, appCfg config.AppConfig, storage service.StorageService, users repo.Users, refresh repo.RefreshTokens) {
	v1 := app.Group("/api/v1")
	v1.Get("/health", handlers.HealthCheck)

	secret := os.Getenv("AUTH_JWT_SECRET")
	accessTTL := time.Duration(parseIntEnv("AUTH_ACCESS_TTL_MIN", 10)) * time.Minute
	refreshTTL := time.Duration(parseIntEnv("AUTH_REFRESH_TTL_MIN", 0)) * time.Minute
	if refreshTTL == 0 {
		refreshTTL = time.Duration(parseIntEnv("AUTH_REFRESH_TTL_DAYS", 7)) * 24 * time.Hour
	}

	authMW := handlers.RequireAuth(secret)
	authHandlers := handlers.NewAuthHandler(users, refresh, secret, accessTTL, refreshTTL)

	sensitive := v1.Group("/auth", limiter.New(limiter.Config{
		Max:        appCfg.RateLimitAuthMax,
		Expiration: time.Duration(appCfg.RateLimitAuthExpire) * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, "too many requests slow down son")
		},
	}))
	sensitive.Post("/login", authHandlers.LoginHandler)
	sensitive.Post("/refresh", authHandlers.RefreshHandler)
	sensitive.Post("/logout", authMW, authHandlers.LogoutHandler)

	userHandlers := handlers.NewUserHandler(users)
	userLimiter := v1.Group("/users", limiter.New(limiter.Config{
		Max:        appCfg.RateLimitUserMax,
		Expiration: time.Duration(appCfg.RateLimitUserExpire) * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, "too many requests guy")
		},
	}))
	userLimiter.Post("", userHandlers.CreateUserHandler)
	userLimiter.Get("/:id", userHandlers.GetUserByIDHandler)
	userLimiter.Patch("/:id", userHandlers.UpdateUserHandler)
	userLimiter.Delete("/:id", userHandlers.DeleteUserHandler)
	userLimiter.Get("", userHandlers.GetAllUsersHandler)

	fileHandlers := handlers.NewFileHandler(storage)
	me := v1.Group("/me", authMW)
	filesLimiter := me.Group("/files", limiter.New(limiter.Config{
		Max:        appCfg.RateLimitFileMax,
		Expiration: time.Duration(appCfg.RateLimitFileExpire) * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, "too many requests guy")
		},
	}))
	me.Get("/files", fileHandlers.GetUserFilesHandler)
	me.Get("/files/:fileID", fileHandlers.GetUserFileByIDHandler)
	filesLimiter.Delete("/:fileID", fileHandlers.DeleteUserFileByIDHandler)
	filesLimiter.Post("/upload", fileHandlers.UploadFileHandler)
	// me.Get("/files/search", fileHandlers.SearchFilesHandler) @@@@@@@@@ v2 @@@@@@@@@
	filesLimiter.Patch("/:fileID/rename", fileHandlers.RenameFileHandler)
}
