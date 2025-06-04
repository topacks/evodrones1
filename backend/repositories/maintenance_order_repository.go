package repositories

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/root/evodrones/backend/models"
)

type MaintenanceOrderRepository struct {
	db *sqlx.DB
}

func NewMaintenanceOrderRepository(db *sqlx.DB) *MaintenanceOrderRepository {
	return &MaintenanceOrderRepository{
		db: db,
	}
}

func (r *MaintenanceOrderRepository) Create(order *models.MaintenanceOrder) error {
	query := `
		INSERT INTO maintenance_orders (
			drone_id, user_id, status, description, start_date, created_at, updated_at
		) VALUES (
			:drone_id, :user_id, :status, :description, :start_date, :created_at, :updated_at
		) RETURNING id
	`

	order.Status = models.MaintenanceOrderStatusPending
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRow(order).Scan(&order.ID)
}

func (r *MaintenanceOrderRepository) GetByID(id string) (*models.MaintenanceOrderWithDetails, error) {
	query := `
		SELECT 
			mo.*,
			d.id as "drone.id", d.model as "drone.model", d.serial as "drone.serial",
			d.status as "drone.status", d.last_check as "drone.last_check",
			d.next_check as "drone.next_check", d.created_at as "drone.created_at",
			d.updated_at as "drone.updated_at",
			u.id as "user.id", u.name as "user.name", u.email as "user.email",
			u.role as "user.role", u.created_at as "user.created_at",
			u.updated_at as "user.updated_at"
		FROM maintenance_orders mo
		LEFT JOIN drones d ON d.id = mo.drone_id
		LEFT JOIN users u ON u.id = mo.user_id
		WHERE mo.id = $1
	`

	var order models.MaintenanceOrderWithDetails
	err := r.db.Get(&order, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *MaintenanceOrderRepository) GetByUserID(userID string) ([]*models.MaintenanceOrderWithDetails, error) {
	query := `
		SELECT 
			mo.*,
			d.id as "drone.id", d.model as "drone.model", d.serial as "drone.serial",
			d.status as "drone.status", d.last_check as "drone.last_check",
			d.next_check as "drone.next_check", d.created_at as "drone.created_at",
			d.updated_at as "drone.updated_at",
			u.id as "user.id", u.name as "user.name", u.email as "user.email",
			u.role as "user.role", u.created_at as "user.created_at",
			u.updated_at as "user.updated_at"
		FROM maintenance_orders mo
		LEFT JOIN drones d ON d.id = mo.drone_id
		LEFT JOIN users u ON u.id = mo.user_id
		WHERE mo.user_id = $1
		ORDER BY mo.created_at DESC
	`

	var orders []*models.MaintenanceOrderWithDetails
	err := r.db.Select(&orders, query, userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *MaintenanceOrderRepository) GetByDroneID(droneID string) ([]*models.MaintenanceOrderWithDetails, error) {
	query := `
		SELECT 
			mo.*,
			d.id as "drone.id", d.model as "drone.model", d.serial as "drone.serial",
			d.status as "drone.status", d.last_check as "drone.last_check",
			d.next_check as "drone.next_check", d.created_at as "drone.created_at",
			d.updated_at as "drone.updated_at",
			u.id as "user.id", u.name as "user.name", u.email as "user.email",
			u.role as "user.role", u.created_at as "user.created_at",
			u.updated_at as "user.updated_at"
		FROM maintenance_orders mo
		LEFT JOIN drones d ON d.id = mo.drone_id
		LEFT JOIN users u ON u.id = mo.user_id
		WHERE mo.drone_id = $1
		ORDER BY mo.created_at DESC
	`

	var orders []*models.MaintenanceOrderWithDetails
	err := r.db.Select(&orders, query, droneID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *MaintenanceOrderRepository) Update(id string, input *models.UpdateMaintenanceOrderInput) error {
	query := `
		UPDATE maintenance_orders
		SET 
			status = COALESCE($1, status),
			description = COALESCE($2, description),
			end_date = COALESCE($3, end_date),
			updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(query,
		input.Status,
		input.Description,
		input.EndDate,
		time.Now(),
		id,
	)

	return err
}

func (r *MaintenanceOrderRepository) Delete(id string) error {
	query := `DELETE FROM maintenance_orders WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *MaintenanceOrderRepository) List(offset, limit int) ([]*models.MaintenanceOrderWithDetails, error) {
	query := `
		SELECT 
			mo.*,
			d.id as "drone.id", d.model as "drone.model", d.serial as "drone.serial",
			d.status as "drone.status", d.last_check as "drone.last_check",
			d.next_check as "drone.next_check", d.created_at as "drone.created_at",
			d.updated_at as "drone.updated_at",
			u.id as "user.id", u.name as "user.name", u.email as "user.email",
			u.role as "user.role", u.created_at as "user.created_at",
			u.updated_at as "user.updated_at"
		FROM maintenance_orders mo
		LEFT JOIN drones d ON d.id = mo.drone_id
		LEFT JOIN users u ON u.id = mo.user_id
		ORDER BY mo.created_at DESC
		LIMIT $1 OFFSET $2
	`

	var orders []*models.MaintenanceOrderWithDetails
	err := r.db.Select(&orders, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *MaintenanceOrderRepository) GetByStatus(status models.MaintenanceOrderStatus) ([]*models.MaintenanceOrderWithDetails, error) {
	query := `
		SELECT 
			mo.*,
			d.id as "drone.id", d.model as "drone.model", d.serial as "drone.serial",
			d.status as "drone.status", d.last_check as "drone.last_check",
			d.next_check as "drone.next_check", d.created_at as "drone.created_at",
			d.updated_at as "drone.updated_at",
			u.id as "user.id", u.name as "user.name", u.email as "user.email",
			u.role as "user.role", u.created_at as "user.created_at",
			u.updated_at as "user.updated_at"
		FROM maintenance_orders mo
		LEFT JOIN drones d ON d.id = mo.drone_id
		LEFT JOIN users u ON u.id = mo.user_id
		WHERE mo.status = $1
		ORDER BY mo.created_at DESC
	`

	var orders []*models.MaintenanceOrderWithDetails
	err := r.db.Select(&orders, query, status)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
