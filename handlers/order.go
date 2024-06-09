package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Order struct {
	ID         int       `json:"id"`
	CustomerID int       `json:"customer_id"`
	Date       time.Time `json:"date"`
}

func OrdersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			HandleGetOrders(db, w, r)
		case "POST":
			action := r.FormValue("action")
			if action == "delete" {
				HandleDeleteOrder(db, w, r)
			} else {
				HandlePostOrder(db, w, r)
			}
		}
	}
}

func AddOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
        <html>
        <head><title>Add Order</title></head>
        <body>
        <form action="/orders" method="post">
            <label>Customer ID: <input type="number" name="customerID" /></label><br/>
            <label>Date (YYYY-MM-DD): <input type="date" name="date" /></label><br/>
            <input type="submit" value="Add Order" />
        </form>
        </body>
        </html>
        `)
	}
}

func DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
        <html>
        <head>
            <title>Delete Order</title>
            <script src="https://unpkg.com/htmx.org@1.9.12"></script>
        </head>
        <body>
        <form hx-post="/orders" hx-target="#response">
            <input type="hidden" name="action" value="delete" />
            <label>Order ID: <input type="number" name="ID" /></label><br/>
            <input type="submit" value="Delete Order" />
        </form>
        <div id="response"></div>
        </body>
        </html>
        `)
	}
}

func HandleGetOrders(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	orders, err := FetchOrders(db)
	if err != nil {
		http.Error(w, "Failed to fetch orders: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<html><head><title>Orders List</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<nav><a href="/">Home</a> | <a href="/addOrder">Add Order</a> | <a href="/deleteOrder">Delete Order</a></nav>`)
	fmt.Fprintf(w, `<h1>Orders</h1><ul>`)
	for _, o := range orders {
		fmt.Fprintf(w, `<li>Order ID: %d, Customer ID: %d, Date: %s</li>`, o.ID, o.CustomerID, o.Date.Format("2006-01-02"))
	}
	fmt.Fprintf(w, `</ul></body></html>`)
}

func HandlePostOrder(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var o Order
	o.CustomerID, _ = strconv.Atoi(r.FormValue("customerID"))
	o.Date, _ = time.Parse("2006-01-02", r.FormValue("date"))
	_, err := AddOrder(db, o.CustomerID, o.Date)
	if err != nil {
		http.Error(w, "Failed to add order: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<html><head><title>Add Order</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<h1>New order added</h1></body></html>`)
}

func HandleDeleteOrder(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("ID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return
	}

	err = DeleteOrder(db, id)
	if err != nil {
		if err.Error() == fmt.Sprintf("order with ID %d not found", id) {
			fmt.Fprintf(w, ``)                              // weird
			http.Error(w, err.Error(), http.StatusNotFound) // 404 Not Found
			return
		} else {
			fmt.Fprintf(w, ``)
			http.Error(w, "Failed to delete order: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	fmt.Fprintf(w, `<html><head><title>Delete Order</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<h1>Order deleted successfully</h1></body></html>`)
}

func FetchOrders(db *sql.DB) ([]Order, error) {
	rows, err := db.Query("SELECT ID, CustomerID, Date FROM CustomerOrder")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []Order
	for rows.Next() {
		var o Order
		var dateInt int64
		if err := rows.Scan(&o.ID, &o.CustomerID, &dateInt); err != nil {
			return nil, err
		}
		o.Date = time.Unix(dateInt, 0)
		orders = append(orders, o)
	}
	return orders, nil
}

func AddOrder(db *sql.DB, customerID int, date time.Time) (int, error) {
	result, err := db.Exec("INSERT INTO CustomerOrder (CustomerID, Date) VALUES (?, ?)", customerID, date.Unix())
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func DeleteOrder(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM CustomerOrder WHERE ID = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order with ID %d not found", id)
	}

	return nil
}
