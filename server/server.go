package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"library-app/db"
	"library-app/handlers"
	"log"
	"mime"
	"net/http"
)

var (
	store = sessions.NewCookieStore([]byte("your-secret-key"))
)

func init() {
	// Initialize database
	db.InitDB()

	// Add MIME types
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
}

func main() {
	r := mux.NewRouter()

	// Static files with proper MIME types
	fileServer := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	// Serve static files from root path for convenience
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))

	// Auth middleware
	auth := handlers.AuthMiddleware

	// Pages
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("GET", "POST")
	r.Handle("/dashboard", auth(http.HandlerFunc(handlers.DashboardHandler))).Methods("GET")
	r.Handle("/add-book", auth(http.HandlerFunc(handlers.AddBookPageHandler))).Methods("GET")
	r.Handle("/checkout", auth(http.HandlerFunc(handlers.CheckoutPageHandler))).Methods("GET")

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.Handle("/books", auth(http.HandlerFunc(handlers.AddBook))).Methods("POST")
	api.Handle("/books/{id}", auth(http.HandlerFunc(handlers.GetBook))).Methods("GET")
	api.Handle("/checkout", auth(http.HandlerFunc(handlers.ProcessCheckout))).Methods("POST")
	api.Handle("/members/search", auth(http.HandlerFunc(handlers.SearchMembers))).Methods("GET")
	api.Handle("/books/search", auth(http.HandlerFunc(handlers.SearchBooks))).Methods("GET")

	// Start server
	log.Printf("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
