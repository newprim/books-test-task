package main

import (
	"log"

	"github.com/newprim/books-test-task/config"
	"github.com/newprim/books-test-task/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
