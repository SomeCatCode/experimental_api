package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type App struct {
	router http.Handler

	db *mongo.Client
}

func New() *App {
	// client, err := mongo.Connect(options.Client().ApplyURI(""))
	// if err != nil {
	// 	panic(fmt.Sprintf("failed to connect to MongoDB: %v", err))
	// }

	return &App{
		router: loadRoutes(),
		db:     nil,
	}
}

func (a *App) Start(ctx context.Context) error {
	var err error

	server := &http.Server{
		Addr:    ":8080",
		Handler: a.router,
	}

	// Channel to capture errors from ListenAndServe
	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("error starting server: %w", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(timeout)
	}
}
