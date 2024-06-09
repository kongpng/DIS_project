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
			HandlePostOrder(db, w, r)
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

func DeleteOrderHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			HandleDeleteOrder(db, w, r)
		}
	}
}

func HandleGetOrders(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	orders, err := FetchOrders(db)
	if err != nil {
		http.Error(w, "Failed to fetch orders: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<html><head><title>Orders List</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<nav><a href="/">Home</a> | <a href="/addOrder">Add Order</a></nav>`)
	fmt.Fprintf(w, `<h1>Orders</h1><ul>`)
	for _, o := range orders {
		fmt.Fprintf(w, `<li>Order ID: %d, Customer ID: %d, Date: %s</li>`, o.ID, o.CustomerID, o.Date)
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
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}
	if err := DeleteOrder(db, id); err != nil {
		http.Error(w, "Failed to delete order: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<html><head><title>Delete Order</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<h1>Order deleted successfully</h1></body></html>`)
}

func FetchOrders(db *sql.DB) ([]Order, error) {
	rows, err := db.Query("SELECT ID, CustomerID, Date FROM Orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.Date); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func AddOrder(db *sql.DB, customerID int, date time.Time) (int, error) {
	result, err := db.Exec("INSERT INTO Orders (CustomerID, Date) VALUES (?, ?)", customerID, date)
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
	_, err := db.Exec("DELETE FROM Orders WHERE ID = ?", id)
	return err
}
