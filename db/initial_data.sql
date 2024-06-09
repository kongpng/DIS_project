-- initial_data.sql

-- Insert some initial customers
INSERT INTO Customer (Name, Address, Email) VALUES ('John Doe', '123 Elm Street', 'john.doe@example.com');
INSERT INTO Customer (Name, Address, Email) VALUES ('Jane Smith', '456 Oak Avenue', 'jane.smith@example.com');
INSERT INTO Customer (Name, Address, Email) VALUES ('Alice Johnson', '789 Birch Street', 'alice.johnson@example.com');
INSERT INTO Customer (Name, Address, Email) VALUES ('Bob Brown', '101 Maple Lane', 'bob.brown@example.com');

-- Insert some initial restaurants
INSERT INTO Restaurant (Name, Address, Open, Cuisine) VALUES ('Pizza Palace', '789 Pine Road', 1, 'Italian');
INSERT INTO Restaurant (Name, Address, Open, Cuisine) VALUES ('Sushi World', '101 Maple Lane', 1, 'Japanese');
INSERT INTO Restaurant (Name, Address, Open, Cuisine) VALUES ('Taco Town', '234 Oak Avenue', 1, 'Mexican');
INSERT INTO Restaurant (Name, Address, Open, Cuisine) VALUES ('Burger Barn', '567 Elm Street', 1, 'American');

-- Insert some initial menus
INSERT INTO Menu (RestaurantID) VALUES (1);
INSERT INTO Menu (RestaurantID) VALUES (2);
INSERT INTO Menu (RestaurantID) VALUES (3);
INSERT INTO Menu (RestaurantID) VALUES (4);

-- Insert some initial dishes
INSERT INTO Dishes (MenuID, Name, Price, Vegan, Shellfish, Nuts) VALUES (1, 'Margherita Pizza', 10, 0, 0, 0);
INSERT INTO Dishes (MenuID, Name, Price, Vegan, Shellfish, Nuts) VALUES (1, 'Pepperoni Pizza', 12, 0, 0, 0);
INSERT INTO Dishes (MenuID, Name, Price, Vegan, Shellfish, Nuts) VALUES (2, 'California Roll', 8, 0, 1, 0);
INSERT INTO Dishes (MenuID, Name, Price, Vegan, Shellfish, Nuts) VALUES (2, 'Spicy Tuna Roll', 9, 0, 1, 0);
INSERT INTO Dishes (MenuID, Name, Price, Vegan, Shellfish, Nuts) VALUES (3, 'Beef Taco', 3, 0, 0, 0);
INSERT INTO Dishes (MenuID, Name, Price, Vegan, Shellfish, Nuts) VALUES (3, 'Veggie Taco', 2, 1, 0, 0);
INSERT INTO Dishes (MenuID, Name, Price, Vegan, Shellfish, Nuts) VALUES (4, 'Cheeseburger', 5, 0, 0, 0);
INSERT INTO Dishes (MenuID, Name, Price, Vegan, Shellfish, Nuts) VALUES (4, 'Veggie Burger', 6, 1, 0, 0);

-- Insert some initial orders
INSERT INTO CustomerOrder (CustomerID, Date) VALUES (1, strftime('%s', '2024-09-01'));
INSERT INTO CustomerOrder (CustomerID, Date) VALUES (2, strftime('%s', '2023-03-02'));
INSERT INTO CustomerOrder (CustomerID, Date) VALUES (3, strftime('%s', '2022-06-03'));
INSERT INTO CustomerOrder (CustomerID, Date) VALUES (4, strftime('%s', '2021-02-04'));

-- Insert some initial deliveries
INSERT INTO DeliveryService (OrderID, DeliveryTime, DeliveryCost) VALUES (1, 1200, 5);
INSERT INTO DeliveryService (OrderID, DeliveryTime, DeliveryCost) VALUES (2, 1300, 4);
INSERT INTO DeliveryService (OrderID, DeliveryTime, DeliveryCost) VALUES (3, 1400, 6);
INSERT INTO DeliveryService (OrderID, DeliveryTime, DeliveryCost) VALUES (4, 1500, 3);
