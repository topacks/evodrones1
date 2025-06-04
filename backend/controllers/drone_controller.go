package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/root/evodrones/backend/models"
	"github.com/root/evodrones/backend/services"
)

type DroneController struct {
	droneService *services.DroneService
}

func NewDroneController(droneService *services.DroneService) *DroneController {
	return &DroneController{
		droneService: droneService,
	}
}

func (c *DroneController) Create(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)

	var input models.CreateDroneInput
	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	// Removendo a validação que usa o pacote validate não definido

	drone, err := c.droneService.Create(userID, input)
	if err != nil {
		if err == services.ErrUnauthorized {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create drone")
	}

	return ctx.Status(fiber.StatusCreated).JSON(drone)
}

func (c *DroneController) GetByID(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)

	droneID := ctx.Params("id")

	drone, err := c.droneService.GetByID(userID, droneID)
	if err != nil {
		switch err {
		case services.ErrDroneNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case services.ErrUnauthorized:
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "failed to get drone")
		}
	}

	return ctx.JSON(drone)
}

func (c *DroneController) GetByUserID(ctx *fiber.Ctx) error {
	// Removendo a variável userRole que não é usada
	userID := ctx.Locals("user_id").(string)

	drones, err := c.droneService.GetByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get drones")
	}

	return ctx.JSON(drones)
}

func (c *DroneController) Update(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)

	droneID := ctx.Params("id")

	var input models.UpdateDroneInput
	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	// Removendo a validação que usa o pacote validate não definido

	drone, err := c.droneService.Update(userID, droneID, input)
	if err != nil {
		switch err {
		case services.ErrDroneNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case services.ErrUnauthorized:
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "failed to update drone")
		}
	}

	return ctx.JSON(drone)
}

func (c *DroneController) Delete(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)

	droneID := ctx.Params("id")

	if err := c.droneService.Delete(userID, droneID); err != nil {
		switch err {
		case services.ErrDroneNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case services.ErrUnauthorized:
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "failed to delete drone")
		}
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (c *DroneController) List(ctx *fiber.Ctx) error {
	offset, _ := strconv.Atoi(ctx.Query("offset", "0"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	drones, err := c.droneService.List(offset, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to list drones")
	}

	return ctx.JSON(drones)
}

func (c *DroneController) GetByStatus(ctx *fiber.Ctx) error {
	status := ctx.Params("status")
	if status == "" {
		return fiber.NewError(fiber.StatusBadRequest, "status is required")
	}

	drones, err := c.droneService.GetByStatus(status)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get drones by status")
	}

	return ctx.JSON(drones)
}

func (c *DroneController) GetDueForCheck(ctx *fiber.Ctx) error {
	drones, err := c.droneService.GetDueForCheck()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get drones due for check")
	}

	return ctx.JSON(drones)
}
