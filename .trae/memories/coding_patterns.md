# EVO DRONES - Padrões de Codificação

## Convenções de Nomenclatura
- **Structs**: PascalCase (User, DroneItem, MaintenanceOrder)
- **Funções**: camelCase (getUserByID, createDrone, updateOrderStatus)
- **Constantes**: UPPER_SNAKE_CASE (JWT_SECRET, DB_MAX_CONNECTIONS)
- **Arquivos**: snake_case.go (auth_controller.go, drone_service.go)
- **Packages**: lowercase (controllers, services, repositories)
- **Endpoints**: kebab-case (/api/maintenance-orders, /api/reports)

## Padrões de Controller ATUALIZADOS
```go
func (c *ControllerName) HandlerName(ctx *fiber.Ctx) error {
    // 1. Extração de dados do contexto (auth)
    userID := ctx.Locals("user_id").(string)
    userRole := ctx.Locals("role").(string)
    
    // 2. Validação de entrada
    var input InputStruct
    if err := ctx.BodyParser(&input); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "dados inválidos")
    }
    
    // 3. Chamada do service com contexto de usuário
    result, err := c.service.MethodName(&input, userID, userRole)
    if err != nil {
        // Tratamento específico de erros de negócio
        switch err {
        case services.ErrNotFound:
            return fiber.NewError(fiber.StatusNotFound, "recurso não encontrado")
        case services.ErrUnauthorized:
            return fiber.NewError(fiber.StatusForbidden, "acesso negado")
        default:
            return fiber.NewError(fiber.StatusInternalServerError, "erro interno")
        }
    }
    
    // 4. Resposta com status apropriado
    return ctx.Status(fiber.StatusCreated).JSON(result) // ou StatusOK
}
```

## Padrões de Service ATUALIZADOS
```go
func (s *ServiceName) MethodName(input *InputStruct, userID, userRole string) (*Model, error) {
    // 1. Validação de autorização baseada em role
    if userRole != "admin" && userRole != "technician" {
        return nil, ErrUnauthorized
    }
    
    // 2. Validação de negócio e dependências
    if err := s.validateDependencies(input); err != nil {
        return nil, err
    }
    
    // 3. Criação do modelo com timestamps
    model := &Model{
        UserID:    userID,
        Field1:    input.Field1,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    // 4. Chamada do repository
    if err := s.repository.Create(model); err != nil {
        return nil, err
    }
    
    return model, nil
}
```

## Padrões de Repository ATUALIZADOS
```go
func (r *RepositoryName) Create(model *Model) error {
    query := `
        INSERT INTO table_name (
            field1, field2, user_id, created_at, updated_at
        ) VALUES (
            :field1, :field2, :user_id, :created_at, :updated_at
        ) RETURNING id
    `
    
    stmt, err := r.db.PrepareNamed(query)
    if err != nil {
        return err
    }
    defer stmt.Close()
    
    return stmt.QueryRow(model).Scan(&model.ID)
}

func (r *RepositoryName) GetByID(id string) (*ModelWithDetails, error) {
    query := `
        SELECT 
            m.*,
            u.name as "user.name", u.email as "user.email"
        FROM table_name m
        LEFT JOIN users u ON m.user_id = u.id
        WHERE m.id = $1
    `
    
    var result ModelWithDetails
    err := r.db.Get(&result, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // Retorna nil para não encontrado
        }
        return nil, err
    }
    
    return &result, nil
}
```

## Tratamento de Erros ATUALIZADO
- Usar erros customizados definidos em `/services/errors.go`
- Logging estruturado com zerolog
- Respostas HTTP com fiber.NewError
- Mensagens em português para o usuário
- Nunca expor detalhes internos ao cliente

```go
// Erros customizados em services/errors.go
var (
    ErrUserNotFound             = errors.New("user not found")
    ErrMaintenanceOrderNotFound = errors.New("maintenance order not found")
    ErrUnauthorized            = errors.New("unauthorized access")
    ErrInvalidInput            = errors.New("invalid input data")
)
```

## Padrões de Middleware ATUALIZADOS
```go
// AuthMiddleware - Validação de JWT
func AuthMiddleware(cfg *config.JWTConfig) fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return fiber.NewError(fiber.StatusUnauthorized, "token required")
        }
        
        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        claims, err := ValidateToken(tokenString, cfg)
        if err != nil {
            return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
        }
        
        c.Locals("user_id", claims.UserID)
        c.Locals("role", claims.Role)
        return c.Next()
    }
}

// AuthorizeResource - Autorização baseada em recursos
func AuthorizeResource(resourceType, paramName string, action Action, authService *services.AuthorizationService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userID := c.Locals("user_id").(string)
        userRole := c.Locals("role").(string)
        
        var resourceID string
        if paramName != "" {
            resourceID = c.Params(paramName)
        }
        
        if !authService.CanPerformAction(userID, userRole, resourceType, resourceID, action) {
            return fiber.NewError(fiber.StatusForbidden, "access denied")
        }
        
        return c.Next()
    }
}
```

## Padrões de Teste
```go
func TestFunctionName(t *testing.T) {
    // Arrange
    mockService := &MockService{}
    controller := NewController(mockService)
    
    // Act
    result, err := controller.Method(params)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    mockService.AssertExpectations(t)
}
```

## Validação de Entrada
- Sempre validar dados de entrada
- Usar struct tags para validação automática
- Sanitizar dados antes de processar

```go
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}
```

## Logging Patterns
```go
log.Info().Str("user_id", userID).Msg("User created successfully")
log.Error().Err(err).Str("operation", "create_user").Msg("Failed to create user")
```

## Database Patterns
- Sempre usar prepared statements
- Usar transações para operações múltiplas
- Implementar connection pooling
- UUIDs como primary keys

```go
// Transação
tx, err := r.db.Beginx()
if err != nil {
    return err
}
defer tx.Rollback()

// Operações...

return tx.Commit()
```