package models

import (
	"time"
)

type Book struct {
	ID              int       `json:"id" db:"book_id"`
	ISBN            string    `json:"isbn" db:"isbn"`
	Title           string    `json:"title" db:"title"`
	Author          string    `json:"author" db:"author"`
	Publisher       string    `json:"publisher" db:"publisher"`
	PublicationYear int       `json:"publication_year" db:"publication_year"`
	Category        string    `json:"category" db:"category"`
	Language        string    `json:"language" db:"language"`
	Pages           int       `json:"pages" db:"pages"`
	ShelfLocation   string    `json:"shelf_location" db:"shelf_location"`
	Status          string    `json:"status" db:"status"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	AvailableCopies int       `json:"available_copies" db:"available_copies"`
	TotalCopies     int       `json:"total_copies" db:"total_copies"`
}

type BookInventory struct {
	ID              int `json:"id" db:"inventory_id"`
	BookID          int `json:"book_id" db:"book_id"`
	AvailableCopies int `json:"available_copies" db:"available_copies"`
	TotalCopies     int `json:"total_copies" db:"total_copies"`
}

type AddBookRequest struct {
	ISBN            string `json:"isbn" validate:"required,min=10,max=13"`
	Title           string `json:"title" validate:"required"`
	Author          string `json:"author"`
	Publisher       string `json:"publisher"`
	PublicationYear int    `json:"publication_year" validate:"min=1000,max=2024"`
	Category        string `json:"category"`
	Language        string `json:"language"`
	Pages           int    `json:"pages" validate:"min=1"`
	ShelfLocation   string `json:"shelf_location"`
	AvailableCopies int    `json:"available_copies" validate:"min=0"`
	TotalCopies     int    `json:"total_copies" validate:"min=1"`
}
