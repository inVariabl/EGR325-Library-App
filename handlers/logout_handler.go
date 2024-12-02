package handlers

import (
	"library-app/db"
	"log"
	"net/http"
	"strings"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get the session
	session, err := store.Get(r, "library-session")
	if err != nil {
		log.Printf("Error getting session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Log the logout activity before clearing the session
	if adminID, ok := session.Values["admin_id"].(int); ok {
		_, err = db.DB.Exec(`
            INSERT INTO activity_log (admin_id, action, details)
            VALUES (?, 'logout', 'User logged out')
        `, adminID)
		if err != nil {
			log.Printf("Error logging activity: %v", err)
		}
	}

	// Clear session values
	session.Values["authenticated"] = false
	session.Values["admin_id"] = nil
	session.Values["username"] = nil
	session.Values["role"] = nil

	// Set the session to expire
	session.Options.MaxAge = -1

	// Save the session
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Error saving session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check if request accepts JSON
	acceptHeader := r.Header.Get("Accept")
	if strings.Contains(acceptHeader, "application/json") {
		// Send JSON response for API requests
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "message": "Logged out successfully", "data": {"redirect": "/login"}}`))
	} else {
		// Redirect for browser requests
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
