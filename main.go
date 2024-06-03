package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func connect_to_db() error {
	var count = 0

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	for {
		var err error
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Println("Failed to open database connection:", err)
			return err
		}

		err = db.Ping()
		if err == nil {
			log.Println("Successfully connected to the database.")
			return nil
		}

		log.Println("Failed to ping the database. Retrying...", err)
		db.Close()

		count++
		if count >= 3 {
			log.Println("Exiting, too many attempts.")
			return err
		}

		time.Sleep(2 * time.Second)
	}
}

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {

	err := connect_to_db()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/", greet)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
