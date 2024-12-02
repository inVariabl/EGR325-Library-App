package models

import "time"

type Return struct {
	ReturnID      int       `json:"return_id" db:"return_id"`
	CheckoutID    int       `json:"checkout_id" db:"checkout_id"`
	ReturnDate    time.Time `json:"return_date" db:"return_date"`
	BookCondition string    `json:"book_condition" db:"book_condition"`
	Notes         string    `json:"notes" db:"notes"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type ReturnRequest struct {
	CheckoutID int    `json:"checkout_id" validate:"required"`
	Condition  string `json:"condition" validate:"required,oneof=good damaged lost"`
	Notes      string `json:"notes"`
}
