package permissions

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// ExampleController demonstrates how to use the new RBAC system in controllers
type ExampleController struct {
	permissionService *Service
	db                *sqlx.DB
}

// NewExampleController creates a new example controller
func NewExampleController(permissionService *Service, db *sqlx.DB) *ExampleController {
	return &ExampleController{
		permissionService: permissionService,
		db:                db,
	}
}

// ExampleDroneCreate demonstrates how to create a drone with permission checking
func (ec *ExampleController) ExampleDroneCreate(c *fiber.Ctx) error {
	// Extract user information from context (set by auth middleware)
	userID := c.Locals("user_id").(string)
	userRole := c.Locals("role").(string)

	// Build permission context
	permissionContext := map[string]interface{}{
		"user_id":   userID,
		"user_role": userRole,
		"client_ip": c.IP(),
		"method":    c.Method(),
		"path":      c.Path(),
	}

	// Check permission using the new RBAC system
	allowed, err := ec.permissionService.CheckPermission(
		context.Background(),
		userID,
		ResourceDrone,
		ActionCreate,
		permissionContext,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "erro ao verificar permissões")
	}

	if !allowed {
		return fiber.NewError(fiber.StatusForbidden, "acesso negado")
	}

	// Parse request body
	type CreateDroneRequest struct {
		Name         string `json:"name" validate:"required,min=3,max=100"`
		Model        string `json:"model" validate:"required,min=3,max=50"`
		SerialNumber string `json:"serial_number" validate:"required,min=5,max=50"`
	}

	var req CreateDroneRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "dados inválidos")
	}

	// Validate request
	if err := validateStruct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Create drone logic here...
	log.Printf("Creating drone: %+v by user: %s", req, userID)

	// Audit the action
	auditContext := map[string]interface{}{
		"action":      "drone_created",
		"drone_name":  req.Name,
		"drone_model": req.Model,
	}
	ec.auditAction(userID, "drone_create", "drone", "", nil, auditContext, c)

	return c.JSON(fiber.Map{
		"message": "Drone criado com sucesso",
		"data":    req,
	})
}

// ExampleDroneUpdate demonstrates how to update a drone with ownership checking
func (ec *ExampleController) ExampleDroneUpdate(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userRole := c.Locals("role").(string)
	droneID := c.Params("id")

	// Get drone owner information
	var droneOwnerID string
	err := ec.db.Get(&droneOwnerID, "SELECT user_id FROM drones WHERE id = $1", droneID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "drone não encontrado")
	}

	// Build permission context with ownership information
	permissionContext := map[string]interface{}{
		"user_id":           userID,
		"user_role":         userRole,
		"resource_id":       droneID,
		"resource_owner_id": droneOwnerID,
		"client_ip":         c.IP(),
		"method":            c.Method(),
		"path":              c.Path(),
	}

	// Check permission
	allowed, err := ec.permissionService.CheckPermission(
		context.Background(),
		userID,
		ResourceDrone,
		ActionUpdate,
		permissionContext,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "erro ao verificar permissões")
	}

	if !allowed {
		return fiber.NewError(fiber.StatusForbidden, "acesso negado")
	}

	// Parse request body
	type UpdateDroneRequest struct {
		Name   *string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
		Model  *string `json:"model,omitempty" validate:"omitempty,min=3,max=50"`
		Status *string `json:"status,omitempty" validate:"omitempty,oneof=active inactive maintenance"`
	}

	var req UpdateDroneRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "dados inválidos")
	}

	// Validate request
	if err := validateStruct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Update drone logic here...
	log.Printf("Updating drone %s: %+v by user: %s", droneID, req, userID)

	// Audit the action
	auditContext := map[string]interface{}{
		"action":   "drone_updated",
		"drone_id": droneID,
		"changes":  req,
	}
	ec.auditAction(userID, "drone_update", "drone", droneID, nil, auditContext, c)

	return c.JSON(fiber.Map{
		"message": "Drone atualizado com sucesso",
		"data":    req,
	})
}

// ExampleUserRoleManagement demonstrates role management operations
func (ec *ExampleController) ExampleUserRoleManagement(c *fiber.Ctx) error {
	currentUserID := c.Locals("user_id").(string)
	targetUserID := c.Params("userId")

	// Check if current user can manage roles
	allowed, err := ec.permissionService.HasPermission(
		context.Background(),
		currentUserID,
		ResourceUser,
		ActionUpdate,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "erro ao verificar permissões")
	}

	if !allowed {
		return fiber.NewError(fiber.StatusForbidden, "acesso negado")
	}

	// Parse request
	type RoleAssignmentRequest struct {
		Action string `json:"action" validate:"required,oneof=assign revoke"`
		RoleID string `json:"role_id" validate:"required,uuid"`
	}

	var req RoleAssignmentRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "dados inválidos")
	}

	// Validate request
	if err := validateStruct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Perform role assignment/revocation
	switch req.Action {
	case "assign":
		err = ec.permissionService.AssignRole(context.Background(), targetUserID, req.RoleID, currentUserID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "erro ao atribuir role")
		}
		ec.auditAction(currentUserID, "role_assigned", "user", targetUserID, nil, map[string]interface{}{
			"target_user_id": targetUserID,
			"role_id":        req.RoleID,
		}, c)

	case "revoke":
		err = ec.permissionService.RevokeRole(context.Background(), targetUserID, req.RoleID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "erro ao revogar role")
		}
		ec.auditAction(currentUserID, "role_revoked", "user", targetUserID, nil, map[string]interface{}{
			"target_user_id": targetUserID,
			"role_id":        req.RoleID,
		}, c)
	}

	return c.JSON(fiber.Map{
		"message": "Operação realizada com sucesso",
	})
}

// ExampleListWithPermissionFiltering demonstrates how to list resources with permission filtering
func (ec *ExampleController) ExampleListWithPermissionFiltering(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userRole := c.Locals("role").(string)

	// Check list permission
	allowed, err := ec.permissionService.HasPermission(
		context.Background(),
		userID,
		ResourceDrone,
		ActionList,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "erro ao verificar permissões")
	}

	if !allowed {
		return fiber.NewError(fiber.StatusForbidden, "acesso negado")
	}

	// Build query based on user role and permissions
	var query string
	var args []interface{}

	// Check if user is admin (can see all drones)
	userRoles, err := ec.permissionService.GetUserRoles(context.Background(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "erro ao obter roles")
	}

	isAdmin := false
	for _, role := range userRoles {
		if role.Name == RoleAdmin {
			isAdmin = true
			break
		}
	}

	if isAdmin {
		// Admin can see all drones
		query = "SELECT id, name, model, serial_number, status, user_id FROM drones ORDER BY name"
	} else {
		// Non-admin users can only see their own drones
		query = "SELECT id, name, model, serial_number, status, user_id FROM drones WHERE user_id = $1 ORDER BY name"
		args = append(args, userID)
	}

	// Execute query
	type DroneResponse struct {
		ID           string `json:"id" db:"id"`
		Name         string `json:"name" db:"name"`
		Model        string `json:"model" db:"model"`
		SerialNumber string `json:"serial_number" db:"serial_number"`
		Status       string `json:"status" db:"status"`
		UserID       string `json:"user_id" db:"user_id"`
	}

	var drones []DroneResponse
	err = ec.db.Select(&drones, query, args...)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "erro ao buscar drones")
	}

	log.Printf("User %s (%s) listed %d drones", userID, userRole, len(drones))

	return c.JSON(fiber.Map{
		"data":  drones,
		"count": len(drones),
	})
}

// Helper functions

// validateStruct validates a struct using validation tags
func validateStruct(s interface{}) error {
	// In a real implementation, you would use a validation library like go-playground/validator
	// For now, we'll just return nil
	return nil
}

// auditAction logs an action for audit purposes
func (ec *ExampleController) auditAction(userID, action, resourceType, resourceID string, oldValue, newValue interface{}, c *fiber.Ctx) {
	// In a real implementation, you would insert into the permission_audit table
	log.Printf("AUDIT: User %s performed %s on %s:%s from IP %s", userID, action, resourceType, resourceID, c.IP())
}

// ExampleMiddlewareSetup demonstrates how to set up routes with the new middleware
func ExampleMiddlewareSetup(app *fiber.App, permissionService *Service) {
	// Create API group with authentication
	api := app.Group("/api")

	// Note: AuthMiddleware should be applied first to set user context
	// api.Use(middleware.AuthMiddleware(&cfg.JWT))

	// Drone routes with new permission middleware
	drones := api.Group("/drones")
	drones.Post("/", CreateDroneMiddleware(permissionService))                 // Create
	drones.Get("/", ListResourcesMiddleware(permissionService, ResourceDrone)) // List
	drones.Get("/:id", ReadDroneMiddleware(permissionService))                 // Read
	drones.Put("/:id", UpdateDroneMiddleware(permissionService))               // Update
	drones.Delete("/:id", DeleteDroneMiddleware(permissionService))            // Delete

	// Maintenance order routes
	maintenance := api.Group("/maintenance-orders")
	maintenance.Post("/", CreateMaintenanceOrderMiddleware(permissionService))
	maintenance.Get("/", ListResourcesMiddleware(permissionService, ResourceMaintenanceOrder))
	maintenance.Get("/:id", ReadMaintenanceOrderMiddleware(permissionService))
	maintenance.Put("/:id", UpdateMaintenanceOrderMiddleware(permissionService))
	maintenance.Delete("/:id", DeleteMaintenanceOrderMiddleware(permissionService))

	// Technical report routes
	reports := api.Group("/technical-reports")
	reports.Post("/", CreateTechnicalReportMiddleware(permissionService))
	reports.Get("/", ListResourcesMiddleware(permissionService, ResourceTechnicalReport))
	reports.Get("/:id", ReadTechnicalReportMiddleware(permissionService))
	reports.Put("/:id", UpdateTechnicalReportMiddleware(permissionService))
	reports.Delete("/:id", DeleteTechnicalReportMiddleware(permissionService))

	// Admin-only routes
	admin := api.Group("/admin", AdminOnlyMiddleware(permissionService))
	admin.Get("/users", ListResourcesMiddleware(permissionService, ResourceUser))
	admin.Post("/users/:userId/roles", PermissionMiddleware(permissionService, ResourceUser, ActionUpdate))
	admin.Get("/audit", PermissionMiddleware(permissionService, "system", "audit"))

	// Role management routes
	roles := api.Group("/roles", TechnicianOrAdminMiddleware(permissionService))
	roles.Get("/", PermissionMiddleware(permissionService, "role", ActionList))
	roles.Get("/:id", PermissionMiddleware(permissionService, "role", ActionRead))

	// Admin-only role management
	roles.Post("/", AdminOnlyMiddleware(permissionService))
	roles.Put("/:id", AdminOnlyMiddleware(permissionService))
	roles.Delete("/:id", AdminOnlyMiddleware(permissionService))
}
