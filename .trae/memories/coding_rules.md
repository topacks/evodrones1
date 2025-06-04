# EVO DRONES - Regras de Codificação

## 1. Padrões de Autenticação e Autorização

### 1.1 Middleware Chain Obrigatório
```go
// Todas as rotas protegidas DEVEM seguir este padrão:
routes.Use(middleware.AuthMiddleware(&cfg.JWT))
routes.Use(middleware.AuthorizeResource(resourceType, paramName, action, authService))
```

### 1.2 Contexto de Usuário nos Controllers
```go
// SEMPRE extrair userID e userRole do contexto:
userID := ctx.Locals("user_id").(string)
userRole := ctx.Locals("role").(string)

// SEMPRE passar para o service:
result, err := c.service.MethodName(input, userID, userRole)
```

### 1.3 Validação de Autorização nos Services
```go
// SEMPRE validar autorização no início dos métodos de service:
if userRole != "admin" && userRole != "technician" {
    return nil, ErrUnauthorized
}

// Para operações específicas do usuário:
if userRole != "admin" && userID != resourceUserID {
    return nil, ErrUnauthorized
}
```

## 2. Padrões de Estrutura de Dados

### 2.1 Models com Validação
```go
type CreateInputStruct struct {
    Field1 string `json:"field1" validate:"required,min=5,max=100"`
    Field2 string `json:"field2" validate:"required,min=10,max=1000"`
}

type UpdateInputStruct struct {
    Field1 *string `json:"field1,omitempty" validate:"omitempty,min=5,max=100"`
    Field2 *string `json:"field2,omitempty" validate:"omitempty,min=10,max=1000"`
}
```

### 2.2 Models com Relacionamentos
```go
type ModelWithDetails struct {
    Model
    RelatedModel *RelatedModel `json:"related_model,omitempty"`
    User         *User         `json:"user,omitempty"`
}
```

## 3. Padrões de Repository

### 3.1 Queries com JOIN
```go
// SEMPRE usar aliases para evitar conflitos:
query := `
    SELECT 
        m.*,
        u.name as "user.name", u.email as "user.email",
        r.field as "related.field"
    FROM main_table m
    LEFT JOIN users u ON m.user_id = u.id
    LEFT JOIN related_table r ON m.related_id = r.id
    WHERE m.id = $1
`
```

### 3.2 Prepared Statements
```go
// SEMPRE usar PrepareNamed para INSERTs:
stmt, err := r.db.PrepareNamed(query)
if err != nil {
    return err
}
defer stmt.Close()

return stmt.QueryRow(model).Scan(&model.ID)
```

## 4. Tratamento de Erros Específicos

### 4.1 Erros de Negócio
```go
// Em services/errors.go:
var (
    ErrUserNotFound             = errors.New("user not found")
    ErrMaintenanceOrderNotFound = errors.New("maintenance order not found")
    ErrTechnicalReportNotFound  = errors.New("technical report not found")
    ErrUnauthorized            = errors.New("unauthorized access")
    ErrInvalidInput            = errors.New("invalid input data")
)
```

### 4.2 Tratamento nos Controllers
```go
// SEMPRE usar switch para erros específicos:
switch err {
case services.ErrUserNotFound:
    return fiber.NewError(fiber.StatusNotFound, "usuário não encontrado")
case services.ErrMaintenanceOrderNotFound:
    return fiber.NewError(fiber.StatusNotFound, "ordem de manutenção não encontrada")
case services.ErrUnauthorized:
    return fiber.NewError(fiber.StatusForbidden, "acesso negado")
default:
    return fiber.NewError(fiber.StatusInternalServerError, "erro interno")
}
```

## 5. Padrões de Rotas

### 5.1 Estrutura de Grupos
```go
// Padrão para agrupamento de rotas:
resources := app.Group("/api/resources", middleware.AuthMiddleware(&cfg.JWT))

// Rotas com autorização específica:
resources.Post("/", 
    middleware.AuthorizeResource(services.ResourceType, "", services.ActionCreate, authService), 
    controller.Create)
resources.Get("/:id", 
    middleware.AuthorizeResource(services.ResourceType, "id", services.ActionRead, authService), 
    controller.GetByID)
```

### 5.2 Endpoints Padronizados
```
POST   /api/resources          - Create
GET    /api/resources          - List
GET    /api/resources/:id      - GetByID
PUT    /api/resources/:id      - Update
DELETE /api/resources/:id      - Delete

// Endpoints específicos:
GET    /api/resources/maintenance/:maintenanceID  - GetByMaintenanceID
GET    /api/resources/user/:userID               - GetByUserID
```

## 6. Validações Obrigatórias

### 6.1 Antes de Modificar Código
- [ ] Verificar se struct é usado em outros services
- [ ] Verificar se campos são referenciados em queries SQL
- [ ] Verificar se mudança afeta serialização JSON
- [ ] Verificar se validações existentes são suficientes

### 6.2 Antes de Adicionar Endpoints
- [ ] Definir middleware de autorização apropriado
- [ ] Criar validações de entrada
- [ ] Implementar tratamento de erros específicos
- [ ] Documentar no mapa de dependências

## 7. Regras de Segurança

### 7.1 NUNCA Expor
- Detalhes de erros internos ao cliente
- Informações de outros usuários sem autorização
- Dados sensíveis em logs
- Queries SQL em mensagens de erro

### 7.2 SEMPRE Validar
- Autorização em todos os endpoints protegidos
- Entrada de dados com tags de validação
- Existência de recursos relacionados
- Propriedade de recursos por usuário

## 8. Padrões de Teste

### 8.1 Estrutura de Testes
```go
func TestServiceMethod(t *testing.T) {
    // Setup
    service := setupTestService()
    
    // Test cases
    tests := []struct {
        name     string
        input    InputStruct
        userID   string
        userRole string
        want     *Model
        wantErr  error
    }{
        // casos de teste
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := service.Method(tt.input, tt.userID, tt.userRole)
            assert.Equal(t, tt.wantErr, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

## 9. Convenções de Commit

### 9.1 Prefixos Obrigatórios
- `feat:` - Nova funcionalidade
- `fix:` - Correção de bug
- `refactor:` - Refatoração sem mudança de funcionalidade
- `test:` - Adição ou modificação de testes
- `docs:` - Atualização de documentação
- `auth:` - Mudanças relacionadas à autenticação/autorização

### 9.2 Exemplos
```
feat: adicionar endpoint GetByUserID para technical reports
fix: corrigir validação de autorização em technical_report_service
refactor: extrair userRole do contexto em todos os controllers
auth: implementar middleware AuthorizeResource
```

## 10. Checklist de Implementação

### 10.1 Para Novos Controllers
- [ ] Extrair userID e userRole do contexto
- [ ] Validar entrada com BodyParser
- [ ] Chamar service com parâmetros corretos
- [ ] Implementar tratamento de erros específicos
- [ ] Retornar status HTTP apropriado

### 10.2 Para Novos Services
- [ ] Validar autorização baseada em role
- [ ] Verificar dependências (user, maintenance order, etc.)
- [ ] Implementar lógica de negócio
- [ ] Chamar repository com dados validados
- [ ] Retornar erros customizados apropriados

### 10.3 Para Novos Repositories
- [ ] Usar prepared statements para INSERTs
- [ ] Implementar JOINs com aliases corretos
- [ ] Tratar sql.ErrNoRows adequadamente
- [ ] Retornar nil para recursos não encontrados
- [ ] Usar timestamps automáticos

## 11. Regras de Servidor

### 11.1 Inicialização do Servidor
- SEMPRE usar o script start.sh para iniciar o servidor Go
- NUNCA usar servidores HTTP alternativos (como servidores baseados em Node.js ou outros)
- Acessar o frontend através do servidor Go em http://localhost:8080
- O servidor Go é responsável por servir todos os arquivos estáticos

### 11.2 Desenvolvimento Frontend
- Usar o servidor Go para desenvolvimento (./start.sh)
- Acessar as páginas através do servidor Go em http://localhost:8080
- NUNCA usar servidores de desenvolvimento alternativos
- Todas as requisições de API devem ser feitas para o servidor Go

### 11.3 Proibições Explícitas
- NUNCA usar servidores HTTP baseados em Node.js, Python ou qualquer outra tecnologia que não seja Go
- NUNCA iniciar servidores HTTP alternativos em portas como 3000, 8000 ou similares
- NUNCA usar comandos como "npm run dev" ou similares para iniciar servidores de desenvolvimento
- SEMPRE usar o servidor Go para servir todos os arquivos estáticos e APIs

## 12. Convenções de Nomenclatura

### 12.1 Estruturas e Tipos
- **Structs**: PascalCase (User, DroneItem, MaintenanceOrder)
- **Funções**: camelCase (getUserByID, createDrone, updateOrderStatus)
- **Constantes**: UPPER_SNAKE_CASE (JWT_SECRET, DB_MAX_CONNECTIONS)
- **Arquivos**: snake_case.go (auth_controller.go, drone_service.go)
- **Packages**: lowercase (controllers, services, repositories)
- **Endpoints**: kebab-case (/api/maintenance-orders, /api/reports)

### 12.2 Database
- **Tabelas**: snake_case (maintenance_orders, technical_reports)
- **Colunas**: snake_case (created_at, user_id, drone_id)
- **Índices**: idx_table_column (idx_users_email, idx_drones_user_id)
- **Foreign Keys**: fk_table_reftable (fk_orders_users, fk_reports_orders)

## 13. Padrões de Documentação

### 13.1 Comentários de Código
```go
// UserRepository implementa operações de banco de dados para usuários
type UserRepository struct {
    db *sqlx.DB
}

// GetByID retorna um usuário pelo seu ID ou nil se não encontrado
// Retorna erro em caso de falha no banco de dados
func (r *UserRepository) GetByID(id string) (*User, error) {
    // implementação
}
```

### 13.2 Documentação de API
```go
// Create godoc
// @Summary Criar novo recurso
// @Description Cria um novo recurso com os dados fornecidos
// @Tags recursos
// @Accept json
// @Produce json
// @Param input body CreateResourceInput true "Dados do recurso"
// @Success 201 {object} Resource
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/resources [post]
func (c *ResourceController) Create(ctx *fiber.Ctx) error {
    // implementação
}
```

## 14. Padrões de Configuração

### 14.1 Variáveis de Ambiente
```go
// Prioridade de configuração:
// 1. Variáveis de ambiente
// 2. Arquivo config.yaml
// 3. Valores padrão

viper.SetDefault("server.port", "8080")
viper.SetDefault("database.host", "localhost")
viper.AutomaticEnv()
```

### 14.2 Configuração de Logs
```go
// Configuração de logs estruturados
log := zerolog.New(os.Stdout).With().Timestamp().Logger()
log.Info().Str("component", "api").Msg("Iniciando servidor")
```

## 15. Padrões de Transação

### 15.1 Uso de Transações
```go
// Para operações que afetam múltiplas tabelas:
db.WithTransaction(func(tx *sqlx.Tx) error {
    // Operação 1
    if err := repository1.CreateWithTx(tx, model1); err != nil {
        return err
    }
    
    // Operação 2
    if err := repository2.CreateWithTx(tx, model2); err != nil {
        return err
    }
    
    return nil
})
```