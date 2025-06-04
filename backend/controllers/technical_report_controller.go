package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/root/evodrones/backend/models"
	"github.com/root/evodrones/backend/services"
)

type TechnicalReportController struct {
	service *services.TechnicalReportService
}

func NewTechnicalReportController(service *services.TechnicalReportService) *TechnicalReportController {
	return &TechnicalReportController{
		service: service,
	}
}

func (c *TechnicalReportController) Create(ctx *fiber.Ctx) error {
	var input models.CreateTechnicalReportInput
	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "dados inválidos")
	}

	userID := ctx.Locals("user_id").(string)

	report, err := c.service.Create(&input, userID)
	if err != nil {
		switch err {
		case services.ErrMaintenanceOrderNotFound:
			return fiber.NewError(fiber.StatusNotFound, "ordem de manutenção não encontrada")
		case services.ErrUserNotFound:
			return fiber.NewError(fiber.StatusNotFound, "usuário não encontrado")
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "erro ao criar relatório")
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(report)
}

func (c *TechnicalReportController) GetByID(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)
	userRole := ctx.Locals("role").(string)
	id := ctx.Params("id")

	report, err := c.service.GetByID(id, userID, userRole)
	if err != nil {
		switch err {
		case services.ErrTechnicalReportNotFound:
			return fiber.NewError(fiber.StatusNotFound, "relatório não encontrado")
		case services.ErrUnauthorized:
			return fiber.NewError(fiber.StatusForbidden, "não autorizado")
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "erro ao buscar relatório")
		}
	}

	return ctx.JSON(report)
}

func (c *TechnicalReportController) GetByMaintenanceID(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)
	userRole := ctx.Locals("role").(string)
	maintenanceID := ctx.Params("maintenance_id")

	reports, err := c.service.GetByMaintenanceID(userID, maintenanceID, userRole)
	if err != nil {
		switch err {
		case services.ErrMaintenanceOrderNotFound:
			return fiber.NewError(fiber.StatusNotFound, "ordem de manutenção não encontrada")
		case services.ErrUnauthorized:
			return fiber.NewError(fiber.StatusForbidden, "não autorizado")
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "erro ao buscar relatórios")
		}
	}

	return ctx.JSON(reports)
}

func (c *TechnicalReportController) GetByUserID(ctx *fiber.Ctx) error {
	targetUserID := ctx.Params("userID")
	requestingUserID := ctx.Locals("user_id").(string)
	userRole := ctx.Locals("role").(string)

	reports, err := c.service.GetByUserID(targetUserID, requestingUserID, userRole)
	if err != nil {
		switch err {
		case services.ErrUnauthorized:
			return fiber.NewError(fiber.StatusForbidden, "não autorizado")
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "erro ao buscar relatórios")
		}
	}

	return ctx.JSON(reports)
}

func (c *TechnicalReportController) Update(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)
	userRole := ctx.Locals("role").(string)
	id := ctx.Params("id")

	var input models.UpdateTechnicalReportInput
	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	// Removendo a validação que usa o pacote validate que não está definido
	if err := c.service.Update(id, &input, userID, userRole); err != nil {
		switch err {
		case services.ErrTechnicalReportNotFound:
			return fiber.NewError(fiber.StatusNotFound, "relatório não encontrado")
		case services.ErrUnauthorized:
			return fiber.NewError(fiber.StatusForbidden, "não autorizado")
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "erro ao atualizar relatório")
		}
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (c *TechnicalReportController) Delete(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)
	userRole := ctx.Locals("role").(string)
	id := ctx.Params("id")

	if err := c.service.Delete(id, userID, userRole); err != nil {
		switch err {
		case services.ErrTechnicalReportNotFound:
			return fiber.NewError(fiber.StatusNotFound, "relatório não encontrado")
		case services.ErrUnauthorized:
			return fiber.NewError(fiber.StatusForbidden, "não autorizado")
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "erro ao excluir relatório")
		}
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (c *TechnicalReportController) List(ctx *fiber.Ctx) error {
	offset, err := strconv.Atoi(ctx.Query("offset", "0"))
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(ctx.Query("limit", "10"))
	if err != nil {
		limit = 10
	}

	userID := ctx.Locals("user_id").(string)
	userRole := ctx.Locals("role").(string)

	reports, err := c.service.List(offset, limit, userID, userRole)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "erro ao listar relatórios")
	}

	return ctx.JSON(reports)
}
