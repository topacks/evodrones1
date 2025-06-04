package permissions

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/root/evodrones/backend/config"
	"github.com/root/evodrones/backend/controllers"
)

// Integration contém as funções para integrar o sistema RBAC com o projeto existente

// SetupPermissionRoutes configura as rotas com o novo sistema de permissões
func SetupPermissionRoutes(app *fiber.App, cfg *config.Config, droneController *controllers.DroneController, maintenanceController *controllers.MaintenanceOrderController, reportController *controllers.TechnicalReportController, permissionService *Service) {
	// Grupo de rotas da API
	api := app.Group("/api")

	// Rotas de Drones
	drones := api.Group("/drones")
	drones.Post("/", CreateDroneMiddleware(permissionService), droneController.Create)
	drones.Get("/", ListResourcesMiddleware(permissionService, "drone"), droneController.List)
	drones.Get("/:id", ReadDroneMiddleware(permissionService), droneController.GetByID)
	drones.Put("/:id", UpdateDroneMiddleware(permissionService), droneController.Update)
	drones.Delete("/:id", DeleteDroneMiddleware(permissionService), droneController.Delete)

	// Rotas de Ordens de Manutenção
	maintenance := api.Group("/maintenance-orders")
	maintenance.Post("/", CreateMaintenanceOrderMiddleware(permissionService), maintenanceController.Create)
	maintenance.Get("/", ListResourcesMiddleware(permissionService, "maintenance_order"), maintenanceController.List)
	maintenance.Get("/:id", ReadMaintenanceOrderMiddleware(permissionService), maintenanceController.GetByID)
	maintenance.Put("/:id", UpdateMaintenanceOrderMiddleware(permissionService), maintenanceController.Update)
	maintenance.Delete("/:id", DeleteMaintenanceOrderMiddleware(permissionService), maintenanceController.Delete)

	// Rotas de Relatórios Técnicos
	reports := api.Group("/technical-reports")
	reports.Post("/", CreateTechnicalReportMiddleware(permissionService), reportController.Create)
	reports.Get("/", ListResourcesMiddleware(permissionService, "technical_report"), reportController.List)
	reports.Get("/:id", ReadTechnicalReportMiddleware(permissionService), reportController.GetByID)
	reports.Put("/:id", UpdateTechnicalReportMiddleware(permissionService), reportController.Update)
	reports.Delete("/:id", DeleteTechnicalReportMiddleware(permissionService), reportController.Delete)

	// Rotas de Usuários (apenas admin)
	users := api.Group("/users")
	users.Post("/", AdminOnlyMiddleware(permissionService), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Usuário criado"})
	})
	users.Get("/", AdminOnlyMiddleware(permissionService), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Lista de usuários"})
	})
	users.Get("/:id", AdminOnlyMiddleware(permissionService), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Detalhes do usuário"})
	})
	users.Put("/:id", AdminOnlyMiddleware(permissionService), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Usuário atualizado"})
	})
	users.Delete("/:id", AdminOnlyMiddleware(permissionService), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Usuário excluído"})
	})

	// Rotas de Gerenciamento de Roles
	roles := api.Group("/roles")
	roles.Post("/", AdminOnlyMiddleware(permissionService), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Role criado"})
	})
	roles.Get("/", AdminOnlyMiddleware(permissionService), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Lista de roles"})
	})
	roles.Put("/:id", AdminOnlyMiddleware(permissionService), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Role atualizado"})
	})
	roles.Delete("/:id", AdminOnlyMiddleware(permissionService), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Role excluído"})
	})

	// Rotas de Atribuição de Roles
	userRoles := api.Group("/user-roles")
	userRoles.Post("/assign", AdminOnlyMiddleware(permissionService), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Role atribuído ao usuário"})
	})
	userRoles.Post("/revoke", AdminOnlyMiddleware(permissionService), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Role revogado do usuário"})
	})
	userRoles.Get("/user/:userID", AdminOnlyMiddleware(permissionService), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Roles do usuário"})
	})
}

// MigrateFromOldSystem migra do sistema de autorização antigo para o novo RBAC
func MigrateFromOldSystem(db *sqlx.DB, permissionService *Service) error {
	ctx := context.Background()

	log.Println("Iniciando migração do sistema antigo...")

	// 1. Migrar usuários existentes
	query := `
		SELECT id, role 
		FROM users 
		WHERE id NOT IN (
			SELECT DISTINCT user_id 
			FROM user_roles
		)
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	migratedCount := 0
	for rows.Next() {
		var userID, oldRole string
		if err := rows.Scan(&userID, &oldRole); err != nil {
			continue
		}

		// Mapear role antigo para novo
		var newRoleID string
		switch oldRole {
		case "admin":
			newRoleID = RoleAdmin
		case "technician":
			newRoleID = RoleTechnician
		default:
			newRoleID = RoleViewer
		}

		// Atribuir novo role
		if err := permissionService.AssignRole(ctx, userID, newRoleID, "system"); err != nil {
			log.Printf("Erro ao migrar usuário %s: %v", userID, err)
			continue
		}

		migratedCount++
	}

	log.Printf("Migração concluída: %d usuários migrados", migratedCount)
	return nil
}

// ValidateIntegration valida se a integração está funcionando corretamente
func ValidateIntegration(permissionService *Service) error {
	ctx := context.Background()

	// Verificar se roles padrão existem
	roles, err := permissionService.ListRoles(ctx)
	if err != nil {
		return err
	}

	if len(roles) == 0 {
		return ErrRoleNotFound
	}

	// Verificar se permissões padrão existem
	hasAdminRole := false
	hasTechnicianRole := false
	hasViewerRole := false

	for _, role := range roles {
		switch role.Name {
		case RoleAdmin:
			hasAdminRole = true
		case RoleTechnician:
			hasTechnicianRole = true
		case RoleViewer:
			hasViewerRole = true
		}
	}

	if !hasAdminRole || !hasTechnicianRole || !hasViewerRole {
		return ErrRoleNotFound
	}

	log.Println("Integração validada com sucesso")
	return nil
}

// GetPermissionContext extrai o contexto de permissão do Fiber Context
func GetPermissionContext(c *fiber.Ctx) map[string]interface{} {
	context := make(map[string]interface{})

	// Informações do usuário
	if userID := c.Locals("user_id"); userID != nil {
		context["user_id"] = userID
	}

	if userRole := c.Locals("role"); userRole != nil {
		context["user_role"] = userRole
	}

	// Informações da requisição
	context["client_ip"] = c.IP()
	context["user_agent"] = c.Get("User-Agent")
	context["method"] = c.Method()
	context["path"] = c.Path()

	// Parâmetros da rota
	if resourceID := c.Params("id"); resourceID != "" {
		context["resource_id"] = resourceID
	}

	// Headers customizados
	if department := c.Get("X-User-Department"); department != "" {
		context["user_department"] = department
	}

	return context
}

// LogPermissionAction registra uma ação de permissão para auditoria
func LogPermissionAction(c *fiber.Ctx, action, resource, resourceID string, allowed bool) {
	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = uid.(string)
	}

	status := "DENIED"
	if allowed {
		status = "ALLOWED"
	}

	log.Printf("PERMISSION_AUDIT: user=%s action=%s resource=%s resource_id=%s status=%s ip=%s",
		userID, action, resource, resourceID, status, c.IP())
}

// ExtractResourceOwnerID extrai o ID do proprietário do recurso
func ExtractResourceOwnerID(c *fiber.Ctx, resourceType, resourceID string, db *sqlx.DB) (string, error) {
	var ownerID string
	var query string

	switch resourceType {
	case ResourceDrone:
		query = "SELECT user_id FROM drones WHERE id = $1"
	case ResourceMaintenanceOrder:
		query = "SELECT user_id FROM maintenance_orders WHERE id = $1"
	case ResourceTechnicalReport:
		query = "SELECT user_id FROM technical_reports WHERE id = $1"
	case ResourceUser:
		// Para usuários, o próprio ID é o owner
		return resourceID, nil
	default:
		return "", ErrInvalidResource
	}

	err := db.GetContext(context.Background(), &ownerID, query, resourceID)
	if err != nil {
		return "", err
	}

	return ownerID, nil
}
