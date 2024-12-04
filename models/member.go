// models/member.go
package models

import "time"

type Member struct {
	MemberID       int       `json:"member_id" db:"member_id"`
	Name           string    `json:"name" db:"name"`
	Email          string    `json:"email" db:"email"`
	PhoneNumber    string    `json:"phone_number" db:"phone_number"`
	MembershipDate time.Time `json:"membership_date" db:"membership_date"`
	Address        string    `json:"address" db:"address"`
}

type MemberResponse struct {
	Member        Member    `json:"member"`
	CheckoutCount int       `json:"checkout_count"`
	LastCheckout  time.Time `json:"last_checkout,omitempty"`
}

type CreateMemberRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

type UpdateMemberRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}
