package controllers

import (
	"server/src/configs"
	"server/src/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	type RegisterRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	user := models.User{Username: req.Username, Password: string(hashedPassword)}
	result := configs.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusCreated).SendString("User created successfully")
}

func Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var user models.User
	result := configs.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
	}

	token, err := models.GenerateToken(user.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{"token": token})
}