package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/root/evodrones/backend/config"
	"github.com/root/evodrones/backend/controllers"
	"github.com/root/evodrones/backend/middleware"
	"github.com/root/evodrones/backend/services"
)

func SetupRoutes(router fiber.Router, cfg *config.Config, droneController *controllers.DroneController, maintenanceOrderController *controllers.MaintenanceOrderController, authorizationService *services.AuthorizationService) {
	// Grupo de rotas autenticadas
	api := router.Group("/") // Já está no prefixo /api

	// Rotas de drones
	drones := api.Group("/drones")
	drones.Post("/", middleware.AuthorizeResource(services.ResourceTypeDrone, "", services.ActionCreate, authorizationService), droneController.Create)
	drones.Get("/", middleware.AuthorizeResource(services.ResourceTypeDrone, "", services.ActionList, authorizationService), droneController.List)
	drones.Get("/user", middleware.AuthorizeResource(services.ResourceTypeDrone, "", services.ActionList, authorizationService), droneController.GetByUserID)
	drones.Get("/status/:status", middleware.AuthorizeResource(services.ResourceTypeDrone, "", services.ActionList, authorizationService), droneController.GetByStatus)
	drones.Get("/due-check", middleware.AuthorizeResource(services.ResourceTypeDrone, "", services.ActionList, authorizationService), droneController.GetDueForCheck)
	drones.Get("/:id", middleware.AuthorizeResource(services.ResourceTypeDrone, "id", services.ActionRead, authorizationService), droneController.GetByID)
	drones.Put("/:id", middleware.AuthorizeResource(services.ResourceTypeDrone, "id", services.ActionUpdate, authorizationService), droneController.Update)
	drones.Delete("/:id", middleware.AuthorizeResource(services.ResourceTypeDrone, "id", services.ActionDelete, authorizationService), droneController.Delete)

	// Rotas de ordens de manutenção
	maintenanceOrders := api.Group("/maintenance-orders")
	maintenanceOrders.Post("/", middleware.AuthorizeResource(services.ResourceTypeMaintenanceOrder, "", services.ActionCreate, authorizationService), maintenanceOrderController.Create)
	maintenanceOrders.Get("/", middleware.AuthorizeResource(services.ResourceTypeMaintenanceOrder, "", services.ActionList, authorizationService), maintenanceOrderController.List)
	maintenanceOrders.Get("/user", middleware.AuthorizeResource(services.ResourceTypeMaintenanceOrder, "", services.ActionList, authorizationService), maintenanceOrderController.GetByUserID)
	maintenanceOrders.Get("/drone/:drone_id", middleware.AuthorizeResource(services.ResourceTypeMaintenanceOrder, "drone_id", services.ActionList, authorizationService), maintenanceOrderController.GetByDroneID)
	maintenanceOrders.Get("/status/:status", middleware.AuthorizeResource(services.ResourceTypeMaintenanceOrder, "", services.ActionList, authorizationService), maintenanceOrderController.GetByStatus)
	maintenanceOrders.Get("/:id", middleware.AuthorizeResource(services.ResourceTypeMaintenanceOrder, "id", services.ActionRead, authorizationService), maintenanceOrderController.GetByID)
	maintenanceOrders.Put("/:id", middleware.AuthorizeResource(services.ResourceTypeMaintenanceOrder, "id", services.ActionUpdate, authorizationService), maintenanceOrderController.Update)
	maintenanceOrders.Delete("/:id", middleware.AuthorizeResource(services.ResourceTypeMaintenanceOrder, "id", services.ActionDelete, authorizationService), maintenanceOrderController.Delete)
}
