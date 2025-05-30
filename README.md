# Frappuccino: Coffee Shop Management System

## Overview

Frappuccino is a coffee shop management system built in Go with a PostgreSQL database backend. It allows users to manage orders, menu items, inventory, and provides detailed reporting and analytics. This project refactors an existing coffee shop system that previously used JSON-based storage to a scalable PostgreSQL relational database.

With this system, coffee shop managers and staff can efficiently track orders, manage ingredients, update menu items, and generate useful reports to improve business decisions.

## Features

- **Order Management**: Create, retrieve, update, and close orders.
- **Menu Management**: Add, update, and delete menu items.
- **Inventory Management**: Manage inventory levels and track ingredient usage.
- **Reporting**: Generate reports on sales, popular items, and inventory leftovers.
- **Aggregation & Analytics**: Get aggregated data on ordered items, sales, and inventory trends.

## Prerequisites

Before using this project, ensure you have the following installed:

- **Docker**: Required to run the application and PostgreSQL database in containers.
- **Go**: Required for building and running the project if you want to modify or extend it.

If you don’t have Docker and Go installed, please refer to the following installation guides:

- [Docker Installation](https://docs.docker.com/get-docker/)
- [Go Installation](https://go.dev/doc/install)

## Setup and Installation

### 1. Clone the Repository

Clone this repository to your local machine using the following command:

``` bash
git clone https://github.com/your-repo/frappuccino.git
cd frappuccino 
```

 2. Configure Docker

The project uses Docker to run both the application and the PostgreSQL database. All configurations are ready for containerization, so you only need to follow the steps below.

Ensure you are in the root folder of the project.

Run the following command to build and start the containers:


```bash
docker compose up --build
```
This command will:

Build the application and database containers.

Set up the PostgreSQL database and initialize it with the necessary tables via the init.sql file.

Start both the application and database containers.

3. Running the Application
Once the Docker containers are up, the application will be available at:

API: http://localhost:8090


4. Accessing the API
You can interact with the system using the API endpoints described below. You can test them using a tool like Postman or directly through curl.

The project provides a RESTful API to manage orders, menu items, and inventory, and generate various reports.

API Endpoints
Orders API
-POST /orders: Create a new order.

-GET /orders: Retrieve all orders.

-GET /orders/{id}: Retrieve a specific order by ID.

-PUT /orders/{id}: Update an existing order.

-DELETE /orders/{id}: Delete an order.

-POST /orders/{id}/close: Close an order.

Example Request to Create an Order

POST /orders
Content-Type: application/json
```bash
{
  "customer_name": "John Doe",
  "status": "pending",
  "items": [
    { "product_id": 1, "quantity": 2 },
    { "product_id": 2, "quantity": 1 }
  ]
}
```
Menu Items API
-POST /menu: Add a new menu item.

-GET /menu: Retrieve all menu items.

-GET /menu/{id}: Retrieve a specific menu item by ID.

-PUT /menu/{id}: Update a menu item.

DELETE /menu/{id}: Delete a menu item.

Inventory API
-POST /inventory: Add a new inventory item.

-GET /inventory: Retrieve all inventory items.

-GET /inventory/{id}: Retrieve a specific inventory item by ID.

-PUT /inventory/{id}: Update an inventory item.

-DELETE /inventory/{id}: Delete an inventory item.

Reporting and Aggregation Endpoints
-GET /reports/total-sales: Get the total sales amount for the specified period.

-GET /reports/popular-items: Get a list of popular menu items based on sales.


Number of Ordered Items
-GET /orders/numberOfOrderedItems?startDate={startDate}&endDate={endDate}

Returns a list of ordered items and their quantities for a specified time period.

Full Text Search Report
-GET /reports/search?q={query}&filter={filter}&minPrice={minPrice}&maxPrice={maxPrice}

Search through orders, menu items, and customers with partial matching and ranking.

Parameters:

q: Search query string (required).

filter: Comma-separated list of filters: orders, menu, or all (optional).

Get Leftovers
-GET /inventory/getLeftOvers?sortBy={value}&page={page}&pageSize={pageSize}

Returns leftover inventory items, with sorting and pagination.

