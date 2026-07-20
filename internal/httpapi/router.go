package httpapi

import (
	"net/http"

	openapiui "github.com/PeterTakahashi/gin-openapi/openapiui"
	"github.com/eduardoquea3/finance-go/internal/auth"
	"github.com/eduardoquea3/finance-go/internal/httpapi/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, tokens *auth.TokenService) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.GET("/docs/*any", openapiui.WrapHandler(openapiui.Config{
		SpecURL:      "/docs/openapi.json",
		SpecFilePath: "./docs/swagger.json",
		Title:        "Finance Go API",
		Theme:        "light",
	}))

	router.GET("/health", healthHandler(db))

	api := router.Group("/api/v1")
	api.Use(middleware.JWT(tokens))
	api.GET("/me", meHandler)

	return router
}

// healthHandler returns the API and database health status.
// @Summary Health check
// @Description Checks whether the API can reach the database.
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /health [get]
func healthHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.Exec("SELECT 1").Error; err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// meHandler returns the subject from the authenticated JWT.
// @Summary Current user
// @Description Returns the subject from the bearer token.
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/me [get]
func meHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"subject": c.GetString("jwt_subject")})
}
