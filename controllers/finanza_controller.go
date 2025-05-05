package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pdm-backend/repositories"
	"pdm-backend/services"
)

type FinanzaHandler struct {
	FinanceRepo *repositories.FinanzaRepository
}

func NewFinanzaHandler(financeRepo *repositories.FinanzaRepository) *FinanzaHandler {
	return &FinanzaHandler{FinanceRepo: financeRepo}
}

func (h *FinanzaHandler) GetDashboardSummary(c *gin.Context) {

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	inicioMes, finMes, httpCode, jsonResponse, ok := services.ParseMonthAndYear(c)
	if !ok {
		c.JSON(httpCode, jsonResponse)
		return
	}

	resumenFinanciero, err := h.FinanceRepo.GetFinanceSummary(userClaims.FinanzaId, inicioMes, finMes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Hubo un error al traer un resumen de los datos"})
		return
	}

	resumenEgresos, err := h.FinanceRepo.GetEgresoSummary(userClaims.FinanzaId, inicioMes, finMes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Hubo un error al traer un resumen de los egresos"})
		return
	}

	resumenAhorro, err := h.FinanceRepo.GetSavingsSummary(userClaims.FinanzaId, inicioMes, finMes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Hubo un error al traer un resumen de los ahorros"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"finanza_principal": gin.H{
			"resumen_financiero": resumenFinanciero,
			"resumen_egresos":    resumenEgresos,
			"resumen_ahorros":    resumenAhorro,
		},
	})
}

func (h *FinanzaHandler) GetDashboardData(c *gin.Context) {

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	inicioMes, finMes, httpCode, jsonResponse, ok := services.ParseMonthAndYear(c)
	if !ok {
		c.JSON(httpCode, jsonResponse)
		return
	}

	resultado, err := h.FinanceRepo.GetDataSummary(inicioMes, finMes, userClaims.FinanzaId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al conseguir los datos del dashboard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "categorias": resultado})
}

func (h *FinanzaHandler) CreateTransaction(c *gin.Context) {
}
