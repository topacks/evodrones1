package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/root/evodrones/backend/models"
	"github.com/root/evodrones/backend/services"
)

type MaintenanceOrderController struct {
	maintenanceOrderService *services.MaintenanceOrderService
}

func NewMaintenanceOrderController(maintenanceOrderService *services.MaintenanceOrderService) *MaintenanceOrderController {
	return &MaintenanceOrderController{
		maintenanceOrderService: maintenanceOrderService,
	}
}

func (c *MaintenanceOrderController) Create(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)

	var input models.CreateMaintenanceOrderInput
	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	order, err := c.maintenanceOrderService.Create(userID, input)
	if err != nil {
		switch err {
		case services.ErrUnauthorized:
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		case services.ErrDroneNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "failed to create maintenance order")
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(order)
}

func (c *MaintenanceOrderController) GetByID(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)

	orderID := ctx.Params("id")

	order, err := c.maintenanceOrderService.GetByID(userID, orderID)
	if err != nil {
		switch err {
		case services.ErrMaintenanceOrderNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case services.ErrUnauthorized:
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "failed to get maintenance order")
		}
	}

	return ctx.JSON(order)
}

func (c *MaintenanceOrderController) GetByUserID(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)

	orders, err := c.maintenanceOrderService.GetByUserID(userID)
	if err != nil {
		if err == services.ErrUnauthorized {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get maintenance orders")
	}

	return ctx.JSON(orders)
}

func (c *MaintenanceOrderController) GetByDroneID(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)

	droneID := ctx.Params("id")

	orders, err := c.maintenanceOrderService.GetByDroneID(userID, droneID)
	if err != nil {
		switch err {
		case services.ErrDroneNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case services.ErrUnauthorized:
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "failed to get maintenance orders")
		}
	}

	return ctx.JSON(orders)
}

func (c *MaintenanceOrderController) Update(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)

	orderID := ctx.Params("id")

	var input models.UpdateMaintenanceOrderInput
	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	order, err := c.maintenanceOrderService.Update(userID, orderID, input)
	if err != nil {
		switch err {
		case services.ErrMaintenanceOrderNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case services.ErrUnauthorized:
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "failed to update maintenance order")
		}
	}

	return ctx.JSON(order)
}

func (c *MaintenanceOrderController) Delete(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)

	orderID := ctx.Params("id")

	if err := c.maintenanceOrderService.Delete(userID, orderID); err != nil {
		switch err {
		case services.ErrMaintenanceOrderNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case services.ErrUnauthorized:
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "failed to delete maintenance order")
		}
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (c *MaintenanceOrderController) List(ctx *fiber.Ctx) error {
	offset, _ := strconv.Atoi(ctx.Query("offset", "0"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	orders, err := c.maintenanceOrderService.List(offset, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to list maintenance orders")
	}

	return ctx.JSON(orders)
}

func (c *MaintenanceOrderController) GetByStatus(ctx *fiber.Ctx) error {
	status := models.MaintenanceOrderStatus(ctx.Params("status"))
	if status == "" {
		return fiber.NewError(fiber.StatusBadRequest, "status is required")
	}

	orders, err := c.maintenanceOrderService.GetByStatus(status)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get maintenance orders by status")
	}

	return ctx.JSON(orders)
}
