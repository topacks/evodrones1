package permissions

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// PermissionMiddleware creates a unified permission middleware
func PermissionMiddleware(service *Service, resource, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract user information from context (set by auth middleware)
		userID := c.Locals("user_id")
		if userID == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "usuário não autenticado")
		}

		userIDStr, ok := userID.(string)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "ID de usuário inválido")
		}

		// Build permission context
		permissionContext := buildPermissionContext(c)

		// Check permission
		allowed, err := service.CheckPermission(context.Background(), userIDStr, resource, action, permissionContext)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "erro ao verificar permissões")
		}

		if !allowed {
			return fiber.NewError(fiber.StatusForbidden, "acesso negado")
		}

		return c.Next()
	}
}

// ResourcePermissionMiddleware creates middleware for resource-specific permissions
func ResourcePermissionMiddleware(service *Service, resource, action, resourceIDParam string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract user information from context
		userID := c.Locals("user_id")
		if userID == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "usuário não autenticado")
		}

		userIDStr, ok := userID.(string)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "ID de usuário inválido")
		}

		// Get resource ID from URL parameters
		resourceID := c.Params(resourceIDParam)
		if resourceID == "" {
			return fiber.NewError(fiber.StatusBadRequest, "ID do recurso não fornecido")
		}

		// Build permission context with resource information
		permissionContext := buildPermissionContext(c)
		permissionContext["resource_id"] = resourceID

		// Add resource owner information if available
		resourceOwnerID, err := getResourceOwnerID(service, resource, resourceID)
		if err == nil && resourceOwnerID != "" {
			permissionContext["resource_owner_id"] = resourceOwnerID
		}

		// Check permission
		allowed, err := service.CheckPermission(context.Background(), userIDStr, resource, action, permissionContext)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "erro ao verificar permissões")
		}

		if !allowed {
			return fiber.NewError(fiber.StatusForbidden, "acesso negado")
		}

		return c.Next()
	}
}

// RoleMiddleware creates middleware that requires specific roles
func RoleMiddleware(service *Service, requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract user information from context
		userID := c.Locals("user_id")
		if userID == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "usuário não autenticado")
		}

		userIDStr, ok := userID.(string)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "ID de usuário inválido")
		}

		// Get user roles
		userRoles, err := service.GetUserRoles(context.Background(), userIDStr)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "erro ao obter roles do usuário")
		}

		// Check if user has any of the required roles
		hasRequiredRole := false
		for _, userRole := range userRoles {
			for _, requiredRole := range requiredRoles {
				if userRole.Name == requiredRole {
					hasRequiredRole = true
					break
				}
			}
			if hasRequiredRole {
				break
			}
		}

		if !hasRequiredRole {
			return fiber.NewError(fiber.StatusForbidden, "role insuficiente")
		}

		return c.Next()
	}
}

// AdminOnlyMiddleware creates middleware that requires admin role
func AdminOnlyMiddleware(service *Service) fiber.Handler {
	return RoleMiddleware(service, RoleAdmin)
}

// TechnicianOrAdminMiddleware creates middleware that requires technician or admin role
func TechnicianOrAdminMiddleware(service *Service) fiber.Handler {
	return RoleMiddleware(service, RoleTechnician, RoleAdmin)
}

// OwnershipMiddleware creates middleware that checks resource ownership
func OwnershipMiddleware(service *Service, resource, resourceIDParam string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract user information from context
		userID := c.Locals("user_id")
		if userID == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "usuário não autenticado")
		}

		userIDStr, ok := userID.(string)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "ID de usuário inválido")
		}

		// Check if user is admin (admins can access everything)
		userRoles, err := service.GetUserRoles(context.Background(), userIDStr)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "erro ao obter roles do usuário")
		}

		for _, role := range userRoles {
			if role.Name == RoleAdmin {
				return c.Next() // Admin can access everything
			}
		}

		// Get resource ID from URL parameters
		resourceID := c.Params(resourceIDParam)
		if resourceID == "" {
			return fiber.NewError(fiber.StatusBadRequest, "ID do recurso não fornecido")
		}

		// Check resource ownership
		resourceOwnerID, err := getResourceOwnerID(service, resource, resourceID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "erro ao verificar propriedade do recurso")
		}

		if resourceOwnerID != userIDStr {
			return fiber.NewError(fiber.StatusForbidden, "acesso negado - recurso não pertence ao usuário")
		}

		return c.Next()
	}
}

// Helper functions

// buildPermissionContext builds the permission context from the request
func buildPermissionContext(c *fiber.Ctx) map[string]interface{} {
	context := make(map[string]interface{})

	// Add user information
	if userID := c.Locals("user_id"); userID != nil {
		context["user_id"] = userID
	}
	if userRole := c.Locals("role"); userRole != nil {
		context["user_role"] = userRole
	}

	// Add request information
	context["method"] = c.Method()
	context["path"] = c.Path()
	context["client_ip"] = c.IP()
	context["user_agent"] = c.Get("User-Agent")

	// Add query parameters
	queryParams := make(map[string]string)
	c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
		queryParams[string(key)] = string(value)
	})
	context["query_params"] = queryParams

	// Add headers (selective)
	headers := make(map[string]string)
	c.Request().Header.VisitAll(func(key, value []byte) {
		keyStr := strings.ToLower(string(key))
		// Only include safe headers
		if keyStr == "content-type" || keyStr == "accept" || keyStr == "accept-language" {
			headers[keyStr] = string(value)
		}
	})
	context["headers"] = headers

	return context
}

// getResourceOwnerID retrieves the owner ID for a specific resource
func getResourceOwnerID(service *Service, resource, resourceID string) (string, error) {
	// This is a simplified implementation
	// In a real application, you would query the database to get the resource owner
	// For now, we'll implement basic logic for known resources

	ctx := context.Background()

	switch resource {
	case ResourceDrone:
		// Query drone owner from database
		var ownerID string
		err := service.db.GetContext(ctx, &ownerID, "SELECT user_id FROM drones WHERE id = $1", resourceID)
		if err != nil {
			return "", err
		}
		return ownerID, nil

	case ResourceMaintenanceOrder:
		// Query maintenance order owner from database
		var ownerID string
		err := service.db.GetContext(ctx, &ownerID, "SELECT user_id FROM maintenance_orders WHERE id = $1", resourceID)
		if err != nil {
			return "", err
		}
		return ownerID, nil

	case ResourceTechnicalReport:
		// Query technical report owner from database
		var ownerID string
		err := service.db.GetContext(ctx, &ownerID, "SELECT user_id FROM technical_reports WHERE id = $1", resourceID)
		if err != nil {
			return "", err
		}
		return ownerID, nil

	case ResourceUser:
		// For user resources, the resource ID is the owner ID
		return resourceID, nil

	default:
		return "", nil // Unknown resource type
	}
}

// Convenience middleware creators for common patterns

// CreateDroneMiddleware creates middleware for drone creation
func CreateDroneMiddleware(service *Service) fiber.Handler {
	return PermissionMiddleware(service, ResourceDrone, ActionCreate)
}

// ReadDroneMiddleware creates middleware for drone reading
func ReadDroneMiddleware(service *Service) fiber.Handler {
	return ResourcePermissionMiddleware(service, ResourceDrone, ActionRead, "id")
}

// UpdateDroneMiddleware creates middleware for drone updating
func UpdateDroneMiddleware(service *Service) fiber.Handler {
	return ResourcePermissionMiddleware(service, ResourceDrone, ActionUpdate, "id")
}

// DeleteDroneMiddleware creates middleware for drone deletion
func DeleteDroneMiddleware(service *Service) fiber.Handler {
	return ResourcePermissionMiddleware(service, ResourceDrone, ActionDelete, "id")
}

// CreateMaintenanceOrderMiddleware creates middleware for maintenance order creation
func CreateMaintenanceOrderMiddleware(service *Service) fiber.Handler {
	return PermissionMiddleware(service, ResourceMaintenanceOrder, ActionCreate)
}

// ReadMaintenanceOrderMiddleware creates middleware for maintenance order reading
func ReadMaintenanceOrderMiddleware(service *Service) fiber.Handler {
	return ResourcePermissionMiddleware(service, ResourceMaintenanceOrder, ActionRead, "id")
}

// UpdateMaintenanceOrderMiddleware creates middleware for maintenance order updating
func UpdateMaintenanceOrderMiddleware(service *Service) fiber.Handler {
	return ResourcePermissionMiddleware(service, ResourceMaintenanceOrder, ActionUpdate, "id")
}

// DeleteMaintenanceOrderMiddleware creates middleware for maintenance order deletion
func DeleteMaintenanceOrderMiddleware(service *Service) fiber.Handler {
	return ResourcePermissionMiddleware(service, ResourceMaintenanceOrder, ActionDelete, "id")
}

// CreateTechnicalReportMiddleware creates middleware for technical report creation
func CreateTechnicalReportMiddleware(service *Service) fiber.Handler {
	return PermissionMiddleware(service, ResourceTechnicalReport, ActionCreate)
}

// ReadTechnicalReportMiddleware creates middleware for technical report reading
func ReadTechnicalReportMiddleware(service *Service) fiber.Handler {
	return ResourcePermissionMiddleware(service, ResourceTechnicalReport, ActionRead, "id")
}

// UpdateTechnicalReportMiddleware creates middleware for technical report updating
func UpdateTechnicalReportMiddleware(service *Service) fiber.Handler {
	return ResourcePermissionMiddleware(service, ResourceTechnicalReport, ActionUpdate, "id")
}

// DeleteTechnicalReportMiddleware creates middleware for technical report deletion
func DeleteTechnicalReportMiddleware(service *Service) fiber.Handler {
	return ResourcePermissionMiddleware(service, ResourceTechnicalReport, ActionDelete, "id")
}

// ListResourcesMiddleware creates middleware for listing resources
func ListResourcesMiddleware(service *Service, resource string) fiber.Handler {
	return PermissionMiddleware(service, resource, ActionList)
}
