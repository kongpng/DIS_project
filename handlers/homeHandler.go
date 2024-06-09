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

}
