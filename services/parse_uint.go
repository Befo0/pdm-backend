package services

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseUint(c *gin.Context) (*uint, int, gin.H) {
	idParam := c.Param("id")

	idUint, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return nil, http.StatusBadRequest, gin.H{"success": false, "message": "El id no es un numero valido"}
	}
	idCategoria := uint(idUint)

	return &idCategoria, 0, nil
}
