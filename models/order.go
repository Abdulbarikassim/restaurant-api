package models



// order Item

type OrderItem struct {
	MenuItemId int     `json:"menu_item_id"`
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
}

// All the order items

type Order struct {
	ID            int      `json:"id"`
	CustomerPhone string      `json:"customer_phone"`
	Status        string      `json:"status"`
	TotalAmount   float64     `json:"total_amount"`
	Items         []OrderItem `json:"items"`
}
