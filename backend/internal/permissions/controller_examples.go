package permissions

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Exemplos práticos de como integrar o sistema RBAC com controllers existentes

// DroneControllerWithRBAC demonstra como integrar RBAC com o controller de drones
type DroneControllerWithRBAC struct {
	db                *sqlx.DB
	permissionService *Service
}

// NewDroneControllerWithRBAC cria uma nova instância do controller com RBAC
func NewDroneControllerWithRBAC(db *sqlx.DB, permissionService *Service) *DroneControllerWithRBAC {
	return &DroneControllerWithRBAC{
		db:                db,
		permissionService: permissionService,
	}
}

// CreateDrone exemplo de criação de drone com verificação RBAC
func (c *DroneControllerWithRBAC) CreateDrone(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)

	// Estrutura de entrada
	type CreateDroneInput struct {
		Name         string `json:"name" validate:"required,min=3,max=100"`
		Model        string `json:"model" validate:"required,min=2,max=50"`
		SerialNumber string `json:"serial_number" validate:"required,min=5,max=50"`
		Description  string `json:"description,omitempty" validate:"max=500"`
	}

	var input CreateDroneInput
	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "dados inválidos")
	}

	// Verificação adicional de permissão (opcional, já verificada no middleware)
	allowed, err := c.permissionService.HasPermission(
		context.Background(),
		userID,
		ResourceDrone,
		ActionCreate,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "erro interno")
	}
	if !allowed {
		return fiber.NewError(fiber.StatusForbidden, "acesso negado")
	}

	// Lógica de criação do drone
	query := `
		INSERT INTO drones (name, model, serial_number, description, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at
	`

	var droneID string
	var createdAt string
	err = c.db.QueryRowContext(
		context.Background(),
		query,
		input.Name,
		input.Model,
		input.SerialNumber,
		input.Description,
		userID,
	).Scan(&droneID, &createdAt)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "erro ao criar drone")
	}

	// Resposta
	response := map[string]interface{}{
		"id":            droneID,
		"name":          input.Name,
		"model":         input.Model,
		"serial_number": input.SerialNumber,
		"description":   input.Description,
		"user_id":       userID,
		"created_at":    createdAt,
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}

// GetDrone exemplo de leitura de drone com verificação de propriedade
func (c *DroneControllerWithRBAC) GetDrone(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)
	droneID := ctx.Params("id")

	// Construir contexto de permissão
	permissionContext := map[string]interface{}{
		"user_id":     userID,
		"resource_id": droneID,
		"client_ip":   ctx.IP(),
	}

	// Verificar permissão com contexto
	allowed, err := c.permissionService.CheckPermission(
		context.Background(),
		userID,
		ResourceDrone,
		ActionRead,
		permissionContext,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "erro interno")
	}
	if !allowed {
		return fiber.NewError(fiber.StatusForbidden, "acesso negado")
	}

	// Buscar drone
	query := `
		SELECT 
			d.id, d.name, d.model, d.serial_number, d.description, 
			d.user_id, d.created_at, d.updated_at,
			u.name as "user.name", u.email as "user.email"
		FROM drones d
		LEFT JOIN users u ON d.user_id = u.id
		WHERE d.id = $1
	`

	type DroneResponse struct {
		ID           string `json:"id" db:"id"`
		Name         string `json:"name" db:"name"`
		Model        string `json:"model" db:"model"`
		SerialNumber string `json:"serial_number" db:"serial_number"`
		Description  string `json:"description" db:"description"`
		UserID       string `json:"user_id" db:"user_id"`
		CreatedAt    string `json:"created_at" db:"created_at"`
		UpdatedAt    string `json:"updated_at" db:"updated_at"`
		User         *struct {
			Name  string `json:"name" db:"user.name"`
			Email string `json:"email" db:"user.email"`
		} `json:"user,omitempty"`
	}

	var drone DroneResponse
	err = c.db.GetContext(context.Background(), &drone, query, droneID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "drone não encontrado")
	}

	return ctx.JSON(drone)
}

// ListDrones exemplo de listagem com filtragem baseada em permissões
func (c *DroneControllerWithRBAC) ListDrones(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)

	// Verificar permissão de listagem
	allowed, err := c.permissionService.HasPermission(
		context.Background(),
		userID,
		ResourceDrone,
		ActionList,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "erro interno")
	}
	if !allowed {
		return fiber.NewError(fiber.StatusForbidden, "acesso negado")
	}

	// Verificar se usuário é admin (pode ver todos os drones)
	userRoles, err := c.permissionService.GetUserRoles(context.Background(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "erro interno")
	}

	isAdmin := false
	for _, role := range userRoles {
		if role.ID == RoleAdmin {
			isAdmin = true
			break
		}
	}

	// Construir query baseada em permissões
	var query string
	var args []interface{}

	if isAdmin {
		// Admin vê todos os drones
		query = `
			SELECT 
				d.id, d.name, d.model, d.serial_number, 
				d.user_id, d.created_at,
				u.name as "user.name", u.email as "user.email"
			FROM drones d
			LEFT JOIN users u ON d.user_id = u.id
			ORDER BY d.created_at DESC
		`
	} else {
		// Outros usuários veem apenas seus próprios drones
		query = `
			SELECT 
				d.id, d.name, d.model, d.serial_number, 
				d.user_id, d.created_at,
				u.name as "user.name", u.email as "user.email"
			FROM drones d
			LEFT JOIN users u ON d.user_id = u.id
			WHERE d.user_id = $1
			ORDER BY d.created_at DESC
		`
		args = append(args, userID)
	}

	// Paginação
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page - 1) * limit

	query += " LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	type DroneListItem struct {
		ID           string `json:"id" db:"id"`
		Name         string `json:"name" db:"name"`
		Model        string `json:"model" db:"model"`
		SerialNumber string `json:"serial_number" db:"serial_number"`
		UserID       string `json:"user_id" db:"user_id"`
		CreatedAt    string `json:"created_at" db:"created_at"`
		User         *struct {
			Name  string `json:"name" db:"user.name"`
			Email string `json:"email" db:"user.email"`
		} `json:"user,omitempty"`
	}

	var drones []DroneListItem
	err = c.db.SelectContext(context.Background(), &drones, query, args...)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "erro ao buscar drones")
	}

	// Resposta com metadados
	response := map[string]interface{}{
		"data": drones,
		"meta": map[string]interface{}{
			"page":     page,
			"limit":    limit,
			"total":    len(drones),
			"is_admin": isAdmin,
			"user_id":  userID,
		},
	}

	return ctx.JSON(response)
}

// UpdateDrone exemplo de atualização com verificação de propriedade
func (c *DroneControllerWithRBAC) UpdateDrone(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)
	droneID := ctx.Params("id")

	// Estrutura de entrada para atualização
	type UpdateDroneInput struct {
		Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
		Model       *string `json:"model,omitempty" validate:"omitempty,min=2,max=50"`
		Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
	}

	var input UpdateDroneInput
	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "dados inválidos")
	}

	// Verificar se drone existe e obter proprietário
	var droneOwnerID string
	err := c.db.GetContext(
		context.Background(),
		&droneOwnerID,
		"SELECT user_id FROM drones WHERE id = $1",
		droneID,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "drone não encontrado")
	}

	// Construir contexto de permissão com informações de propriedade
	permissionContext := map[string]interface{}{
		"user_id":        userID,
		"resource_id":    droneID,
		"resource_owner": droneOwnerID,
		"client_ip":      ctx.IP(),
	}

	// Verificar permissão com contexto
	allowed, err := c.permissionService.CheckPermission(
		context.Background(),
		userID,
		ResourceDrone,
		ActionUpdate,
		permissionContext,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "erro interno")
	}
	if !allowed {
		return fiber.NewError(fiber.StatusForbidden, "acesso negado")
	}

	// Construir query de atualização dinamicamente
	updateFields := []string{}
	args := []interface{}{}
	argIndex := 1

	if input.Name != nil {
		updateFields = append(updateFields, "name = $"+strconv.Itoa(argIndex))
		args = append(args, *input.Name)
		argIndex++
	}

	if input.Model != nil {
		updateFields = append(updateFields, "model = $"+strconv.Itoa(argIndex))
		args = append(args, *input.Model)
		argIndex++
	}

	if input.Description != nil {
		updateFields = append(updateFields, "description = $"+strconv.Itoa(argIndex))
		args = append(args, *input.Description)
		argIndex++
	}

	if len(updateFields) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "nenhum campo para atualizar")
	}

	// Adicionar updated_at e WHERE clause
	updateFields = append(updateFields, "updated_at = NOW()")
	args = append(args, droneID)

	query := "UPDATE drones SET " +
		string(updateFields[0])
	for i := 1; i < len(updateFields); i++ {
		query += ", " + updateFields[i]
	}
	query += " WHERE id = $" + strconv.Itoa(argIndex) + " RETURNING updated_at"

	var updatedAt string
	err = c.db.QueryRowContext(context.Background(), query, args...).Scan(&updatedAt)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "erro ao atualizar drone")
	}

	// Resposta
	response := map[string]interface{}{
		"id":         droneID,
		"updated_at": updatedAt,
		"message":    "Drone atualizado com sucesso",
	}

	return ctx.JSON(response)
}

// DeleteDrone exemplo de exclusão com verificação de propriedade
func (c *DroneControllerWithRBAC) DeleteDrone(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)
	droneID := ctx.Params("id")

	// Verificar se drone existe e obter proprietário
	var droneOwnerID string
	err := c.db.GetContext(
		context.Background(),
		&droneOwnerID,
		"SELECT user_id FROM drones WHERE id = $1",
		droneID,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "drone não encontrado")
	}

	// Construir contexto de permissão
	permissionContext := map[string]interface{}{
		"user_id":        userID,
		"resource_id":    droneID,
		"resource_owner": droneOwnerID,
		"client_ip":      ctx.IP(),
	}

	// Verificar permissão
	allowed, err := c.permissionService.CheckPermission(
		context.Background(),
		userID,
		ResourceDrone,
		ActionDelete,
		permissionContext,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "erro interno")
	}
	if !allowed {
		return fiber.NewError(fiber.StatusForbidden, "acesso negado")
	}

	// Excluir drone
	result, err := c.db.ExecContext(
		context.Background(),
		"DELETE FROM drones WHERE id = $1",
		droneID,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "erro ao excluir drone")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "drone não encontrado")
	}

	return ctx.JSON(map[string]interface{}{
		"message": "Drone excluído com sucesso",
		"id":      droneID,
	})
}

// SetupDroneRoutesWithRBAC configura as rotas de drone com RBAC
func SetupDroneRoutesWithRBAC(app *fiber.App, db *sqlx.DB, permissionService *Service) {
	controller := NewDroneControllerWithRBAC(db, permissionService)

	drones := app.Group("/api/drones")

	// Aplicar middleware de autenticação em todas as rotas
	// drones.Use(middleware.AuthMiddleware(&cfg.JWT)) // Aplicado globalmente

	// Rotas com middleware RBAC específico
	drones.Post("/",
		CreateDroneMiddleware(permissionService),
		controller.CreateDrone)

	drones.Get("/",
		ListResourcesMiddleware(permissionService, "drone"),
		controller.ListDrones)

	drones.Get("/:id",
		ReadDroneMiddleware(permissionService),
		controller.GetDrone)

	drones.Put("/:id",
		UpdateDroneMiddleware(permissionService),
		controller.UpdateDrone)

	drones.Delete("/:id",
		DeleteDroneMiddleware(permissionService),
		controller.DeleteDrone)
}
