package handlers

import (
    "encoding/json"
    "golang.org/x/crypto/bcrypt"
    "library-app/db"
    "library-app/models"
    "log"
    "net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        renderTemplate(w, "login", PageData{Title: "Login"})
        return
    }

    // Set content type for JSON responses
    w.Header().Set("Content-Type", "application/json")

    var loginRequest models.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
        log.Printf("Error decoding login request: %v", err)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "success": false,
            "message": "Invalid request format",
        })
        return
    }

    var admin models.Admin
    err := db.DB.Get(&admin, "SELECT * FROM admin WHERE username = ?", loginRequest.Username)
    if err != nil {
        log.Printf("User not found or database error: %v", err)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "success": false,
            "message": "Invalid credentials",
        })
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(loginRequest.Password)); err != nil {
        log.Printf("Password mismatch for user: %s", loginRequest.Username)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "success": false,
            "message": "Invalid credentials",
        })
        return
    }

    // Create session
    session, _ := store.Get(r, "library-session")
    session.Values["authenticated"] = true
    session.Values["admin_id"] = admin.ID
    session.Values["username"] = admin.Username
    session.Values["role"] = admin.Role

    if err := session.Save(r, w); err != nil {
        log.Printf("Error saving session: %v", err)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "success": false,
            "message": "Error creating session",
        })
        return
    }

    // Update last login
    _, err = db.DB.Exec("UPDATE admin SET last_login = CURRENT_TIMESTAMP WHERE admin_id = ?", admin.ID)
    if err != nil {
        log.Printf("Error updating last login: %v", err)
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "message": "Login successful",
        "data": map[string]string{
            "redirect": "/dashboard",
        },
    })
}
