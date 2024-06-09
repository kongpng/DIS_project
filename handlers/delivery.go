package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

type Delivery struct {
	ID           int `json:"id"`
	OrderID      int `json:"order_id"`
	DeliveryTime int `json:"delivery_time"`
	DeliveryCost int `json:"delivery_cost"`
}

func DeliveryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			HandleGetDeliveries(db, w, r)
		case "POST":
			action := r.FormValue("action")
			if action == "delete" {
				HandleDeleteDelivery(db, w, r)
			} else {
				HandlePostDelivery(db, w, r)
			}
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
                <label>Delivery Time (HH:MM): <input type="time" name="deliveryTime" /></label><br/>
                <label>Delivery Cost: <input type="number" name="deliveryCost" /></label><br/>
                <input type="submit" value="Add Delivery" />
            </form>
            </body>
            </html>
            `)
		}
	}
}

func DeleteDeliveryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
        <html>
        <head>
            <title>Delete Delivery</title>
            <script src="https://unpkg.com/htmx.org@1.9.12"></script>
        </head>
        <body>
        <form hx-post="/deliveries" hx-target="#response">
            <input type="hidden" name="action" value="delete" />
            <label>Delivery ID: <input type="number" name="ID" /></label><br/>
            <input type="submit" value="Delete Delivery" />
        </form>
        <div id="response"></div>
        </body>
        </html>
        `)
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
	fmt.Fprintf(w, `<nav><a href="/">Home</a> | <a href="/addDelivery">Add Delivery</a> | <a href="/deleteDelivery">Delete Delivery</a></nav>`)
	fmt.Fprintf(w, `<h1>Deliveries</h1><ul>`)
	for _, d := range deliveries {
		fmt.Fprintf(w, `<li>Delivery ID: %d, Order ID: %d, Time: %d, Cost: %d</li>`, d.ID, d.OrderID, d.DeliveryTime, d.DeliveryCost)
	}
	fmt.Fprintf(w, `</ul></body></html>`)
}

func HandlePostDelivery(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var d Delivery
	d.OrderID, _ = strconv.Atoi(r.FormValue("orderID"))
	d.DeliveryTime, _ = strconv.Atoi(r.FormValue("deliveryTime"))
	d.DeliveryCost, _ = strconv.Atoi(r.FormValue("deliveryCost"))
	_, err := AddDelivery(db, d.OrderID, d.DeliveryTime, d.DeliveryCost)
	if err != nil {
		http.Error(w, "Failed to add delivery: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<html><head><title>Add Delivery</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<h1>New delivery added</h1></body></html>`)
}

func HandleDeleteDelivery(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("ID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return
	}
	err = DeleteDelivery(db, id)
	if err != nil {
		if err.Error() == fmt.Sprintf("delivery with ID %d not found", id) {
			fmt.Fprintf(w, ``)                              // weird
			http.Error(w, err.Error(), http.StatusNotFound) // 404 Not Found
			return
		} else {
			fmt.Fprintf(w, ``)
			http.Error(w, "Failed to delete delivery: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	fmt.Fprintf(w, `<html><head><title>Delete Delivery</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<h1>Delivery deleted successfully</h1></body></html>`)
}

func FetchDeliveries(db *sql.DB) ([]Delivery, error) {
	rows, err := db.Query("SELECT ID, OrderID, DeliveryTime, DeliveryCost FROM DeliveryService")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var deliveries []Delivery
	for rows.Next() {
		var d Delivery
		if err := rows.Scan(&d.ID, &d.OrderID, &d.DeliveryTime, &d.DeliveryCost); err != nil {
			return nil, err
		}
		deliveries = append(deliveries, d)
	}
	return deliveries, nil
}

func AddDelivery(db *sql.DB, orderID, deliveryTime, deliveryCost int) (int, error) {
	result, err := db.Exec("INSERT INTO DeliveryService (OrderID, DeliveryTime, DeliveryCost) VALUES (?, ?, ?)", orderID, deliveryTime, deliveryCost)
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
	result, err := db.Exec("DELETE FROM DeliveryService WHERE ID = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("delivery service with ID %d not found", id)
	}

	return nil
}
