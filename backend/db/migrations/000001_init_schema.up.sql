-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tabela de usuários
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    role VARCHAR(20) NOT NULL CHECK (role IN ('client', 'technician')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de drones
CREATE TABLE drones (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    model VARCHAR(255) NOT NULL,
    serial_number VARCHAR(255) UNIQUE NOT NULL,
    color VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de itens do drone
CREATE TABLE drone_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    drone_id UUID NOT NULL REFERENCES drones(id),
    item_type VARCHAR(50) NOT NULL CHECK (item_type IN ('bag', 'radio_control', 'batteries', 'propellers', 'charging_hub', 'strap', 'drone')),
    quantity INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de ordens de manutenção
CREATE TABLE maintenance_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    drone_id UUID NOT NULL REFERENCES drones(id),
    user_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(50) NOT NULL CHECK (status IN ('created', 'analyzed', 'in_maintenance', 'approved_paid', 'completed')),
    defect_description TEXT,
    broken_parts JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de relatórios técnicos
CREATE TABLE technical_reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    maintenance_order_id UUID NOT NULL REFERENCES maintenance_orders(id),
    technician_id UUID NOT NULL REFERENCES users(id),
    report TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de custos de manutenção
CREATE TABLE maintenance_costs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    maintenance_order_id UUID NOT NULL REFERENCES maintenance_orders(id),
    labor_cost DECIMAL(10,2) NOT NULL,
    parts JSONB,
    total_cost DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de mídias de manutenção
CREATE TABLE maintenance_media (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    maintenance_order_id UUID NOT NULL REFERENCES maintenance_orders(id),
    type VARCHAR(20) NOT NULL CHECK (type IN ('photo', 'video')),
    url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índices
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_drones_user_id ON drones(user_id);
CREATE INDEX idx_drones_serial_number ON drones(serial_number);
CREATE INDEX idx_drone_items_drone_id ON drone_items(drone_id);
CREATE INDEX idx_maintenance_orders_drone_id ON maintenance_orders(drone_id);
CREATE INDEX idx_maintenance_orders_user_id ON maintenance_orders(user_id);
CREATE INDEX idx_maintenance_orders_status ON maintenance_orders(status);
CREATE INDEX idx_technical_reports_maintenance_order_id ON technical_reports(maintenance_order_id);
CREATE INDEX idx_maintenance_costs_maintenance_order_id ON maintenance_costs(maintenance_order_id);
CREATE INDEX idx_maintenance_media_maintenance_order_id ON maintenance_media(maintenance_order_id); 