package controller

import (
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
	database *mongo.Database
}

func NewHealthController(database *mongo.Database) *HealthController {
	return &HealthController{database: database}
}

func (c *HealthController) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", c.GetHealth)
}

// GetHealth godoc
// @Summary      Health Check
// @Description  Überprüft den Gesundheitszustand des Dienstes
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200 {object} object "Health Check Response"
// @Failure      500 {object} object "ErrorResponse JSON"
// @Router       /health [get]
func (c *HealthController) GetHealth(ctx *gin.Context) {
	// Check DB Connection
	if err := c.database.Client().Ping(ctx.Request.Context(), nil); err != nil {
		ctx.JSON(500, gin.H{
			"status":  "error",
			"message": "Database connection failed",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "ok",
		"message": "Service is running",
	})
}
