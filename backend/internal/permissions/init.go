package permissions

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

// InitializePermissionSystem initializes the RBAC permission system
func InitializePermissionSystem(db *sqlx.DB) (*Service, error) {
	// Create policy evaluator
	policyEvaluator := NewDefaultPolicyEvaluator()

	// Create permission service
	service := NewService(db, policyEvaluator)

	// Verify database tables exist
	err := verifyRBACTables(db)
	if err != nil {
		return nil, fmt.Errorf("failed to verify RBAC tables: %w", err)
	}

	// Migrate existing users to new role system
	err = migrateExistingUsers(db, service)
	if err != nil {
		log.Printf("Warning: failed to migrate existing users: %v", err)
		// Don't fail initialization, just log the warning
	}

	log.Println("RBAC permission system initialized successfully")
	return service, nil
}

// verifyRBACTables checks if all required RBAC tables exist
func verifyRBACTables(db *sqlx.DB) error {
	requiredTables := []string{
		"roles",
		"permissions",
		"role_permissions",
		"user_roles",
		"policies",
		"permission_policies",
		"permission_audit",
	}

	for _, table := range requiredTables {
		var exists bool
		query := `SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)`
		err := db.Get(&exists, query, table)
		if err != nil {
			return fmt.Errorf("failed to check table %s: %w", table, err)
		}
		if !exists {
			return fmt.Errorf("required table %s does not exist", table)
		}
	}

	return nil
}

// migrateExistingUsers migrates existing users from the old role system to the new RBAC system
func migrateExistingUsers(db *sqlx.DB, service *Service) error {
	ctx := context.Background()

	// Get all existing users
	type ExistingUser struct {
		ID   string `db:"id"`
		Role string `db:"role"`
	}

	var users []ExistingUser
	query := `SELECT id, role FROM users WHERE role IS NOT NULL`
	err := db.SelectContext(ctx, &users, query)
	if err != nil {
		return fmt.Errorf("failed to get existing users: %w", err)
	}

	// Get role IDs for mapping
	roleMap := make(map[string]string)
	roles, err := service.ListRoles(ctx)
	if err != nil {
		return fmt.Errorf("failed to get roles: %w", err)
	}

	for _, role := range roles {
		roleMap[role.Name] = role.ID
	}

	// Migrate each user
	migratedCount := 0
	for _, user := range users {
		// Map old role names to new role names
		newRoleName := mapOldRoleToNew(user.Role)
		roleID, exists := roleMap[newRoleName]
		if !exists {
			log.Printf("Warning: role %s not found for user %s", newRoleName, user.ID)
			continue
		}

		// Check if user already has this role assigned
		var hasRole bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM user_roles WHERE user_id = $1 AND role_id = $2 AND is_active = true)`
		err := db.GetContext(ctx, &hasRole, checkQuery, user.ID, roleID)
		if err != nil {
			log.Printf("Warning: failed to check existing role for user %s: %v", user.ID, err)
			continue
		}

		if hasRole {
			continue // User already has this role
		}

		// Assign the role
		err = service.AssignRole(ctx, user.ID, roleID, "system")
		if err != nil {
			log.Printf("Warning: failed to assign role to user %s: %v", user.ID, err)
			continue
		}

		migratedCount++
	}

	log.Printf("Successfully migrated %d users to new RBAC system", migratedCount)
	return nil
}

// mapOldRoleToNew maps old role names to new RBAC role names
func mapOldRoleToNew(oldRole string) string {
	switch oldRole {
	case "admin", "administrator":
		return RoleAdmin
	case "technician", "tech":
		return RoleTechnician
	case "user", "viewer", "client":
		return RoleViewer
	default:
		// Default to viewer for unknown roles
		return RoleViewer
	}
}

// CreateDefaultPolicies creates default policies for the system
func CreateDefaultPolicies(service *Service) error {
	ctx := context.Background()

	// Create business hours policy
	businessHoursConditions, _ := json.Marshal(map[string]interface{}{
		"allowed_hours": []map[string]string{{"start": "08:00", "end": "18:00"}},
		"allowed_days":  []int{1, 2, 3, 4, 5},
		"timezone":      "America/Sao_Paulo",
	})
	businessHoursPolicy := Policy{
		ID:         generateUUID(),
		Name:       "business_hours_only",
		Type:       PolicyTypeTimeRestriction,
		Resource:   "*",
		Action:     "*",
		Conditions: businessHoursConditions,
		Effect:     PolicyEffectAllow,
		Priority:   100,
		IsActive:   true,
	}

	// Create resource ownership policy
	ownershipConditions, _ := json.Marshal(map[string]interface{}{
		"require_ownership":   true,
		"allowed_owner_roles": []string{"admin", "technician"},
	})
	ownershipPolicy := Policy{
		ID:         generateUUID(),
		Name:       "resource_ownership_required",
		Type:       PolicyTypeResourceOwnership,
		Resource:   "*",
		Action:     "*",
		Conditions: ownershipConditions,
		Effect:     PolicyEffectAllow,
		Priority:   200,
		IsActive:   true,
	}

	// Create IP restriction policy for admin access
	adminIPConditions, _ := json.Marshal(map[string]interface{}{
		"allowed_ips": []string{"192.168.1.*", "10.0.0.*", "127.0.0.1"},
	})
	adminIPPolicy := Policy{
		ID:         generateUUID(),
		Name:       "admin_ip_restriction",
		Type:       PolicyTypeIPRestriction,
		Resource:   "*",
		Action:     "admin",
		Conditions: adminIPConditions,
		Effect:     PolicyEffectAllow,
		Priority:   300,
		IsActive:   false, // Disabled by default
	}

	policies := []Policy{businessHoursPolicy, ownershipPolicy, adminIPPolicy}

	for _, policy := range policies {
		err := createPolicyIfNotExists(service, ctx, policy)
		if err != nil {
			log.Printf("Warning: failed to create policy %s: %v", policy.Name, err)
		}
	}

	return nil
}

// createPolicyIfNotExists creates a policy if it doesn't already exist
func createPolicyIfNotExists(service *Service, ctx context.Context, policy Policy) error {
	// Check if policy already exists
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM policies WHERE name = $1)`
	err := service.db.GetContext(ctx, &exists, query, policy.Name)
	if err != nil {
		return fmt.Errorf("failed to check policy existence: %w", err)
	}

	if exists {
		return nil // Policy already exists
	}

	// Create the policy
	insertQuery := `
		INSERT INTO policies (id, name, type, resource, action, conditions, effect, priority, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	now := time.Now()
	_, err = service.db.ExecContext(ctx, insertQuery,
		policy.ID, policy.Name, policy.Type, policy.Resource, policy.Action,
		policy.Conditions, policy.Effect, policy.Priority, policy.IsActive, now, now)
	if err != nil {
		return fmt.Errorf("failed to create policy: %w", err)
	}

	log.Printf("Created policy: %s", policy.Name)
	return nil
}

// generateUUID generates a new UUID (simplified version)
func generateUUID() string {
	// In a real implementation, you would use a proper UUID library
	// For now, we'll use a placeholder
	return "uuid_placeholder"
}

// SetupDefaultRoleAssignments sets up default role assignments for system users
func SetupDefaultRoleAssignments(service *Service) error {
	ctx := context.Background()

	// Find admin users (users with admin role in old system)
	var adminUsers []string
	query := `SELECT id FROM users WHERE role = 'admin' OR email LIKE '%admin%'`
	err := service.db.SelectContext(ctx, &adminUsers, query)
	if err != nil {
		return fmt.Errorf("failed to find admin users: %w", err)
	}

	// Get admin role ID
	adminRole, err := service.GetRole(ctx, RoleAdmin)
	if err != nil {
		return fmt.Errorf("failed to get admin role: %w", err)
	}

	// Assign admin role to admin users
	for _, userID := range adminUsers {
		err := service.AssignRole(ctx, userID, adminRole.ID, "system")
		if err != nil {
			log.Printf("Warning: failed to assign admin role to user %s: %v", userID, err)
		}
	}

	log.Printf("Set up default role assignments for %d admin users", len(adminUsers))
	return nil
}

// ValidatePermissionSystem validates that the permission system is working correctly
func ValidatePermissionSystem(service *Service) error {
	ctx := context.Background()

	// Check that default roles exist
	defaultRoles := []string{RoleAdmin, RoleTechnician, RoleViewer}
	for _, roleName := range defaultRoles {
		roles, err := service.ListRoles(ctx)
		if err != nil {
			return fmt.Errorf("failed to list roles: %w", err)
		}

		found := false
		for _, role := range roles {
			if role.Name == roleName {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("default role %s not found", roleName)
		}
	}

	// Check that basic permissions exist
	basicPermissions := []struct {
		resource string
		action   string
	}{
		{ResourceDrone, ActionRead},
		{ResourceMaintenanceOrder, ActionCreate},
		{ResourceTechnicalReport, ActionUpdate},
		{ResourceUser, ActionList},
	}

	for _, perm := range basicPermissions {
		var exists bool
		query := `SELECT EXISTS(SELECT 1 FROM permissions WHERE resource = $1 AND action = $2)`
		err := service.db.GetContext(ctx, &exists, query, perm.resource, perm.action)
		if err != nil {
			return fmt.Errorf("failed to check permission %s:%s: %w", perm.resource, perm.action, err)
		}

		if !exists {
			return fmt.Errorf("basic permission %s:%s not found", perm.resource, perm.action)
		}
	}

	log.Println("Permission system validation completed successfully")
	return nil
}
