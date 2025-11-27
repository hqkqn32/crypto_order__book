package repository

import (
	"context"
	"crypto-orderbook/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create creates a new order
func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	query := `
		INSERT INTO orders (user_id, order_type, price, amount, status, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query,
		order.UserID,
		order.OrderType,
		order.Price,
		order.Amount,
		order.Status,
	).Scan(&order.ID, &order.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

// GetAll retrieves all active orders
func (r *OrderRepository) GetAll(ctx context.Context) ([]models.Order, error) {
	query := `
		SELECT o.id, o.user_id, u.username, o.order_type, o.price, o.amount, o.status, o.created_at
		FROM orders o
		JOIN users u ON o.user_id = u.id
		WHERE o.status = 'active'
		ORDER BY o.created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Username,
			&order.OrderType,
			&order.Price,
			&order.Amount,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// GetOrderBook retrieves buy and sell orders separately
func (r *OrderRepository) GetOrderBook(ctx context.Context) (*models.OrderBook, error) {
	query := `
		SELECT o.id, o.user_id, u.username, o.order_type, o.price, o.amount, o.status, o.created_at
		FROM orders o
		JOIN users u ON o.user_id = u.id
		WHERE o.status = 'active'
		ORDER BY 
			CASE WHEN o.order_type = 'buy' THEN o.price END DESC,
			CASE WHEN o.order_type = 'sell' THEN o.price END ASC,
			o.created_at ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get order book: %w", err)
	}
	defer rows.Close()

	orderBook := &models.OrderBook{
		BuyOrders:  []models.Order{},
		SellOrders: []models.Order{},
	}

	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Username,
			&order.OrderType,
			&order.Price,
			&order.Amount,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		if order.OrderType == "buy" {
			orderBook.BuyOrders = append(orderBook.BuyOrders, order)
		} else {
			orderBook.SellOrders = append(orderBook.SellOrders, order)
		}
	}

	return orderBook, nil
}

// GetByUserID retrieves all orders for a specific user
func (r *OrderRepository) GetByUserID(ctx context.Context, userID int64) ([]models.Order, error) {
	query := `
		SELECT o.id, o.user_id, u.username, o.order_type, o.price, o.amount, o.status, o.created_at
		FROM orders o
		JOIN users u ON o.user_id = u.id
		WHERE o.user_id = $1
		ORDER BY o.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Username,
			&order.OrderType,
			&order.Price,
			&order.Amount,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// Delete deletes an order (soft delete by updating status)
func (r *OrderRepository) Delete(ctx context.Context, orderID int64, userID int64) error {
	query := `
		UPDATE orders
		SET status = 'cancelled'
		WHERE id = $1 AND user_id = $2 AND status = 'active'
	`

	result, err := r.db.Exec(ctx, query, orderID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("order not found or already cancelled")
	}

	return nil
}
