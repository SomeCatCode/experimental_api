package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SomeCatCode/experimental_api/model"
	"github.com/SomeCatCode/experimental_api/repository/organisation"
	"github.com/go-chi/chi/v5"
)

type Organisation struct {
	Repo *organisation.MongoRepository
}

func (o *Organisation) Create(w http.ResponseWriter, r *http.Request) {
	var obj model.Organisation

	// Decode the request body into the Organisation object
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set timestamps
	now := time.Now().UTC()
	obj.CreatedAt = &now
	obj.UpdatedAt = &now

	// Write the organisation ID if it is not set
	err = o.Repo.Insert(r.Context(), obj)
	if err != nil {
		http.Error(w, "Failed to create organisation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(obj)
}

func (o *Organisation) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	organisation, err := o.Repo.FindByID(r.Context(), idParam)
	if err != nil {
		http.Error(w, "Failed to retrieve organisation", http.StatusInternalServerError)
		return
	}
	if organisation == nil {
		http.Error(w, "Organisation not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(organisation)
}

func (o *Organisation) List(w http.ResponseWriter, r *http.Request) {
}

func (o *Organisation) UpdateByID(w http.ResponseWriter, r *http.Request) {
}

func (o *Organisation) DeleteByID(w http.ResponseWriter, r *http.Request) {
}
