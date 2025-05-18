package controllers

import (
	"net/http"
	"pdm-backend/repositories"
	"pdm-backend/services"
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

func (h *FinanzaHandler) GetDashboardSummary(c *gin.Context) {

	var resumenFinanciero, resumenEgresos, resumenAhorro gin.H
	errCh := make(chan error, 3)

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

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

	go func() {
		resumen, err := h.FinanceRepo.GetFinanceSummary(userClaims.FinanzaId, inicioMes, finMes)
		if err == nil {
			resumenFinanciero = resumen
		}
		errCh <- err

	}()

	go func() {
		resumen, err := h.FinanceRepo.GetEgresoSummary(userClaims.FinanzaId, inicioMes, finMes)
		if err == nil {
			resumenEgresos = resumen
		}
		errCh <- err
	}()

	go func() {
		resumen, err := h.FinanceRepo.GetSavingsSummary(userClaims.FinanzaId, mes, anio)
		if err == nil {
			resumenAhorro = resumen
		}
		errCh <- err
	}()

	for i := 0; i < 3; i++ {
		if err := <-errCh; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Hubo un error al conseguir los datos del dashboard"})
			return
		}
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
