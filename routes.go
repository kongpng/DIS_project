package main

import (
	"database/sql"
	"net/http"
	"project/handlers"
)

func setupRoutes(db *sql.DB) {
	// Home page
	http.HandleFunc("/", handlers.HomeHandler)

	// Customer handlers
	http.HandleFunc("/customers", handlers.CustomersHandler(db))
	http.HandleFunc("/addCustomer", handlers.AddCustomerHandler)
	http.HandleFunc("/deleteCustomer", handlers.DeleteCustomerHandler)

	// Order handlers
	http.HandleFunc("/orders", handlers.OrdersHandler(db))
	http.HandleFunc("/addOrder", handlers.AddOrderHandler)
	http.HandleFunc("/deleteOrder", handlers.DeleteOrderHandler)

	// Delivery handlers
	http.HandleFunc("/deliveries", handlers.DeliveryHandler(db))
	http.HandleFunc("/addDelivery", handlers.AddDeliveryHandler(db))
	http.HandleFunc("/deleteDelivery", handlers.DeleteDeliveryHandler)

	// Restaurant handlers
	http.HandleFunc("/restaurants", handlers.RestaurantsHandler(db))
	http.HandleFunc("/addRestaurant", handlers.AddRestaurantHandler)
	http.HandleFunc("/deleteRestaurant", handlers.DeleteRestaurantHandler)

	// Menu handlers
	http.HandleFunc("/menus", handlers.MenusHandler(db))
	http.HandleFunc("/addMenu", handlers.AddMenuHandler)
	http.HandleFunc("/deleteMenu", handlers.DeleteMenuHandler)

	// Dish handlers
	http.HandleFunc("/dishes", handlers.DishesHandler(db))
	http.HandleFunc("/addDish", handlers.AddDishHandler)
	http.HandleFunc("/deleteDish", handlers.DeleteDishHandler)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Catch-all handler for undefined routes
	//	http.HandleFunc("/catch-all", handlers.CatchAllHandler)
}
