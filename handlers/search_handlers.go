package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
)

func SearchCustomersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Missing search parameter: name", http.StatusBadRequest)
			return
		}

		//
		pattern := "%" + name + "%"

		query := "SELECT * FROM Customer WHERE Name LIKE ?"
		rows, err := db.Query(query, pattern)
		if err != nil {
			http.Error(w, "Failed to execute search query: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			http.Error(w, "Failed to get columns: "+err.Error(), http.StatusInternalServerError)
			return
		}

		results := make([]map[string]interface{}, 0)
		for rows.Next() {
			columnsData := make([]interface{}, len(columns))
			columnsPointers := make([]interface{}, len(columns))
			for i := range columnsData {
				columnsPointers[i] = &columnsData[i]
			}
			if err := rows.Scan(columnsPointers...); err != nil {
				http.Error(w, "Failed to scan row: "+err.Error(), http.StatusInternalServerError)
				return
			}
			rowMap := make(map[string]interface{})
			for i, col := range columns {
				rowMap[col] = columnsData[i]
			}
			results = append(results, rowMap)
		}

		fmt.Fprintf(w, `<html><head><title>Search Results</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
		fmt.Fprintf(w, `<h1>Search Results</h1><ul>`)
		for _, result := range results {
			fmt.Fprintf(w, `<li>`)
			for col, val := range result {
				fmt.Fprintf(w, `%s: %v, `, col, val)
			}
			fmt.Fprintf(w, `</li>`)
		}
		fmt.Fprintf(w, `</ul></body></html>`)
	}
}
