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
	Router   http.Handler
	Database *mongo.Database
	Config   Config
}

func New(config Config) *App {
	app := &App{
		Config: config,
	}
	return app
}

func (app *App) Start(ctx context.Context) error {
	var err error

	// Datenbank laden
	err = app.loadDatabase(ctx)
	if err != nil {
		return fmt.Errorf("Fehler beim Laden der Datenbank: %w", err)
	}
	defer app.Database.Client().Disconnect(ctx)

	// Router initialisieren
	app.loadRoutes()

	// Server initialisieren
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Config.Port),
		Handler: app.Router,
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

func (app *App) loadDatabase(ctx context.Context) error {
	var err error
	var client *mongo.Client

	client, err = mongo.Connect(options.Client().ApplyURI(app.Config.MongoUri))
	if err != nil {
		return fmt.Errorf("MongoDB-Verbindung fehlgeschlagen: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("MongoDB-Ping fehlgeschlagen: %w", err)
	}

	app.Database = client.Database(app.Config.MongoDb)
	return nil
}
