package models

import "time"

type Order struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username,omitempty"` // For display purposes
	OrderType string    `json:"order_type"`         // "buy" or "sell"
	Price     float64   `json:"price"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"` // "active", "filled", "cancelled"
	CreatedAt time.Time `json:"created_at"`
}

type CreateOrderRequest struct {
	OrderType string  `json:"order_type" validate:"required,oneof=buy sell"`
	Price     float64 `json:"price" validate:"required,gt=0"`
	Amount    float64 `json:"amount" validate:"required,gt=0"`
}

type OrderBook struct {
	BuyOrders  []Order `json:"buy_orders"`
	SellOrders []Order `json:"sell_orders"`
}
