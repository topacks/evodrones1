package repositories

import (
	"database/sql"
	"time"

	"github.com/root/evodrones/backend/db"
	"github.com/root/evodrones/backend/models"
)

type UserRepository struct {
	db *db.DB
}

func NewUserRepository(db *db.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, name, phone, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	return r.db.QueryRow(query,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.Phone,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, name, phone, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.Get(user, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, name, phone, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.Get(user, query, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET email = $1, password_hash = $2, name = $3, phone = $4, role = $5, updated_at = $6
		WHERE id = $7
	`

	user.UpdatedAt = time.Now()

	_, err := r.db.Exec(query,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.Phone,
		user.Role,
		user.UpdatedAt,
		user.ID,
	)

	return err
}

func (r *UserRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *UserRepository) List(offset, limit int) ([]*models.User, error) {
	query := `
		SELECT id, email, password_hash, name, phone, role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	users := []*models.User{}
	err := r.db.Select(&users, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return users, nil
}
