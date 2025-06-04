package permissions

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DefaultPolicyEvaluator implements the PolicyEvaluator interface
type DefaultPolicyEvaluator struct{}

// NewDefaultPolicyEvaluator creates a new default policy evaluator
func NewDefaultPolicyEvaluator() *DefaultPolicyEvaluator {
	return &DefaultPolicyEvaluator{}
}

// EvaluatePolicy evaluates a single policy against an access request
func (pe *DefaultPolicyEvaluator) EvaluatePolicy(ctx context.Context, policy Policy, req AccessRequest) (bool, error) {
	return pe.evaluatePolicy(ctx, policy, req.Context)
}

// EvaluatePolicies evaluates multiple policies against an access request
func (pe *DefaultPolicyEvaluator) EvaluatePolicies(ctx context.Context, policies []Policy, req AccessRequest) (*AccessResult, error) {
	result := &AccessResult{
		Allowed:         true,
		AppliedPolicies: make([]string, 0),
	}

	// If no policies, allow by default
	if len(policies) == 0 {
		result.Reason = "No policies defined"
		return result, nil
	}

	// Evaluate each policy
	for _, policy := range policies {
		allowed, err := pe.EvaluatePolicy(ctx, policy, req)
		if err != nil {
			result.Allowed = false
			result.Reason = fmt.Sprintf("Policy evaluation failed: %v", err)
			return result, nil
		}

		result.AppliedPolicies = append(result.AppliedPolicies, policy.Name)

		if policy.Effect == PolicyEffectDeny && allowed {
			result.Allowed = false
			result.Reason = fmt.Sprintf("Denied by policy: %s", policy.Name)
			return result, nil
		}

		if policy.Effect == PolicyEffectAllow && !allowed {
			result.Allowed = false
			result.Reason = fmt.Sprintf("Not allowed by policy: %s", policy.Name)
			return result, nil
		}
	}

	result.Reason = "All policies passed"
	return result, nil
}

// Evaluate evaluates a permission against the provided context (legacy method)
func (pe *DefaultPolicyEvaluator) Evaluate(ctx context.Context, permission Permission, context map[string]interface{}) (bool, error) {
	// Convert to new format for compatibility
	req := AccessRequest{
		Resource: permission.Resource,
		Action:   permission.Action,
		Context:  context,
	}

	// Create policies from permission (if any)
	var policies []Policy
	// Note: This is a simplified conversion - in practice, permissions might have associated policies

	result, err := pe.EvaluatePolicies(ctx, policies, req)
	if err != nil {
		return false, err
	}

	return result.Allowed, nil
}

// evaluatePolicy evaluates a single policy
func (pe *DefaultPolicyEvaluator) evaluatePolicy(ctx context.Context, policy Policy, context map[string]interface{}) (bool, error) {
	switch policy.Type {
	case PolicyTypeTimeRestriction:
		return pe.evaluateTimeRestriction(policy, context)
	case PolicyTypeResourceOwnership:
		return pe.evaluateResourceOwnership(policy, context)
	case PolicyTypeIPRestriction:
		return pe.evaluateIPRestriction(policy, context)
	case PolicyTypeCustom:
		return pe.evaluateCustomPolicy(policy, context)
	default:
		return false, fmt.Errorf("unknown policy type: %s", policy.Type)
	}
}

// evaluateTimeRestriction evaluates time-based access restrictions
func (pe *DefaultPolicyEvaluator) evaluateTimeRestriction(policy Policy, context map[string]interface{}) (bool, error) {
	var conditions TimeRestrictionConditions
	err := json.Unmarshal(policy.Conditions, &conditions)
	if err != nil {
		return false, fmt.Errorf("failed to parse time restriction conditions: %w", err)
	}

	now := time.Now()
	currentTime := now.Format("15:04")
	currentWeekday := int(now.Weekday())

	// Check allowed hours
	if conditions.AllowedHours != nil {
		allowed := false
		for _, timeRange := range conditions.AllowedHours {
			if currentTime >= timeRange.Start && currentTime <= timeRange.End {
				allowed = true
				break
			}
		}
		if !allowed {
			return false, nil
		}
	}

	// Check allowed days
	if len(conditions.AllowedDays) > 0 {
		allowed := false
		for _, day := range conditions.AllowedDays {
			if day == currentWeekday {
				allowed = true
				break
			}
		}
		if !allowed {
			return false, nil
		}
	}

	// Check timezone
	if conditions.Timezone != "" {
		loc, err := time.LoadLocation(conditions.Timezone)
		if err != nil {
			return false, fmt.Errorf("invalid timezone: %w", err)
		}
		localTime := now.In(loc)
		localTimeStr := localTime.Format("15:04")
		localWeekday := int(localTime.Weekday())

		// Re-check with local time
		if conditions.AllowedHours != nil {
			allowed := false
			for _, timeRange := range conditions.AllowedHours {
				if localTimeStr >= timeRange.Start && localTimeStr <= timeRange.End {
					allowed = true
					break
				}
			}
			if !allowed {
				return false, nil
			}
		}

		if len(conditions.AllowedDays) > 0 {
			allowed := false
			for _, day := range conditions.AllowedDays {
				if day == localWeekday {
					allowed = true
					break
				}
			}
			if !allowed {
				return false, nil
			}
		}
	}

	return true, nil
}

// evaluateResourceOwnership evaluates resource ownership restrictions
func (pe *DefaultPolicyEvaluator) evaluateResourceOwnership(policy Policy, context map[string]interface{}) (bool, error) {
	var conditions ResourceOwnershipConditions
	err := json.Unmarshal(policy.Conditions, &conditions)
	if err != nil {
		return false, fmt.Errorf("failed to parse resource ownership conditions: %w", err)
	}

	// Get user ID from context
	userID, ok := context["user_id"].(string)
	if !ok {
		return false, fmt.Errorf("user_id not found in context")
	}

	// Check if user owns the resource
	if conditions.RequireOwnership {
		resourceOwnerID, ok := context["resource_owner_id"].(string)
		if !ok {
			return false, fmt.Errorf("resource_owner_id not found in context")
		}
		if userID != resourceOwnerID {
			return false, nil
		}
	}

	// Check allowed owner roles
	if len(conditions.AllowedOwnerRoles) > 0 {
		userRole, ok := context["user_role"].(string)
		if !ok {
			return false, fmt.Errorf("user_role not found in context")
		}

		allowed := false
		for _, allowedRole := range conditions.AllowedOwnerRoles {
			if userRole == allowedRole {
				allowed = true
				break
			}
		}
		if !allowed {
			return false, nil
		}
	}

	return true, nil
}

// evaluateIPRestriction evaluates IP-based access restrictions
func (pe *DefaultPolicyEvaluator) evaluateIPRestriction(policy Policy, context map[string]interface{}) (bool, error) {
	var conditions IPRestrictionConditions
	err := json.Unmarshal(policy.Conditions, &conditions)
	if err != nil {
		return false, fmt.Errorf("failed to parse IP restriction conditions: %w", err)
	}

	// Get client IP from context
	clientIP, ok := context["client_ip"].(string)
	if !ok {
		return false, fmt.Errorf("client_ip not found in context")
	}

	// Check allowed IPs
	if len(conditions.AllowedIPs) > 0 {
		allowed := false
		for _, allowedIP := range conditions.AllowedIPs {
			if pe.matchesIPPattern(clientIP, allowedIP) {
				allowed = true
				break
			}
		}
		if !allowed {
			return false, nil
		}
	}

	// Check blocked IPs
	if len(conditions.BlockedIPs) > 0 {
		for _, blockedIP := range conditions.BlockedIPs {
			if pe.matchesIPPattern(clientIP, blockedIP) {
				return false, nil
			}
		}
	}

	return true, nil
}

// evaluateCustomPolicy evaluates custom policy conditions
func (pe *DefaultPolicyEvaluator) evaluateCustomPolicy(policy Policy, context map[string]interface{}) (bool, error) {
	var conditions CustomPolicyConditions
	err := json.Unmarshal(policy.Conditions, &conditions)
	if err != nil {
		return false, fmt.Errorf("failed to parse custom policy conditions: %w", err)
	}

	// Evaluate each rule
	for _, rule := range conditions.Rules {
		allowed, err := pe.evaluateCustomRule(rule, context)
		if err != nil {
			return false, fmt.Errorf("custom rule evaluation failed: %w", err)
		}
		if !allowed {
			return false, nil
		}
	}

	return true, nil
}

// evaluateCustomRule evaluates a single custom rule
func (pe *DefaultPolicyEvaluator) evaluateCustomRule(rule CustomRule, context map[string]interface{}) (bool, error) {
	// Get the field value from context
	fieldValue, exists := context[rule.Field]
	if !exists {
		if rule.Required {
			return false, fmt.Errorf("required field %s not found in context", rule.Field)
		}
		return true, nil // Optional field not present, allow
	}

	// Convert field value to string for comparison
	fieldStr := fmt.Sprintf("%v", fieldValue)

	// Evaluate based on operator
	switch rule.Operator {
	case "equals":
		return fieldStr == rule.Value, nil
	case "not_equals":
		return fieldStr != rule.Value, nil
	case "contains":
		return strings.Contains(fieldStr, rule.Value), nil
	case "not_contains":
		return !strings.Contains(fieldStr, rule.Value), nil
	case "starts_with":
		return strings.HasPrefix(fieldStr, rule.Value), nil
	case "ends_with":
		return strings.HasSuffix(fieldStr, rule.Value), nil
	case "greater_than":
		return pe.compareNumeric(fieldStr, rule.Value, ">")
	case "less_than":
		return pe.compareNumeric(fieldStr, rule.Value, "<")
	case "greater_equal":
		return pe.compareNumeric(fieldStr, rule.Value, ">=")
	case "less_equal":
		return pe.compareNumeric(fieldStr, rule.Value, "<=")
	case "in":
		return pe.checkInList(fieldStr, rule.Value), nil
	case "not_in":
		return !pe.checkInList(fieldStr, rule.Value), nil
	default:
		return false, fmt.Errorf("unknown operator: %s", rule.Operator)
	}
}

// Helper methods

// matchesIPPattern checks if an IP matches a pattern (supports CIDR and wildcards)
func (pe *DefaultPolicyEvaluator) matchesIPPattern(ip, pattern string) bool {
	// Simple exact match
	if ip == pattern {
		return true
	}

	// Wildcard support (e.g., 192.168.1.*)
	if strings.Contains(pattern, "*") {
		patternParts := strings.Split(pattern, ".")
		ipParts := strings.Split(ip, ".")

		if len(patternParts) != len(ipParts) {
			return false
		}

		for i, part := range patternParts {
			if part != "*" && part != ipParts[i] {
				return false
			}
		}
		return true
	}

	// TODO: Add CIDR support if needed
	// For now, just exact match
	return false
}

// compareNumeric compares two numeric strings
func (pe *DefaultPolicyEvaluator) compareNumeric(field, value, operator string) (bool, error) {
	fieldNum, err := strconv.ParseFloat(field, 64)
	if err != nil {
		return false, fmt.Errorf("field value is not numeric: %s", field)
	}

	valueNum, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false, fmt.Errorf("comparison value is not numeric: %s", value)
	}

	switch operator {
	case ">":
		return fieldNum > valueNum, nil
	case "<":
		return fieldNum < valueNum, nil
	case ">=":
		return fieldNum >= valueNum, nil
	case "<=":
		return fieldNum <= valueNum, nil
	default:
		return false, fmt.Errorf("unknown numeric operator: %s", operator)
	}
}

// checkInList checks if a value is in a comma-separated list
func (pe *DefaultPolicyEvaluator) checkInList(field, listValue string) bool {
	values := strings.Split(listValue, ",")
	for _, value := range values {
		if strings.TrimSpace(value) == field {
			return true
		}
	}
	return false
}

// Policy condition structures

// TimeRestrictionConditions defines time-based access restrictions
type TimeRestrictionConditions struct {
	AllowedHours []TimeRange `json:"allowed_hours,omitempty"`
	AllowedDays  []int       `json:"allowed_days,omitempty"` // 0=Sunday, 1=Monday, etc.
	Timezone     string      `json:"timezone,omitempty"`
}

// TimeRange defines a time range
type TimeRange struct {
	Start string `json:"start"` // Format: "HH:MM"
	End   string `json:"end"`   // Format: "HH:MM"
}

// ResourceOwnershipConditions defines resource ownership restrictions
type ResourceOwnershipConditions struct {
	RequireOwnership     bool     `json:"require_ownership"`
	AllowedOwnerRoles    []string `json:"allowed_owner_roles,omitempty"`
	AllowDelegatedAccess bool     `json:"allow_delegated_access,omitempty"`
}

// IPRestrictionConditions defines IP-based access restrictions
type IPRestrictionConditions struct {
	AllowedIPs []string `json:"allowed_ips,omitempty"`
	BlockedIPs []string `json:"blocked_ips,omitempty"`
}

// CustomPolicyConditions defines custom policy rules
type CustomPolicyConditions struct {
	Rules []CustomRule `json:"rules"`
}

// CustomRule defines a single custom rule
type CustomRule struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
	Required bool   `json:"required,omitempty"`
}
