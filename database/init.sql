-- Criação das extensões necessárias
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Criação dos tipos ENUM
CREATE TYPE item_type AS ENUM ('bag', 'radio_control', 'batteries', 'propellers', 'charging_hub', 'strap', 'drone');
CREATE TYPE order_status AS ENUM ('created', 'analyzed', 'in_maintenance', 'approved_paid', 'completed');
CREATE TYPE user_role AS ENUM ('client', 'technician');
CREATE TYPE media_type AS ENUM ('photo', 'video');

-- Tabela de usuários
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    role user_role NOT NULL DEFAULT 'client',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Índice para busca por email
CREATE INDEX idx_users_email ON users(email);

-- Tabela de drones
CREATE TABLE drones (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    serial_number VARCHAR(100) NOT NULL,
    color VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Índice para busca de drones por usuário
CREATE INDEX idx_drones_user_id ON drones(user_id);

-- Tabela de itens do drone
CREATE TABLE drone_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    drone_id UUID NOT NULL REFERENCES drones(id) ON DELETE CASCADE,
    item_type item_type NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(drone_id, item_type)
);

-- Índice para busca de itens por drone
CREATE INDEX idx_drone_items_drone_id ON drone_items(drone_id);

-- Tabela de ordens de manutenção
CREATE TABLE maintenance_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    drone_id UUID NOT NULL REFERENCES drones(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status order_status NOT NULL DEFAULT 'created',
    defect_description TEXT NOT NULL,
    broken_parts JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Índices para busca de ordens
CREATE INDEX idx_maintenance_orders_drone_id ON maintenance_orders(drone_id);
CREATE INDEX idx_maintenance_orders_user_id ON maintenance_orders(user_id);
CREATE INDEX idx_maintenance_orders_status ON maintenance_orders(status);

-- Tabela de relatórios técnicos
CREATE TABLE technical_reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    maintenance_order_id UUID NOT NULL REFERENCES maintenance_orders(id) ON DELETE CASCADE,
    technician_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    report TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Índice para busca de relatórios por ordem
CREATE INDEX idx_technical_reports_maintenance_order_id ON technical_reports(maintenance_order_id);
CREATE INDEX idx_technical_reports_technician_id ON technical_reports(technician_id);

-- Tabela de custos de manutenção
CREATE TABLE maintenance_costs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    maintenance_order_id UUID NOT NULL REFERENCES maintenance_orders(id) ON DELETE CASCADE,
    labor_cost DECIMAL(10, 2) NOT NULL CHECK (labor_cost >= 0),
    parts JSONB NOT NULL,
    total_cost DECIMAL(10, 2) NOT NULL CHECK (total_cost >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Índice para busca de custos por ordem
CREATE INDEX idx_maintenance_costs_maintenance_order_id ON maintenance_costs(maintenance_order_id);

-- Tabela de mídias de manutenção
CREATE TABLE maintenance_media (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    maintenance_order_id UUID NOT NULL REFERENCES maintenance_orders(id) ON DELETE CASCADE,
    type media_type NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Índice para busca de mídias por ordem
CREATE INDEX idx_maintenance_media_maintenance_order_id ON maintenance_media(maintenance_order_id);

-- Tabela de comentários
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    maintenance_order_id UUID NOT NULL REFERENCES maintenance_orders(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Índice para busca de comentários por ordem
CREATE INDEX idx_comments_maintenance_order_id ON comments(maintenance_order_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);

-- Função para atualizar o timestamp de updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers para atualizar o timestamp de updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_drones_updated_at BEFORE UPDATE ON drones FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_drone_items_updated_at BEFORE UPDATE ON drone_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_maintenance_orders_updated_at BEFORE UPDATE ON maintenance_orders FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_technical_reports_updated_at BEFORE UPDATE ON technical_reports FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_maintenance_costs_updated_at BEFORE UPDATE ON maintenance_costs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Inserir um usuário técnico padrão para testes
INSERT INTO users (email, password_hash, name, phone, role) VALUES 
('tecnico@evodrones.com', '$2a$10$1qAz2wSx3eDc4rFv5tGb5eT3TEwmLYbUAYNldGLEGt4WCMeRMqlte', 'Técnico Padrão', '(11) 99999-9999', 'technician');

-- Inserir um usuário cliente padrão para testes
INSERT INTO users (email, password_hash, name, phone, role) VALUES 
('cliente@evodrones.com', '$2a$10$1qAz2wSx3eDc4rFv5tGb5eT3TEwmLYbUAYNldGLEGt4WCMeRMqlte', 'Cliente Padrão', '(11) 88888-8888', 'client');

-- Nota: O hash da senha corresponde a 'senha123' usando bcrypt