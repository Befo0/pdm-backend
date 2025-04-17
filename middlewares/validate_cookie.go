package middlewares

import (
	"net/http"
	"pdm-backend/services"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "token no proporcionado"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Formato de token inválido"})
			c.Abort()
			return
		}

		_, claims, err := services.ValidateJWT(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"succcess": false, "message": "Token invalido o expirado"})
			c.Abort()
			return
		}

		c.Set("claims", claims)

		c.Next()
	}
}
