package middleware

import (
	"net/http"
	"strings"

	"github.com/eduardoquea3/finance-go/internal/platform/token"
	"github.com/gin-gonic/gin"
)

func Auth(tokens *token.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		parts := strings.Fields(c.GetHeader("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		claims, err := tokens.Parse(parts[1], "access")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("jwt_subject", claims.Subject)
		c.Next()
	}
}
