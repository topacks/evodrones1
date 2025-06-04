package permissions

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

// Service implements the unified RBAC permission system
type Service struct {
	db              *sqlx.DB
	policyEvaluator PolicyEvaluator
	cache           *permissionCache
	mu              sync.RWMutex
}

// permissionCache provides intelligent caching for permissions
type permissionCache struct {
	userPermissions map[string]map[string]bool // userID -> permission -> allowed
	userRoles       map[string][]string        // userID -> roles
	lastUpdate      map[string]time.Time       // userID -> last cache update
	cacheTTL        time.Duration
	mu              sync.RWMutex
}

// NewService creates a new permission service instance
func NewService(db *sqlx.DB, policyEvaluator PolicyEvaluator) *Service {
	return &Service{
		db:              db,
		policyEvaluator: policyEvaluator,
		cache: &permissionCache{
			userPermissions: make(map[string]map[string]bool),
			userRoles:       make(map[string][]string),
			lastUpdate:      make(map[string]time.Time),
			cacheTTL:        5 * time.Minute,
		},
	}
}

// CheckPermission verifies if a user has permission to perform an action on a resource
func (s *Service) CheckPermission(ctx context.Context, userID, resource, action string, context map[string]interface{}) (bool, error) {
	// Check cache first
	permissionKey := fmt.Sprintf("%s:%s", resource, action)
	if allowed, found := s.getCachedPermission(userID, permissionKey); found {
		return allowed, nil
	}

	// Get user roles
	userRoles, err := s.GetUserRoles(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Check permissions for each role
	allowed := false
	for _, role := range userRoles {
		rolePermissions, err := s.getRolePermissions(ctx, role.Name)
		if err != nil {
			continue
		}

		for _, perm := range rolePermissions {
			if perm.Resource == resource && perm.Action == action {
				// Check policy conditions if any
				if s.policyEvaluator != nil {
					// Create a dummy policy for evaluation
					policy := Policy{
						Resource: perm.Resource,
						Action:   perm.Action,
						Effect:   PolicyEffectAllow,
					}
					// Create AccessRequest from context
					accessReq := AccessRequest{
						UserID:   userID,
						Resource: resource,
						Action:   action,
						Context:  context,
					}
					policyAllowed, err := s.policyEvaluator.EvaluatePolicy(ctx, policy, accessReq)
					if err != nil {
						continue
					}
					allowed = policyAllowed
				} else {
					allowed = true
				}
				break
			}
		}
		if allowed {
			break
		}
	}

	// Cache the result
	s.cachePermission(userID, permissionKey, allowed)

	return allowed, nil
}

// HasPermission is a simplified version of CheckPermission without context
func (s *Service) HasPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	return s.CheckPermission(ctx, userID, resource, action, nil)
}

// GetUserRoles retrieves all roles assigned to a user
func (s *Service) GetUserRoles(ctx context.Context, userID string) ([]Role, error) {
	// Check cache first
	if roles, found := s.getCachedRoles(userID); found {
		return s.convertStringRolesToRoles(ctx, roles)
	}

	query := `
		SELECT r.id, r.name, r.description, r.parent_role_id, r.is_active, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1 AND r.is_active = true AND ur.is_active = true
		ORDER BY r.name
	`

	var roles []Role
	err := s.db.SelectContext(ctx, &roles, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Include inherited roles from parent roles
	allRoles := make([]Role, 0, len(roles))
	processedRoles := make(map[string]bool)

	for _, role := range roles {
		inheritedRoles, err := s.getInheritedRoles(ctx, role.ID, processedRoles)
		if err != nil {
			continue
		}
		allRoles = append(allRoles, inheritedRoles...)
	}

	// Cache the role names
	roleNames := make([]string, len(allRoles))
	for i, role := range allRoles {
		roleNames[i] = role.Name
	}
	s.cacheRoles(userID, roleNames)

	return allRoles, nil
}

// GetUserPermissions retrieves all permissions for a user (aggregated from all roles)
func (s *Service) GetUserPermissions(ctx context.Context, userID string) ([]Permission, error) {
	userRoles, err := s.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	permissionMap := make(map[string]Permission)

	for _, role := range userRoles {
		rolePermissions, err := s.getRolePermissions(ctx, role.Name)
		if err != nil {
			continue
		}

		for _, perm := range rolePermissions {
			key := fmt.Sprintf("%s:%s", perm.Resource, perm.Action)
			permissionMap[key] = perm
		}
	}

	permissions := make([]Permission, 0, len(permissionMap))
	for _, perm := range permissionMap {
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

// AssignRole assigns a role to a user
func (s *Service) AssignRole(ctx context.Context, userID, roleID string, assignedBy string) error {
	// Check if role exists and is active
	var roleExists bool
	err := s.db.GetContext(ctx, &roleExists, "SELECT EXISTS(SELECT 1 FROM roles WHERE id = $1 AND is_active = true)", roleID)
	if err != nil {
		return fmt.Errorf("failed to check role existence: %w", err)
	}
	if !roleExists {
		return fmt.Errorf("role not found or inactive")
	}

	// Check if assignment already exists
	var assignmentExists bool
	err = s.db.GetContext(ctx, &assignmentExists, "SELECT EXISTS(SELECT 1 FROM user_roles WHERE user_id = $1 AND role_id = $2)", userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to check existing assignment: %w", err)
	}
	if assignmentExists {
		return fmt.Errorf("role already assigned to user")
	}

	// Create the assignment
	query := `
		INSERT INTO user_roles (user_id, role_id, assigned_by, assigned_at, is_active)
		VALUES ($1, $2, $3, NOW(), true)
	`
	_, err = s.db.ExecContext(ctx, query, userID, roleID, assignedBy)
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	// Clear cache for this user
	s.clearUserCache(userID)

	return nil
}

// RevokeRole removes a role from a user
func (s *Service) RevokeRole(ctx context.Context, userID, roleID string) error {
	query := `
		UPDATE user_roles 
		SET is_active = false, revoked_at = NOW()
		WHERE user_id = $1 AND role_id = $2 AND is_active = true
	`
	result, err := s.db.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to revoke role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("role assignment not found or already revoked")
	}

	// Clear cache for this user
	s.clearUserCache(userID)

	return nil
}

// CreateRole creates a new role
func (s *Service) CreateRole(ctx context.Context, role *Role) error {
	query := `
		INSERT INTO roles (name, description, parent_role_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := s.db.QueryRowContext(ctx, query, role.Name, role.Description, role.ParentRoleID, role.IsActive).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

// UpdateRole updates an existing role
func (s *Service) UpdateRole(ctx context.Context, roleID string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	// Build dynamic update query
	setClauses := make([]string, 0, len(updates))
	args := make([]interface{}, 0, len(updates)+1)
	argIndex := 1

	for field, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, argIndex))
		args = append(args, value)
		argIndex++
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	args = append(args, roleID)

	query := fmt.Sprintf("UPDATE roles SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argIndex+1)

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("role not found")
	}

	return nil
}

// DeleteRole soft deletes a role
func (s *Service) DeleteRole(ctx context.Context, roleID string) error {
	query := `UPDATE roles SET is_active = false, updated_at = NOW() WHERE id = $1 AND is_active = true`
	result, err := s.db.ExecContext(ctx, query, roleID)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("role not found or already deleted")
	}

	return nil
}

// ListRoles retrieves all active roles
func (s *Service) ListRoles(ctx context.Context) ([]Role, error) {
	query := `
		SELECT id, name, description, parent_role_id, is_active, created_at, updated_at
		FROM roles
		WHERE is_active = true
		ORDER BY name
	`

	var roles []Role
	err := s.db.SelectContext(ctx, &roles, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	return roles, nil
}

// GetRole retrieves a specific role by ID
func (s *Service) GetRole(ctx context.Context, roleID string) (*Role, error) {
	query := `
		SELECT id, name, description, parent_role_id, is_active, created_at, updated_at
		FROM roles
		WHERE id = $1 AND is_active = true
	`

	var role Role
	err := s.db.GetContext(ctx, &role, query, roleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return &role, nil
}

// Private helper methods

func (s *Service) getRolePermissions(ctx context.Context, roleName string) ([]Permission, error) {
	query := `
		SELECT p.id, p.resource, p.action, p.description, p.created_at, p.updated_at
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN roles r ON rp.role_id = r.id
		WHERE r.name = $1 AND r.is_active = true AND rp.is_active = true
	`

	var permissions []Permission
	err := s.db.SelectContext(ctx, &permissions, query, roleName)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}

	return permissions, nil
}

func (s *Service) getInheritedRoles(ctx context.Context, roleID string, processed map[string]bool) ([]Role, error) {
	if processed[roleID] {
		return nil, nil // Avoid circular dependencies
	}
	processed[roleID] = true

	// Get the role itself
	role, err := s.GetRole(ctx, roleID)
	if err != nil {
		return nil, err
	}

	roles := []Role{*role}

	// Get parent role if exists
	if role.ParentRoleID != nil {
		parentRoles, err := s.getInheritedRoles(ctx, *role.ParentRoleID, processed)
		if err != nil {
			return roles, nil // Return current role even if parent fails
		}
		roles = append(roles, parentRoles...)
	}

	return roles, nil
}

func (s *Service) convertStringRolesToRoles(ctx context.Context, roleNames []string) ([]Role, error) {
	if len(roleNames) == 0 {
		return []Role{}, nil
	}

	placeholders := make([]string, len(roleNames))
	args := make([]interface{}, len(roleNames))
	for i, name := range roleNames {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = name
	}

	query := fmt.Sprintf(`
		SELECT id, name, description, parent_role_id, is_active, created_at, updated_at
		FROM roles
		WHERE name IN (%s) AND is_active = true
		ORDER BY name
	`, strings.Join(placeholders, ","))

	var roles []Role
	err := s.db.SelectContext(ctx, &roles, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to convert role names to roles: %w", err)
	}

	return roles, nil
}

// Cache management methods

func (s *Service) getCachedPermission(userID, permissionKey string) (bool, bool) {
	s.cache.mu.RLock()
	defer s.cache.mu.RUnlock()

	if s.isCacheExpired(userID) {
		return false, false
	}

	if userPerms, exists := s.cache.userPermissions[userID]; exists {
		if allowed, found := userPerms[permissionKey]; found {
			return allowed, true
		}
	}

	return false, false
}

func (s *Service) cachePermission(userID, permissionKey string, allowed bool) {
	s.cache.mu.Lock()
	defer s.cache.mu.Unlock()

	if s.cache.userPermissions[userID] == nil {
		s.cache.userPermissions[userID] = make(map[string]bool)
	}

	s.cache.userPermissions[userID][permissionKey] = allowed
	s.cache.lastUpdate[userID] = time.Now()
}

func (s *Service) getCachedRoles(userID string) ([]string, bool) {
	s.cache.mu.RLock()
	defer s.cache.mu.RUnlock()

	if s.isCacheExpired(userID) {
		return nil, false
	}

	if roles, exists := s.cache.userRoles[userID]; exists {
		return roles, true
	}

	return nil, false
}

func (s *Service) cacheRoles(userID string, roles []string) {
	s.cache.mu.Lock()
	defer s.cache.mu.Unlock()

	s.cache.userRoles[userID] = roles
	s.cache.lastUpdate[userID] = time.Now()
}

func (s *Service) clearUserCache(userID string) {
	s.cache.mu.Lock()
	defer s.cache.mu.Unlock()

	delete(s.cache.userPermissions, userID)
	delete(s.cache.userRoles, userID)
	delete(s.cache.lastUpdate, userID)
}

func (s *Service) isCacheExpired(userID string) bool {
	lastUpdate, exists := s.cache.lastUpdate[userID]
	if !exists {
		return true
	}
	return time.Since(lastUpdate) > s.cache.cacheTTL
}

// ClearAllCache clears the entire permission cache
func (s *Service) ClearAllCache() {
	s.cache.mu.Lock()
	defer s.cache.mu.Unlock()

	s.cache.userPermissions = make(map[string]map[string]bool)
	s.cache.userRoles = make(map[string][]string)
	s.cache.lastUpdate = make(map[string]time.Time)
}
