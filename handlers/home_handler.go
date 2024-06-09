package handlers

import (
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
	<html>
	<head>
		<title>Welcome</title>
		<link rel='stylesheet' type='text/css' href='/static/style.css'>
	</head>
	<body>
		<h1>Dinner Dashers</h1>
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
		<h2>Search Customers</h2>
		<form action="/searchCustomers" method="get">
			<label>Name: <input type="text" name="name" /></label><br/>
			<label>Address: <input type="text" name="address" /></label><br/>
			<label>Email: <input type="email" name="email" /></label><br/>
			<input type="submit" value="Search" />
		</form>
		<h2>Search Restaurants</h2>
		<form action="/searchRestaurants" method="get">
			<label>Name: <input type="text" name="name" /></label><br/>
			<label>Address: <input type="text" name="address" /></label><br/>
			<label>Open: <input type="checkbox" name="open" /></label><br/>
			<label>Cuisine: <input type="text" name="cuisine" /></label><br/>
			<input type="submit" value="Search" />
		</form>
		<h2>Search Menus</h2>
		<form action="/searchMenus" method="get">
			<label>Restaurant ID: <input type="number" name="restaurantID" /></label><br/>
			<input type="submit" value="Search" />
		</form>
		<h2>Search Dishes</h2>
		<form action="/searchDishes" method="get">
			<label>Menu ID: <input type="number" name="menuID" /></label><br/>
			<label>Name: <input type="text" name="name" /></label><br/>
			<label>Price: <input type="number" name="price" /></label><br/>
			<label>Vegan: <input type="checkbox" name="vegan" /></label><br/>
			<label>Shellfish: <input type="checkbox" name="shellfish" /></label><br/>
			<label>Nuts: <input type="checkbox" name="nuts" /></label><br/>
			<input type="submit" value="Search" />
		</form>
		<h2>Search Orders</h2>
		<form action="/searchOrders" method="get">
			<label>Customer ID: <input type="number" name="customerID" /></label><br/>
			<label>Date (YYYY-MM-DD): <input type="date" name="date" /></label><br/>
			<input type="submit" value="Search" />
		</form>
		<h2>Search Deliveries</h2>
		<form action="/searchDeliveries" method="get">
			<label>Order ID: <input type="number" name="orderID" /></label><br/>
			<label>Delivery Time (HHMM): <input type="text" name="deliveryTime" /></label><br/>
			<label>Delivery Cost: <input type="number" name="deliveryCost" /></label><br/>
			<input type="submit" value="Search" />
		</form>
	</body>
	</html>
	`)
}
