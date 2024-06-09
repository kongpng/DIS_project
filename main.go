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
    http.HandleFunc("/customers", customerHandler(db))
    http.HandleFunc("/orders", orderHandler(db))
    http.HandleFunc("/deliveries", deliveryHandler(db))
    http.HandleFunc("/restaurants", restaurantHandler(db))
    http.HandleFunc("/menus", menuHandler(db))
    http.HandleFunc("/dishes", dishHandler(db))
}

func customerHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
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

// Define the remaining necessary handlers as shown in the previous messages...
// Example of how to handle DELETE requests for an Order
func handleDeleteOrder(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid order ID", http.StatusBadRequest)
        return
    }
    err = deleteOrder(db, id)
    if err != nil {
        http.Error(w, "Failed to delete order", http.StatusInternalServerError)
        return
    }
    fmt.Fprintln(w, "Order deleted successfully")
}

// Example of how to handle GET requests for Deliveries
func handleGetDeliveries(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT ID, OrderID, DTime, DCost FROM DeliveryService")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var deliveries []Delivery
    for rows.Next() {
        var delivery Delivery
        if err := rows.Scan(&delivery.ID, &delivery.OrderID, &delivery.DTime, &delivery.DCost); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        deliveries = append(deliveries, delivery)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(deliveries)
}

// Complete the remaining handlers
func handleDeleteDelivery(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid delivery ID", http.StatusBadRequest)
        return
    }
    err = deleteDelivery(db, id)
    if err != nil {
        http.Error(w, "Failed to delete delivery", http.StatusInternalServerError)
        return
    }
    fmt.Fprintln(w, "Delivery deleted successfully")
}

func handleDeleteRestaurant(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
        return
    }
    err = deleteRestaurant(db, id)
    if err != nil {
        http.Error(w, "Failed to delete restaurant", http.StatusInternalServerError)
        return
    }
    fmt.Fprintln(w, "Restaurant deleted successfully")
}

func handleDeleteMenu(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid menu ID", http.StatusBadRequest)
        return
    }
    err = deleteMenu(db, id)
    if err != nil {
        http.Error(w, "Failed to delete menu", http.StatusInternalServerError)
        return
    }
    fmt.Fprintln(w, "Menu deleted successfully")
}

func handleDeleteDish(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid dish ID", http.StatusBadRequest)
        return
    }
    err = deleteDish(db, id)
    if err != nil {
        http.Error(w, "Failed to delete dish", http.StatusInternalServerError)
        return
    }
    fmt.Fprintln(w, "Dish deleted successfully")
}

// Assume additional GET, POST, and other CRUD operations are defined similarly
