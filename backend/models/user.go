package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	UserRoleAdmin      = "admin"
	UserRoleUser       = "user"
	UserRoleTechnician = "technician"
)

type User struct {
	ID           string    `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Name         string    `db:"name" json:"name"`
	Phone        string    `db:"phone" json:"phone"`
	Role         string    `db:"role" json:"role"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type CreateUserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=admin user technician"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

func (u *User) HashPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}

func NewUser(input CreateUserInput) (*User, error) {
	user := &User{
		Email: input.Email,
		Name:  input.Name,
		Phone: input.Phone,
		Role:  input.Role,
	}

	if err := user.HashPassword(input.Password); err != nil {
		return nil, err
	}

	return user, nil
}
