package permissions

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

// InitializePermissionSystemComplete inicializa completamente o sistema de permissões RBAC
// Esta é a função principal que deve ser chamada no main.go do projeto
func InitializePermissionSystemComplete(db *sqlx.DB) (*PermissionManager, error) {
	log.Println("🚀 Inicializando Sistema de Permissões RBAC...")

	// 1. Criar configuração padrão
	config := DefaultConfig()

	// 2. Criar gerenciador de permissões
	permissionManager, err := NewPermissionManager(db, config)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar gerenciador de permissões: %w", err)
	}

	// 3. Inicializar sistema completo
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := permissionManager.Initialize(ctx); err != nil {
		return nil, fmt.Errorf("erro ao inicializar sistema: %w", err)
	}

	// 4. Migrar do sistema antigo
	if err := MigrateFromOldSystem(db, permissionManager.Service); err != nil {
		log.Printf("⚠️  Warning: Erro na migração do sistema antigo: %v", err)
	}

	// 5. Validar integração
	if err := ValidateIntegration(permissionManager.Service); err != nil {
		return nil, fmt.Errorf("erro na validação da integração: %w", err)
	}

	log.Println("✅ Sistema de Permissões RBAC inicializado com sucesso!")
	return permissionManager, nil
}

// QuickSetup é uma função de configuração rápida para desenvolvimento
func QuickSetup(db *sqlx.DB) (*Service, error) {
	permissionManager, err := InitializePermissionSystemComplete(db)
	if err != nil {
		return nil, err
	}

	return permissionManager.Service, nil
}

// SetupWithCustomConfig inicializa o sistema com configuração customizada
func SetupWithCustomConfig(db *sqlx.DB, config *Config) (*PermissionManager, error) {
	if config == nil {
		config = DefaultConfig()
	}

	log.Printf("🔧 Configuração personalizada: Cache=%v, Auditoria=%v, TTL=%v",
		config.EnableCache, config.EnableAudit, config.CacheTTL)

	// Criar gerenciador com configuração customizada
	permissionManager, err := NewPermissionManager(db, config)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar gerenciador: %w", err)
	}

	// Inicializar
	ctx, cancel := context.WithTimeout(context.Background(), config.DatabaseTimeout)
	defer cancel()

	if err := permissionManager.Initialize(ctx); err != nil {
		return nil, fmt.Errorf("erro na inicialização: %w", err)
	}

	return permissionManager, nil
}

// CreateTestPermissionSystem cria um sistema de permissões para testes
func CreateTestPermissionSystem(db *sqlx.DB) (*Service, error) {
	config := &Config{
		CacheTTL:        1 * time.Minute, // Cache mais curto para testes
		EnableAudit:     false,           // Desabilitar auditoria em testes
		EnableCache:     true,
		DatabaseTimeout: 10 * time.Second,
		MaxCacheSize:    100,
	}

	permissionManager, err := NewPermissionManager(db, config)
	if err != nil {
		return nil, err
	}

	// Inicialização mínima para testes
	// TODO: Implement VerifyDatabaseTables and SetupDefaultRoles
	// if err := VerifyDatabaseTables(db); err != nil {
	//	return nil, err
	// }

	// if err := SetupDefaultRoles(permissionManager.Service); err != nil {
	//	return nil, err
	// }

	return permissionManager.Service, nil
}

// GetSystemStatus retorna o status atual do sistema de permissões
func GetSystemStatus(permissionManager *PermissionManager) map[string]interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status := map[string]interface{}{
		"system":  "RBAC Permission System",
		"version": "1.0.0",
		"status":  "unknown",
	}

	// Verificar saúde do sistema
	if err := permissionManager.HealthCheck(ctx); err != nil {
		status["status"] = "error"
		status["error"] = err.Error()
	} else {
		status["status"] = "healthy"
	}

	// Adicionar estatísticas do cache
	cacheStats := permissionManager.GetCacheStats()
	status["cache"] = cacheStats

	// Adicionar configuração
	config := permissionManager.GetConfig()
	status["config"] = map[string]interface{}{
		"cache_enabled":     config.EnableCache,
		"audit_enabled":     config.EnableAudit,
		"cache_ttl_minutes": config.CacheTTL.Minutes(),
		"max_cache_size":    config.MaxCacheSize,
	}

	// Contar roles e permissões
	service := permissionManager.GetService()
	if roles, err := service.ListRoles(ctx); err == nil {
		status["total_roles"] = len(roles)
	}

	return status
}

// ResetPermissionSystem reseta completamente o sistema (USE COM CUIDADO!)
func ResetPermissionSystem(db *sqlx.DB) error {
	log.Println("⚠️  ATENÇÃO: Resetando sistema de permissões...")

	// Limpar todas as tabelas RBAC
	tables := []string{
		"permission_audit",
		"permission_policies",
		"user_roles",
		"role_permissions",
		"policies",
		"permissions",
		"roles",
	}

	for _, table := range tables {
		if _, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			log.Printf("Erro ao limpar tabela %s: %v", table, err)
		}
	}

	log.Println("🗑️  Sistema de permissões resetado")
	return nil
}

// ExportPermissionData exporta dados do sistema para backup
func ExportPermissionData(service *Service) (map[string]interface{}, error) {
	ctx := context.Background()

	exportData := make(map[string]interface{})

	// Exportar roles
	roles, err := service.ListRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao exportar roles: %w", err)
	}
	exportData["roles"] = roles

	// Exportar user_roles
	var userRoles []UserRole
	query := `SELECT user_id, role_id, assigned_by, assigned_at FROM user_roles`
	if err := service.db.SelectContext(ctx, &userRoles, query); err != nil {
		return nil, fmt.Errorf("erro ao exportar user_roles: %w", err)
	}
	exportData["user_roles"] = userRoles

	// Exportar políticas
	var policies []Policy
	query = `SELECT id, name, type, rules, created_at, updated_at FROM policies`
	if err := service.db.SelectContext(ctx, &policies, query); err != nil {
		return nil, fmt.Errorf("erro ao exportar policies: %w", err)
	}
	exportData["policies"] = policies

	exportData["exported_at"] = time.Now()
	exportData["version"] = "1.0.0"

	return exportData, nil
}

// ImportPermissionData importa dados do sistema de um backup
func ImportPermissionData(service *Service, data map[string]interface{}) error {
	ctx := context.Background()

	log.Println("📥 Importando dados do sistema de permissões...")

	// Importar roles
	if rolesData, ok := data["roles"]; ok {
		if roles, ok := rolesData.([]Role); ok {
			for _, role := range roles {
				if err := service.CreateRole(ctx, &role); err != nil {
					log.Printf("Erro ao importar role %s: %v", role.ID, err)
				}
			}
		}
	}

	// Importar user_roles
	if userRolesData, ok := data["user_roles"]; ok {
		if userRoles, ok := userRolesData.([]UserRole); ok {
			for _, userRole := range userRoles {
				if err := service.AssignRole(ctx, userRole.UserID, userRole.RoleID, userRole.GrantedBy); err != nil {
					log.Printf("Erro ao importar user_role %s->%s: %v", userRole.UserID, userRole.RoleID, err)
				}
			}
		}
	}

	log.Println("✅ Importação concluída")
	return nil
}

// MonitorPermissionSystem monitora o sistema e reporta métricas
func MonitorPermissionSystem(permissionManager *PermissionManager, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("📊 Iniciando monitoramento do sistema (intervalo: %v)", interval)

	for {
		select {
		case <-ticker.C:
			status := GetSystemStatus(permissionManager)
			log.Printf("📈 Status: %s | Cache: %d entradas | Roles: %v",
				status["status"],
				status["cache"].(map[string]interface{})["total_entries"],
				status["total_roles"])
		}
	}
}
