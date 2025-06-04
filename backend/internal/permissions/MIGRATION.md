# Guia de Migração - Sistema RBAC EVO DRONES

Este documento fornece instruções detalhadas para migrar do sistema de autorização atual para o novo sistema RBAC.

## 📋 Pré-requisitos

- Go 1.19+
- PostgreSQL 12+
- Backup do banco de dados atual
- Acesso administrativo ao sistema

## 🚀 Processo de Migração

### Passo 1: Backup do Sistema Atual

```bash
# Fazer backup do banco de dados
pg_dump evodrones > backup_pre_rbac_$(date +%Y%m%d_%H%M%S).sql

# Fazer backup dos arquivos de código
cp -r backend/middleware backup_middleware_$(date +%Y%m%d_%H%M%S)
cp -r backend/services backup_services_$(date +%Y%m%d_%H%M%S)
```

### Passo 2: Executar Migração do Banco

```bash
# Executar script de migração SQL
psql -d evodrones -f backend/db/migrations/001_create_rbac_tables.sql

# Verificar se tabelas foram criadas
psql -d evodrones -c "\dt rbac*"
```

### Passo 3: Atualizar main.go

**Antes:**
```go
package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
    "your-project/config"
    "your-project/db"
    "your-project/routes"
)

func main() {
    cfg := config.Load()
    database := db.Connect(cfg.Database)
    
    app := fiber.New()
    routes.Setup(app, database, cfg)
    
    log.Fatal(app.Listen(":8080"))
}
```

**Depois:**
```go
package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
    "your-project/config"
    "your-project/db"
    "your-project/routes"
    "your-project/internal/permissions"
)

func main() {
    cfg := config.Load()
    database := db.Connect(cfg.Database)
    
    // Inicializar sistema RBAC
    permissionManager, err := permissions.InitializePermissionSystem(database)
    if err != nil {
        log.Fatal("Erro ao inicializar sistema de permissões:", err)
    }
    
    app := fiber.New()
    
    // Configurar rotas com novo sistema
    routes.Setup(app, database, cfg, permissionManager.GetService())
    
    log.Fatal(app.Listen(":8080"))
}
```

### Passo 4: Atualizar Rotas

**Antes (routes/routes.go):**
```go
func Setup(app *fiber.App, db *sqlx.DB, cfg *config.Config) {
    api := app.Group("/api")
    
    // Middleware antigo
    api.Use(middleware.AuthMiddleware(&cfg.JWT))
    
    // Rotas com middleware antigo
    api.Post("/drones", 
        middleware.RoleMiddleware("admin"), 
        droneController.Create)
}
```

**Depois (routes/routes.go):**
```go
func Setup(app *fiber.App, db *sqlx.DB, cfg *config.Config, permissionService *permissions.Service) {
    api := app.Group("/api")
    
    // Middleware de autenticação (mantém o mesmo)
    api.Use(middleware.AuthMiddleware(&cfg.JWT))
    
    // Rotas com novo middleware RBAC
    api.Post("/drones", 
        permissions.CreateDroneMiddleware(permissionService), 
        droneController.Create)
    
    api.Get("/drones", 
        permissions.ListDronesMiddleware(permissionService), 
        droneController.List)
    
    api.Get("/drones/:id", 
        permissions.ReadDroneMiddleware(permissionService), 
        droneController.GetByID)
    
    api.Put("/drones/:id", 
        permissions.UpdateDroneMiddleware(permissionService), 
        droneController.Update)
    
    api.Delete("/drones/:id", 
        permissions.DeleteDroneMiddleware(permissionService), 
        droneController.Delete)
}
```

### Passo 5: Atualizar Controllers

**Antes:**
```go
func (c *DroneController) Create(ctx *fiber.Ctx) error {
    userID := ctx.Locals("user_id").(string)
    userRole := ctx.Locals("role").(string)
    
    // Validação manual de role
    if userRole != "admin" {
        return fiber.NewError(fiber.StatusForbidden, "acesso negado")
    }
    
    // Lógica do controller...
}
```

**Depois:**
```go
func (c *DroneController) Create(ctx *fiber.Ctx) error {
    userID := ctx.Locals("user_id").(string)
    
    // Permissão já verificada pelo middleware RBAC
    // Apenas implementar lógica de negócio
    
    var input models.CreateDroneInput
    if err := ctx.BodyParser(&input); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "dados inválidos")
    }
    
    drone, err := c.service.Create(input, userID)
    if err != nil {
        return handleServiceError(err)
    }
    
    return ctx.Status(fiber.StatusCreated).JSON(drone)
}
```

### Passo 6: Remover Código Antigo

**Arquivos para remover/modificar:**

1. **middleware/authorization.go** - Substituído pelo novo sistema
2. **services/authorization_service.go** - Substituído pelo novo sistema
3. **Middleware RoleMiddleware** em auth.go - Substituído pelos novos middlewares

**Manter:**
- **middleware/auth.go** - AuthMiddleware (validação JWT)
- Estruturas de User em models/user.go
- Lógica de login/logout

### Passo 7: Verificar Migração

```bash
# Executar testes
go test ./internal/permissions/...

# Verificar se servidor inicia
go run cmd/api/main.go

# Testar endpoints
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/drones
```

## 🔄 Mapeamento de Roles

| Role Antigo | Role Novo | Permissões |
|-------------|-----------|------------|
| `admin` | `admin` | Todas as permissões |
| `technician` | `technician` | Manutenção, relatórios, leitura de drones |
| `user` | `viewer` | Apenas leitura |

## 📊 Verificação de Permissões

### Testando Permissões por Role

```sql
-- Verificar roles de um usuário
SELECT u.email, r.name as role_name
FROM users u
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
WHERE u.email = 'admin@example.com';

-- Verificar permissões de um usuário
SELECT DISTINCT p.resource, p.action, p.description
FROM users u
JOIN user_roles ur ON u.id = ur.user_id
JOIN role_permissions rp ON ur.role_id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE u.email = 'tech@example.com';
```

### Testando via API

```bash
# Login como admin
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password"}' | jq -r '.token')

# Testar criação de drone (deve funcionar)
curl -X POST http://localhost:8080/api/drones \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Drone Test","model":"DJI"}'

# Login como technician
TECH_TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"tech@example.com","password":"password"}' | jq -r '.token')

# Testar criação de usuário (deve falhar)
curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer $TECH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"New User","email":"new@example.com"}'
```

## 🛠️ Troubleshooting

### Problema: "permission denied"

**Causa:** Usuário não tem role ou permissão necessária

**Solução:**
```sql
-- Verificar roles do usuário
SELECT * FROM user_roles WHERE user_id = 'user-id';

-- Atribuir role se necessário
INSERT INTO user_roles (user_id, role_id, assigned_by) 
VALUES ('user-id', 'admin', 'system');
```

### Problema: "role not found"

**Causa:** Roles padrão não foram criados

**Solução:**
```sql
-- Executar novamente a migração
\i backend/db/migrations/001_create_rbac_tables.sql
```

### Problema: Cache não funciona

**Causa:** Configuração incorreta

**Solução:**
```go
// Verificar configuração
config := &permissions.Config{
    EnableCache: true,
    CacheTTL: 5 * time.Minute,
}
```

### Problema: Auditoria não registra

**Causa:** Auditoria desabilitada

**Solução:**
```go
config.EnableAudit = true
```

## 📈 Monitoramento

### Logs de Auditoria

```sql
-- Ver últimas ações
SELECT 
    pa.action,
    pa.resource_type,
    pa.allowed,
    u.email,
    pa.created_at
FROM permission_audit pa
JOIN users u ON pa.user_id = u.id
ORDER BY pa.created_at DESC
LIMIT 50;
```

### Métricas do Sistema

```go
// Obter status do sistema
status := permissions.GetSystemStatus(permissionManager)
log.Printf("Sistema: %+v", status)
```

## 🔄 Rollback (Se Necessário)

### Rollback Completo

```bash
# 1. Parar aplicação
pkill -f "go run"

# 2. Restaurar backup do banco
psql -d evodrones < backup_pre_rbac_YYYYMMDD_HHMMSS.sql

# 3. Restaurar código
rm -rf backend/middleware
cp -r backup_middleware_YYYYMMDD_HHMMSS backend/middleware

# 4. Restaurar main.go original
git checkout HEAD~1 -- cmd/api/main.go

# 5. Reiniciar aplicação
go run cmd/api/main.go
```

### Rollback Parcial (Manter Tabelas)

```sql
-- Apenas desabilitar sistema RBAC
UPDATE roles SET active = false;
UPDATE permissions SET active = false;
```

## ✅ Checklist de Migração

- [ ] Backup do banco de dados criado
- [ ] Backup do código criado
- [ ] Migração SQL executada com sucesso
- [ ] Tabelas RBAC criadas
- [ ] Roles padrão inseridos
- [ ] Permissões padrão inseridas
- [ ] main.go atualizado
- [ ] Rotas atualizadas
- [ ] Controllers atualizados
- [ ] Middleware antigo removido
- [ ] Testes executados com sucesso
- [ ] Servidor inicia sem erros
- [ ] Endpoints funcionam corretamente
- [ ] Permissões testadas por role
- [ ] Auditoria funcionando
- [ ] Cache funcionando
- [ ] Monitoramento configurado

## 📞 Suporte

Em caso de problemas durante a migração:

1. Verificar logs da aplicação
2. Verificar logs do PostgreSQL
3. Consultar seção de Troubleshooting
4. Executar rollback se necessário
5. Consultar documentação completa em `README.md`

---

**⚠️ IMPORTANTE:** Sempre teste a migração em ambiente de desenvolvimento antes de aplicar em produção!