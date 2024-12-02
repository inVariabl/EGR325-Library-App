package handlers

import (
	"library-app/db"
	"log"
	"net/http"
)

type ActivityLog struct {
	CreatedAt string `json:"created_at" db:"created_at"`
	Action    string `json:"action" db:"action"`
	Details   string `json:"details" db:"details"`
	Username  string `json:"username" db:"username"`
}

// ActivityLogHandler handles the AJAX request for recent activities
func ActivityLogHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("ActivityLogHandler called")

	// Set content type
	w.Header().Set("Content-Type", "application/json")

	var activities []ActivityLog

	// First, check if the activity_log table exists and has data
	var count int
	err := db.DB.Get(&count, "SELECT COUNT(*) FROM activity_log")
	if err != nil {
		log.Printf("Error checking activity_log table: %v", err)
		sendErrorResponse(w, "Error accessing activity log", http.StatusInternalServerError)
		return
	}
	log.Printf("Found %d activities in the database", count)

	// Get activities with user information
	err = db.DB.Select(&activities, `
        SELECT
            DATE_FORMAT(al.created_at, '%Y-%m-%d %H:%i:%s') as created_at,
            al.action,
            al.details,
            a.username
        FROM activity_log al
        JOIN admin a ON al.admin_id = a.admin_id
        ORDER BY al.created_at DESC
        LIMIT 10
    `)

	if err != nil {
		log.Printf("Error fetching activities: %v", err)
		sendErrorResponse(w, "Error fetching activity log", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully fetched %d activities", len(activities))
	for i, activity := range activities {
		log.Printf("Activity %d: Time=%s, Action=%s, User=%s",
			i+1, activity.CreatedAt, activity.Action, activity.Username)
	}

	sendSuccessResponse(w, "", map[string]interface{}{
		"activities": activities,
	})
}
