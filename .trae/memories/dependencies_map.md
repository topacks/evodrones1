# EVO DRONES - Mapa de Dependências

## Dependências Críticas do Backend

### Controllers → Services (ATUALIZADO)
- `auth_controller.go` → `auth_service.go`
  - Login, Register, RefreshToken, Logout
- `drone_controller.go` → `drone_service.go`
  - CreateDrone, GetDrones, UpdateDrone, DeleteDrone
- `maintenance_order_controller.go` → `maintenance_order_service.go`
  - CreateOrder, GetOrders, UpdateOrderStatus, AssignTechnician
- `technical_report_controller.go` → `technical_report_service.go`
  - Create, GetByID, GetByMaintenanceID, GetByUserID, Update, Delete, List
  - **IMPORTANTE**: Todos os métodos recebem userID e userRole do contexto

### Services → Repositories (ATUALIZADO)
- `auth_service.go` → `user_repository.go`
  - GetUserByEmail, CreateUser, UpdateUser, GetUserByID
- `drone_service.go` → `drone_repository.go`
  - CreateDrone, GetDroneByID, GetDronesByUser, UpdateDrone, DeleteDrone
- `maintenance_order_service.go` → `maintenance_order_repository.go`
  - CreateOrder, GetOrderByID, GetOrdersByStatus, UpdateOrderStatus
- `technical_report_service.go` → `technical_report_repository.go`, `maintenance_order_repository.go`, `user_repository.go`
  - **MÚLTIPLAS DEPENDÊNCIAS**: Validação cruzada de dados
  - CreateReport, GetReportByID, GetReportsByOrder, UpdateReport
  - Validações de existência de MaintenanceOrder e User

### Models Compartilhados
- `user.go` - Usado por:
  - auth_service.go (autenticação)
  - maintenance_order_service.go (atribuição de técnicos)
  - technical_report_service.go (criação de relatórios)
  
- `drone.go` - Usado por:
  - drone_service.go (gestão de drones)
  - maintenance_order_service.go (ordens de manutenção)
  
- `maintenance_order.go` - Usado por:
  - maintenance_order_service.go (gestão de ordens)
  - technical_report_service.go (vinculação de relatórios)
  
- `technical_report.go` - Usado por:
  - technical_report_service.go (gestão de relatórios)
  
- `drone_item.go` - Usado por:
  - drone_service.go (componentes do drone)
  - technical_report_service.go (itens reportados)

### Middleware Dependencies
- `auth.go` → `config/config.go` (JWT settings)
- Todos os controllers protegidos → `middleware/auth.go`
- Routes protegidas:
  - `/api/drones/*` (requer autenticação)
  - `/api/maintenance-orders/*` (requer autenticação + role)
  - `/api/technical-reports/*` (requer autenticação + role)

### Configuration Dependencies
- `config.go` → `config.yaml` (arquivo de configuração)
- `db.go` → `config.go` (configurações de banco)
- Todos os services → `config.go` (configurações gerais)

### Route Dependencies
- `routes.go` → todos os controllers
- `auth_routes.go` → `auth_controller.go`
- `technical_report_routes.go` → `technical_report_controller.go`
- Middleware aplicado em grupos de rotas

## Dependências do Frontend

### CSS Dependencies
- `main.css` → todas as páginas HTML
- Bootstrap 5 → layout e componentes
- Font Awesome → ícones
- Google Fonts (Inter, Roboto) → tipografia

### JavaScript Dependencies
- Bootstrap JS → funcionalidades interativas
- AOS (Animate On Scroll) → animações
- Custom JS → validações e API calls

### Component Dependencies
- `navbar/` → todas as páginas autenticadas
- `footer/` → todas as páginas
- `modals/` → páginas que requerem confirmações

## Regras de Modificação

### ❌ NUNCA ALTERAR
1. **Nomes de structs existentes**:
   - User, Drone, MaintenanceOrder, TechnicalReport, DroneItem
   
2. **Assinaturas de funções públicas**:
   - Métodos de controllers expostos em rotas
   - Métodos de services usados por controllers
   - Métodos de repositories usados por services
   
3. **Campos de modelos sem verificação**:
   - Verificar uso em queries SQL
   - Verificar serialização JSON
   - Verificar validações

### ✅ VERIFICAÇÃO OBRIGATÓRIA

#### Antes de modificar qualquer struct:
```bash
# Verificar uso em todo o backend
grep -r "StructName" backend/
grep -r "fieldName" backend/
```

#### Antes de modificar funções:
```bash
# Verificar chamadas da função
grep -r "FunctionName" backend/
# Verificar testes
grep -r "TestFunctionName" backend/
```

#### Antes de modificar rotas:
```bash
# Verificar definições de rotas
grep -r "endpoint-path" backend/routes/
# Verificar middleware aplicado
grep -r "middleware" backend/routes/
```

### ✅ VALIDAÇÃO PÓS-MODIFICAÇÃO

#### Testes obrigatórios:
```bash
# Executar todos os testes
go test ./...

# Verificar build
go build ./cmd/api

# Validar sintaxe
go vet ./...
go fmt ./...
```

#### Verificação de dependências:
```bash
# Verificar imports
go mod tidy
go mod verify

# Verificar dependências circulares
go list -deps ./...
```

## Fluxo de Dados Críticos

### Autenticação
1. `auth_controller.go` → `auth_service.go` → `user_repository.go`
2. JWT gerado em `auth_service.go`
3. Validação em `middleware/auth.go`
4. Contexto de usuário propagado para controllers

### Criação de Ordem de Manutenção
1. `maintenance_order_controller.go` → `maintenance_order_service.go`
2. Validação de drone em `drone_repository.go`
3. Validação de usuário em `user_repository.go`
4. Criação em `maintenance_order_repository.go`

### Relatório Técnico
1. `technical_report_controller.go` → `technical_report_service.go`
2. Validação de ordem em `maintenance_order_repository.go`
3. Validação de técnico em `user_repository.go`
4. Criação em `technical_report_repository.go`

## Pontos de Atenção

### Segurança
- Middleware de autenticação em todas as rotas protegidas
- Validação de roles para operações sensíveis
- Prepared statements em todas as queries
- Sanitização de inputs

### Performance
- Connection pooling configurado
- Índices de banco otimizados
- Queries eficientes com LIMIT/OFFSET
- Cache de configurações

### Manutenibilidade
- Separação clara de responsabilidades
- Interfaces bem definidas
- Testes unitários abrangentes
- Documentação atualizada