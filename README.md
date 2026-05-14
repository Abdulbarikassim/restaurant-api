#  Restaurant Management API

A RESTful backend API built with Go (Gin framework) and PostgreSQL for managing restaurant menu items and customer orders.

This API handles menu management, order processing, and order status tracking with persistent database storage.

---

##  Features

### Menu Management
- Create menu items
- Retrieve all menu items
- Retrieve a single menu item
- Update menu items
- Delete menu items

### Order Management
- Create customer orders
- Retrieve all orders
- Retrieve order by ID
- Update order status
- Automatic total calculation
- Item price validation from database

---

## 🛠️ Tech Stack

- **Language:** Go
- **Framework:** Gin
- **Database:** PostgreSQL
- **Environment Management:** godotenv
- **Data Format:** JSON

---

## 📁 Project Structure
```
restaurant-api/
│
├── db/                  # Database connection logic
│   └── postgres.go
│├── models/                  #  models
│    └── menu.go
     └── order.go
│├── handlers/                  # handler for the project logic
│   └── menu_handler.go
│   └── order_handler.go
│├── routes/                  # all routes
│   └── routes.go
│
├── main.go             # Application entry point
├── go.mod
├── go.sum
├── .env
└── README.md
```
## Installation & set up

1.Clone the repository

```bash
git clone https://github.com/Abdulbarikassim/restaurant-api.git

cd restaurant-api
```

2. Instaling dependecies
```
go mody tidy
```

## API Endpoint

For this project, below are the API endpoint:

| Method | Endpoint            | Description         |
| ------ | ------------------- | ------------------- |
| GET    | `/menu`             | Get all menu items  |
| GET    | `/menu/:id`         | Get menu item by ID |
| POST   | `/menu`             | Create menu item    |
| PUT    | `/menu/:id`         | Update menu item    |
| DELETE | `/menu/:id`         | Delete menu item    |
| GET    | `/orders`           | Get all orders      |
| GET    | `/orders/:id`       | Get order by ID     |
| POST   | `/orders`           | Create new order    |
| PUT    | `/order/:id/status` | Update order status |
