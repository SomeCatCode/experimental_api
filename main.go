package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/SomeCatCode/experimental_api/application"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	app := application.New(application.LoadConfig())
	err := app.Start(ctx)
	if err != nil {
		fmt.Printf("Error starting application: %v\n", err)
		return
	}
}
