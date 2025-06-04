package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/root/evodrones/backend/controllers"
)

func SetupAuthRoutes(router fiber.Router, authController *controllers.AuthController) {
	// Rotas públicas de autenticação
	auth := router.Group("/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
}
