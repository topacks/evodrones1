package repositories

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/root/evodrones/backend/models"
)

type TechnicalReportRepository struct {
	db *sqlx.DB
}

func NewTechnicalReportRepository(db *sqlx.DB) *TechnicalReportRepository {
	return &TechnicalReportRepository{
		db: db,
	}
}

func (r *TechnicalReportRepository) Create(report *models.TechnicalReport) error {
	query := `
		INSERT INTO technical_reports (
			maintenance_id, user_id, title, description, findings, recommendations,
			created_at, updated_at
		) VALUES (
			:maintenance_id, :user_id, :title, :description, :findings, :recommendations,
			:created_at, :updated_at
		) RETURNING id
	`

	report.CreatedAt = time.Now()
	report.UpdatedAt = time.Now()

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRow(report).Scan(&report.ID)
}

func (r *TechnicalReportRepository) GetByID(id string) (*models.TechnicalReportWithDetails, error) {
	query := `
		SELECT 
			tr.*,
			mo.id as "maintenance_order.id", mo.drone_id as "maintenance_order.drone_id",
			mo.user_id as "maintenance_order.user_id", mo.status as "maintenance_order.status",
			mo.description as "maintenance_order.description", mo.start_date as "maintenance_order.start_date",
			mo.end_date as "maintenance_order.end_date", mo.created_at as "maintenance_order.created_at",
			mo.updated_at as "maintenance_order.updated_at",
			d.id as "maintenance_order.drone.id", d.model as "maintenance_order.drone.model", 
			d.serial as "maintenance_order.drone.serial", d.status as "maintenance_order.drone.status",
			d.last_check as "maintenance_order.drone.last_check", d.next_check as "maintenance_order.drone.next_check",
			d.created_at as "maintenance_order.drone.created_at", d.updated_at as "maintenance_order.drone.updated_at",
			u.id as "user.id", u.name as "user.name", u.email as "user.email",
			u.role as "user.role", u.created_at as "user.created_at",
			u.updated_at as "user.updated_at"
		FROM technical_reports tr
		LEFT JOIN maintenance_orders mo ON mo.id = tr.maintenance_id
		LEFT JOIN drones d ON d.id = mo.drone_id
		LEFT JOIN users u ON u.id = tr.user_id
		WHERE tr.id = $1
	`

	var report models.TechnicalReportWithDetails
	err := r.db.Get(&report, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (r *TechnicalReportRepository) GetByMaintenanceID(maintenanceID string) ([]*models.TechnicalReportWithDetails, error) {
	query := `
		SELECT 
			tr.*,
			mo.id as "maintenance_order.id", mo.drone_id as "maintenance_order.drone_id",
			mo.user_id as "maintenance_order.user_id", mo.status as "maintenance_order.status",
			mo.description as "maintenance_order.description", mo.start_date as "maintenance_order.start_date",
			mo.end_date as "maintenance_order.end_date", mo.created_at as "maintenance_order.created_at",
			mo.updated_at as "maintenance_order.updated_at",
			d.id as "maintenance_order.drone.id", d.model as "maintenance_order.drone.model", 
			d.serial as "maintenance_order.drone.serial", d.status as "maintenance_order.drone.status",
			d.last_check as "maintenance_order.drone.last_check", d.next_check as "maintenance_order.drone.next_check",
			d.created_at as "maintenance_order.drone.created_at", d.updated_at as "maintenance_order.drone.updated_at",
			u.id as "user.id", u.name as "user.name", u.email as "user.email",
			u.role as "user.role", u.created_at as "user.created_at",
			u.updated_at as "user.updated_at"
		FROM technical_reports tr
		LEFT JOIN maintenance_orders mo ON mo.id = tr.maintenance_id
		LEFT JOIN drones d ON d.id = mo.drone_id
		LEFT JOIN users u ON u.id = tr.user_id
		WHERE tr.maintenance_id = $1
		ORDER BY tr.created_at DESC
	`

	var reports []*models.TechnicalReportWithDetails
	err := r.db.Select(&reports, query, maintenanceID)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *TechnicalReportRepository) GetByUserID(userID string) ([]*models.TechnicalReportWithDetails, error) {
	query := `
		SELECT 
			tr.*,
			mo.id as "maintenance_order.id", mo.drone_id as "maintenance_order.drone_id",
			mo.user_id as "maintenance_order.user_id", mo.status as "maintenance_order.status",
			mo.description as "maintenance_order.description", mo.start_date as "maintenance_order.start_date",
			mo.end_date as "maintenance_order.end_date", mo.created_at as "maintenance_order.created_at",
			mo.updated_at as "maintenance_order.updated_at",
			d.id as "maintenance_order.drone.id", d.model as "maintenance_order.drone.model", 
			d.serial as "maintenance_order.drone.serial", d.status as "maintenance_order.drone.status",
			d.last_check as "maintenance_order.drone.last_check", d.next_check as "maintenance_order.drone.next_check",
			d.created_at as "maintenance_order.drone.created_at", d.updated_at as "maintenance_order.drone.updated_at",
			u.id as "user.id", u.name as "user.name", u.email as "user.email",
			u.role as "user.role", u.created_at as "user.created_at",
			u.updated_at as "user.updated_at"
		FROM technical_reports tr
		LEFT JOIN maintenance_orders mo ON mo.id = tr.maintenance_id
		LEFT JOIN drones d ON d.id = mo.drone_id
		LEFT JOIN users u ON u.id = tr.user_id
		WHERE tr.user_id = $1
		ORDER BY tr.created_at DESC
	`

	var reports []*models.TechnicalReportWithDetails
	err := r.db.Select(&reports, query, userID)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *TechnicalReportRepository) Update(id string, input *models.UpdateTechnicalReportInput) error {
	query := `
		UPDATE technical_reports
		SET 
			title = COALESCE($1, title),
			description = COALESCE($2, description),
			findings = COALESCE($3, findings),
			recommendations = COALESCE($4, recommendations),
			updated_at = $5
		WHERE id = $6
	`

	_, err := r.db.Exec(query,
		input.Title,
		input.Description,
		input.Findings,
		input.Recommendations,
		time.Now(),
		id,
	)

	return err
}

func (r *TechnicalReportRepository) Delete(id string) error {
	query := `DELETE FROM technical_reports WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *TechnicalReportRepository) List(offset, limit int) ([]*models.TechnicalReportWithDetails, error) {
	query := `
		SELECT 
			tr.*,
			mo.id as "maintenance_order.id", mo.drone_id as "maintenance_order.drone_id",
			mo.user_id as "maintenance_order.user_id", mo.status as "maintenance_order.status",
			mo.description as "maintenance_order.description", mo.start_date as "maintenance_order.start_date",
			mo.end_date as "maintenance_order.end_date", mo.created_at as "maintenance_order.created_at",
			mo.updated_at as "maintenance_order.updated_at",
			d.id as "maintenance_order.drone.id", d.model as "maintenance_order.drone.model", 
			d.serial as "maintenance_order.drone.serial", d.status as "maintenance_order.drone.status",
			d.last_check as "maintenance_order.drone.last_check", d.next_check as "maintenance_order.drone.next_check",
			d.created_at as "maintenance_order.drone.created_at", d.updated_at as "maintenance_order.drone.updated_at",
			u.id as "user.id", u.name as "user.name", u.email as "user.email",
			u.role as "user.role", u.created_at as "user.created_at",
			u.updated_at as "user.updated_at"
		FROM technical_reports tr
		LEFT JOIN maintenance_orders mo ON mo.id = tr.maintenance_id
		LEFT JOIN drones d ON d.id = mo.drone_id
		LEFT JOIN users u ON u.id = tr.user_id
		ORDER BY tr.created_at DESC
		LIMIT $1 OFFSET $2
	`

	var reports []*models.TechnicalReportWithDetails
	err := r.db.Select(&reports, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return reports, nil
}
