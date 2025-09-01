package app

import (
	"context"
	"fmt"
	"mortgage-calculator/internal/cache"
	"mortgage-calculator/internal/calculator"
	"mortgage-calculator/internal/config"
	"mortgage-calculator/internal/controller"
	"mortgage-calculator/internal/middleware"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type App struct {
	server *http.Server
}

func NewApp(cfg *config.Config) (*App, error) {
	// Initialize dependencies
	calc := calculator.NewCalculator()
	cache := cache.NewInMemoryCache()
	controller := controller.NewMortgageController(calc, cache)

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logging) // Add logging middleware
	controller.RegisterRoutes(r)

	// Create server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{server: server}, nil
}

func (a *App) Run() error {
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
