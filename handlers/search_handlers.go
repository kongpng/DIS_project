package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Customer Search Handler
func SearchCustomersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleCustomerSearch(db, w, r)
	}
}

func handleCustomerSearch(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	name := params.Get("name")
	address := params.Get("address")
	email := params.Get("email")

	// Validate email if provided
	if email != "" && !ValidateEmail(email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT ID, Name, Address, Email FROM Customer WHERE 1=1")

	if name != "" {
		queryBuilder.WriteString(" AND Name LIKE '%" + name + "%'")
	}
	if address != "" {
		queryBuilder.WriteString(" AND Address LIKE '%" + address + "%'")
	}
	if email != "" {
		queryBuilder.WriteString(" AND Email LIKE '%" + email + "%'")
	}

	query := queryBuilder.String()
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Failed to execute search query: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var c Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Address, &c.Email); err != nil {
			http.Error(w, "Failed to scan customer: "+err.Error(), http.StatusInternalServerError)
			return
		}
		customers = append(customers, c)
	}

	fmt.Fprintf(w, `<html><head><title>Customer Search Results</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<h1>Customer Search Results</h1><ul>`)
	for _, c := range customers {
		fmt.Fprintf(w, `<li>ID: %d, Name: %s, Address: %s, Email: %s</li>`, c.ID, c.Name, c.Address, c.Email)
	}
	fmt.Fprintf(w, `</ul></body></html>`)
}

// Restaurant Search Handler
func SearchRestaurantsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleRestaurantSearch(db, w, r)
	}
}

func handleRestaurantSearch(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	name := params.Get("name")
	address := params.Get("address")
	open := params.Get("open")
	cuisine := params.Get("cuisine")

	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT ID, Name, Address, Open, Cuisine FROM Restaurant WHERE 1=1")

	if name != "" {
		queryBuilder.WriteString(" AND Name LIKE '%" + name + "%'")
	}
	if address != "" {
		queryBuilder.WriteString(" AND Address LIKE '%" + address + "%'")
	}
	if open != "" {
		queryBuilder.WriteString(" AND Open = " + open)
	}
	if cuisine != "" {
		queryBuilder.WriteString(" AND Cuisine LIKE '%" + cuisine + "%'")
	}

	query := queryBuilder.String()
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Failed to execute search query: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var restaurants []Restaurant
	for rows.Next() {
		var r Restaurant
		if err := rows.Scan(&r.ID, &r.Name, &r.Address, &r.Open, &r.Cuisine); err != nil {
			http.Error(w, "Failed to scan restaurant: "+err.Error(), http.StatusInternalServerError)
			return
		}
		restaurants = append(restaurants, r)
	}

	fmt.Fprintf(w, `<html><head><title>Restaurant Search Results</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<h1>Restaurant Search Results</h1><ul>`)
	for _, r := range restaurants {
		fmt.Fprintf(w, `<li>ID: %d, Name: %s, Address: %s, Open: %t, Cuisine: %s</li>`, r.ID, r.Name, r.Address, r.Open, r.Cuisine)
	}
	fmt.Fprintf(w, `</ul></body></html>`)
}

// Menu Search Handler
func SearchMenusHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleMenuSearch(db, w, r)
	}
}

func handleMenuSearch(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	restaurantIDStr := params.Get("restaurantID")

	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT ID, RestaurantID FROM Menu WHERE 1=1")

	if restaurantIDStr != "" {
		restaurantID, err := strconv.Atoi(restaurantIDStr)
		if err != nil {
			http.Error(w, "Invalid Restaurant ID", http.StatusBadRequest)
			return
		}
		queryBuilder.WriteString(" AND RestaurantID = " + strconv.Itoa(restaurantID))
	}

	query := queryBuilder.String()
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Failed to execute search query: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var menus []Menu
	for rows.Next() {
		var m Menu
		if err := rows.Scan(&m.ID, &m.RestaurantID); err != nil {
			http.Error(w, "Failed to scan menu: "+err.Error(), http.StatusInternalServerError)
			return
		}
		menus = append(menus, m)
	}

	fmt.Fprintf(w, `<html><head><title>Menu Search Results</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<h1>Menu Search Results</h1><ul>`)
	for _, m := range menus {
		fmt.Fprintf(w, `<li>ID: %d, Restaurant ID: %d</li>`, m.ID, m.RestaurantID)
	}
	fmt.Fprintf(w, `</ul></body></html>`)
}

// Dish Search Handler
func SearchDishesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleDishSearch(db, w, r)
	}
}

func handleDishSearch(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	menuIDStr := params.Get("menuID")
	name := params.Get("name")
	priceStr := params.Get("price")
	vegan := params.Get("vegan")
	shellfish := params.Get("shellfish")
	nuts := params.Get("nuts")

	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT ID, MenuID, Name, Price, Vegan, Shellfish, Nuts FROM Dishes WHERE 1=1")

	if menuIDStr != "" {
		menuID, err := strconv.Atoi(menuIDStr)
		if err != nil {
			http.Error(w, "Invalid Menu ID", http.StatusBadRequest)
			return
		}
		queryBuilder.WriteString(" AND MenuID = " + strconv.Itoa(menuID))
	}
	if name != "" {
		queryBuilder.WriteString(" AND Name LIKE '%" + name + "%'")
	}
	if priceStr != "" {
		price, err := strconv.Atoi(priceStr)
		if err != nil {
			http.Error(w, "Invalid Price", http.StatusBadRequest)
			return
		}
		queryBuilder.WriteString(" AND Price = " + strconv.Itoa(price))
	}
	if vegan != "" {
		queryBuilder.WriteString(" AND Vegan = " + vegan)
	}
	if shellfish != "" {
		queryBuilder.WriteString(" AND Shellfish = " + shellfish)
	}
	if nuts != "" {
		queryBuilder.WriteString(" AND Nuts = " + nuts)
	}

	query := queryBuilder.String()
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Failed to execute search query: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var dishes []Dish
	for rows.Next() {
		var d Dish
		if err := rows.Scan(&d.ID, &d.MenuID, &d.Name, &d.Price, &d.Vegan, &d.Shellfish, &d.Nuts); err != nil {
			http.Error(w, "Failed to scan dish: "+err.Error(), http.StatusInternalServerError)
			return
		}
		dishes = append(dishes, d)
	}

	fmt.Fprintf(w, `<html><head><title>Dish Search Results</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<h1>Dish Search Results</h1><ul>`)
	for _, d := range dishes {
		fmt.Fprintf(w, `<li>ID: %d, Menu ID: %d, Name: %s, Price: %d, Vegan: %t, Shellfish: %t, Nuts: %t</li>`, d.ID, d.MenuID, d.Name, d.Price, d.Vegan, d.Shellfish, d.Nuts)
	}
	fmt.Fprintf(w, `</ul></body></html>`)
}

// Order Search Handler
func SearchOrdersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleOrderSearch(db, w, r)
	}
}

func handleOrderSearch(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	customerIDStr := params.Get("customerID")
	date := params.Get("date")

	// Validate date if provided
	if date != "" && !ValidateDate(date) {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT ID, CustomerID, Date FROM CustomerOrder WHERE 1=1")

	if customerIDStr != "" {
		customerID, err := strconv.Atoi(customerIDStr)
		if err != nil {
			http.Error(w, "Invalid Customer ID", http.StatusBadRequest)
			return
		}
		queryBuilder.WriteString(" AND CustomerID = " + strconv.Itoa(customerID))
	}
	if date != "" {
		queryBuilder.WriteString(" AND Date = " + date)
	}

	query := queryBuilder.String()
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Failed to execute search query: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		var dateInt int64
		if err := rows.Scan(&o.ID, &o.CustomerID, &dateInt); err != nil {
			http.Error(w, "Failed to scan order: "+err.Error(), http.StatusInternalServerError)
			return
		}
		o.Date = time.Unix(dateInt, 0)
		orders = append(orders, o)
	}

	fmt.Fprintf(w, `<html><head><title>Order Search Results</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<h1>Order Search Results</h1><ul>`)
	for _, o := range orders {
		fmt.Fprintf(w, `<li>Order ID: %d, Customer ID: %d, Date: %s</li>`, o.ID, o.CustomerID, o.Date.Format("2006-01-02"))
	}
	fmt.Fprintf(w, `</ul></body></html>`)
}

// Delivery Search Handler
func SearchDeliveriesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleDeliverySearch(db, w, r)
	}
}

func handleDeliverySearch(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	orderIDStr := params.Get("orderID")
	deliveryTimeStr := params.Get("deliveryTime")
	deliveryCostStr := params.Get("deliveryCost")

	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT ID, OrderID, DeliveryTime, DeliveryCost FROM DeliveryService WHERE 1=1")

	if orderIDStr != "" {
		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil {
			http.Error(w, "Invalid Order ID", http.StatusBadRequest)
			return
		}
		queryBuilder.WriteString(" AND OrderID = " + strconv.Itoa(orderID))
	}
	if deliveryTimeStr != "" {
		deliveryTime, err := strconv.Atoi(deliveryTimeStr)
		if err != nil {
			http.Error(w, "Invalid Delivery Time", http.StatusBadRequest)
			return
		}
		queryBuilder.WriteString(" AND DeliveryTime = " + strconv.Itoa(deliveryTime))
	}
	if deliveryCostStr != "" {
		deliveryCost, err := strconv.Atoi(deliveryCostStr)
		if err != nil {
			http.Error(w, "Invalid Delivery Cost", http.StatusBadRequest)
			return
		}
		queryBuilder.WriteString(" AND DeliveryCost = " + strconv.Itoa(deliveryCost))
	}

	query := queryBuilder.String()
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Failed to execute search query: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var deliveries []Delivery
	for rows.Next() {
		var d Delivery
		if err := rows.Scan(&d.ID, &d.OrderID, &d.DeliveryTime, &d.DeliveryCost); err != nil {
			http.Error(w, "Failed to scan delivery: "+err.Error(), http.StatusInternalServerError)
			return
		}
		deliveries = append(deliveries, d)
	}

	fmt.Fprintf(w, `<html><head><title>Delivery Search Results</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<h1>Delivery Search Results</h1><ul>`)
	for _, d := range deliveries {
		fmt.Fprintf(w, `<li>Delivery ID: %d, Order ID: %d, Delivery Time: %d, Delivery Cost: %d</li>`, d.ID, d.OrderID, d.DeliveryTime, d.DeliveryCost)
	}
	fmt.Fprintf(w, `</ul></body></html>`)
}
