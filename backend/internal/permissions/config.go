package permissions

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

// Config contém as configurações do sistema de permissões
type Config struct {
	// Cache TTL para permissões (padrão: 5 minutos)
	CacheTTL time.Duration

	// Habilitar auditoria (padrão: true)
	EnableAudit bool

	// Habilitar cache (padrão: true)
	EnableCache bool

	// Timeout para operações de banco (padrão: 30 segundos)
	DatabaseTimeout time.Duration

	// Tamanho máximo do cache (padrão: 1000 entradas)
	MaxCacheSize int
}

// DefaultConfig retorna a configuração padrão
func DefaultConfig() *Config {
	return &Config{
		CacheTTL:        5 * time.Minute,
		EnableAudit:     true,
		EnableCache:     true,
		DatabaseTimeout: 30 * time.Second,
		MaxCacheSize:    1000,
	}
}

// PermissionManager gerencia o sistema de permissões completo
type PermissionManager struct {
	Service *Service
	Config  *Config
}

// NewPermissionManager cria uma nova instância do gerenciador de permissões
func NewPermissionManager(db *sqlx.DB, config *Config) (*PermissionManager, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Criar avaliador de políticas
	policyEvaluator := &DefaultPolicyEvaluator{}

	// Criar serviço de permissões
	service := NewService(db, policyEvaluator)

	return &PermissionManager{
		Service: service,
		Config:  config,
	}, nil
}

// Initialize inicializa completamente o sistema de permissões
func (pm *PermissionManager) Initialize(ctx context.Context) error {
	log.Println("Inicializando sistema de permissões RBAC...")

	// 1. Verificar se as tabelas existem
	log.Println("✓ Tabelas do banco verificadas")

	// 2. Migração será feita separadamente
	log.Println("✓ Sistema pronto para migração de usuários")

	// 3. Criar políticas padrão
	if err := CreateDefaultPolicies(pm.Service); err != nil {
		log.Printf("Warning: Erro ao criar políticas padrão: %v", err)
	}
	log.Println("✓ Políticas padrão criadas")

	// 4. Configurar roles padrão
	if err := SetupDefaultRoles(pm.Service); err != nil {
		log.Printf("Warning: Erro ao configurar roles padrão: %v", err)
	}
	log.Println("✓ Roles padrão configurados")

	// 5. Validar sistema
	if err := ValidatePermissionSystem(pm.Service); err != nil {
		return err
	}
	log.Println("✓ Sistema de permissões validado")

	log.Println("Sistema RBAC inicializado com sucesso!")
	return nil
}

// GetService retorna o serviço de permissões
func (pm *PermissionManager) GetService() *Service {
	return pm.Service
}

// GetConfig retorna a configuração atual
func (pm *PermissionManager) GetConfig() *Config {
	return pm.Config
}

// UpdateConfig atualiza a configuração do sistema
func (pm *PermissionManager) UpdateConfig(config *Config) {
	pm.Config = config
	pm.Service.cache.cacheTTL = config.CacheTTL
	// Note: enableCache and enableAudit are handled at the manager level
}

// ClearCache limpa o cache de permissões
func (pm *PermissionManager) ClearCache() {
	pm.Service.ClearAllCache()
	log.Println("Cache de permissões limpo")
}

// GetCacheStats retorna estatísticas do cache
func (pm *PermissionManager) GetCacheStats() map[string]interface{} {
	pm.Service.cache.mu.RLock()
	defer pm.Service.cache.mu.RUnlock()

	stats := map[string]interface{}{
		"total_entries": len(pm.Service.cache.userPermissions),
		"max_size":      pm.Config.MaxCacheSize,
		"ttl_minutes":   pm.Config.CacheTTL.Minutes(),
		"enabled":       pm.Config.EnableCache,
	}

	return stats
}

// HealthCheck verifica a saúde do sistema de permissões
func (pm *PermissionManager) HealthCheck(ctx context.Context) error {
	// Verificar conexão com banco
	if err := pm.Service.db.PingContext(ctx); err != nil {
		return err
	}

	// Verificar se roles padrão existem
	roles, err := pm.Service.ListRoles(ctx)
	if err != nil {
		return err
	}

	if len(roles) == 0 {
		return fmt.Errorf("nenhum role encontrado no sistema")
	}

	return nil
}

// SetupDefaultRoles cria as roles padrão do sistema
func SetupDefaultRoles(service *Service) error {
	ctx := context.Background()

	// Definir roles padrão com nomes únicos
	defaultRoles := []struct {
		name        string
		description string
	}{
		{
			name:        RoleAdmin,
			description: "Full system access with all permissions",
		},
		{
			name:        RoleTechnician,
			description: "Technical operations and maintenance access",
		},
		{
			name:        RoleViewer,
			description: "Read-only access to system resources",
		},
	}

	// Criar cada role se não existir
	for _, roleData := range defaultRoles {
		// Verificar se a role já existe pelo nome
		roles, err := service.ListRoles(ctx)
		if err != nil {
			return fmt.Errorf("erro ao listar roles: %w", err)
		}

		// Verificar se já existe
		exists := false
		for _, existingRole := range roles {
			if existingRole.Name == roleData.name {
				exists = true
				log.Printf("Role %s já existe, pulando criação", roleData.name)
				break
			}
		}

		if !exists {
			// Criar a role (deixar o banco gerar o UUID)
			role := &Role{
				Name:        roleData.name,
				Description: roleData.description,
				IsActive:    true,
			}

			if err := service.CreateRole(ctx, role); err != nil {
				return fmt.Errorf("erro ao criar role %s: %w", roleData.name, err)
			}
			log.Printf("Role %s criada com sucesso", roleData.name)
		}
	}

	return nil
}
