package services

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func ParseMonthAndYear(c *gin.Context) (inicioMes, finMes time.Time, httpCode int, jsonResponse gin.H, ok bool) {
	mesString := c.Query("mes")
	anioString := c.Query("anio")

	mes, err := strconv.Atoi(mesString)
	if err != nil || mes < 1 || mes > 12 {
		return time.Time{}, time.Time{}, http.StatusBadRequest, gin.H{"error": "Mes inválido"}, false
	}

	anio, err := strconv.Atoi(anioString)
	if err != nil || anio < 1900 {
		return time.Time{}, time.Time{}, http.StatusBadRequest, gin.H{"error": "Año inválido"}, false
	}

	inicioMes = time.Date(anio, time.Month(mes), 1, 0, 0, 0, 0, time.UTC)
	finMes = inicioMes.AddDate(0, 1, 0)

	return inicioMes, finMes, 0, nil, true
}
