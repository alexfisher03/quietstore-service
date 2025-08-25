package v1

import (
	"os"
	"strconv"
	"time"

	"github.com/alexfisher03/quietstore-service/QuietStore/internal/handlers"
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/repo"
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/service"
	"github.com/gofiber/fiber/v2"
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

func RegisterRoutes(app *fiber.App, storage service.StorageService, users repo.Users, refresh repo.RefreshTokens) {
	v1 := app.Group("/api/v1")
	v1.Get("/health", handlers.HealthCheck)

	secret := os.Getenv("AUTH_JWT_SECRET")

	accessMin := parseIntEnv("AUTH_ACCESS_TTL_MIN", 1)
	accessTTL := time.Duration(accessMin) * time.Minute

	var refreshTTL time.Duration
	if mins := os.Getenv("AUTH_REFRESH_TTL_MIN"); mins != "" {
		m := parseIntEnv("AUTH_REFRESH_TTL_MIN", 0)
		if m > 0 {
			refreshTTL = time.Duration(m) * time.Minute
		}
	}
	if refreshTTL == 0 {
		days := parseIntEnv("AUTH_REFRESH_TTL_DAYS", 7)
		refreshTTL = time.Duration(days) * 24 * time.Hour
	}

	authHandlers := handlers.NewAuthHandler(users, refresh, secret, accessTTL, refreshTTL)
	v1.Post("/auth/login", authHandlers.LoginHandler)
	v1.Post("/auth/refresh", authHandlers.RefreshHandler)

	authMW := handlers.RequireAuth(secret)

	userHandlers := handlers.NewUserHandler(users)
	v1.Get("/users/:id", userHandlers.GetUserByIDHandler)
	v1.Post("/users", userHandlers.CreateUserHandler)
	v1.Patch("/users/:id", userHandlers.UpdateUserHandler)
	v1.Delete("/users/:id", userHandlers.DeleteUserHandler)
	v1.Get("/users", userHandlers.GetAllUsersHandler)

	fileHandlers := handlers.NewFileHandler(storage)
	me := v1.Group("/me", authMW)
	me.Get("/files", fileHandlers.GetUserFilesHandler)
	me.Get("/files/:fileID", fileHandlers.GetUserFileByIDHandler)
	me.Delete("/files/:fileID", fileHandlers.DeleteUserFileByIDHandler)
	me.Post("/files/upload", fileHandlers.UploadFileHandler)
}
