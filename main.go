package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "io/ioutil"
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
    // Handle root
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Welcome to the API Server")
    })

    // Existing handlers
    http.HandleFunc("/customers", customerHandler(db))
    http.HandleFunc("/orders", orderHandler(db))
    http.HandleFunc("/deliveries", deliveryHandler(db))
    http.HandleFunc("/restaurants", restaurantHandler(db))
    http.HandleFunc("/menus", menuHandler(db))
    http.HandleFunc("/dishes", dishHandler(db))

    // Catch-all handler for undefined routes
    http.HandleFunc("/catch-all", func(w http.ResponseWriter, r *http.Request) {
        http.Error(w, "This route is not defined: "+r.URL.Path, http.StatusNotFound)
    })
}


func customerHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received %s request for customers", r.Method)
        switch r.Method {
        case "GET":
            handleListCustomers(db, w, r)
        case "POST":
            handleAddCustomer(db, w, r)
        case "DELETE":
            handleDeleteCustomer(db, w, r)
        default:
            http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
        }
    }
}
func handleListCustomers(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    customers, err := fetchCustomers(db)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error fetching customers: %v", err), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(customers)
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
func handleAddCustomer(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var c Customer
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Cannot read body", http.StatusBadRequest)
        return
    }
    if err := json.Unmarshal(body, &c); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    id, err := addCustomer(db, c.Name, c.Address, c.Email)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error adding customer: %v", err), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "New customer added with ID: %d", id)
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
func handleDeleteCustomer(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid customer ID", http.StatusBadRequest)
        return
    }
    if err := deleteCustomer(db, id); err != nil {
        http.Error(w, fmt.Sprintf("Error deleting customer: %v", err), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "Customer deleted successfully")
}

func deleteCustomer(db *sql.DB, id int) error {
    _, err := db.Exec("DELETE FROM Customer WHERE ID = ?", id)
    return err
}

func handleGetOrders(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    orders, err := fetchOrders(db)
    if err != nil {
        http.Error(w, "Failed to fetch orders: "+err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(orders)
}

func handleAddOrder(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var o Order
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Cannot read body", http.StatusBadRequest)
        return
    }
    if err := json.Unmarshal(body, &o); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    id, err := addOrder(db, o.CustomerID, o.Date)
    if err != nil {
        http.Error(w, "Failed to add order: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "New order added with ID: %d", id)
}

func handleDeleteOrder(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid order ID", http.StatusBadRequest)
        return
    }
    if err := deleteOrder(db, id); err != nil {
        http.Error(w, "Failed to delete order: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "Order deleted successfully")
}

// Fetch, Add, and Delete functions for Orders
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
func orderHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            handleGetOrders(db, w, r)
        case "POST":
            handleAddOrder(db, w, r)
        case "DELETE":
            handleDeleteOrder(db, w, r)
        default:
            http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
        }
    }
}
func handleGetDeliveries(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    deliveries, err := fetchDeliveries(db)
    if err != nil {
        http.Error(w, "Failed to fetch deliveries: "+err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(deliveries)
}

func handleAddDelivery(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var d Delivery
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Cannot read body", http.StatusBadRequest)
        return
    }
    if err := json.Unmarshal(body, &d); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    id, err := addDelivery(db, d.OrderID, d.DTime, d.DCost)
    if err != nil {
        http.Error(w, "Failed to add delivery: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "New delivery added with ID: %d", id)
}

func handleDeleteDelivery(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid delivery ID", http.StatusBadRequest)
        return
    }
    if err := deleteDelivery(db, id); err != nil {
        http.Error(w, "Failed to delete delivery: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "Delivery deleted successfully")
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

func deliveryHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            handleGetDeliveries(db, w, r)
        case "POST":
            handleAddDelivery(db, w, r)
        case "DELETE":
            handleDeleteDelivery(db, w, r)
        default:
            http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
        }
    }
}

func handleGetRestaurants(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    restaurants, err := fetchRestaurants(db)
    if err != nil {
        http.Error(w, "Failed to fetch restaurants: "+err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(restaurants)
}

func handleAddRestaurant(db *sql.DB, w http.ResponseWriter, req *http.Request) {
    var restaurant Restaurant
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        http.Error(w, "Cannot read body", http.StatusBadRequest)
        return
    }
    if err := json.Unmarshal(body, &restaurant); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    id, err := addRestaurant(db, restaurant.Name, restaurant.Address, restaurant.Open, restaurant.Cuisine)
    if err != nil {
        http.Error(w, "Failed to add restaurant: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "New restaurant added with ID: %d", id)
}


func handleDeleteRestaurant(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
        return
    }
    if err := deleteRestaurant(db, id); err != nil {
        http.Error(w, "Failed to delete restaurant: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "Restaurant deleted successfully")
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

func restaurantHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            handleGetRestaurants(db, w, r)
        case "POST":
            handleAddRestaurant(db, w, r)
        case "DELETE":
            handleDeleteRestaurant(db, w, r)
        default:
            http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
        }
    }
}

func menuHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            handleGetMenus(db, w, r)
        case "POST":
            handleAddMenu(db, w, r)
        case "DELETE":
            handleDeleteMenu(db, w, r)
        default:
            http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
        }
    }
}
// Menu Handlers
func handleGetMenus(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    menus, err := fetchMenus(db)
    if err != nil {
        http.Error(w, "Failed to fetch menus: "+err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(menus)
}

func handleAddMenu(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var m Menu
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Cannot read body", http.StatusBadRequest)
        return
    }
    if err := json.Unmarshal(body, &m); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    id, err := addMenu(db, m.RestaurantID)
    if err != nil {
        http.Error(w, "Failed to add menu: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "New menu added with ID: %d", id)
}

func handleDeleteMenu(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid menu ID", http.StatusBadRequest)
        return
    }
    if err := deleteMenu(db, id); err != nil {
        http.Error(w, "Failed to delete menu: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "Menu deleted successfully")
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
    json.NewEncoder(w).Encode(dishes)
}

func handleAddDish(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var d Dish
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Cannot read body", http.StatusBadRequest)
        return
    }
    if err := json.Unmarshal(body, &d); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    id, err := addDish(db, d.MenuID, d.Name, d.Price, d.Vegan, d.Shellfish, d.Nuts)
    if err != nil {
        http.Error(w, "Failed to add dish: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "New dish added with ID: %d", id)
}

func handleDeleteDish(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid dish ID", http.StatusBadRequest)
        return
    }
    if err := deleteDish(db, id); err != nil {
        http.Error(w, "Failed to delete dish: "+err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "Dish deleted successfully")
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

func dishHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            handleGetDishes(db, w, r)
        case "POST":
            handleAddDish(db, w, r)
        case "DELETE":
            handleDeleteDish(db, w, r)
        default:
            http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
        }
    }
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