package main

import (
	"log"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	// Connect to database
	db, err := sqlx.Connect("mysql", "libraryuser:password@tcp(localhost:3306)/library_db?parseTime=true")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Create password hash
	password := "admin123" // Change this password!
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}

	// Insert admin user
	result, err := db.Exec(`
		INSERT INTO admin (
			username,
			password_hash,
			email,
			first_name,
			last_name,
			role,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`,
		"admin",
		string(hashedPassword),
		"admin@library.com",
		"Admin",
		"User",
		"superadmin",
	)

	if err != nil {
		log.Fatalf("Error creating admin user: %v", err)
	}

	id, _ := result.LastInsertId()
	log.Printf("Admin user created successfully with ID: %d", id)
	log.Printf("Username: admin")
	log.Printf("Password: %s", password)
}
