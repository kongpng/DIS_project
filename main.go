package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, err := sql.Open("sqlite3", "db/Data.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    setupRoutes(db)

    log.Println("Server starting on http://localhost:8080/")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupRoutes(db *sql.DB) {
    // Home page
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, `
        <html>
        <head>
            <title>Welcome</title>
            <link rel='stylesheet' type='text/css' href='/static/style.css'>
        </head>
        <body>
            <h1>Welcome to the API Server</h1>
            <nav>
                <ul>
                    <li><a href="/customers">Customers</a></li>
                    <li><a href="/orders">Orders</a></li>
                    <li><a href="/deliveries">Deliveries</a></li>
                    <li><a href="/restaurants">Restaurants</a></li>
                    <li><a href="/menus">Menus</a></li>
                    <li><a href="/dishes">Dishes</a></li>
                </ul>
            </nav>
        </body>
        </html>
        `)
    })

    // Customer handlers
    http.HandleFunc("/customers", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            handleGetCustomers(db, w, r)
        case "POST":
            handlePostCustomer(db, w, r)
        }
    })

    http.HandleFunc("/addCustomer", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "GET" {
            fmt.Fprintf(w, `
            <html>
            <head><title>Add Customer</title></head>
            <body>
            <form action="/customers" method="post">
                <label>Name: <input type="text" name="name" /></label><br/>
                <label>Address: <input type="text" name="address" /></label><br/>
                <label>Email: <input type="email" name="email" /></label><br/>
                <input type="submit" value="Add Customer" />
            </form>
            </body>
            </html>
            `)
        }
    })

    http.HandleFunc("/deleteCustomer", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            handleDeleteCustomer(db, w, r)
        }
    })

    // Order handlers
    http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            handleGetOrders(db, w, r)
        case "POST":
            handlePostOrder(db, w, r)
        }
    })

    http.HandleFunc("/addOrder", func(w http.ResponseWriter, r *http.Request) {
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
    })

    http.HandleFunc("/deleteOrder", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            handleDeleteOrder(db, w, r)
        }
    })

    // Delivery handlers
    http.HandleFunc("/deliveries", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            handleGetDeliveries(db, w, r)
        case "POST":
            handlePostDelivery(db, w, r)
        }
    })

    http.HandleFunc("/addDelivery", func(w http.ResponseWriter, r *http.Request) {
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
    })

    http.HandleFunc("/deleteDelivery", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            handleDeleteDelivery(db, w, r)
        }
    })

    // Restaurant handlers
    http.HandleFunc("/restaurants", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            handleGetRestaurants(db, w, r)
        case "POST":
            handlePostRestaurant(db, w, r)
        }
    })

    http.HandleFunc("/addRestaurant", func(w http.ResponseWriter, r *http.Request) {
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
    })

    http.HandleFunc("/deleteRestaurant", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            handleDeleteRestaurant(db, w, r)
        }
    })

    // Menu handlers
    http.HandleFunc("/menus", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            handleGetMenus(db, w, r)
        case "POST":
            handlePostMenu(db, w, r)
        }
    })

    http.HandleFunc("/addMenu", func(w http.ResponseWriter, r *http.Request) {
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
    })

    http.HandleFunc("/deleteMenu", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            handleDeleteMenu(db, w, r)
        }
    })

    // Dish handlers
    http.HandleFunc("/dishes", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            handleGetDishes(db, w, r)
        case "POST":
            handlePostDish(db, w, r)
        }
    })

    http.HandleFunc("/addDish", func(w http.ResponseWriter, r *http.Request) {
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
    })

    http.HandleFunc("/deleteDish", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            handleDeleteDish(db, w, r)
        }
    })

    // Serve static files
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    // Catch-all handler for undefined routes
    http.HandleFunc("/catch-all", func(w http.ResponseWriter, r *http.Request) {
        http.Error(w, "This route is not defined: "+r.URL.Path, http.StatusNotFound)
    })
}

// Customer Handlers
func handleGetCustomers(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    customers, err := fetchCustomers(db)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error fetching customers: %v", err), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, `<html><head><title>Customers List</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
    fmt.Fprintf(w, `<nav><a href="/">Home</a> | <a href="/addCustomer">Add Customer</a></nav>`)
    fmt.Fprintf(w, `<h1>Customers</h1><ul>`)
    for _, c := range customers {
        fmt.Fprintf(w, `<li>ID: %d, Name: %s, Address: %s, Email: %s</li>`, c.ID, c.Name, c.Address, c.Email)
    }
    fmt.Fprintf(w, `</ul></body></html>`)
}

func handlePostCustomer(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var c Customer
    c.Name = r.FormValue("name")
    c.Address = r.FormValue("address")
    c.Email = r.FormValue("email")
    _, err := addCustomer(db, c.Name, c.Address, c.Email)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error adding customer: %v", err), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, `<html><head><title>Add Customer</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
    fmt.Fprintf(w, `<h1>New customer added</h1></body></html>`)
}

func handleDeleteCustomer(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    idStr := r.FormValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid customer ID", http.StatusBadRequest)
        return
    }
    if err := deleteCustomer(db, id); err != nil {
        http.Error(w, fmt.Sprintf("Error deleting customer: %v", err), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, `<html><head><title>Delete Customer</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
    fmt.Fprintf(w, `<h1>Customer deleted successfully</h1></body></html>`)
}

func fetchCustomers(db *sql.DB) ([]Customer, error) {
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

func addCustomer(db *sql.DB, name, address, email string) (int, error) {
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

func deleteCustomer(db *sql.DB, id int) error {
    _, err := db.Exec("DELETE FROM Customer WHERE ID = ?", id)
    return err
}

// Order Handlers
func handleGetOrders(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    orders, err := fetchOrders(db)
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

func handlePostOrder(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var o Order
    o.CustomerID, _ = strconv.Atoi(r.FormValue("customerID"))
    o.Date, _ = time.Parse("2006-01-02", r.FormValue("date"))
    _, err := addOrder(db, o.CustomerID, o.Date)
    if err != nil {
        http.Error(w, "Failed to add order: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, `<html><head><title>Add Order</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
    fmt.Fprintf(w, `<h1>New order added</h1></body></html>`)
}

func handleDeleteOrder(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    idStr := r.FormValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid order ID", http.StatusBadRequest)
        return
    }
    if err := deleteOrder(db, id); err != nil {
        http.Error(w, "Failed to delete order: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, `<html><head><title>Delete Order</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
    fmt.Fprintf(w, `<h1>Order deleted successfully</h1></body></html>`)
}

func fetchOrders(db *sql.DB) ([]Order, error) {
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

func addOrder(db *sql.DB, customerID int, date time.Time) (int, error) {
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

func deleteOrder(db *sql.DB, id int) error {
    _, err := db.Exec("DELETE FROM Orders WHERE ID = ?", id)
    return err
}

// Delivery Handlers
func handleGetDeliveries(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    deliveries, err := fetchDeliveries(db)
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

func handlePostDelivery(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var d Delivery
    d.OrderID, _ = strconv.Atoi(r.FormValue("orderID"))
    d.DTime, _ = strconv.Atoi(r.FormValue("dTime"))
    d.DCost, _ = strconv.Atoi(r.FormValue("dCost"))
    _, err := addDelivery(db, d.OrderID, d.DTime, d.DCost)
    if err != nil {
        http.Error(w, "Failed to add delivery: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, `<html><head><title>Add Delivery</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
    fmt.Fprintf(w, `<h1>New delivery added</h1></body></html>`)
}

func handleDeleteDelivery(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    idStr := r.FormValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid delivery ID", http.StatusBadRequest)
        return
    }
    if err := deleteDelivery(db, id); err != nil {
        http.Error(w, "Failed to delete delivery: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, `<html><head><title>Delete Delivery</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
    fmt.Fprintf(w, `<h1>Delivery deleted successfully</h1></body></html>`)
}

func fetchDeliveries(db *sql.DB) ([]Delivery, error) {
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

func addDelivery(db *sql.DB, orderID, dTime, dCost int) (int, error) {
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

func deleteDelivery(db *sql.DB, id int) error {
    _, err := db.Exec("DELETE FROM DeliveryService WHERE ID = ?", id)
    return err
}

// Restaurant Handlers
func handleGetRestaurants(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    restaurants, err := fetchRestaurants(db)
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

func handlePostRestaurant(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var res Restaurant
    res.Name = r.FormValue("name")
    res.Address = r.FormValue("address")
    res.Open = r.FormValue("open") == "on"
    res.Cuisine = r.FormValue("cuisine")
    _, err := addRestaurant(db, res.Name, res.Address, res.Open, res.Cuisine)
    if err != nil {
        http.Error(w, "Failed to add restaurant: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, `<html><head><title>Add Restaurant</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
    fmt.Fprintf(w, `<h1>New restaurant added</h1></body></html>`)
}

func handleDeleteRestaurant(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    idStr := r.FormValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
        return
    }
    if err := deleteRestaurant(db, id); err != nil {
        http.Error(w, "Failed to delete restaurant: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, `<html><head><title>Delete Restaurant</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
    fmt.Fprintf(w, `<h1>Restaurant deleted successfully</h1></body></html>`)
}

func fetchRestaurants(db *sql.DB) ([]Restaurant, error) {
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

func addRestaurant(db *sql.DB, name, address string, open bool, cuisine string) (int, error) {
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

func deleteRestaurant(db *sql.DB, id int) error {
    _, err := db.Exec("DELETE FROM Restaurant WHERE ID = ?", id)
    return err
}

// Menu Handlers
func handleGetMenus(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    menus, err := fetchMenus(db)
    if err != nil {
        http.Error(w, "Failed to fetch menus: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, `<html><head><title>Menus List</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
    fmt.Fprintf(w, `<nav><a href="/">Home</a> | <a href="/addMenu">Add Menu</a></nav>`)
    fmt.Fprintf(w, `<h1>Menus</h1><ul>`)
    for _, m := range menus {
        fmt.Fprintf(w, `<li>ID: %d, Restaurant ID: %d</li>`, m.ID, m.RestaurantID)
    }
    fmt.Fprintf(w, `</ul></body></html>`)
}

func handlePostMenu(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var m Menu
    m.RestaurantID, _ = strconv.Atoi(r.FormValue("restaurantID"))
    _, err := addMenu(db, m.RestaurantID)
    if err != nil {
        http.Error(w, "Failed to add menu: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, `<html><head><title>Add Menu</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
    fmt.Fprintf(w, `<h1>New menu added</h1></body></html>`)
}

func handleDeleteMenu(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    idStr := r.FormValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid menu ID", http.StatusBadRequest)
        return
    }
    if err := deleteMenu(db, id); err != nil {
        http.Error(w, "Failed to delete menu: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, `<html><head><title>Delete Menu</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
    fmt.Fprintf(w, `<h1>Menu deleted successfully</h1></body></html>`)
}

func fetchMenus(db *sql.DB) ([]Menu, error) {
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

func addMenu(db *sql.DB, restaurantID int) (int, error) {
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

func deleteMenu(db *sql.DB, id int) error {
    _, err := db.Exec("DELETE FROM Menu WHERE ID = ?", id)
    return err
}

// Dish Handlers
func handleGetDishes(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    dishes, err := fetchDishes(db)
    if err != nil {
        http.Error(w, "Failed to fetch dishes: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, `<html><head><title>Dishes List</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
    fmt.Fprintf(w, `<nav><a href="/">Home</a> | <a href="/addDish">Add Dish</a></nav>`)
    fmt.Fprintf(w, `<h1>Dishes</h1><ul>`)
    for _, d := range dishes {
        fmt.Fprintf(w, `<li>ID: %d, Menu ID: %d, Name: %s, Price: %d, Vegan: %t, Shellfish: %t, Nuts: %t</li>`, d.ID, d.MenuID, d.Name, d.Price, d.Vegan, d.Shellfish, d.Nuts)
    }
    fmt.Fprintf(w, `</ul></body></html>`)
}

func handlePostDish(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var d Dish
    d.MenuID, _ = strconv.Atoi(r.FormValue("menuID"))
    d.Name = r.FormValue("name")
    d.Price, _ = strconv.Atoi(r.FormValue("price"))
    d.Vegan = r.FormValue("vegan") == "on"
    d.Shellfish = r.FormValue("shellfish") == "on"
    d.Nuts = r.FormValue("nuts") == "on"
    _, err := addDish(db, d.MenuID, d.Name, d.Price, d.Vegan, d.Shellfish, d.Nuts)
    if err != nil {
        http.Error(w, "Failed to add dish: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, `<html><head><title>Add Dish</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
    fmt.Fprintf(w, `<h1>New dish added</h1></body></html>`)
}

func handleDeleteDish(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    idStr := r.FormValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid dish ID", http.StatusBadRequest)
        return
    }
    if err := deleteDish(db, id); err != nil {
        http.Error(w, "Failed to delete dish: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, `<html><head><title>Delete Dish</title><link rel="stylesheet" type="text/css" href="/static/style.css"></head><body>`)
    fmt.Fprintf(w, `<h1>Dish deleted successfully</h1></body></html>`)
}

func fetchDishes(db *sql.DB) ([]Dish, error) {
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

func addDish(db *sql.DB, menuID int, name string, price int, vegan, shellfish, nuts bool) (int, error) {
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

func deleteDish(db *sql.DB, id int) error {
    _, err := db.Exec("DELETE FROM Dishes WHERE ID = ?", id)
    return err
}

// Define data structures
type Customer struct {
    ID      int    `json:"id"`
    Name    string `json:"name"`
    Address string `json:"address"`
    Email   string `json:"email"`
}

type Order struct {
    ID         int       `json:"id"`
    CustomerID int       `json:"customer_id"`
    Date       time.Time `json:"date"`
}

type Delivery struct {
    ID      int `json:"id"`
    OrderID int `json:"order_id"`
    DTime   int `json:"dtime"`
    DCost   int `json:"dcost"`
}

type Restaurant struct {
    ID      int    `json:"id"`
    Name    string `json:"name"`
    Address string `json:"address"`
    Open    bool   `json:"open"`
    Cuisine string `json:"cuisine"`
}

type Menu struct {
    ID           int `json:"id"`
    RestaurantID int `json:"restaurant_id"`
}

type Dish struct {
    ID          int    `json:"id"`
    MenuID      int    `json:"menu_id"`
    Name        string `json:"name"`
    Price       int    `json:"price"`
    Vegan       bool   `json:"vegan"`
    Shellfish   bool   `json:"shellfish"`
    Nuts        bool   `json:"nuts"`
}
