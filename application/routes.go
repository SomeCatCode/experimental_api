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

	// routes
	router.Get("/health", a.handleHealthCheck)
	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/", a.handleDefault)
		r.Route("/organisation", a.loadOrganisationRoutes)
	})

	a.router = router
}

func (a *App) handleDefault(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (a *App) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	err := a.db.Client().Ping(r.Context(), nil)
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
			Database: a.db,
		},
	}

	router.Post("/", modelHandler.Create)           // Create organisation
	router.Get("/", modelHandler.List)              // List organisations
	router.Get("/{id}", modelHandler.GetByID)       // Get organisation by ID
	router.Put("/{id}", modelHandler.UpdateByID)    // Update organisation by ID
	router.Delete("/{id}", modelHandler.DeleteByID) // Delete organisation by ID
}
