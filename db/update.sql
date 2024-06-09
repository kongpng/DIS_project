-- update.sql

-- Drop existing tables if they exist
DROP TABLE IF EXISTS Customer;
DROP TABLE IF EXISTS CustomerOrder;
DROP TABLE IF EXISTS DeliveryService;
DROP TABLE IF EXISTS Restaurant;
DROP TABLE IF EXISTS Menu;
DROP TABLE IF EXISTS Dishes;

-- Table for Customer
CREATE TABLE Customer (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT NOT NULL,
    Address TEXT NOT NULL,
    Email TEXT NOT NULL
);

-- Table for CustomerOrder
CREATE TABLE CustomerOrder (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    CustomerID INTEGER NOT NULL,
    Date INTEGER NOT NULL,
    FOREIGN KEY (CustomerID) REFERENCES Customer(ID)
);

-- Table for DeliveryService
CREATE TABLE DeliveryService (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    OrderID INTEGER NOT NULL,
    DeliveryTime INTEGER NOT NULL,
    DeliveryCost INTEGER NOT NULL,
    FOREIGN KEY (OrderID) REFERENCES CustomerOrder(ID)
);

-- Table for Restaurant
CREATE TABLE Restaurant (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT NOT NULL,
    Address TEXT NOT NULL,
    Open BOOLEAN NOT NULL,
    Cuisine TEXT NOT NULL
);

-- Table for Menu
CREATE TABLE Menu (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    RestaurantID INTEGER NOT NULL,
    FOREIGN KEY (RestaurantID) REFERENCES Restaurant(ID)
);

-- Table for Dishes
CREATE TABLE Dishes (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    MenuID INTEGER NOT NULL,
    Name TEXT NOT NULL,
    Price INTEGER NOT NULL,
    Vegan BOOLEAN NOT NULL,
    Shellfish BOOLEAN NOT NULL,
    Nuts BOOLEAN NOT NULL,
    FOREIGN KEY (MenuID) REFERENCES Menu(ID)
);
