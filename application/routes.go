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

	// Middleware for logging
	router.Use(middleware.Logger)

	// Default route
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Model routes
	router.Route("/organisation", a.loadOrganisationRoutes)
	a.router = router
}

func (a *App) loadOrganisationRoutes(router chi.Router) {
	modelHandler := &handler.Organisation{
		Repo: &organisation.MongoRepository{
			Database: a.db,
		},
	}

	router.Post("/", modelHandler.Create)           // Create organisation
	router.Get("/", modelHandler.List)              // List organisations
	router.Get("/{id}", modelHandler.GetByID)       // Get organisation by ID
	router.Put("/{id}", modelHandler.UpdateByID)    // Update organisation by ID
	router.Delete("/{id}", modelHandler.DeleteByID) // Delete organisation by ID
}
