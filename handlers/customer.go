package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

type Customer struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Email   string `json:"email"`
}

// Customer Handlers

func CustomersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			HandleGetCustomers(db, w, r)
		case "POST":
			action := r.FormValue("action")
			if action == "delete" {
				HandleDeleteCustomer(db, w, r)
			} else {
				HandlePostCustomer(db, w, r)
			}
		}
	}
}

func HandleGetCustomers(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	customers, err := FetchCustomers(db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching customers: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<html><head><title>Customers List</title><link rel="stylesheet" type="text/css" href="/static/style.css"><script src="https://unpkg.com/htmx.org@1.9.12"></script></head><body>`)
	fmt.Fprintf(w, `<nav><a href="/">Home</a> | <a href="/addCustomer">Add Customer</a> | <a href="/deleteCustomer">Delete Customer</a></nav>`)
	fmt.Fprintf(w, `<h1>Customers</h1><ul id="customerList">`)
	for _, c := range customers {
		fmt.Fprintf(w, `<li>ID: %d, Name: %s, Address: %s, Email: %s</li>`, c.ID, c.Name, c.Address, c.Email)
	}
	fmt.Fprintf(w, `</ul></body></html>`)
}

func HandlePostCustomer(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var c Customer
	c.Name = r.FormValue("name")
	c.Address = r.FormValue("address")
	c.Email = r.FormValue("email")
	id, err := AddCustomer(db, c.Name, c.Address, c.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding customer: %v", err), http.StatusInternalServerError)
		return
	}
	c.ID = id
	fmt.Fprintf(w, `<div>New customer added: ID: %d, Name: %s, Address: %s, Email: %s</div>`, c.ID, c.Name, c.Address, c.Email)
}

func HandleDeleteCustomer(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	idStr := r.FormValue("ID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return
	}

	err = DeleteCustomer(db, id)
	if err != nil {
		if err.Error() == fmt.Sprintf("customer with ID %d not found", id) {
			fmt.Fprintf(w, ``)                              // weird
			http.Error(w, err.Error(), http.StatusNotFound) // 404 Not Found
			return
		} else {
			fmt.Fprintf(w, ``)
			http.Error(w, "Failed to delete customer: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	fmt.Fprintf(w, `<html><head><title>Delete Customer</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
	fmt.Fprintf(w, `<h1>Customer deleted successfully</h1></body></html>`)
}

func FetchCustomers(db *sql.DB) ([]Customer, error) {
	rows, err := db.Query("SELECT ID, Name, Address, Email FROM Customer")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var customers []Customer
	for rows.Next() {
		var c Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Address, &c.Email); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}

func AddCustomer(db *sql.DB, name, address, email string) (int, error) {
	result, err := db.Exec("INSERT INTO Customer (Name, Address, Email) VALUES (?, ?, ?)", name, address, email)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func DeleteCustomer(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM Customer WHERE ID = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("customer service with ID %d not found", id)
	}

	return nil
}

func AddCustomerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
        <html>
        <head>
            <title>Add Customer</title>
            <script src="https://unpkg.com/htmx.org@1.9.12"></script>
        </head>
        <body>
        <form hx-post="/customers" hx-target="#response">
            <label>Name: <input type="text" name="name" /></label><br/>
            <label>Address: <input type="text" name="address" /></label><br/>
            <label>Email: <input type="email" name="email" /></label><br/>
            <input type="submit" value="Add Customer" />
        </form>
        <div id="response"></div>
        </body>
        </html>
        `)
	}
}

func DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
        <html>
        <head>
            <title>Delete Customer</title>
            <script src="https://unpkg.com/htmx.org@1.9.12"></script>
        </head>
        <body>
        <form hx-post="/customers" hx-target="#response">
            <input type="hidden" name="action" value="delete" />
            <label>Customer ID: <input type="number" name="ID" /></label><br/>
            <input type="submit" value="Delete Customer" />
        </form>
        <div id="response"></div>
        </body>
        </html>
        `)
	}
}
