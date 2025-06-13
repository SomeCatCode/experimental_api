package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/SomeCatCode/experimental_api/handler"
	"github.com/SomeCatCode/experimental_api/repository/organisation"
)

func (app *App) loadRoutes() {
	router := chi.NewRouter()

	// routes
	router.Get("/health", app.handleHealthCheck)
	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/", app.handleDefault)
		r.Route("/organisation", app.loadOrganisationRoutes)
	})

	app.Router = router
}

func (a *App) handleDefault(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (a *App) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	err := a.Database.Client().Ping(r.Context(), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("UNHEALTHY"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (a *App) loadOrganisationRoutes(router chi.Router) {
	modelHandler := &handler.Organisation{
		Repo: &organisation.MongoRepository{
			Collection: "organisations",
			Database:   a.Database,
		},
	}

	router.Post("/", modelHandler.Create)           // Create organisation
	router.Get("/", modelHandler.List)              // List organisations
	router.Get("/{id}", modelHandler.GetByID)       // Get organisation by ID
	router.Put("/{id}", modelHandler.UpdateByID)    // Update organisation by ID
	router.Delete("/{id}", modelHandler.DeleteByID) // Delete organisation by ID
}
