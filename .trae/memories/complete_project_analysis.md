# EVO DRONES - Análise Completa do Projeto

## 📋 Visão Geral do Projeto

**EVO DRONES** é uma plataforma completa de gestão de drones e manutenção, desenvolvida com arquitetura moderna e padrões de segurança robustos.

### Tecnologias Principais
- **Backend**: Go 1.21+ com Fiber Framework
- **Frontend**: Bootstrap 5 + JavaScript Vanilla (migração para Svelte planejada)
- **Database**: PostgreSQL 14+ com SQLX
- **Autenticação**: JWT com middleware customizado
- **Containerização**: Docker Compose
- **Autorização**: Sistema baseado em recursos e ações

---

## 🏗️ Arquitetura do Sistema

### Padrão Arquitetural
- **Clean Architecture**: Controller → Service → Repository
- **Dependency Injection**: Injeção manual nos construtores
- **Middleware Chain**: Auth → Authorization → Handler
- **Error Handling**: Erros customizados centralizados

### Estrutura de Diretórios

```
evodrones/
├── backend/                    # API Go com Fiber
│   ├── cmd/api/               # Ponto de entrada
│   ├── config/                # Configurações (Viper)
│   ├── controllers/           # Handlers HTTP
│   ├── middleware/            # Auth, CORS, Authorization
│   ├── models/                # Estruturas de dados
│   ├── repositories/          # Camada de dados
│   ├── services/              # Lógica de negócio
│   ├── routes/                # Definição de rotas
│   ├── db/                    # Conexão e migrações
│   └── static/                # Frontend servido pelo Go
├── database/                  # Scripts SQL
├── docker-compose.yml         # Orquestração
└── .trae/                     # Documentação e regras
```

---

## 🔐 Sistema de Autenticação e Autorização

### Roles do Sistema
- **admin**: Acesso total ao sistema
- **technician**: Gestão de manutenções e relatórios
- **client**: Gestão de próprios drones e ordens

### Fluxo de Autenticação
1. **Login** → Validação de credenciais
2. **Token Generation** → Access Token (15min) + Refresh Token (24h)
3. **Middleware Chain** → AuthMiddleware → AuthorizeResource
4. **Context Injection** → userID e userRole disponíveis nos controllers

### Sistema de Autorização
```go
// Recursos do Sistema
ResourceTypeDrone
ResourceTypeMaintenanceOrder
ResourceTypeTechnicalReport
ResourceTypeUser

// Ações Disponíveis
ActionCreate, ActionRead, ActionUpdate, ActionDelete, ActionList
```

---

## 📊 Modelos de Dados

### Entidades Principais

#### User
```go
type User struct {
    ID           string    `db:"id" json:"id"`
    Email        string    `db:"email" json:"email"`
    PasswordHash string    `db:"password_hash" json:"-"`
    Name         string    `db:"name" json:"name"`
    Phone        string    `db:"phone" json:"phone"`
    Role         string    `db:"role" json:"role"` // admin, technician, client
    CreatedAt    time.Time `db:"created_at" json:"created_at"`
    UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
```

#### Drone
```go
type Drone struct {
    ID        string    `db:"id" json:"id"`
    UserID    string    `db:"user_id" json:"user_id"`
    Model     string    `db:"model" json:"model"`
    Serial    string    `db:"serial" json:"serial"`
    Status    string    `db:"status" json:"status"` // active, maintenance, inactive
    LastCheck time.Time `db:"last_check" json:"last_check"`
    NextCheck time.Time `db:"next_check" json:"next_check"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
```

#### MaintenanceOrder
```go
type MaintenanceOrder struct {
    ID          string                 `json:"id" db:"id"`
    DroneID     string                 `json:"drone_id" db:"drone_id"`
    UserID      string                 `json:"user_id" db:"user_id"`
    Status      MaintenanceOrderStatus `json:"status" db:"status"`
    Description string                 `json:"description" db:"description"`
    StartDate   time.Time              `json:"start_date" db:"start_date"`
    EndDate     *time.Time             `json:"end_date,omitempty" db:"end_date"`
    CreatedAt   time.Time              `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// Status: pending, in_progress, completed, cancelled
```

### Relacionamentos
- **User** → **Drone** (1:N)
- **User** → **MaintenanceOrder** (1:N)
- **Drone** → **MaintenanceOrder** (1:N)
- **MaintenanceOrder** → **TechnicalReport** (1:N)
- **Drone** → **DroneItem** (1:N)

---

## 🛠️ Padrões de Implementação

### Controller Pattern
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
        // Tratamento específico de erros
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
    return ctx.Status(fiber.StatusCreated).JSON(result)
}
```

### Service Pattern
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

### Repository Pattern
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
```

---

## 🌐 API Endpoints

### Estrutura de Rotas
```
POST   /api/auth/login           - Autenticação
POST   /api/auth/register        - Registro de usuário
POST   /api/auth/refresh         - Renovação de token

POST   /api/drones               - Criar drone
GET    /api/drones               - Listar drones
GET    /api/drones/:id           - Obter drone por ID
PUT    /api/drones/:id           - Atualizar drone
DELETE /api/drones/:id           - Deletar drone
GET    /api/drones/user          - Drones do usuário
GET    /api/drones/status/:status - Drones por status

POST   /api/maintenance-orders   - Criar ordem
GET    /api/maintenance-orders   - Listar ordens
GET    /api/maintenance-orders/:id - Obter ordem por ID
PUT    /api/maintenance-orders/:id - Atualizar ordem
DELETE /api/maintenance-orders/:id - Deletar ordem

POST   /api/technical-reports    - Criar relatório
GET    /api/technical-reports    - Listar relatórios
GET    /api/technical-reports/:id - Obter relatório
```

### Middleware Chain
```go
// Todas as rotas protegidas seguem este padrão:
api.Use(middleware.AuthMiddleware(&cfg.JWT))
api.Use(middleware.AuthorizeResource(resourceType, paramName, action, authService))
```

---

## 🗄️ Banco de Dados

### Configuração
- **PostgreSQL 14+**
- **Connection Pool**: 25 max connections, 5min lifetime
- **Extensions**: uuid-ossp para UUIDs
- **Migrations**: Scripts SQL em `/database/init.sql`

### Tabelas Principais
```sql
-- Usuários com roles
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    role user_role NOT NULL DEFAULT 'client',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Drones com relacionamento ao usuário
CREATE TABLE drones (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    serial_number VARCHAR(100) NOT NULL,
    color VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Ordens de manutenção
CREATE TABLE maintenance_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    drone_id UUID NOT NULL REFERENCES drones(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status order_status NOT NULL DEFAULT 'created',
    defect_description TEXT NOT NULL,
    broken_parts JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

---

## 🎨 Frontend

### Tecnologias
- **Bootstrap 5**: Framework CSS
- **JavaScript Vanilla**: Lógica de interação
- **Font Awesome**: Ícones
- **AOS**: Animações on scroll

### Estrutura de Páginas
```
static/
├── assets/
│   ├── css/main.css           # Estilos principais
│   └── js/                    # Scripts JavaScript
├── components/
│   ├── navbar/                # Barra de navegação
│   ├── footer/                # Rodapé
│   └── modals/                # Modais reutilizáveis
├── pages/
│   ├── auth/                  # Login e registro
│   ├── dashboard/             # Dashboards por role
│   └── public/                # Páginas públicas
└── index.html                 # Página inicial
```

### Padrão de Acesso
- **Servidor Go**: Porta 8080
- **Frontend**: Servido pelo Go Fiber
- **URLs**: http://localhost:8080/pages/...
- **API**: http://localhost:8080/api/...

---

## 🔧 Configuração e Deploy

### Configuração (config.yaml)
```yaml
server:
  port: "8080"
  read_timeout: "10s"
  write_timeout: "10s"

database:
  host: "localhost"
  port: "5432"
  user: "evodrones"
  password: "evodrones_password"
  database: "evodrones_db"
  ssl_mode: "disable"

jwt:
  secret_key: "evodrones_jwt_secret_key"
  access_token_ttl: "15m"
  refresh_token_ttl: "24h"
  signing_method: "HS256"
```

### Docker Compose
```yaml
services:
  postgres:
    image: postgres:14-alpine
    environment:
      POSTGRES_USER: evodrones
      POSTGRES_PASSWORD: evodrones_password
      POSTGRES_DB: evodrones_db
    ports:
      - "5432:5432"

  backend:
    build: ./backend
    ports:
      - "8080:8080"
    depends_on:
      - postgres
```

### Inicialização
```bash
# Desenvolvimento
./start.sh

# Docker
docker-compose up -d
```

---

## 📋 Funcionalidades Implementadas

### ✅ Autenticação e Autorização
- [x] Sistema de login/registro
- [x] JWT com refresh tokens
- [x] Middleware de autenticação
- [x] Sistema de roles (admin, technician, client)
- [x] Autorização baseada em recursos

### ✅ Gestão de Drones
- [x] CRUD completo de drones
- [x] Relacionamento com usuários
- [x] Status de drones (active, maintenance, inactive)
- [x] Filtros por status e usuário

### ✅ Ordens de Manutenção
- [x] CRUD de ordens de manutenção
- [x] Workflow de status (pending → in_progress → completed)
- [x] Relacionamento drone-usuário
- [x] Descrição de defeitos

### ✅ Relatórios Técnicos
- [x] Criação de relatórios por técnicos
- [x] Vinculação com ordens de manutenção
- [x] Sistema de permissões

### ✅ Interface Web
- [x] Dashboard por role
- [x] Páginas de autenticação
- [x] Interface responsiva
- [x] Componentes reutilizáveis

---

## 🚀 Roadmap e Melhorias

### 🔄 Em Desenvolvimento
- [ ] Migração para Svelte no frontend
- [ ] Sistema de notificações
- [ ] Upload de imagens/documentos
- [ ] Relatórios em PDF

### 🎯 Próximas Features
- [ ] Dashboard analytics
- [ ] Sistema de custos
- [ ] Integração com APIs externas
- [ ] Mobile app

---

## 🛡️ Segurança

### Medidas Implementadas
- **Prepared Statements**: Prevenção de SQL Injection
- **Password Hashing**: bcrypt para senhas
- **JWT Tokens**: Autenticação stateless
- **CORS**: Configuração adequada
- **Input Validation**: Validação em todos os endpoints
- **Role-Based Access**: Controle granular de permissões

### Boas Práticas
- Nunca expor detalhes de erros internos
- Validação de autorização em todas as camadas
- Logs estruturados para auditoria
- Timeouts configurados para conexões

---

## 📚 Convenções e Padrões

### Nomenclatura
- **Structs**: PascalCase (User, DroneItem)
- **Funções**: camelCase (getUserByID, createDrone)
- **Constantes**: UPPER_SNAKE_CASE (JWT_SECRET)
- **Arquivos**: snake_case.go (auth_controller.go)
- **Endpoints**: kebab-case (/api/maintenance-orders)

### Commits
- `feat:` - Nova funcionalidade
- `fix:` - Correção de bug
- `refactor:` - Refatoração
- `test:` - Testes
- `docs:` - Documentação
- `auth:` - Autenticação/Autorização

---

## 🔍 Dependências Críticas

### Backend Dependencies
```go
// Core
github.com/gofiber/fiber/v2     // Web framework
github.com/jmoiron/sqlx         // Database toolkit
github.com/lib/pq               // PostgreSQL driver

// Auth
github.com/golang-jwt/jwt/v5    // JWT tokens
golang.org/x/crypto             // Password hashing

// Config
github.com/spf13/viper          // Configuration
github.com/google/uuid          // UUID generation

// Testing
github.com/stretchr/testify     // Test assertions
```

### Interdependências
- Controllers → Services → Repositories
- Middleware → Config (JWT settings)
- Services → Multiple Repositories (validação cruzada)
- Models compartilhados entre camadas

---

## 📖 Conclusão

O projeto **EVO DRONES** implementa uma arquitetura sólida e escalável para gestão de drones e manutenção, com foco em:

- **Segurança**: Sistema robusto de autenticação e autorização
- **Escalabilidade**: Arquitetura limpa e bem estruturada
- **Manutenibilidade**: Padrões consistentes e documentação completa
- **Performance**: Otimizações de banco e connection pooling
- **Usabilidade**: Interface moderna e responsiva

A base está sólida para expansões futuras e novas funcionalidades.