package models

import "time"

type CheckoutResponse struct {
	CheckoutID   int        `db:"checkout_id" json:"checkout_id"`
	BookID       int        `db:"book_id" json:"book_id"`
	MemberID     int        `db:"member_id" json:"member_id"`
	CheckoutDate time.Time  `db:"checkout_date" json:"checkout_date"`
	DueDate      time.Time  `db:"due_date" json:"due_date"`
	ReturnDate   *time.Time `db:"return_date" json:"return_date,omitempty"`
	Notes        string     `db:"notes" json:"notes"`

	// Book details
	ISBN          string `db:"isbn" json:"isbn"`
	Title         string `db:"title" json:"title"`
	Author        string `db:"author" json:"author"`
	ShelfLocation string `db:"shelf_location" json:"shelf_location"`

	// Member details
	Name        string `db:"name" json:"member_name"` // note the json tag matches what frontend expects
	Email       string `db:"email" json:"email"`
	PhoneNumber string `db:"phone_number" json:"phone_number"`
}
