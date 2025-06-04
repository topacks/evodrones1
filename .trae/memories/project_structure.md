# EVO DRONES - Estrutura do Projeto

## Arquitetura Geral
- **Backend**: Go 1.21+ com Fiber framework
- **Frontend**: Bootstrap 5 + JavaScript vanilla (migração para Svelte planejada)
- **Database**: PostgreSQL 14+ com sqlx
- **Containerização**: Docker Compose
- **Autenticação**: JWT com middleware customizado
- **Autorização**: Sistema baseado em recursos e ações

## Estrutura de Diretórios

### Backend (/backend)
- `cmd/api/` - Ponto de entrada da aplicação
- `config/` - Configurações (Viper) com JWT settings
- `controllers/` - Handlers HTTP com validação de entrada
- `middleware/` - Auth, CORS e autorização por recursos
- `models/` - Estruturas de dados com validações
- `repositories/` - Camada de dados com prepared statements
- `services/` - Lógica de negócio e regras de autorização
- `routes/` - Definição de rotas com middleware aplicado
- `db/` - Conexão e migrações
- `static/` - Arquivos estáticos servidos pelo backend

### Frontend (/frontend)
- `assets/` - CSS, JS, imagens
- `components/` - Componentes reutilizáveis (navbar, footer, modals)
- `pages/` - Páginas da aplicação organizadas por funcionalidade
- `src/` - Código fonte (preparado para Svelte)
- Estrutura atual: HTML + Bootstrap + JS vanilla

### Database (/database)
- `init.sql` - Schema inicial e dados de seed

## Tecnologias Críticas
- **JWT**: golang-jwt/jwt/v5 para tokens
- **Database**: sqlx para queries otimizadas
- **Validation**: Validações customizadas nos models
- **Authorization**: Sistema de recursos e ações
- **Logging**: Zerolog para logs estruturados
- **Docker**: Compose para orquestração

## Padrões de Arquitetura
- **Clean Architecture**: Controller → Service → Repository
- **Dependency Injection**: Injeção manual nos construtores
- **Middleware Chain**: Auth → Authorization → Handler
- **Error Handling**: Erros customizados em services/errors.go
- **Security**: Prepared statements, validação de entrada
- **Role-Based Access**: Admin, Technician, Client roles

## Configurações Importantes
- **Database**: 25 max connections, 5min lifetime
- **JWT**: Access token (15min) + Refresh token (7 dias)
- **CORS**: Configurado para desenvolvimento
- **Logging**: Zerolog com níveis configuráveis
- **Validation**: Tags de validação nos structs

## Endpoints Principais
- `/api/auth/*` - Autenticação (login, register, refresh)
- `/api/drones/*` - Gestão de drones (CRUD completo)
- `/api/maintenance-orders/*` - Ordens de manutenção
- `/api/reports/*` - Relatórios técnicos (novo padrão)
- Todos protegidos por AuthMiddleware + AuthorizeResource

## Modelos de Dados
- **User**: UUID, roles (admin/technician/client), timestamps
- **Drone**: UUID, status, specifications, user relationship
- **MaintenanceOrder**: UUID, status workflow, drone relationship
- **TechnicalReport**: UUID, linked to maintenance, user relationship
- **DroneItem**: UUID, components and parts
- Todos com validações e relacionamentos bem definidos

## Sistema de Autorização
- **Recursos**: TechnicalReport, MaintenanceOrder, Drone, User
- **Ações**: Create, Read, Update, Delete
- **Middleware**: AuthorizeResource para controle granular
- **Service Layer**: Validações adicionais baseadas em roles

## Servidor e Acesso
- **Inicialização**: Usar exclusivamente o script start.sh
- **Porta**: Servidor Go na porta 8080
- **Acesso Frontend**: http://localhost:8080/pages/...
- **Acesso API**: http://localhost:8080/api/...
- **Arquivos Estáticos**: Servidos diretamente pelo Go Fiber
- **Desenvolvimento**: Sempre usar o servidor Go, nunca servidores alternativos