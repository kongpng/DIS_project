package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

type Restaurant struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Open    bool   `json:"open"`
	Cuisine string `json:"cuisine"`
}

func RestaurantsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			HandleGetRestaurants(db, w, r)
		case "POST":
			HandlePostRestaurant(db, w, r)
		}
	}
}

func AddRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
        <html>
        <head><title>Add Restaurant</title></head>
        <body>
        <form action="/restaurants" method="post">
            <label>Name: <input type="text" name="name" /></label><br/>
            <label>Address: <input type="text" name="address" /></label><br/>
            <label>Open: <input type="checkbox" name="open" /></label><br/>
            <label>Cuisine: <input type="text" name="cuisine" /></label><br/>
            <input type="submit" value="Add Restaurant" />
        </form>
        </body>
        </html>
        `)
	}
}

func DeleteRestaurantHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			HandleDeleteRestaurant(db, w, r)
		}
	}
}

// Restaurant Handlers
func HandleGetRestaurants(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	restaurants, err := FetchRestaurants(db)
	if err != nil {
		http.Error(w, "Failed to fetch restaurants: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<html><head><title>Restaurants List</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<nav><a href="/">Home</a> | <a href="/addRestaurant">Add Restaurant</a></nav>`)
	fmt.Fprintf(w, `<h1>Restaurants</h1><ul>`)
	for _, r := range restaurants {
		fmt.Fprintf(w, `<li>ID: %d, Name: %s, Address: %s, Open: %t, Cuisine: %s</li>`, r.ID, r.Name, r.Address, r.Open, r.Cuisine)
	}
	fmt.Fprintf(w, `</ul></body></html>`)
}

func HandlePostRestaurant(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var res Restaurant
	res.Name = r.FormValue("name")
	res.Address = r.FormValue("address")
	res.Open = r.FormValue("open") == "on"
	res.Cuisine = r.FormValue("cuisine")
	_, err := AddRestaurant(db, res.Name, res.Address, res.Open, res.Cuisine)
	if err != nil {
		http.Error(w, "Failed to add restaurant: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<html><head><title>Add Restaurant</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<h1>New restaurant added</h1></body></html>`)
}

func HandleDeleteRestaurant(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}
	if err := DeleteRestaurant(db, id); err != nil {
		http.Error(w, "Failed to delete restaurant: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<html><head><title>Delete Restaurant</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<h1>Restaurant deleted successfully</h1></body></html>`)
}

func FetchRestaurants(db *sql.DB) ([]Restaurant, error) {
	rows, err := db.Query("SELECT ID, Name, Address, Open, Cuisine FROM Restaurant")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var restaurants []Restaurant
	for rows.Next() {
		var r Restaurant
		if err := rows.Scan(&r.ID, &r.Name, &r.Address, &r.Open, &r.Cuisine); err != nil {
			return nil, err
		}
		restaurants = append(restaurants, r)
	}
	return restaurants, nil
}

func AddRestaurant(db *sql.DB, name, address string, open bool, cuisine string) (int, error) {
	result, err := db.Exec("INSERT INTO Restaurant (Name, Address, Open, Cuisine) VALUES (?, ?, ?, ?)", name, address, open, cuisine)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func DeleteRestaurant(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM Restaurant WHERE ID = ?", id)
	return err
}
