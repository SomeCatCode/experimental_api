package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type App struct {
	router http.Handler
	db     *mongo.Database
	config *Config
}

func New() *App {
	config, err := loadConfig()
	if err != nil {
		panic(fmt.Sprintf("Fehler beim Laden der Konfiguration: %v", err))
	}

	app := &App{
		db:     nil,
		config: config,
	}
	app.loadRoutes()
	return app
}

func (a *App) Start(ctx context.Context) error {
	var err error

	if err = a.ConnectMongo(ctx); err != nil {
		return err
	}
	defer a.db.Client().Disconnect(ctx)

	server := &http.Server{
		Addr:    a.config.Port,
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

func (a *App) ConnectMongo(ctx context.Context) error {
	client, err := mongo.Connect(options.Client().ApplyURI(a.config.MongoUri))
	if err != nil {
		return fmt.Errorf("MongoDB-Verbindung fehlgeschlagen: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("MongoDB-Ping fehlgeschlagen: %w", err)
	}

	a.db = client.Database(a.config.MongoDb)
	return nil
}
