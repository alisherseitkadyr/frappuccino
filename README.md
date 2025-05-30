Frappuccino: Coffee Shop Management System
Overview
Frappuccino is a coffee shop management system built in Go with a PostgreSQL database backend. It allows users to manage orders, menu items, inventory, and provides detailed reporting and analytics. This project refactors an existing coffee shop system that previously used JSON-based storage to a scalable PostgreSQL relational database.

With this system, coffee shop managers and staff can efficiently track orders, manage ingredients, update menu items, and generate useful reports to improve business decisions.

Features
Order Management: Create, retrieve, update, and close orders.

Menu Management: Add, update, and delete menu items.

Inventory Management: Manage inventory levels and track ingredient usage.

Reporting: Generate reports on sales, popular items, and inventory leftovers.

Aggregation & Analytics: Get aggregated data on ordered items, sales, and inventory trends.

Prerequisites
Before using this project, ensure you have the following installed:

Docker: Required to run the application and PostgreSQL database in containers.

Go: Required for building and running the project if you want to modify or extend it.

If you don’t have Docker and Go installed, please refer to the following installation guides:

Setup and Installation
1. Clone the Repository
Clone this repository to your local machine using the following command:

`git clone https://github.com/your-repo/frappuccino.git`
`cd frappuccino`

2. Configure Docker
The project uses Docker to run both the application and the PostgreSQL database. All configurations are ready for containerization, so you only need to follow the steps below.

Ensure you are in the root folder of the project.

Run the following command to build and start the containers:

`docker compose up --build`
This command will:

Build the application and database containers.

Set up the PostgreSQL database and initialize it with the necessary tables via the init.sql file.

Start both the application and database containers.

3. Running the Application
Once the Docker containers are up, the application will be available at:

API: http://localhost:8090

Database: Accessible via Docker container db, using the following connection credentials:

Host: db

Port: 5432

User: latte

Password: latte

Database: frappuccino

4. Accessing the API
You can interact with the system using the API endpoints described below. You can test them using a tool like Postman or directly through curl.

The project provides a RESTful API to manage orders, menu items, and inventory, and generate various reports.

API Endpoints
Orders API
POST /orders: Create a new order.

GET /orders: Retrieve all orders.

GET /orders/{id}: Retrieve a specific order by ID.

PUT /orders/{id}: Update an existing order.

DELETE /orders/{id}: Delete an order.

POST /orders/{id}/close: Close an order.

Example Request to Create an Order
bash
Copy
Edit
POST /orders
Content-Type: application/json

{
  "customer_name": "John Doe",
  "status": "pending",
  "items": [
    { "product_id": 1, "quantity": 2 },
    { "product_id": 2, "quantity": 1 }
  ]
}
Menu Items API
POST /menu: Add a new menu item.

GET /menu: Retrieve all menu items.

GET /menu/{id}: Retrieve a specific menu item by ID.

PUT /menu/{id}: Update a menu item.

DELETE /menu/{id}: Delete a menu item.

Inventory API
POST /inventory: Add a new inventory item.

GET /inventory: Retrieve all inventory items.

GET /inventory/{id}: Retrieve a specific inventory item by ID.

PUT /inventory/{id}: Update an inventory item.

DELETE /inventory/{id}: Delete an inventory item.

Reporting and Aggregation Endpoints
GET /reports/total-sales: Get the total sales amount for the specified period.

GET /reports/popular-items: Get a list of popular menu items based on sales.

New Aggregation Endpoints
1. Number of Ordered Items
GET /orders/numberOfOrderedItems?startDate={startDate}&endDate={endDate}

Returns a list of ordered items and their quantities for a specified time period.

Parameters:

startDate: Start date in YYYY-MM-DD format (optional).

endDate: End date in YYYY-MM-DD format (optional).

Response Example:

json
Copy
Edit
{
  "latte": 109,
  "muffin": 56,
  "espresso": 120,
  "raff": 0
}
2. Full Text Search Report
GET /reports/search?q={query}&filter={filter}&minPrice={minPrice}&maxPrice={maxPrice}

Search through orders, menu items, and customers with partial matching and ranking.

Parameters:

q: Search query string (required).

filter: Comma-separated list of filters: orders, menu, or all (optional).

minPrice: Minimum price (optional).

maxPrice: Maximum price (optional).

Response Example:

json
Copy
Edit
{
  "menu_items": [
    {
      "id": "12",
      "name": "Double Chocolate Cake",
      "description": "Rich chocolate layer cake",
      "price": 15.99,
      "relevance": 0.89
    },
    {
      "id": "15",
      "name": "Chocolate Cheesecake",
      "description": "Creamy cheesecake with chocolate",
      "price": 12.99,
      "relevance": 0.75
    }
  ],
  "orders": [
    {
      "id": "1234",
      "customer_name": "Alice Brown",
      "items": ["Chocolate Cake", "Coffee"],
      "total": 18.99,
      "relevance": 0.68
    }
  ],
  "total_matches": 3
}
3. Ordered Items by Period
GET /reports/orderedItemsByPeriod?period={day|month}&month={month}

Get the number of orders for the specified period, grouped by day or month.

Response Example:

json
Copy
Edit
{
  "period": "day",
  "month": "october",
  "orderedItems": [
    { "1": 109 },
    { "2": 234 },
    ...
  ]
}
4. Get Leftovers
GET /inventory/getLeftOvers?sortBy={value}&page={page}&pageSize={pageSize}

Returns leftover inventory items, with sorting and pagination.

Response Example:

json
Copy
Edit
{
  "currentPage": 1,
  "hasNextPage": true,
  "pageSize": 4,
  "totalPages": 10,
  "data": [
    {
      "name": "croissant",
      "quantity": 109,
      "price": 950
    },
    {
      "name": "sugar",
      "quantity": 93,
      "price": 50
    }
  ]
}