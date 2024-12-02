package handlers

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"library-app/db"
	"library-app/models"
	"log"
	"net/http"
)

var validate = validator.New()

func AddBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var bookRequest models.AddBookRequest
	if err := json.NewDecoder(r.Body).Decode(&bookRequest); err != nil {
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := validate.Struct(bookRequest); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Start transaction
	tx, err := db.DB.Beginx()
	if err != nil {
		sendErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Insert book
	result, err := tx.Exec(`
		INSERT INTO book (isbn, title, author, publisher, publication_year,
		category, language, pages, shelf_location, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		bookRequest.ISBN, bookRequest.Title, bookRequest.Author,
		bookRequest.Publisher, bookRequest.PublicationYear,
		bookRequest.Category, bookRequest.Language, bookRequest.Pages,
		bookRequest.ShelfLocation, "available")

	if err != nil {
		log.Printf("Error inserting book: %v", err)
		sendErrorResponse(w, "Error adding book", http.StatusInternalServerError)
		return
	}

	bookID, err := result.LastInsertId()
	if err != nil {
		sendErrorResponse(w, "Error getting book ID", http.StatusInternalServerError)
		return
	}

	// Insert inventory
	_, err = tx.Exec(`
		INSERT INTO book_inventory (book_id, available_copies, total_copies)
		VALUES (?, ?, ?)`,
		bookID, bookRequest.AvailableCopies, bookRequest.TotalCopies)

	if err != nil {
		log.Printf("Error inserting inventory: %v", err)
		sendErrorResponse(w, "Error adding book inventory", http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		sendErrorResponse(w, "Error committing transaction", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, "Book added successfully", nil)
}

// GetBook retrieves a single book by ID
func GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]

	var book models.Book
	err := db.DB.Get(&book, `
        SELECT b.*, bi.available_copies, bi.total_copies
        FROM book b
        LEFT JOIN book_inventory bi ON b.book_id = bi.book_id
        WHERE b.book_id = ?`, bookID)

	if err != nil {
		log.Printf("Error getting book: %v", err)
		sendErrorResponse(w, "Book not found", http.StatusNotFound)
		return
	}

	sendSuccessResponse(w, "", book)
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var books []models.Book
	err := db.DB.Select(&books, `
		SELECT b.*, bi.available_copies, bi.total_copies
		FROM book b
		LEFT JOIN book_inventory bi ON b.book_id = bi.book_id
		ORDER BY b.created_at DESC`)

	if err != nil {
		log.Printf("Error getting books: %v", err)
		sendErrorResponse(w, "Error retrieving books", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, "", books)
}

func AddBookPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/add-book.html")
}
