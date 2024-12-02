package handlers

import (
	"database/sql"
	"library-app/db"
	"log"
	"net/http"
)

// SearchBooks handles book search requests
func SearchBooks(w http.ResponseWriter, r *http.Request) {
	log.Printf("SearchBooks handler called with method: %s", r.Method)
	log.Printf("Full URL path: %s", r.URL.Path)

	query := r.URL.Query().Get("q")
	log.Printf("Search query: %s", query)

	if query == "" {
		log.Printf("Empty search query")
		sendErrorResponse(w, "Search query is required", http.StatusBadRequest)
		return
	}

	log.Printf("Executing book search query")
	var books []struct {
		ID              int    `db:"book_id" json:"id"`
		ISBN            string `db:"isbn" json:"isbn"`
		Title           string `db:"title" json:"title"`
		Author          string `db:"author" json:"author"`
		Publisher       string `db:"publisher" json:"publisher"`
		Category        string `db:"category" json:"category"`
		Language        string `db:"language" json:"language"`
		Pages           int    `db:"pages" json:"pages"`
		ShelfLocation   string `db:"shelf_location" json:"shelf_location"`
		Status          string `db:"status" json:"status"`
		AvailableCopies int    `db:"available_copies" json:"available_copies"`
		TotalCopies     int    `db:"total_copies" json:"total_copies"`
	}

	// Debug: Print the SQL query that will be executed
	searchQuery := `
        SELECT
            b.book_id, b.isbn, b.title, b.author, b.publisher,
            b.category, b.language, b.pages, b.shelf_location, b.status,
            COALESCE(bi.available_copies, 0) as available_copies,
            COALESCE(bi.total_copies, 0) as total_copies
        FROM book b
        LEFT JOIN book_inventory bi ON b.book_id = bi.book_id
        WHERE
            LOWER(b.title) LIKE LOWER(?) OR
            LOWER(b.isbn) LIKE LOWER(?) OR
            LOWER(b.author) LIKE LOWER(?)`

	log.Printf("Executing SQL query: %s", searchQuery)
	log.Printf("With parameters: %q, %q, %q", "%"+query+"%", "%"+query+"%", "%"+query+"%")

	// Try to ping the database first
	if err := db.DB.Ping(); err != nil {
		log.Printf("Database connection error: %v", err)
		sendErrorResponse(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	err := db.DB.Select(&books, searchQuery,
		"%"+query+"%", "%"+query+"%", "%"+query+"%")

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No books found for query: %s", query)
			sendSuccessResponse(w, "", []interface{}{})
			return
		}
		log.Printf("Error searching books: %v", err)
		sendErrorResponse(w, "Error searching books", http.StatusInternalServerError)
		return
	}

	log.Printf("Found %d books for query: %s", len(books), query)
	for i, book := range books {
		log.Printf("Book %d: %s by %s", i+1, book.Title, book.Author)
	}

	sendSuccessResponse(w, "", books)
}
