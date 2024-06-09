package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

type Delivery struct {
	ID      int `json:"id"`
	OrderID int `json:"order_id"`
	DTime   int `json:"dtime"`
	DCost   int `json:"dcost"`
}

func DeliveryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			HandleGetDeliveries(db, w, r)
		case "POST":
			HandlePostDelivery(db, w, r)
		}
	}
}

func AddDeliveryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Fprintf(w, `
            <html>
            <head><title>Add Delivery</title></head>
            <body>
            <form action="/deliveries" method="post">
                <label>Order ID: <input type="number" name="orderID" /></label><br/>
                <label>Delivery Time (HH:MM): <input type="time" name="dTime" /></label><br/>
                <label>Delivery Cost: <input type="number" name="dCost" /></label><br/>
                <input type="submit" value="Add Delivery" />
            </form>
            </body>
            </html>
            `)
		}
	}
}

func DeleteDeliveryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			HandleDeleteDelivery(db, w, r)
		}
	}
}

// Delivery Handlers
func HandleGetDeliveries(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	deliveries, err := FetchDeliveries(db)
	if err != nil {
		http.Error(w, "Failed to fetch deliveries: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<html><head><title>Deliveries</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<nav><a href="/">Home</a> | <a href="/addDelivery">Add Delivery</a></nav>`)
	fmt.Fprintf(w, `<h1>Deliveries</h1><ul>`)
	for _, d := range deliveries {
		fmt.Fprintf(w, `<li>Delivery ID: %d, Order ID: %d, Time: %d, Cost: %d</li>`, d.ID, d.OrderID, d.DTime, d.DCost)
	}
	fmt.Fprintf(w, `</ul></body></html>`)
}

func HandlePostDelivery(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var d Delivery
	d.OrderID, _ = strconv.Atoi(r.FormValue("orderID"))
	d.DTime, _ = strconv.Atoi(r.FormValue("dTime"))
	d.DCost, _ = strconv.Atoi(r.FormValue("dCost"))
	_, err := AddDelivery(db, d.OrderID, d.DTime, d.DCost)
	if err != nil {
		http.Error(w, "Failed to add delivery: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<html><head><title>Add Delivery</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<h1>New delivery added</h1></body></html>`)
}

func HandleDeleteDelivery(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid delivery ID", http.StatusBadRequest)
		return
	}
	if err := DeleteDelivery(db, id); err != nil {
		http.Error(w, "Failed to delete delivery: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<html><head><title>Delete Delivery</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<h1>Delivery deleted successfully</h1></body></html>`)
}

func FetchDeliveries(db *sql.DB) ([]Delivery, error) {
	rows, err := db.Query("SELECT ID, OrderID, DTime, DCost FROM DeliveryService")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var deliveries []Delivery
	for rows.Next() {
		var d Delivery
		if err := rows.Scan(&d.ID, &d.OrderID, &d.DTime, &d.DCost); err != nil {
			return nil, err
		}
		deliveries = append(deliveries, d)
	}
	return deliveries, nil
}

func AddDelivery(db *sql.DB, orderID, dTime, dCost int) (int, error) {
	result, err := db.Exec("INSERT INTO DeliveryService (OrderID, DTime, DCost) VALUES (?, ?, ?)", orderID, dTime, dCost)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func DeleteDelivery(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM DeliveryService WHERE ID = ?", id)
	return err
}
