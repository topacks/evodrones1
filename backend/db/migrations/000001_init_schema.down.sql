-- +migrate Down
DROP TABLE IF EXISTS maintenance_media;
DROP TABLE IF EXISTS maintenance_costs;
DROP TABLE IF EXISTS technical_reports;
DROP TABLE IF EXISTS maintenance_orders;
DROP TABLE IF EXISTS drone_items;
DROP TABLE IF EXISTS drones;
DROP TABLE IF EXISTS users; 