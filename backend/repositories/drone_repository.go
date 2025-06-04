package repositories

import (
	"database/sql"
	"time"

	"github.com/root/evodrones/backend/db"
	"github.com/root/evodrones/backend/models"
)

type DroneRepository struct {
	db *db.DB
}

func NewDroneRepository(db *db.DB) *DroneRepository {
	return &DroneRepository{db: db}
}

func (r *DroneRepository) Create(drone *models.Drone) error {
	query := `
		INSERT INTO drones (user_id, model, serial, status, last_check, next_check, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	return r.db.DB.QueryRow(query,
		drone.UserID,
		drone.Model,
		drone.Serial,
		drone.Status,
		drone.LastCheck,
		drone.NextCheck,
		drone.CreatedAt,
		drone.UpdatedAt,
	).Scan(&drone.ID)
}

func (r *DroneRepository) GetByID(id string) (*models.Drone, error) {
	query := `
		SELECT id, user_id, model, serial, status, last_check, next_check, created_at, updated_at
		FROM drones
		WHERE id = $1
	`

	drone := &models.Drone{}
	err := r.db.DB.Get(drone, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return drone, nil
}

func (r *DroneRepository) GetByUserID(userID string) ([]*models.Drone, error) {
	query := `
		SELECT id, user_id, model, serial, status, last_check, next_check, created_at, updated_at
		FROM drones
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	drones := []*models.Drone{}
	err := r.db.DB.Select(&drones, query, userID)
	if err != nil {
		return nil, err
	}

	return drones, nil
}

func (r *DroneRepository) Update(drone *models.Drone) error {
	query := `
		UPDATE drones
		SET model = $1, serial = $2, status = $3, last_check = $4, next_check = $5, updated_at = $6
		WHERE id = $7
	`

	drone.UpdatedAt = time.Now()

	_, err := r.db.DB.Exec(query,
		drone.Model,
		drone.Serial,
		drone.Status,
		drone.LastCheck,
		drone.NextCheck,
		drone.UpdatedAt,
		drone.ID,
	)

	return err
}

func (r *DroneRepository) Delete(id string) error {
	query := `DELETE FROM drones WHERE id = $1`
	_, err := r.db.DB.Exec(query, id)
	return err
}

func (r *DroneRepository) List(offset, limit int) ([]*models.Drone, error) {
	query := `
		SELECT id, user_id, model, serial, status, last_check, next_check, created_at, updated_at
		FROM drones
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	drones := []*models.Drone{}
	err := r.db.DB.Select(&drones, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return drones, nil
}

func (r *DroneRepository) GetByStatus(status string) ([]*models.Drone, error) {
	query := `
		SELECT id, user_id, model, serial, status, last_check, next_check, created_at, updated_at
		FROM drones
		WHERE status = $1
		ORDER BY created_at DESC
	`

	drones := []*models.Drone{}
	err := r.db.DB.Select(&drones, query, status)
	if err != nil {
		return nil, err
	}

	return drones, nil
}

func (r *DroneRepository) GetDueForCheck() ([]*models.Drone, error) {
	query := `
		SELECT id, user_id, model, serial, status, last_check, next_check, created_at, updated_at
		FROM drones
		WHERE next_check <= NOW()
		ORDER BY next_check ASC
	`

	drones := []*models.Drone{}
	err := r.db.DB.Select(&drones, query)
	if err != nil {
		return nil, err
	}

	return drones, nil
}
