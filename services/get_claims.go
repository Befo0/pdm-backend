package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetClaims(c *gin.Context) (*JWTClaims, int, gin.H) {
	claimsInterface, exists := c.Get("claims")

	if !exists {
		return nil, http.StatusUnauthorized, gin.H{"success": false, "message": "No se encontraron claims"}
	}

	userClaims, ok := claimsInterface.(*JWTClaims)
	if !ok {
		return nil, http.StatusInternalServerError, gin.H{"success": false, "message": "Tipo de claim invalido"}
	}

	return userClaims, http.StatusOK, gin.H{}
}
