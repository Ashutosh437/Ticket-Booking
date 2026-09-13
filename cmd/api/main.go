package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"ticket-system/internal/config"
	"ticket-system/internal/handler"
	"ticket-system/internal/middleware"
	"ticket-system/internal/repository"
	"ticket-system/internal/service"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Initialize DB
	db, err := repository.InitDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close()

	// 3. Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	ticketRepo := repository.NewTicketRepository(db)

	// 4. Initialize Services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	ticketService := service.NewTicketService(ticketRepo)

	// 5. Initialize Handlers
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authService)
	ticketHandler := handler.NewTicketHandler(ticketService)

	// 6. Setup Chi Router
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	// Public Routes
	r.Get("/health", healthHandler.HealthCheck)
	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	// Protected Routes (JWT required)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.JWTSecret))

		r.Post("/tickets", ticketHandler.CreateTicket)
		r.Get("/tickets", ticketHandler.GetTickets)
		r.Get("/tickets/{id}", ticketHandler.GetTicketByID)
		r.Patch("/tickets/{id}/status", ticketHandler.UpdateTicketStatus)
	})

	serverAddr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server running on port %s...", cfg.Port)
	if err := http.ListenAndServe(serverAddr, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
