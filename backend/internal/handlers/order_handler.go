package handlers

import (
	"crypto-orderbook/internal/models"
	"crypto-orderbook/internal/repository"
	"crypto-orderbook/internal/websocket"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	orderRepo *repository.OrderRepository
	hub       *websocket.Hub
}

func NewOrderHandler(orderRepo *repository.OrderRepository, hub *websocket.Hub) *OrderHandler {
	return &OrderHandler{
		orderRepo: orderRepo,
		hub:       hub,
	}
}

func (h *OrderHandler) GetOrderBook(c *fiber.Ctx) error {
	orderBook, err := h.orderRepo.GetOrderBook(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get order book"})
	}

	return c.JSON(orderBook)
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int64)
	username := c.Locals("username").(string)

	var req models.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.OrderType != "buy" && req.OrderType != "sell" {
		return c.Status(400).JSON(fiber.Map{"error": "Order type must be 'buy' or 'sell'"})
	}

	if req.Price <= 0 || req.Amount <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Price and amount must be greater than 0"})
	}

	order := &models.Order{
		UserID:    userID,
		Username:  username,
		OrderType: req.OrderType,
		Price:     req.Price,
		Amount:    req.Amount,
		Status:    "active",
	}

	if err := h.orderRepo.Create(c.Context(), order); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create order"})
	}

	// Broadcast to all WebSocket clients
	h.hub.BroadcastOrder(order)

	return c.Status(201).JSON(order)
}

func (h *OrderHandler) GetMyOrders(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int64)

	orders, err := h.orderRepo.GetByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get orders"})
	}

	return c.JSON(orders)
}
