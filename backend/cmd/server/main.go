package main

import (
	"context"
	"crypto-orderbook/internal/config"
	"crypto-orderbook/internal/database"
	"crypto-orderbook/internal/handlers"
	"crypto-orderbook/internal/middleware"
	"crypto-orderbook/internal/repository"
	"crypto-orderbook/internal/websocket"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	ws "github.com/gofiber/websocket/v2"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.RunMigrations(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.Pool)
	orderRepo := repository.NewOrderRepository(db.Pool)

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userRepo, cfg)
	orderHandler := handlers.NewOrderHandler(orderRepo, hub)
	wsHandler := handlers.NewWebSocketHandler(hub)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Crypto Orderbook API",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowCredentials: true,
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now(),
		})
	})

	// Auth routes
	api := app.Group("/api")
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected order routes
	orders := api.Group("/orders", middleware.AuthMiddleware(cfg))
	orders.Get("/", orderHandler.GetOrderBook)
	orders.Post("/", orderHandler.CreateOrder)
	orders.Get("/my", orderHandler.GetMyOrders)

	// WebSocket route
	app.Get("/ws", wsHandler.UpgradeMiddleware(), ws.New(wsHandler.HandleWebSocket))

	// Start server
	port := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("🚀 Server starting on port %s", cfg.Server.Port)

	// Graceful shutdown
	go func() {
		if err := app.Listen(port); err != nil {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
