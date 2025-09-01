package main

import (
	"log"
	"mortgage-calculator/internal/app"
	"mortgage-calculator/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// Create and run application
	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Could not create app: %v", err)
	}

	log.Printf("Starting server on port %d", cfg.Port)
	if err := application.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
