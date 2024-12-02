package handlers

import (
	"net/http"
)

// CheckoutPageHandler serves the checkout page
func CheckoutPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	renderTemplate(w, "checkout", PageData{Title: "Book Checkout"})
}
