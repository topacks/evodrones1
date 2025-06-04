-- Migration: Create RBAC (Role-Based Access Control) Tables
-- Description: Creates the unified permission system tables
-- Author: EVO DRONES System
-- Date: 2024

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    parent_role_id UUID REFERENCES roles(id) ON DELETE SET NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(resource, action)
);

-- Create role_permissions junction table
CREATE TABLE IF NOT EXISTS role_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(role_id, permission_id)
);

-- Create user_roles junction table
CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_by UUID REFERENCES users(id) ON DELETE SET NULL,
    assigned_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, role_id)
);

-- Create policies table for advanced permission rules
CREATE TABLE IF NOT EXISTS policies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(200) NOT NULL UNIQUE,
    description TEXT,
    type VARCHAR(50) NOT NULL, -- time_restriction, resource_ownership, ip_restriction, custom
    conditions JSONB NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create permission_policies junction table
CREATE TABLE IF NOT EXISTS permission_policies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    policy_id UUID NOT NULL REFERENCES policies(id) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(permission_id, policy_id)
);

-- Create audit table for permission changes
CREATE TABLE IF NOT EXISTS permission_audit (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL, -- role_assigned, role_revoked, permission_granted, permission_revoked
    resource_type VARCHAR(100),
    resource_id UUID,
    old_value JSONB,
    new_value JSONB,
    performed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);
CREATE INDEX IF NOT EXISTS idx_roles_parent ON roles(parent_role_id);
CREATE INDEX IF NOT EXISTS idx_roles_active ON roles(is_active);

CREATE INDEX IF NOT EXISTS idx_permissions_resource ON permissions(resource);
CREATE INDEX IF NOT EXISTS idx_permissions_action ON permissions(action);
CREATE INDEX IF NOT EXISTS idx_permissions_resource_action ON permissions(resource, action);

CREATE INDEX IF NOT EXISTS idx_role_permissions_role ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission ON role_permissions(permission_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_active ON role_permissions(is_active);

CREATE INDEX IF NOT EXISTS idx_user_roles_user ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role ON user_roles(role_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_active ON user_roles(is_active);
CREATE INDEX IF NOT EXISTS idx_user_roles_assigned_at ON user_roles(assigned_at);

CREATE INDEX IF NOT EXISTS idx_policies_type ON policies(type);
CREATE INDEX IF NOT EXISTS idx_policies_active ON policies(is_active);

CREATE INDEX IF NOT EXISTS idx_permission_policies_permission ON permission_policies(permission_id);
CREATE INDEX IF NOT EXISTS idx_permission_policies_policy ON permission_policies(policy_id);

CREATE INDEX IF NOT EXISTS idx_permission_audit_user ON permission_audit(user_id);
CREATE INDEX IF NOT EXISTS idx_permission_audit_action ON permission_audit(action);
CREATE INDEX IF NOT EXISTS idx_permission_audit_resource ON permission_audit(resource_type, resource_id);
CREATE INDEX IF NOT EXISTS idx_permission_audit_created_at ON permission_audit(created_at);

-- Create triggers for updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_roles_updated_at BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_permissions_updated_at BEFORE UPDATE ON permissions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_role_permissions_updated_at BEFORE UPDATE ON role_permissions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_roles_updated_at BEFORE UPDATE ON user_roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_policies_updated_at BEFORE UPDATE ON policies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert default roles
INSERT INTO roles (id, name, description) VALUES 
    (uuid_generate_v4(), 'admin', 'Administrador do sistema com acesso total'),
    (uuid_generate_v4(), 'technician', 'Técnico com acesso a manutenções e relatórios'),
    (uuid_generate_v4(), 'viewer', 'Visualizador com acesso somente leitura')
ON CONFLICT (name) DO NOTHING;

-- Insert default permissions
INSERT INTO permissions (resource, action, description) VALUES 
    -- Drone permissions
    ('drone', 'create', 'Criar novos drones'),
    ('drone', 'read', 'Visualizar informações de drones'),
    ('drone', 'update', 'Atualizar informações de drones'),
    ('drone', 'delete', 'Excluir drones'),
    ('drone', 'list', 'Listar drones'),
    
    -- Maintenance Order permissions
    ('maintenance_order', 'create', 'Criar ordens de manutenção'),
    ('maintenance_order', 'read', 'Visualizar ordens de manutenção'),
    ('maintenance_order', 'update', 'Atualizar ordens de manutenção'),
    ('maintenance_order', 'delete', 'Excluir ordens de manutenção'),
    ('maintenance_order', 'list', 'Listar ordens de manutenção'),
    
    -- Technical Report permissions
    ('technical_report', 'create', 'Criar relatórios técnicos'),
    ('technical_report', 'read', 'Visualizar relatórios técnicos'),
    ('technical_report', 'update', 'Atualizar relatórios técnicos'),
    ('technical_report', 'delete', 'Excluir relatórios técnicos'),
    ('technical_report', 'list', 'Listar relatórios técnicos'),
    
    -- User permissions
    ('user', 'create', 'Criar novos usuários'),
    ('user', 'read', 'Visualizar informações de usuários'),
    ('user', 'update', 'Atualizar informações de usuários'),
    ('user', 'delete', 'Excluir usuários'),
    ('user', 'list', 'Listar usuários'),
    
    -- System permissions
    ('system', 'admin', 'Acesso administrativo ao sistema'),
    ('system', 'backup', 'Realizar backup do sistema'),
    ('system', 'restore', 'Restaurar backup do sistema'),
    ('system', 'audit', 'Visualizar logs de auditoria')
ON CONFLICT (resource, action) DO NOTHING;

-- Assign permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign limited permissions to technician role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'technician'
AND p.resource IN ('maintenance_order', 'technical_report', 'drone')
AND p.action IN ('create', 'read', 'update', 'list')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign read-only permissions to viewer role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'viewer'
AND p.action IN ('read', 'list')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Create view for easy permission checking
CREATE OR REPLACE VIEW user_permissions AS
SELECT 
    ur.user_id,
    u.email as user_email,
    r.name as role_name,
    p.resource,
    p.action,
    p.description as permission_description
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
JOIN users u ON ur.user_id = u.id
WHERE ur.is_active = true 
AND r.is_active = true 
AND rp.is_active = true;

-- Create view for role hierarchy
CREATE OR REPLACE VIEW role_hierarchy AS
WITH RECURSIVE role_tree AS (
    -- Base case: roles without parents
    SELECT id, name, description, parent_role_id, 0 as level, ARRAY[name] as path
    FROM roles
    WHERE parent_role_id IS NULL AND is_active = true
    
    UNION ALL
    
    -- Recursive case: roles with parents
    SELECT r.id, r.name, r.description, r.parent_role_id, rt.level + 1, rt.path || r.name
    FROM roles r
    JOIN role_tree rt ON r.parent_role_id = rt.id
    WHERE r.is_active = true
)
SELECT * FROM role_tree ORDER BY level, name;

-- Create function to check user permission
CREATE OR REPLACE FUNCTION check_user_permission(
    p_user_id UUID,
    p_resource VARCHAR,
    p_action VARCHAR
) RETURNS BOOLEAN AS $$
DECLARE
    permission_exists BOOLEAN := FALSE;
BEGIN
    SELECT EXISTS(
        SELECT 1
        FROM user_permissions up
        WHERE up.user_id = p_user_id
        AND up.resource = p_resource
        AND up.action = p_action
    ) INTO permission_exists;
    
    RETURN permission_exists;
END;
$$ LANGUAGE plpgsql;

-- Create function to get user roles
CREATE OR REPLACE FUNCTION get_user_roles(p_user_id UUID)
RETURNS TABLE(role_id UUID, role_name VARCHAR, role_description TEXT) AS $$
BEGIN
    RETURN QUERY
    SELECT r.id, r.name, r.description
    FROM user_roles ur
    JOIN roles r ON ur.role_id = r.id
    WHERE ur.user_id = p_user_id
    AND ur.is_active = true
    AND r.is_active = true
    ORDER BY r.name;
END;
$$ LANGUAGE plpgsql;

-- Create function to audit permission changes
CREATE OR REPLACE FUNCTION audit_permission_change(
    p_user_id UUID,
    p_action VARCHAR,
    p_resource_type VARCHAR DEFAULT NULL,
    p_resource_id UUID DEFAULT NULL,
    p_old_value JSONB DEFAULT NULL,
    p_new_value JSONB DEFAULT NULL,
    p_performed_by UUID DEFAULT NULL,
    p_ip_address INET DEFAULT NULL,
    p_user_agent TEXT DEFAULT NULL
) RETURNS UUID AS $$
DECLARE
    audit_id UUID;
BEGIN
    INSERT INTO permission_audit (
        user_id, action, resource_type, resource_id,
        old_value, new_value, performed_by, ip_address, user_agent
    ) VALUES (
        p_user_id, p_action, p_resource_type, p_resource_id,
        p_old_value, p_new_value, p_performed_by, p_ip_address, p_user_agent
    ) RETURNING id INTO audit_id;
    
    RETURN audit_id;
END;
$$ LANGUAGE plpgsql;

COMMIT;