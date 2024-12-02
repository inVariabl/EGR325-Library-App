package handlers

import (
	"net/http"
)

// HomeHandler redirects to login if not authenticated, dashboard if authenticated
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "library-session")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	// If not authenticated, redirect to login
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
