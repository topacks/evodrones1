package permissions

import (
	"context"
	"errors"
	"time"
)

// Permission representa uma permissão específica no sistema
type Permission struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Resource    string    `json:"resource" db:"resource"`
	Action      string    `json:"action" db:"action"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Role representa um papel/função no sistema
type Role struct {
	ID           string       `json:"id" db:"id"`
	Name         string       `json:"name" db:"name"`
	Description  string       `json:"description" db:"description"`
	ParentRoleID *string      `json:"parent_role_id,omitempty" db:"parent_role_id"`
	IsSystem     bool         `json:"is_system" db:"is_system"`
	IsActive     bool         `json:"is_active" db:"is_active"`
	Permissions  []Permission `json:"permissions,omitempty"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
}

// UserRole representa a associação entre usuário e papel
type UserRole struct {
	ID        string     `json:"id" db:"id"`
	UserID    string     `json:"user_id" db:"user_id"`
	RoleID    string     `json:"role_id" db:"role_id"`
	GrantedBy string     `json:"granted_by" db:"granted_by"`
	GrantedAt time.Time  `json:"granted_at" db:"granted_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" db:"expires_at"`
}

// Policy representa uma política de acesso dinâmica
type Policy struct {
	ID         string       `json:"id" db:"id"`
	Name       string       `json:"name" db:"name"`
	Type       PolicyType   `json:"type" db:"type"`
	Resource   string       `json:"resource" db:"resource"`
	Action     string       `json:"action" db:"action"`
	Conditions []byte       `json:"conditions" db:"conditions"`
	Effect     PolicyEffect `json:"effect" db:"effect"`
	Priority   int          `json:"priority" db:"priority"`
	IsActive   bool         `json:"is_active" db:"is_active"`
	CreatedAt  time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at" db:"updated_at"`
}

// PolicyEffect define o efeito de uma política
type PolicyEffect string

const (
	PolicyEffectAllow PolicyEffect = "allow"
	PolicyEffectDeny  PolicyEffect = "deny"
)

// PolicyType define o tipo de uma política
type PolicyType string

const (
	PolicyTypeTimeRestriction   PolicyType = "time_restriction"
	PolicyTypeResourceOwnership PolicyType = "resource_ownership"
	PolicyTypeIPRestriction     PolicyType = "ip_restriction"
	PolicyTypeCustom            PolicyType = "custom"
)

// Errors
var (
	ErrRoleNotFound     = errors.New("role not found")
	ErrInvalidResource  = errors.New("invalid resource")
	ErrPermissionDenied = errors.New("permission denied")
	ErrUserNotFound     = errors.New("user not found")
)

// AccessRequest representa uma solicitação de acesso
type AccessRequest struct {
	UserID     string                 `json:"user_id"`
	Resource   string                 `json:"resource"`
	Action     string                 `json:"action"`
	ResourceID string                 `json:"resource_id,omitempty"`
	Context    map[string]interface{} `json:"context,omitempty"`
}

// AccessResult representa o resultado de uma verificação de acesso
type AccessResult struct {
	Allowed         bool     `json:"allowed"`
	Reason          string   `json:"reason,omitempty"`
	AppliedRoles    []string `json:"applied_roles,omitempty"`
	AppliedPolicies []string `json:"applied_policies,omitempty"`
}

// PermissionChecker interface para verificação de permissões
type PermissionChecker interface {
	CheckPermission(ctx context.Context, req AccessRequest) (*AccessResult, error)
	HasPermission(ctx context.Context, userID, resource, action string) (bool, error)
	GetUserPermissions(ctx context.Context, userID string) ([]Permission, error)
	GetUserRoles(ctx context.Context, userID string) ([]Role, error)
}

// PolicyEvaluator interface para avaliação de políticas
type PolicyEvaluator interface {
	EvaluatePolicy(ctx context.Context, policy Policy, req AccessRequest) (bool, error)
	EvaluatePolicies(ctx context.Context, policies []Policy, req AccessRequest) (*AccessResult, error)
}

// RoleManager interface para gerenciamento de papéis
type RoleManager interface {
	AssignRole(ctx context.Context, userID, roleID, grantedBy string) error
	RevokeRole(ctx context.Context, userID, roleID string) error
	CreateRole(ctx context.Context, role Role) (*Role, error)
	UpdateRole(ctx context.Context, role Role) (*Role, error)
	DeleteRole(ctx context.Context, roleID string) error
	GetRole(ctx context.Context, roleID string) (*Role, error)
	ListRoles(ctx context.Context) ([]Role, error)
}

// Recursos do sistema
const (
	ResourceDrone            = "drone"
	ResourceMaintenanceOrder = "maintenance_order"
	ResourceTechnicalReport  = "technical_report"
	ResourceUser             = "user"
	ResourceRole             = "role"
	ResourcePermission       = "permission"
	ResourcePolicy           = "policy"
	ResourceDashboard        = "dashboard"
	ResourceReport           = "report"
	ResourceInventory        = "inventory"
	ResourceNotification     = "notification"
)

// Ações do sistema
const (
	ActionCreate  = "create"
	ActionRead    = "read"
	ActionUpdate  = "update"
	ActionDelete  = "delete"
	ActionList    = "list"
	ActionExecute = "execute"
	ActionApprove = "approve"
	ActionReject  = "reject"
	ActionAssign  = "assign"
	ActionExport  = "export"
	ActionImport  = "import"
)

// Papéis do sistema
const (
	RoleSystemAdmin = "system_admin"
	RoleAdmin       = "admin"
	RoleTechnician  = "technician"
	RoleUser        = "user"
	RoleManagerRole = "manager"
	RoleSupervisor  = "supervisor"
	RoleOperator    = "operator"
	RoleViewer      = "viewer"
)
