// main.go
package main

import (
	"github.com/gorilla/mux"
	"library-app/db"
	"library-app/handlers"
	"log"
	"mime"
	"net/http"
)

func init() {
	db.InitDB()
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
}

func main() {
	r := mux.NewRouter()

	// Static files
	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))

	// Auth middleware
	auth := handlers.AuthMiddleware

	// Auth routes
	r.HandleFunc("/login", handlers.LoginHandler).Methods("GET", "POST")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")

	// Pages
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.Handle("/dashboard", auth(http.HandlerFunc(handlers.DashboardHandler))).Methods("GET")
	r.Handle("/add-book", auth(http.HandlerFunc(handlers.AddBookPageHandler))).Methods("GET")
	r.Handle("/checkout", auth(http.HandlerFunc(handlers.CheckoutPageHandler))).Methods("GET")
	r.Handle("/returns", auth(http.HandlerFunc(handlers.ReturnsPageHandler))).Methods("GET")
	r.Handle("/members", auth(http.HandlerFunc(handlers.MembersPageHandler))).Methods("GET")

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Apply auth middleware to all API routes
	api.Use(auth)

	// Dashboard endpoints
	api.HandleFunc("/dashboard/stats", handlers.DashboardStatsHandler).Methods("GET")
	api.HandleFunc("/dashboard/activity", handlers.ActivityLogHandler).Methods("GET")

	// Book endpoints - specific routes must come before parameterized routes
	api.HandleFunc("/books/search", handlers.SearchBooks).Methods("GET")
	api.HandleFunc("/books", handlers.AddBook).Methods("POST")
	api.HandleFunc("/books/{id}", handlers.GetBook).Methods("GET")

	// Member endpoints
	api.HandleFunc("/members/search", handlers.SearchMembers).Methods("GET")
	api.HandleFunc("/members", handlers.GetMembers).Methods("GET")
	api.HandleFunc("/members", handlers.AddMember).Methods("POST")
	api.HandleFunc("/members/{id}", handlers.GetMember).Methods("GET")
	api.HandleFunc("/members/{id}", handlers.UpdateMember).Methods("PUT")

	// Checkout and return endpoints
	api.HandleFunc("/checkouts/search", handlers.SearchCheckouts).Methods("GET")
	api.HandleFunc("/checkout", handlers.ProcessCheckout).Methods("POST")
	api.HandleFunc("/returns", handlers.ProcessReturn).Methods("POST")

	// Debug: Print all registered routes
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			methods, _ := route.GetMethods()
			log.Printf("Route: %s Methods: %v", pathTemplate, methods)
		}
		return nil
	})
	if err != nil {
		log.Printf("Error walking routes: %v", err)
	}

	log.Printf("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
