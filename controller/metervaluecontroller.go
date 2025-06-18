package controller

import (
	"net/http"
	"strconv"

	"github.com/lan1hotspotgmbh/ms_meter/model"
	"github.com/lan1hotspotgmbh/ms_meter/service"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/gin-gonic/gin"
)

type MeterValueController struct {
	service *service.MeterValueService
}

func NewMeterValueController(database *mongo.Database) *MeterValueController {
	service := service.NewMeterValueService(database)
	return &MeterValueController{service: service}
}

func (c *MeterValueController) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/metervalue")
	{
		group.POST("/", c.Create)
		group.GET("/", c.GetAll)
		group.GET("/:id", c.GetByID)
		group.PUT("/:id", c.UpdateByID)
		group.DELETE("/:id", c.DeleteByID)
	}
}

// CreateMeterValue godoc
// @Summary      Neuen Messwert anlegen
// @Description  Legt einen neuen Zählermesswert an
// @Tags         meter-values
// @Accept       json
// @Produce      json
// @Param        value body model.MeterValue true "Messwert-Daten"
// @Success      201 {object} model.MeterValue
// @Failure      400 {object} object
// @Failure      500 {object} object
// @Router       /metervalue/ [post]
func (c *MeterValueController) Create(ctx *gin.Context) {
	var mv model.MeterValue
	if err := ctx.ShouldBindJSON(&mv); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := mv.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.service.Create(ctx.Request.Context(), &mv); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, mv)
}

// GetAllMeterValues godoc
// @Summary      Alle Messwerte abrufen
// @Description  Gibt eine Liste aller Messwerte mit optionalen Filtern und Pagination zurück
// @Tags         meter-values
// @Produce      json
// @Param        filter[MeterId] query string false "Filter nach MeterId"
// @Param        filter[RawData] query string false "Filter nach RawData"
// @Param        filter[search] query string false "Volltextsuche"
// @Param        limit query int false "Maximale Anzahl"
// @Param        start query int false "Offset"
// @Success      200 {array} model.MeterValue
// @Failure      500 {object} object
// @Router       /metervalue/ [get]
func (c *MeterValueController) GetAll(ctx *gin.Context) {
	filters := map[string]string{}
	allowedFields := []string{"MeterId", "RawData"}
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

	values, err := c.service.GetAllWithFilters(ctx.Request.Context(), filters, search, limit, start)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, values)
}

// GetMeterValueByID godoc
// @Summary      Einzelnen Messwert abrufen
// @Description  Gibt einen Messwert anhand seiner ID zurück
// @Tags         meter-values
// @Produce      json
// @Param        id path string true "Messwert ID"
// @Success      200 {object} model.MeterValue
// @Failure      404 {object} object
// @Router       /metervalue/{id} [get]
func (c *MeterValueController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	value, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	ctx.JSON(http.StatusOK, value)
}

// UpdateMeterValueByID godoc
// @Summary      Messwert aktualisieren
// @Description  Aktualisiert einen bestehenden Messwert anhand seiner ID
// @Tags         meter-values
// @Accept       json
// @Produce      json
// @Param        id path string true "Messwert ID"
// @Param        value body model.MeterValue true "Aktualisierte Daten"
// @Success      200 {object} model.MeterValue
// @Failure      400 {object} object
// @Failure      500 {object} object
// @Router       /metervalue/{id} [put]
func (c *MeterValueController) UpdateByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var mv model.MeterValue
	if err := ctx.ShouldBindJSON(&mv); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := mv.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := c.service.UpdateByID(ctx.Request.Context(), id, &mv)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updated)
}

// DeleteMeterValueByID godoc
// @Summary      Messwert löschen
// @Description  Löscht einen Messwert anhand seiner ID
// @Tags         meter-values
// @Param        id path string true "Messwert ID"
// @Success      204 "No Content"
// @Failure      500 {object} object
// @Router       /metervalue/{id} [delete]
func (c *MeterValueController) DeleteByID(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.DeleteByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
