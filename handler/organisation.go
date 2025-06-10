package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SomeCatCode/experimental_api/model"
	"github.com/SomeCatCode/experimental_api/repository/organisation"
)

type Organisation struct {
	Repo *organisation.MongoRepository
}

func (o *Organisation) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Homepage    string  `json:"website_url"`
		Email       string  `json:"contact_email"`
		Phone       string  `json:"contact_phone"`
		Address     string  `json:"address"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		Style       string  `json:"style"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	now := time.Now().UTC()

	organisation := model.Organisation{
		Name:        body.Name,
		Description: body.Description,
		Homepage:    body.Homepage,

		Email:     body.Email,
		Phone:     body.Phone,
		Address:   body.Address,
		Latetude:  body.Latitude,
		Longitude: body.Longitude,
		Style:     body.Style,

		CreatedAt: &now,
		UpdatedAt: &now,
	}

	err := o.Repo.Insert(r.Context(), organisation)
	if err != nil {
		http.Error(w, "Failed to create organisation", http.StatusInternalServerError)
		return
	}
}

func (o *Organisation) List(w http.ResponseWriter, r *http.Request) {
}

func (o *Organisation) GetByID(w http.ResponseWriter, r *http.Request) {
}

func (o *Organisation) UpdateByID(w http.ResponseWriter, r *http.Request) {
}

func (o *Organisation) DeleteByID(w http.ResponseWriter, r *http.Request) {
}
