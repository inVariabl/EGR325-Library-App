package handlers

import (
	"encoding/json"
	"fmt"
	"library-app/db"
	"library-app/models"
	"log"
	"net/http"
	"time"
)

// ReturnsPageHandler serves the returns page
func ReturnsPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "returns", PageData{Title: "Book Return"})
}

// SearchMembers handles member search requests
func SearchMembers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		sendErrorResponse(w, "Search query is required", http.StatusBadRequest)
		return
	}

	var members []models.MemberResponse
	err := db.DB.Select(&members, `
        SELECT
            m.*,
            COUNT(c.checkout_id) as checkout_count,
            MAX(c.checkout_date) as last_checkout
        FROM member m
        LEFT JOIN checkout c ON m.member_id = c.member_id
        WHERE m.name LIKE ? OR m.email LIKE ? OR CAST(m.member_id AS CHAR) LIKE ?
        GROUP BY m.member_id
        LIMIT 10`,
		"%"+query+"%", "%"+query+"%", "%"+query+"%")

	if err != nil {
		log.Printf("Error searching members: %v", err)
		sendErrorResponse(w, "Error searching members", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, "", members)
}

// ProcessCheckout handles book checkout requests
func ProcessCheckout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BookID   int    `json:"book_id"`
		MemberID int    `json:"member_id"`
		DueDate  string `json:"due_date"`
		Notes    string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Start transaction
	tx, err := db.DB.Beginx()
	if err != nil {
		sendErrorResponse(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Check book availability
	var availableCopies int
	err = tx.Get(&availableCopies, `
        SELECT available_copies
        FROM book_inventory
        WHERE book_id = ?
        FOR UPDATE`,
		req.BookID)

	if err != nil {
		sendErrorResponse(w, "Book not found", http.StatusNotFound)
		return
	}

	if availableCopies <= 0 {
		sendErrorResponse(w, "Book is not available", http.StatusBadRequest)
		return
	}

	// Parse due date
	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		sendErrorResponse(w, "Invalid due date format", http.StatusBadRequest)
		return
	}

	// Create checkout record
	session, _ := store.Get(r, "library-session")
	adminID := session.Values["admin_id"].(int)

	result, err := tx.Exec(`
        INSERT INTO checkout (
            book_id,
            member_id,
            checkout_date,
            due_date,
            notes
        ) VALUES (?, ?, CURRENT_TIMESTAMP, ?, ?)`,
		req.BookID, req.MemberID, dueDate, req.Notes)

	if err != nil {
		log.Printf("Error creating checkout: %v", err)
		sendErrorResponse(w, "Error processing checkout", http.StatusInternalServerError)
		return
	}

	// Update book inventory
	_, err = tx.Exec(`
        UPDATE book_inventory
        SET available_copies = available_copies - 1
        WHERE book_id = ?`,
		req.BookID)

	if err != nil {
		log.Printf("Error updating inventory: %v", err)
		sendErrorResponse(w, "Error updating inventory", http.StatusInternalServerError)
		return
	}

	// Get the checkout ID
	checkoutID, _ := result.LastInsertId()

	// Log activity
	_, err = tx.Exec(`
        INSERT INTO activity_log (admin_id, action, details)
        VALUES (?, 'checkout', ?)`,
		adminID,
		json.RawMessage(fmt.Sprintf(`{"checkout_id": %d, "book_id": %d, "member_id": %d}`,
			checkoutID, req.BookID, req.MemberID)))

	if err != nil {
		log.Printf("Error logging activity: %v", err)
		sendErrorResponse(w, "Error logging activity", http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		sendErrorResponse(w, "Error completing checkout", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, "Checkout processed successfully", map[string]interface{}{
		"checkout_id": checkoutID,
	})
}

// ProcessReturn handles book return requests
func ProcessReturn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CheckoutID int    `json:"checkout_id"`
		Condition  string `json:"condition"`
		ReturnDate string `json:"return_date"`
		Notes      string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Start transaction
	tx, err := db.DB.Beginx()
	if err != nil {
		sendErrorResponse(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Get checkout info
	var checkout struct {
		BookID   int  `db:"book_id"`
		Returned bool `db:"returned"`
	}
	err = tx.Get(&checkout, `
        SELECT book_id, return_date IS NOT NULL as returned
        FROM checkout
        WHERE checkout_id = ?`,
		req.CheckoutID)

	if err != nil {
		log.Printf("Error getting checkout info: %v", err)
		sendErrorResponse(w, "Checkout not found", http.StatusNotFound)
		return
	}

	if checkout.Returned {
		sendErrorResponse(w, "This book has already been returned", http.StatusBadRequest)
		return
	}

	// Update checkout with return date
	_, err = tx.Exec(`
        UPDATE checkout
        SET return_date = ?, notes = CONCAT(COALESCE(notes, ''), '\nReturn Notes: ', ?)
        WHERE checkout_id = ?`,
		req.ReturnDate, req.Notes, req.CheckoutID)

	if err != nil {
		log.Printf("Error updating checkout: %v", err)
		sendErrorResponse(w, "Error updating checkout", http.StatusInternalServerError)
		return
	}

	// Create return record
	_, err = tx.Exec(`
        INSERT INTO book_return (
            checkout_id, return_date, book_condition, notes
        ) VALUES (?, ?, ?, ?)`,
		req.CheckoutID, req.ReturnDate, req.Condition, req.Notes)

	if err != nil {
		log.Printf("Error creating return record: %v", err)
		sendErrorResponse(w, "Error recording return", http.StatusInternalServerError)
		return
	}

	// Update book inventory if book isn't lost
	if req.Condition != "lost" {
		_, err = tx.Exec(`
            UPDATE book_inventory
            SET available_copies = available_copies + 1
            WHERE book_id = ?`,
			checkout.BookID)

		if err != nil {
			log.Printf("Error updating inventory: %v", err)
			sendErrorResponse(w, "Error updating inventory", http.StatusInternalServerError)
			return
		}
	}

	// Log the return
	session, _ := store.Get(r, "library-session")
	adminID := session.Values["admin_id"].(int)

	_, err = tx.Exec(`
        INSERT INTO activity_log (admin_id, action, details)
        VALUES (?, 'return', ?)`,
		adminID,
		json.RawMessage(fmt.Sprintf(`{"checkout_id": %d, "condition": "%s"}`,
			req.CheckoutID, req.Condition)))

	if err != nil {
		log.Printf("Error logging activity: %v", err)
		sendErrorResponse(w, "Error logging activity", http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		sendErrorResponse(w, "Error completing return", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, "Return processed successfully", nil)
}

// SearchCheckouts handles searching for active checkouts
func SearchCheckouts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		sendErrorResponse(w, "Search query is required", http.StatusBadRequest)
		return
	}

	log.Printf("Search query: %s", query)

	var checkouts []models.CheckoutResponse
	sqlQuery := `
        SELECT
            c.checkout_id,
            c.checkout_date,
            c.due_date,
            c.return_date,
            c.notes,
            b.book_id,
            b.isbn,
            b.title,
            b.author,
            b.shelf_location,
            m.member_id,
            m.name,
            m.email,
            m.phone_number
        FROM checkout c
        JOIN book b ON c.book_id = b.book_id
        JOIN member m ON c.member_id = m.member_id
        WHERE
            c.return_date IS NULL AND
            (LOWER(b.title) LIKE LOWER(?) OR
             LOWER(b.isbn) LIKE LOWER(?) OR
             LOWER(m.name) LIKE LOWER(?) OR
             CAST(c.checkout_id AS CHAR) LIKE ?)`

	searchPattern := "%" + query + "%"
	log.Printf("Executing query with pattern: %s", searchPattern)

	err := db.DB.Select(&checkouts, sqlQuery,
		searchPattern, searchPattern, searchPattern, searchPattern)

	if err != nil {
		log.Printf("Database error: %v", err)
		sendErrorResponse(w, "Error searching checkouts", http.StatusInternalServerError)
		return
	}

	log.Printf("Found %d checkouts", len(checkouts))
	if len(checkouts) > 0 {
		log.Printf("First checkout details: %+v", checkouts[0])
	}

	if len(checkouts) == 0 {
		// Return empty array instead of error
		log.Printf("No checkouts found for query: %s", query)
		sendSuccessResponse(w, "", []models.CheckoutResponse{})
		return
	}

	sendSuccessResponse(w, "", checkouts)
}
