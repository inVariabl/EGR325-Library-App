package handlers

import (
    "net/http"
    "library-app/db"
    "log"
)

type DashboardStats struct {
    TotalBooks      int `json:"totalBooks"`
    AvailableBooks  int `json:"availableBooks"`
    TotalMembers    int `json:"totalMembers"`
    ActiveCheckouts int `json:"activeCheckouts"`
}

type DashboardData struct {
    Title    string
    Username string
}

// DashboardHandler serves the dashboard page
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
    // Get session
    session, err := store.Get(r, "library-session")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Get username from session
    username, ok := session.Values["username"].(string)
    if !ok {
        username = "Unknown"
    }

    data := DashboardData{
        Title:    "Dashboard",
        Username: username,
    }

    renderTemplate(w, "dashboard", PageData{
        Title: "Dashboard",
        Data:  data,
    })
}

// DashboardStatsHandler handles the AJAX request for dashboard statistics
func DashboardStatsHandler(w http.ResponseWriter, r *http.Request) {
    // Set content type before anything else
    w.Header().Set("Content-Type", "application/json")

    stats := DashboardStats{}

    // Get total books
    err := db.DB.Get(&stats.TotalBooks, `
        SELECT COUNT(*) FROM book`)
    if err != nil {
        log.Printf("Error fetching total books: %v", err)
        sendErrorResponse(w, "Error fetching total books", http.StatusInternalServerError)
        return
    }
    log.Printf("Total books query result: %d", stats.TotalBooks)

    // Get available books
    err = db.DB.Get(&stats.AvailableBooks, `
        SELECT COUNT(*) FROM book WHERE status = 'available'`)
    if err != nil {
        log.Printf("Error fetching available books: %v", err)
        sendErrorResponse(w, "Error fetching available books", http.StatusInternalServerError)
        return
    }
    log.Printf("Available books query result: %d", stats.AvailableBooks)

    // Get total members
    err = db.DB.Get(&stats.TotalMembers, `
        SELECT COUNT(*) FROM member`)
    if err != nil {
        log.Printf("Error fetching total members: %v", err)
        sendErrorResponse(w, "Error fetching total members", http.StatusInternalServerError)
        return
    }
    log.Printf("Total members query result: %d", stats.TotalMembers)

    // Get active checkouts
    err = db.DB.Get(&stats.ActiveCheckouts, `
        SELECT COUNT(*) FROM checkout WHERE return_date IS NULL`)
    if err != nil {
        log.Printf("Error fetching active checkouts: %v", err)
        sendErrorResponse(w, "Error fetching active checkouts", http.StatusInternalServerError)
        return
    }
    log.Printf("Active checkouts query result: %d", stats.ActiveCheckouts)

    // Log the final stats before sending
    log.Printf("Sending stats to client: %+v", stats)

    sendSuccessResponse(w, "", stats)
}
