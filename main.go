package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "db/Data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	setupRoutes(db)

	log.Println("Server starting on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
