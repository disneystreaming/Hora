package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/disneystreaming/Hora/src/controllers"
	"github.com/disneystreaming/Hora/src/models"
)

// AddRoutes adds the defined routes to the given engine
func AddRoutes(s *gin.Engine) {
	// GET
	s.GET("/docs", swaggerRedirect)
	s.GET("/index.html", swaggerRedirect)
	s.GET("/", swaggerRedirect)
	s.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))
	s.GET("/health", health)

	// POST
	s.POST("/validate", validate)
}

// swaggerRedirect will redirect the request to the homepage of swagger docs
func swaggerRedirect(c *gin.Context) {
	c.Redirect(http.StatusFound, "/swagger/index.html")
}

// health will return a health check response
// @Summary Health check
// @Description Simple health check
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse "healthy response"
// @Router /health [get]
func health(c *gin.Context) {
	c.JSON(http.StatusOK, models.HealthResponse{Message: "healthy"})
}

// validate will validate the given candidates against all configured validators
// @Summary Validate
// @Description Validate candidates
// @Tags Validate
// @Accept json
// @Produce json
// @Param payload body []ValidationCandidate true "validation candidates"
// @Success 200 {object} ValidationSummary "validation summary"
// @Failure 400 {object} ErrorResponse "bad request"
// @Failure 500 {object} ErrorResponse "internal server error"
// @Router /validate [post]
func validate(c *gin.Context) {
	rawData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "unable to retrieve payload from request"})
		return
	}

	candidates := []models.ValidationCandidate{}
	err = json.Unmarshal(rawData, &candidates)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "unable to interpret payload"})
		return
	}

	summary, err := controllers.Validate(c, candidates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
