package permissions

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestService cria um service de teste com mock do banco
func setupTestService(t *testing.T) (*Service, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "postgres")

	service := &Service{
		db:              sqlxDB,
		policyEvaluator: &DefaultPolicyEvaluator{},
		permissionCache: make(map[string]*cacheEntry),
		cacheTTL:        5 * time.Minute,
		enableCache:     true,
		enableAudit:     true,
	}

	cleanup := func() {
		db.Close()
	}

	return service, mock, cleanup
}

func TestService_CheckPermission(t *testing.T) {
	service, mock, cleanup := setupTestService(t)
	defer cleanup()

	tests := []struct {
		name      string
		userID    string
		resource  string
		action    string
		context   map[string]interface{}
		mockSetup func(sqlmock.Sqlmock)
		want      bool
		wantErr   bool
	}{
		{
			name:     "admin tem acesso total",
			userID:   "admin-user-id",
			resource: ResourceDrone,
			action:   ActionCreate,
			context:  map[string]interface{}{"user_id": "admin-user-id"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock para buscar roles do usuário
				mock.ExpectQuery(`SELECT r.id, r.name, r.description`).WithArgs("admin-user-id").WillReturnRows(
					sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(RoleAdmin, "Admin", "Administrator"),
				)
				// Mock para buscar permissões
				mock.ExpectQuery(`SELECT p.resource, p.action`).WithArgs(RoleAdmin).WillReturnRows(
					sqlmock.NewRows([]string{"resource", "action"}).AddRow(ResourceDrone, ActionCreate),
				)
				// Mock para auditoria
				mock.ExpectExec(`INSERT INTO permission_audit`).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want:    true,
			wantErr: false,
		},
		{
			name:     "technician não pode criar usuários",
			userID:   "tech-user-id",
			resource: ResourceUser,
			action:   ActionCreate,
			context:  map[string]interface{}{"user_id": "tech-user-id"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock para buscar roles do usuário
				mock.ExpectQuery(`SELECT r.id, r.name, r.description`).WithArgs("tech-user-id").WillReturnRows(
					sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(RoleTechnician, "Technician", "Technician"),
				)
				// Mock para buscar permissões (sem permissão para criar usuários)
				mock.ExpectQuery(`SELECT p.resource, p.action`).WithArgs(RoleTechnician).WillReturnRows(
					sqlmock.NewRows([]string{"resource", "action"}).AddRow(ResourceMaintenanceOrder, ActionCreate),
				)
				// Mock para auditoria
				mock.ExpectExec(`INSERT INTO permission_audit`).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want:    false,
			wantErr: false,
		},
		{
			name:     "viewer pode ler drones",
			userID:   "viewer-user-id",
			resource: ResourceDrone,
			action:   ActionRead,
			context:  map[string]interface{}{"user_id": "viewer-user-id"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock para buscar roles do usuário
				mock.ExpectQuery(`SELECT r.id, r.name, r.description`).WithArgs("viewer-user-id").WillReturnRows(
					sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(RoleViewer, "Viewer", "Viewer"),
				)
				// Mock para buscar permissões
				mock.ExpectQuery(`SELECT p.resource, p.action`).WithArgs(RoleViewer).WillReturnRows(
					sqlmock.NewRows([]string{"resource", "action"}).AddRow(ResourceDrone, ActionRead),
				)
				// Mock para auditoria
				mock.ExpectExec(`INSERT INTO permission_audit`).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			got, err := service.CheckPermission(context.Background(), tt.userID, tt.resource, tt.action, tt.context)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			// Verificar se todas as expectativas do mock foram atendidas
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestService_AssignRole(t *testing.T) {
	service, mock, cleanup := setupTestService(t)
	defer cleanup()

	tests := []struct {
		name       string
		userID     string
		roleID     string
		assignedBy string
		mockSetup  func(sqlmock.Sqlmock)
		wantErr    bool
	}{
		{
			name:       "atribuir role com sucesso",
			userID:     "user-123",
			roleID:     RoleTechnician,
			assignedBy: "admin-456",
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock para verificar se role existe
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM roles WHERE id = \$1`).WithArgs(RoleTechnician).WillReturnRows(
					sqlmock.NewRows([]string{"count"}).AddRow(1),
				)
				// Mock para verificar se usuário existe
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users WHERE id = \$1`).WithArgs("user-123").WillReturnRows(
					sqlmock.NewRows([]string{"count"}).AddRow(1),
				)
				// Mock para inserir user_role
				mock.ExpectExec(`INSERT INTO user_roles`).WithArgs("user-123", RoleTechnician, "admin-456").WillReturnResult(sqlmock.NewResult(1, 1))
				// Mock para auditoria
				mock.ExpectExec(`INSERT INTO permission_audit`).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name:       "erro ao atribuir role inexistente",
			userID:     "user-123",
			roleID:     "role-inexistente",
			assignedBy: "admin-456",
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock para verificar se role existe (não existe)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM roles WHERE id = \$1`).WithArgs("role-inexistente").WillReturnRows(
					sqlmock.NewRows([]string{"count"}).AddRow(0),
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			err := service.AssignRole(context.Background(), tt.userID, tt.roleID, tt.assignedBy)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verificar se todas as expectativas do mock foram atendidas
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestService_GetUserRoles(t *testing.T) {
	service, mock, cleanup := setupTestService(t)
	defer cleanup()

	userID := "user-123"
	expectedRoles := []Role{
		{ID: RoleTechnician, Name: "Technician", Description: "Technician role"},
		{ID: RoleViewer, Name: "Viewer", Description: "Viewer role"},
	}

	// Mock para buscar roles do usuário
	mock.ExpectQuery(`SELECT r.id, r.name, r.description`).WithArgs(userID).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "description"}).
			AddRow(RoleTechnician, "Technician", "Technician role").
			AddRow(RoleViewer, "Viewer", "Viewer role"),
	)

	roles, err := service.GetUserRoles(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedRoles, roles)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_GetUserPermissions(t *testing.T) {
	service, mock, cleanup := setupTestService(t)
	defer cleanup()

	userID := "user-123"
	expectedPermissions := []Permission{
		{Resource: ResourceDrone, Action: ActionRead, Description: "Read drone"},
		{Resource: ResourceMaintenanceOrder, Action: ActionCreate, Description: "Create maintenance order"},
	}

	// Mock para buscar permissões do usuário
	mock.ExpectQuery(`SELECT DISTINCT p.resource, p.action, p.description`).WithArgs(userID).WillReturnRows(
		sqlmock.NewRows([]string{"resource", "action", "description"}).
			AddRow(ResourceDrone, ActionRead, "Read drone").
			AddRow(ResourceMaintenanceOrder, ActionCreate, "Create maintenance order"),
	)

	permissions, err := service.GetUserPermissions(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedPermissions, permissions)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_CreateRole(t *testing.T) {
	service, mock, cleanup := setupTestService(t)
	defer cleanup()

	role := &Role{
		ID:          "custom-role",
		Name:        "Custom Role",
		Description: "A custom role for testing",
	}

	// Mock para inserir role
	mock.ExpectExec(`INSERT INTO roles`).WithArgs(role.ID, role.Name, role.Description).WillReturnResult(sqlmock.NewResult(1, 1))

	err := service.CreateRole(context.Background(), role)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_Cache(t *testing.T) {
	service, mock, cleanup := setupTestService(t)
	defer cleanup()

	userID := "user-123"
	resource := ResourceDrone
	action := ActionRead

	// Primeira chamada - deve consultar o banco
	mock.ExpectQuery(`SELECT r.id, r.name, r.description`).WithArgs(userID).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(RoleViewer, "Viewer", "Viewer"),
	)
	mock.ExpectQuery(`SELECT p.resource, p.action`).WithArgs(RoleViewer).WillReturnRows(
		sqlmock.NewRows([]string{"resource", "action"}).AddRow(ResourceDrone, ActionRead),
	)
	mock.ExpectExec(`INSERT INTO permission_audit`).WillReturnResult(sqlmock.NewResult(1, 1))

	// Primeira verificação
	allowed1, err1 := service.CheckPermission(context.Background(), userID, resource, action, nil)
	assert.NoError(t, err1)
	assert.True(t, allowed1)

	// Segunda chamada - deve usar cache (sem consultas ao banco)
	allowed2, err2 := service.CheckPermission(context.Background(), userID, resource, action, nil)
	assert.NoError(t, err2)
	assert.True(t, allowed2)

	// Verificar se todas as expectativas foram atendidas
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_ClearCache(t *testing.T) {
	service, mock, cleanup := setupTestService(t)
	defer cleanup()

	// Adicionar entrada no cache
	service.permissionCache["test-key"] = &cacheEntry{
		permissions: []Permission{{Resource: ResourceDrone, Action: ActionRead}},
		cachedAt:    time.Now(),
	}

	// Verificar se cache tem entrada
	assert.Len(t, service.permissionCache, 1)

	// Limpar cache
	service.clearCache()

	// Verificar se cache foi limpo
	assert.Len(t, service.permissionCache, 0)
}

func TestService_HasPermission(t *testing.T) {
	service, mock, cleanup := setupTestService(t)
	defer cleanup()

	userID := "user-123"
	resource := ResourceDrone
	action := ActionRead

	// Mock para buscar roles do usuário
	mock.ExpectQuery(`SELECT r.id, r.name, r.description`).WithArgs(userID).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(RoleViewer, "Viewer", "Viewer"),
	)
	// Mock para buscar permissões
	mock.ExpectQuery(`SELECT p.resource, p.action`).WithArgs(RoleViewer).WillReturnRows(
		sqlmock.NewRows([]string{"resource", "action"}).AddRow(ResourceDrone, ActionRead),
	)

	// Testar HasPermission (sem auditoria)
	allowed, err := service.HasPermission(context.Background(), userID, resource, action)

	assert.NoError(t, err)
	assert.True(t, allowed)
	assert.NoError(t, mock.ExpectationsWereMet())
}
