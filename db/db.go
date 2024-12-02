package db

import (
    "log"
    "github.com/jmoiron/sqlx"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sqlx.DB

func InitDB() {
    var err error
    DB, err = sqlx.Connect("mysql", "libraryuser:password@(localhost:3306)/library_db?parseTime=true")
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatalf("Error pinging database: %v", err)
    }

    log.Println("Connected to database successfully")
}
