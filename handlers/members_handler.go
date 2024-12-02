package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"library-app/db"
	"library-app/models"
)

func MembersPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "members", PageData{Title: "Library Members"})
}

func GetMembers(w http.ResponseWriter, r *http.Request) {
	var members []models.Member
	err := db.DB.Select(&members, `
        SELECT * FROM member
        ORDER BY member_id DESC
        LIMIT 50`)

	if err != nil {
		log.Printf("Error getting members: %v", err)
		sendErrorResponse(w, "Error retrieving members", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, "", members)
}

func GetMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	memberID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, "Invalid member ID", http.StatusBadRequest)
		return
	}

	var member models.Member
	err = db.DB.Get(&member, "SELECT * FROM member WHERE member_id = ?", memberID)
	if err != nil {
		log.Printf("Error fetching member: %v", err)
		sendErrorResponse(w, "Member not found", http.StatusNotFound)
		return
	}

	sendSuccessResponse(w, "", member)
}

func UpdateMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	memberID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, "Invalid member ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if email is already in use by another member
	var count int
	err = db.DB.Get(&count,
		"SELECT COUNT(*) FROM member WHERE email = ? AND member_id != ?",
		req.Email, memberID)
	if err != nil {
		log.Printf("Error checking email uniqueness: %v", err)
		sendErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		sendErrorResponse(w, "Email is already in use", http.StatusBadRequest)
		return
	}

	// Update member
	result, err := db.DB.Exec(`
        UPDATE member
        SET name = ?, email = ?, phone_number = ?, address = ?
        WHERE member_id = ?`,
		req.Name, req.Email, req.PhoneNumber, req.Address, memberID)
	if err != nil {
		log.Printf("Error updating member: %v", err)
		sendErrorResponse(w, "Error updating member", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		sendErrorResponse(w, "Error updating member", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		sendErrorResponse(w, "Member not found", http.StatusNotFound)
		return
	}

	// Log the activity
	session, _ := store.Get(r, "library-session")
	adminID := session.Values["admin_id"].(int)

	_, err = db.DB.Exec(`
        INSERT INTO activity_log (admin_id, action, details)
        VALUES (?, 'update_member', ?)`,
		adminID,
		fmt.Sprintf(`{"member_id": %d}`, memberID))
	if err != nil {
		log.Printf("Error logging activity: %v", err)
	}

	sendSuccessResponse(w, "Member updated successfully", nil)
}

func AddMember(w http.ResponseWriter, r *http.Request) {
	var req models.CreateMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec(`
        INSERT INTO member (name, email, phone_number, address, membership_date)
        VALUES (?, ?, ?, ?, CURRENT_DATE)`,
		req.Name, req.Email, req.PhoneNumber, req.Address)

	if err != nil {
		log.Printf("Error adding member: %v", err)
		sendErrorResponse(w, "Error adding member", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	sendSuccessResponse(w, "Member added successfully", map[string]interface{}{
		"member_id": id,
	})
}
