package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

type Menu struct {
	ID           int `json:"id"`
	RestaurantID int `json:"restaurant_id"`
}

func MenusHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			HandleGetMenus(db, w, r)
		case "POST":
			action := r.FormValue("action")
			if action == "delete" {
				HandleDeleteMenu(db, w, r)
			} else {
				HandlePostMenu(db, w, r)
			}
		}
	}
}

func AddMenuHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
        <html>
        <head><title>Add Menu</title></head>
        <body>
        <form action="/menus" method="post">
            <label>Restaurant ID: <input type="number" name="restaurantID" /></label><br/>
            <input type="submit" value="Add Menu" />
        </form>
        </body>
        </html>
        `)
	}
}

func DeleteMenuHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
        <html>
        <head>
            <title>Delete Menu</title>
            <script src="https://unpkg.com/htmx.org@1.5.0"></script>
        </head>
        <body>
        <form hx-post="/menus" hx-target="#response">
            <input type="hidden" name="action" value="delete" />
            <label>Menu ID: <input type="number" name="ID" /></label><br/>
            <input type="submit" value="Delete Menu" />
        </form>
        <div id="response"></div>
        </body>
        </html>
        `)
	}
}

func HandleGetMenus(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	menus, err := FetchMenus(db)
	if err != nil {
		http.Error(w, "Failed to fetch menus: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<html><head><title>Menus List</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<nav><a href="/">Home</a> | <a href="/addMenu">Add Menu</a></nav>| <a href="/deleteMenu">Delete Menu</a></nav>`)
	fmt.Fprintf(w, `<h1>Menus</h1><ul>`)
	for _, m := range menus {
		fmt.Fprintf(w, `<li>ID: %d, Restaurant ID: %d</li>`, m.ID, m.RestaurantID)
	}
	fmt.Fprintf(w, `</ul></body></html>`)
}

func HandlePostMenu(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var m Menu
	m.RestaurantID, _ = strconv.Atoi(r.FormValue("restaurantID"))
	_, err := AddMenu(db, m.RestaurantID)
	if err != nil {
		http.Error(w, "Failed to add menu: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<html><head><title>Add Menu</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<h1>New menu added</h1></body></html>`)
}

func HandleDeleteMenu(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("ID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid menu ID", http.StatusBadRequest)
		return
	}
	if err := DeleteMenu(db, id); err != nil {
		http.Error(w, "Failed to delete menu: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<html><head><title>Delete Menu</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<h1>Menu deleted successfully</h1></body></html>`)
}

func FetchMenus(db *sql.DB) ([]Menu, error) {
	rows, err := db.Query("SELECT ID, RestaurantID FROM Menu")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var menus []Menu
	for rows.Next() {
		var m Menu
		if err := rows.Scan(&m.ID, &m.RestaurantID); err != nil {
			return nil, err
		}
		menus = append(menus, m)
	}
	return menus, nil
}

func AddMenu(db *sql.DB, restaurantID int) (int, error) {
	result, err := db.Exec("INSERT INTO Menu (RestaurantID) VALUES (?)", restaurantID)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func DeleteMenu(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM Menu WHERE ID = ?", id)
	return err
}
