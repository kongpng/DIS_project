package project

import "time"

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

type Menu struct {
	ID           int `json:"id"`
	RestaurantID int `json:"restaurant_id"`
}

type Dish struct {
	ID        int    `json:"id"`
	MenuID    int    `json:"menu_id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Vegan     bool   `json:"vegan"`
	Shellfish bool   `json:"shellfish"`
	Nuts      bool   `json:"nuts"`
}
