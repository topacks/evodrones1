package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/root/evodrones/backend/config"
	"github.com/root/evodrones/backend/controllers"
	"github.com/root/evodrones/backend/middleware"
	"github.com/root/evodrones/backend/services"
)

func SetupTechnicalReportRoutes(router fiber.Router, cfg *config.Config, controller *controllers.TechnicalReportController, authorizationService *services.AuthorizationService) {
	reports := router.Group("/technical-reports")

	reports.Post("/", middleware.AuthorizeResource(services.ResourceTypeTechnicalReport, "", services.ActionCreate, authorizationService), controller.Create)
	reports.Get("/", middleware.AuthorizeResource(services.ResourceTypeTechnicalReport, "", services.ActionList, authorizationService), controller.List)
	reports.Get("/:id", middleware.AuthorizeResource(services.ResourceTypeTechnicalReport, "id", services.ActionRead, authorizationService), controller.GetByID)
	reports.Put("/:id", middleware.AuthorizeResource(services.ResourceTypeTechnicalReport, "id", services.ActionUpdate, authorizationService), controller.Update)
	reports.Delete("/:id", middleware.AuthorizeResource(services.ResourceTypeTechnicalReport, "id", services.ActionDelete, authorizationService), controller.Delete)

	reports.Get("/maintenance/:maintenanceID", middleware.AuthorizeResource(services.ResourceTypeTechnicalReport, "maintenanceID", services.ActionRead, authorizationService), controller.GetByMaintenanceID)
	reports.Get("/user/:userID", middleware.AuthorizeResource(services.ResourceTypeTechnicalReport, "userID", services.ActionRead, authorizationService), controller.GetByUserID)
}
