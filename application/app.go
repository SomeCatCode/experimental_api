package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lan1hotspotgmbh/ms_meter/controller"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	_ "github.com/lan1hotspotgmbh/ms_meter/docs" // Swagger generierter Cod
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	Database *mongo.Database
	Config   Config
	Router   *gin.Engine
}

func New(config Config) *App {
	app := &App{
		Config: config,
	}
	return app
}

func (app *App) Start(ctx context.Context) error {
	// Datenbank laden
	err := app.loadDatabase(ctx)
	if err != nil {
		return fmt.Errorf("Fehler beim Laden der Datenbank: %w", err)
	}
	defer app.Database.Client().Disconnect(ctx)

	//gin.SetMode(gin.ReleaseMode)

	app.Router = gin.Default()
	app.Router.Use(gin.Logger())   // Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release. // By default gin.DefaultWriter = os.Stdout
	app.Router.Use(gin.Recovery()) // Recovery middleware recovers from any panics and writes a 500 if there was one.

	app.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	meterController := controller.NewMeterController(app.Database)
	meterController.RegisterRoutes(app.Router)
	meterValueController := controller.NewMeterValueController(app.Database)
	meterValueController.RegisterRoutes(app.Router)
	healthController := controller.NewHealthController(app.Database)
	healthController.RegisterRoutes(app.Router)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Config.Port),
		Handler: app.Router,
	}
	ch := make(chan error, 1)
	go func() {
		ch <- srv.ListenAndServe()
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(timeout)
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
