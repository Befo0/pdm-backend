package middlewares

import (
	"net/http"
	"pdm-backend/services"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt")

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "No has iniciado sesión"})
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
