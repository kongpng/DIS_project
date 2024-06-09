package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	_ "regexp"

	_ "github.com/mattn/go-sqlite3"
)

func executeSQLFile(db *sql.DB, filepath string) error {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(content))
	return err
}

func main() {

	db, err := sql.Open("sqlite3", "db/Data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Command line flags
	update := flag.Bool("update", false, "Apply update.sql")
	revert := flag.Bool("revert", false, "Apply revert.sql")
	initialData := flag.Bool("initial-data", false, "Apply initial_data.sql")
	flag.Parse()

	if *update {
		err = executeSQLFile(db, "update.sql")
		if err != nil {
			log.Fatalf("Failed to update database: %v", err)
		}
		log.Println("Database updated successfully.")
		return
	}

	if *revert {
		err = executeSQLFile(db, "revert.sql")
		if err != nil {
			log.Fatalf("Failed to revert database: %v", err)
		}
		log.Println("Database reverted successfully.")
		return
	}

	if *initialData {
		err = executeSQLFile(db, "initial_data.sql")
		if err != nil {
			log.Fatalf("Failed to insert initial data: %v", err)
		}
		log.Println("Initial data inserted successfully.")
		return
	}

	setupRoutes(db)

	log.Println("Server starting on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
