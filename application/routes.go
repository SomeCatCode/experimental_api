package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/SomeCatCode/experimental_api/handler"
	"github.com/SomeCatCode/experimental_api/repository/organisation"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger) // Middleware for logging
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Route("/organisation", a.loadOrganisationRoutes)
	a.router = router
}

func (a *App) loadOrganisationRoutes(router chi.Router) {
	orgHandler := &handler.Organisation{
		Repo: &organisation.MongoRepository{
			Database: a.db,
		},
	}

	router.Post("/", orgHandler.Create)           // Create organisation
	router.Get("/", orgHandler.List)              // List organisations
	router.Get("/{id}", orgHandler.GetByID)       // Get organisation by ID
	router.Put("/{id}", orgHandler.UpdateByID)    // Update organisation by ID
	router.Delete("/{id}", orgHandler.DeleteByID) // Delete organisation by ID
}
