package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alexfisher03/quietstore-service/QuietStore/internal/repo"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	users      repo.Users
	refresh    repo.RefreshTokens
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	issuer     string
	audience   string
}

func NewAuthHandler(users repo.Users, refresh repo.RefreshTokens, secret string, accessTTL, refreshTTL time.Duration) *AuthHandler {
	iss := os.Getenv("AUTH_ISSUER")
	if iss == "" {
		iss = "quietstore"
	}
	aud := os.Getenv("AUTH_AUDIENCE")
	if aud == "" {
		aud = "quietstore-api"
	}

	return &AuthHandler{
		users:      users,
		refresh:    refresh,
		jwtSecret:  []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		issuer:     iss,
		audience:   aud,
	}
}

func (h *AuthHandler) LoginHandler(c *fiber.Ctx) error {
	var in struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&in); err != nil || in.Username == "" || in.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid credentials payload")
	}

	u, err := h.users.ByUsername(c.Context(), in.Username)
	if err != nil || u == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(in.Password)); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}

	// Access JWT
	now := time.Now().UTC()
	accessExp := now.Add(h.accessTTL)
	claims := jwt.MapClaims{
		"sub":      u.ID,
		"user_id":  u.ID,
		"username": u.Username,
		"iss":      h.issuer,
		"aud":      h.audience,
		"iat":      now.Unix(),
		"nbf":      now.Unix(),
		"exp":      accessExp.Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessStr, err := tok.SignedString(h.jwtSecret)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to sign token")
	}

	refreshRaw := uuid.NewString()
	refreshHash := sha256.Sum256([]byte(refreshRaw))
	if err := h.refresh.Insert(c.Context(), u.ID, hex.EncodeToString(refreshHash[:]), now.Add(h.refreshTTL)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to persist refresh token")
	}

	return c.JSON(fiber.Map{
		"access_token":  accessStr,
		"token_type":    "Bearer",
		"expires_in":    int(h.accessTTL.Seconds()),
		"refresh_token": refreshRaw,
		"user_id":       u.ID,
		"username":      u.Username,
	})
}

func (h *AuthHandler) RefreshHandler(c *fiber.Ctx) error {
	var in struct {
		UserID       string `json:"user_id"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&in); err != nil || in.UserID == "" || in.RefreshToken == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid refresh payload")
	}

	now := time.Now().UTC()
	hash := sha256.Sum256([]byte(in.RefreshToken))
	valid, err := h.refresh.FindValid(c.Context(), in.UserID, hex.EncodeToString(hash[:]), now)
	if err != nil || !valid {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired refresh token")
	}
	_ = h.refresh.Revoke(c.Context(), in.UserID, hex.EncodeToString(hash[:]))

	accessExp := now.Add(h.accessTTL)
	claims := jwt.MapClaims{
		"sub":     in.UserID,
		"user_id": in.UserID,
		"iss":     h.issuer,
		"aud":     h.audience,
		"iat":     now.Unix(),
		"nbf":     now.Unix(),
		"exp":     accessExp.Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessStr, err := tok.SignedString(h.jwtSecret)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to sign token")
	}

	newRefresh := uuid.NewString()
	newHash := sha256.Sum256([]byte(newRefresh))
	if err := h.refresh.Insert(c.Context(), in.UserID, hex.EncodeToString(newHash[:]), now.Add(h.refreshTTL)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to persist rotated refresh token")
	}

	return c.JSON(fiber.Map{
		"access_token":  accessStr,
		"token_type":    "Bearer",
		"expires_in":    int(h.accessTTL.Seconds()),
		"refresh_token": newRefresh,
	})
}

func RequireAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return fiber.NewError(fiber.StatusUnauthorized, "missing bearer token")
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil || !tok.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}
		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid claims")
		}
		userID, _ := claims["user_id"].(string)
		if userID == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing user id")
		}
		c.Locals("userID", userID)
		return c.Next()
	}
}

func (h *AuthHandler) LogoutHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing user id context")
	}

	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&input); err != nil || input.RefreshToken == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid logout payload")
	}

	hash := sha256.Sum256([]byte(input.RefreshToken))
	if err := h.refresh.Revoke(c.Context(), userID, hex.EncodeToString(hash[:])); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to revoke refresh token")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
