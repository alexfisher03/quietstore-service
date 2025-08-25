package handlers

import (
	"strconv"
	"time"

	"github.com/alexfisher03/quietstore-service/QuietStore/internal/models"
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/repo"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	users repo.Users
}

func NewUserHandler(users repo.Users) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) GetUserByIDHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	u, err := h.users.ByID(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "lookup failed: "+err.Error())
	}
	if u == nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}
	return c.JSON(fiber.Map{
		"id":         u.ID,
		"username":   u.Username,
		"email":      u.Email,
		"created_at": u.CreatedAt,
	})
}

// POST /api/v1/users
func (h *UserHandler) CreateUserHandler(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request payload")
	}
	if input.Username == "" || input.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing username or password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "hash failed")
	}

	u := &models.User{
		ID:        models.GenerateUserID(),
		Username:  input.Username,
		Email:     input.Email,
		Password:  string(hash),
		CreatedAt: time.Now(),
	}
	if err := h.users.Create(c.Context(), u); err != nil {
		return fiber.NewError(fiber.StatusConflict, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":         u.ID,
		"username":   u.Username,
		"email":      u.Email,
		"created_at": u.CreatedAt,
	})
}

// PATCH /api/v1/users/:id
func (h *UserHandler) UpdateUserHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	u, err := h.users.ByID(c.Context(), id)
	if err != nil {
		return fiber.NewError(500, "lookup failed: "+err.Error())
	}
	if u == nil {
		return fiber.NewError(404, "user not found")
	}

	var body struct {
		Username *string `json:"username"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(400, "invalid body")
	}

	if body.Username != nil {
		u.Username = *body.Username
	}
	if body.Email != nil {
		u.Email = *body.Email
	}
	if body.Password != nil && *body.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*body.Password), bcrypt.DefaultCost)
		if err != nil {
			return fiber.NewError(500, "hash failed")
		}
		u.Password = string(hash)
	}

	if err := h.users.Update(c.Context(), u); err != nil {
		return fiber.NewError(500, "update failed: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
	})
}

// DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUserHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.users.Delete(c.Context(), id); err != nil {
		return fiber.NewError(404, "user not found")
	}
	return c.JSON(fiber.Map{"message": "user deleted"})
}

// GET /api/v1/users?limit=50&offset=0
func (h *UserHandler) GetAllUsersHandler(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	users, err := h.users.List(c.Context(), limit, offset)
	if err != nil {
		return fiber.NewError(500, "list failed: "+err.Error())
	}

	out := make([]fiber.Map, 0, len(users))
	for _, u := range users {
		out = append(out, fiber.Map{
			"id":         u.ID,
			"username":   u.Username,
			"email":      u.Email,
			"created_at": u.CreatedAt,
		})
	}
	return c.JSON(out)
}
