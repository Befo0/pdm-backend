package controllers

import (
	"net/http"
	"pdm-backend/repositories"
	"pdm-backend/services"

	"github.com/gin-gonic/gin"
)

type FinanzaHandler struct {
	FinanceRepo *repositories.FinanzaRepository
}

func NewFinanzaHandler(financeRepo *repositories.FinanzaRepository) *FinanzaHandler {
	return &FinanzaHandler{FinanceRepo: financeRepo}
}

func (h *FinanzaHandler) GetDashboardSummary(c *gin.Context) {

	var finanzaId uint

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	id, err := services.GetFinanceId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato del query es incorrecto"})
		return
	}

	finanzaId = userClaims.FinanzaId

	if id != 0 {
		finanzaId = id
	}

	inicioMes, finMes, httpCode, jsonResponse, ok := services.ParseMonthAndYear(c)
	if !ok {
		c.JSON(httpCode, jsonResponse)
		return
	}

	finanzaPrincipal, err := h.FinanceRepo.GetDashboardSummary(finanzaId, inicioMes, finMes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al traer el resumen del dashboard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":           true,
		"finanza_id":        finanzaId,
		"finanza_principal": finanzaPrincipal,
	})
}

func (h *FinanzaHandler) GetDashboardData(c *gin.Context) {

	var finanzaId uint

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	id, err := services.GetFinanceId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato del query es incorrecto"})
		return
	}

	finanzaId = userClaims.FinanzaId

	if id != 0 {
		finanzaId = id
	}

	inicioMes, finMes, httpCode, jsonResponse, ok := services.ParseMonthAndYear(c)
	if !ok {
		c.JSON(httpCode, jsonResponse)
		return
	}

	resultado, err := h.FinanceRepo.GetDataSummary(inicioMes, finMes, finanzaId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al conseguir los datos del dashboard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"finanza_id": finanzaId,
		"categorias": resultado,
	})
}
