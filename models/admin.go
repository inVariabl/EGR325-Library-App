package models

import (
	"database/sql"
	"time"
)

type Admin struct {
	ID           int            `db:"admin_id" json:"id"`
	Username     string         `db:"username" json:"username"`
	PasswordHash string         `db:"password_hash" json:"-"`
	Email        string         `db:"email" json:"email"`
	FirstName    string         `db:"first_name" json:"first_name"`
	LastName     string         `db:"last_name" json:"last_name"`
	PhoneNumber  sql.NullString `db:"phone_number" json:"phone_number,omitempty"` // Made nullable
	Role         string         `db:"role" json:"role"`
	CreatedAt    time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at" json:"updated_at"`
	LastLogin    sql.NullTime   `db:"last_login" json:"last_login,omitempty"` // Made nullable
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateAdminRequest struct {
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required,min=8"`
	Email       string `json:"email" validate:"required,email"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role" validate:"required,oneof=superadmin manager staff"`
}
