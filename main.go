package main

import (
	"context"
	"log"

	"github.com/EthanQuarry/jarbis/internal/app"
)

func main() {
	ctx := context.Background()

	app, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	if err := app.Run(ctx); err != nil {
		log.Fatalf("Error running app: %v", err)
	}
}