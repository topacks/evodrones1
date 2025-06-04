package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/root/evodrones/backend/models"
	"github.com/root/evodrones/backend/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var input models.CreateUserInput
	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	user, err := c.authService.Register(input)
	if err != nil {
		if err.Error() == "user already exists" {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to register user")
	}

	return ctx.Status(fiber.StatusCreated).JSON(user)
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var input models.LoginInput
	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	response, err := c.authService.Login(input)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		case services.ErrInvalidCredentials:
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "failed to login")
		}
	}

	return ctx.JSON(response)
}
