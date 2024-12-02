package handlers

import (
	"github.com/gorilla/sessions"
)

var (
	// Store is the session store for the application
	// Change the secret key in production!
	store = sessions.NewCookieStore([]byte("your-secret-key-123"))
)

func init() {
	// Configure session store options
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
	}
}
