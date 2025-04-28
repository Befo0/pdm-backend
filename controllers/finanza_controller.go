package controllers

import (
	"net/http"
	"pdm-backend/repositories"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FinanzaHandler struct {
	FinanceRepo *repositories.FinanzaRepository
}

func NewFinanzaHandler(financeRepo *repositories.FinanzaRepository) *FinanzaHandler {
	return &FinanzaHandler{FinanceRepo: financeRepo}
}

func (h *FinanzaHandler) GetDashboard(c *gin.Context) {
	mesString := c.Query("mes")
	anioString := c.Query("anio")

	mes, err := strconv.Atoi(mesString)
	if err != nil || mes < 1 || mes > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mes inválido"})
		return
	}

	anio, err := strconv.Atoi(anioString)
	if err != nil || anio < 1900 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Año inválido"})
		return
	}

	inicioMes := time.Date(anio, time.Month(mes), 1, 0, 0, 0, 0, time.UTC)
	finMes := inicioMes.AddDate(0, 1, 0)

	resumen, err := h.FinanceRepo.GetSummary(inicioMes, finMes)

	c.JSON(http.StatusOK, gin.H{
		"finanza_principal": gin.H{
			"nombre": "Finanza principal",
		},
	})
}
