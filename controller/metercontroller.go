package controller

import (
	"net/http"
	"strconv"

	"github.com/lan1hotspotgmbh/ms_meter/model"
	"github.com/lan1hotspotgmbh/ms_meter/service"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/gin-gonic/gin"
)

type MeterController struct {
	service *service.MeterService
}

func NewMeterController(database *mongo.Database) *MeterController {
	service := service.NewMeterService(database)
	return &MeterController{service: service}
}

func (c *MeterController) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/meter")
	{
		group.POST("/", c.Create)
		group.GET("/", c.GetAll)
		group.GET("/:id", c.GetByID)
		group.PUT("/:id", c.UpdateByID)
		group.DELETE("/:id", c.DeleteByID)
	}
}

// CreateMeter godoc
// @Summary      Neuen Zähler anlegen
// @Description  Legt einen neuen Stromzähler an
// @Tags         meters
// @Accept       json
// @Produce      json
// @Param        meter body model.Meter true "Meter Daten"
// @Success      201 {object} model.Meter
// @Failure      400 {object} object  "Fehler bei der Validierung der Eingabedaten"
// @Failure      500 {object} object   "ErrorResponse JSON"
// @Router       /meter/ [post]
func (c *MeterController) Create(ctx *gin.Context) {
	var meter model.Meter
	if err := ctx.ShouldBindJSON(&meter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := meter.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.service.Create(ctx.Request.Context(), &meter); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, meter)
}

// GetAllMeters godoc
// @Summary      Alle Stromzähler abrufen
// @Description  Gibt eine Liste aller Stromzähler zurück mit Filter- und Pagingoptionen
// @Tags         meters
// @Produce      json
// @Param        filter[ID] query string false "Filter nach ID"
// @Param        filter[UUID] query string false "Filter nach UUID"
// @Param        filter[Name] query string false "Filter nach Name"
// @Param        filter[Description] query string false "Filter nach Beschreibung"
// @Param        filter[SerialNumber] query string false "Filter nach Seriennummer"
// @Param        filter[Model] query string false "Filter nach Modell"
// @Param        filter[Manufacturer] query string false "Filter nach Hersteller"
// @Param        filter[search] query string false "Volltextsuche"
// @Param        limit query int false "Anzahl der Einträge"
// @Param        start query int false "Startposition"
// @Success      200 {array} model.Meter
// @Failure      500 {object}  object
// @Router       /meter/ [get]
func (c *MeterController) GetAll(ctx *gin.Context) {
	filters := map[string]string{}
	allowedFields := []string{"id", "uuid", "name", "description", "serial_number", "model", "manufacturer"}
	for _, field := range allowedFields {
		val := ctx.Query("filter[" + field + "]")
		if val != "" {
			filters[field] = val
		}
	}
	search := ctx.Query("filter[search]")
	limitStr := ctx.Query("limit")
	startStr := ctx.Query("start")

	limit := 20
	start := 0
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	if s, err := strconv.Atoi(startStr); err == nil && s >= 0 {
		start = s
	}

	meters, err := c.service.GetAllWithFilters(ctx.Request.Context(), filters, search, limit, start)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, meters)
}

// GetMeterByID godoc
// @Summary      Stromzähler nach ID abrufen
// @Description  Gibt einen einzelnen Stromzähler anhand seiner ID zurück
// @Tags         meters
// @Produce      json
// @Param        id path string true "Meter ID"
// @Success      200 {object} model.Meter
// @Failure      404 {object}  object
// @Router       /meter/{id} [get]
func (c *MeterController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	meter, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	ctx.JSON(http.StatusOK, meter)
}

// UpdateMeterByID godoc
// @Summary      Stromzähler aktualisieren
// @Description  Aktualisiert einen bestehenden Stromzähler anhand seiner ID
// @Tags         meters
// @Accept       json
// @Produce      json
// @Param        id path string true "Meter ID"
// @Param        meter body model.Meter true "Aktualisierte Meter-Daten"
// @Success      200 {object} model.Meter
// @Failure      400 {object} object
// @Failure      500 {object} object
// @Router       /meter/{id} [put]
func (c *MeterController) UpdateByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var meter model.Meter
	if err := ctx.ShouldBindJSON(&meter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := meter.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := c.service.UpdateByID(ctx.Request.Context(), id, &meter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updated)
}

// DeleteMeterByID godoc
// @Summary      Stromzähler löschen
// @Description  Löscht einen Stromzähler anhand seiner ID
// @Tags         meters
// @Param        id path string true "Meter ID"
// @Success      204 "No Content"
// @Failure      500 {object}  object
// @Router       /meter/{id} [delete]
func (c *MeterController) DeleteByID(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.DeleteByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
