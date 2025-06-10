package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/SomeCatCode/experimental_api/handler"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger) // Middleware for logging

	// Routes
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/organisation", loadOrganisationRoutes)

	return router
}

func loadOrganisationRoutes(router chi.Router) {
	orgHandler := &handler.Organisation{}

	router.Post("/", orgHandler.Create)           // Create organisation
	router.Get("/", orgHandler.List)              // List organisations
	router.Get("/{id}", orgHandler.GetByID)       // Get organisation by ID
	router.Put("/{id}", orgHandler.UpdateByID)    // Update organisation by ID
	router.Delete("/{id}", orgHandler.DeleteByID) // Delete organisation by ID
}
