package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

type Dish struct {
	ID        int    `json:"id"`
	MenuID    int    `json:"menu_id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Vegan     bool   `json:"vegan"`
	Shellfish bool   `json:"shellfish"`
	Nuts      bool   `json:"nuts"`
}

func DishesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			HandleGetDishes(db, w, r)
		case "POST":
			action := r.FormValue("action")
			if action == "delete" {
				HandleDeleteDish(db, w, r)
			} else {
				HandlePostDish(db, w, r)
			}
		}
	}
}

func AddDishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
        <html>
        <head><title>Add Dish</title></head>
        <body>
        <form action="/dishes" method="post">
            <label>Menu ID: <input type="number" name="menuID" /></label><br/>
            <label>Name: <input type="text" name="name" /></label><br/>
            <label>Price: <input type="number" name="price" /></label><br/>
            <label>Vegan: <input type="checkbox" name="vegan" /></label><br/>
            <label>Contains Shellfish: <input type="checkbox" name="shellfish" /></label><br/>
            <label>Contains Nuts: <input type="checkbox" name="nuts" /></label><br/>
            <input type="submit" value="Add Dish" />
        </form>
        </body>
        </html>
        `)
	}
}

func DeleteDishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
        <html>
        <head>
            <title>Delete Dish</title>
            <script src="https://unpkg.com/htmx.org@1.9.12"></script>
        </head>
        <body>
        <form hx-post="/dishes" hx-target="#response">
            <input type="hidden" name="action" value="delete" />
            <label>Dish ID: <input type="number" name="ID" /></label><br/>
            <input type="submit" value="Delete Dish" />
        </form>
        <div id="response"></div>
        </body>
        </html>
        `)
	}
}

// Dish Handlers
func HandleGetDishes(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	dishes, err := FetchDishes(db)
	if err != nil {
		http.Error(w, "Failed to fetch dishes: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<html><head><title>Dishes List</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<nav><a href="/">Home</a> | <a href="/addDish">Add Dish</a></nav> | <a href="/deleteDish">Delete Dish</a></nav>`)
	fmt.Fprintf(w, `<h1>Dishes</h1><ul>`)
	for _, d := range dishes {
		fmt.Fprintf(w, `<li>ID: %d, Menu ID: %d, Name: %s, Price: %d, Vegan: %t, Shellfish: %t, Nuts: %t</li>`, d.ID, d.MenuID, d.Name, d.Price, d.Vegan, d.Shellfish, d.Nuts)
	}
	fmt.Fprintf(w, `</ul></body></html>`)
}

func HandlePostDish(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var d Dish
	d.MenuID, _ = strconv.Atoi(r.FormValue("menuID"))
	d.Name = r.FormValue("name")
	d.Price, _ = strconv.Atoi(r.FormValue("price"))
	d.Vegan = r.FormValue("vegan") == "on"
	d.Shellfish = r.FormValue("shellfish") == "on"
	d.Nuts = r.FormValue("nuts") == "on"
	_, err := AddDish(db, d.MenuID, d.Name, d.Price, d.Vegan, d.Shellfish, d.Nuts)
	if err != nil {
		http.Error(w, "Failed to add dish: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<html><head><title>Add Dish</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<h1>New dish added</h1></body></html>`)
}

func HandleDeleteDish(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("ID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return
	}
	err = DeleteDish(db, id)
	if err != nil {
		if err.Error() == fmt.Sprintf("dish with ID %d not found", id) {
			fmt.Fprintf(w, ``)                              // weird
			http.Error(w, err.Error(), http.StatusNotFound) // 404 Not Found
			return
		} else {
			fmt.Fprintf(w, ``)
			http.Error(w, "Failed to delete dish: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintf(w, `<html><head><title>Delete Dish</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<h1>Dish deleted successfully</h1></body></html>`)
}

func FetchDishes(db *sql.DB) ([]Dish, error) {
	rows, err := db.Query("SELECT ID, MenuID, Name, Price, Vegan, Shellfish, Nuts FROM Dishes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var dishes []Dish
	for rows.Next() {
		var d Dish
		if err := rows.Scan(&d.ID, &d.MenuID, &d.Name, &d.Price, &d.Vegan, &d.Shellfish, &d.Nuts); err != nil {
			return nil, err
		}
		dishes = append(dishes, d)
	}
	return dishes, nil
}

func AddDish(db *sql.DB, menuID int, name string, price int, vegan, shellfish, nuts bool) (int, error) {
	result, err := db.Exec("INSERT INTO Dishes (MenuID, Name, Price, Vegan, Shellfish, Nuts) VALUES (?, ?, ?, ?, ?, ?)", menuID, name, price, vegan, shellfish, nuts)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func DeleteDish(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM Dishes WHERE ID = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("dish with ID %d not found", id)
	}

	return nil
}
