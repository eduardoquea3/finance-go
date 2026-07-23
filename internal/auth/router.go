package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(router gin.IRouter, handler *HTTPHandler) {
	auth := router.Group("/api/v1/auth")
	auth.POST("/login", handler.Login)
	auth.POST("/refresh", handler.Refresh)
}
