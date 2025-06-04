package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/root/evodrones/backend/services"
)

// ResourceAuthMiddleware cria um middleware para autorização de acesso a recursos específicos
func ResourceAuthMiddleware(authService *services.AuthorizationService, resourceType, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Obter ID do usuário e role do contexto
		userID := c.Locals("user_id")
		if userID == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "usuário não autenticado")
		}

		userRole := c.Locals("role")
		if userRole == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "role não encontrada")
		}

		// Obter ID do recurso dos parâmetros da rota
		resourceID := c.Params("id")
		if resourceID == "" {
			// Se não houver ID específico, verificar permissão para listar recursos
			canList, err := authService.CanListResources(resourceType, userID.(string), userRole.(string))
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "erro ao verificar permissões")
			}

			if !canList {
				return fiber.NewError(fiber.StatusForbidden, "acesso negado")
			}

			// Armazenar os IDs dos recursos acessíveis no contexto para filtragem posterior
			accessibleIDs, err := authService.GetAccessibleResources(resourceType, userID.(string), userRole.(string), action)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "erro ao obter recursos acessíveis")
			}

			c.Locals("accessible_resources", accessibleIDs)
			return c.Next()
		}

		// Verificar permissão para o recurso específico
		canAccess, err := authService.CanAccessResource(resourceType, resourceID, userID.(string), userRole.(string), action)
		if err != nil {
			switch err {
			case services.ErrResourceNotFound, services.ErrDroneNotFound, services.ErrMaintenanceOrderNotFound, services.ErrTechnicalReportNotFound:
				return fiber.NewError(fiber.StatusNotFound, err.Error())
			default:
				return fiber.NewError(fiber.StatusInternalServerError, "erro ao verificar permissões")
			}
		}

		if !canAccess {
			return fiber.NewError(fiber.StatusForbidden, "acesso negado")
		}

		return c.Next()
	}
}

// OwnershipMiddleware verifica se o usuário é proprietário do recurso ou admin
func OwnershipMiddleware(resourceType string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		if userID == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "usuário não autenticado")
		}

		userRole := c.Locals("role")
		if userRole == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "role não encontrada")
		}

		// Administradores têm acesso total
		if userRole.(string) == "admin" {
			return c.Next()
		}

		// Para usuários comuns, verificar propriedade do recurso
		resourceUserID := c.Locals(resourceType + "_user_id")
		if resourceUserID == nil {
			return fiber.NewError(fiber.StatusInternalServerError, "informação de propriedade não disponível")
		}

		if resourceUserID.(string) != userID.(string) {
			return fiber.NewError(fiber.StatusForbidden, "acesso negado")
		}

		return c.Next()
	}
}

// AuthorizeResource verifica se o usuário tem permissão para acessar um recurso específico
func AuthorizeResource(resourceType, paramName, action string, authService *services.AuthorizationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Obter ID do usuário e role do contexto (definidos pelo AuthMiddleware)
		userID := c.Locals("user_id")
		if userID == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "usuário não autenticado")
		}

		userRole := c.Locals("role")
		if userRole == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "role não encontrada")
		}

		// Obter ID do recurso (se aplicável)
		var resourceID string
		if paramName != "" {
			resourceID = c.Params(paramName)
		}

		// Verificar autorização
		authorized, err := authService.CanAccessResource(
			resourceType,
			resourceID,
			userID.(string),
			userRole.(string),
			action,
		)

		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "erro ao verificar autorização")
		}

		if !authorized {
			return fiber.NewError(fiber.StatusForbidden, "acesso negado ao recurso")
		}

		return c.Next()
	}
}
