package application

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type App struct {
	router http.Handler

	db         *mongo.Database
	collection *mongo.Collection
}

func New() *App {
	app := &App{
		db:         nil,
		collection: nil,
	}
	app.loadRoutes()
	return app
}

func (a *App) ConnectMongo(ctx context.Context) error {
	// load environment file
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("Fehler beim Laden der Umgebungsvariablen: %w", err)
	}

	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")
	if uri == "" || dbName == "" {
		return fmt.Errorf("MONGO_URI oder MONGO_DB nicht gesetzt")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("MongoDB-Verbindung fehlgeschlagen: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("MongoDB-Ping fehlgeschlagen: %w", err)
	}

	a.db = client.Database(dbName)
	return nil
}

func (a *App) Start(ctx context.Context) error {
	var err error

	if err = a.ConnectMongo(ctx); err != nil {
		return err
	}
	defer a.db.Client().Disconnect(ctx)

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
