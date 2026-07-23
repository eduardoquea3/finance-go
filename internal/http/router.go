package http

import (
	"net/http"

	openapiui "github.com/PeterTakahashi/gin-openapi/openapiui"
	"github.com/eduardoquea3/finance-go/internal/auth"
	"github.com/eduardoquea3/finance-go/internal/http/middleware"
	"github.com/eduardoquea3/finance-go/internal/platform/token"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, authHandler *auth.HTTPHandler, tokens *token.Service) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), middleware.Recovery(), middleware.RequestID())
	router.GET("/docs/*any", openapiui.WrapHandler(openapiui.Config{
		SpecURL:      "/docs/openapi.json",
		SpecFilePath: "./docs/swagger.json",
		Title:        "Finance Go API",
		Theme:        "light",
	}))
	router.GET("/health", healthHandler(db))
	auth.RegisterRoutes(router, authHandler)

	api := router.Group("/api/v1")
	api.Use(middleware.Auth(tokens))
	api.GET("/me", meHandler)
	return router
}

func healthHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.Exec("SELECT 1").Error; err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func meHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"subject": c.GetString("jwt_subject")})
}
