package handlers

import (
	"crypto-orderbook/internal/config"
	"crypto-orderbook/internal/models"
	"crypto-orderbook/internal/repository"
	"crypto-orderbook/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

func NewAuthHandler(userRepo *repository.UserRepository, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validate input
	req.Email = strings.TrimSpace(req.Email)
	req.Username = strings.TrimSpace(req.Username)

	if req.Email == "" || req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "All fields are required"})
	}

	if len(req.Password) < 6 {
		return c.Status(400).JSON(fiber.Map{"error": "Password must be at least 6 characters"})
	}

	// Check if email exists
	exists, err := h.userRepo.EmailExists(c.Context(), req.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Server error"})
	}
	if exists {
		return c.Status(400).JSON(fiber.Map{"error": "Email already exists"})
	}

	// Check if username exists
	exists, err = h.userRepo.UsernameExists(c.Context(), req.Username)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Server error"})
	}
	if exists {
		return c.Status(400).JSON(fiber.Map{"error": "Username already exists"})
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Server error"})
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashedPassword,
	}

	if err := h.userRepo.Create(c.Context(), user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Username, h.cfg.JWT.Secret, h.cfg.JWT.ExpireHours)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.Status(201).JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	req.Email = strings.TrimSpace(req.Email)

	if req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Email and password are required"})
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(c.Context(), req.Email)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Check password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Username, h.cfg.JWT.Secret, h.cfg.JWT.ExpireHours)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	})
}
