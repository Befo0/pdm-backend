package controllers

import (
	"pdm-backend/repositories"
	"pdm-backend/services"

	"github.com/gin-gonic/gin"
)

type TransaccionHandler struct {
	TransaccionRepo *repositories.TransaccionRepository
}

func NewTransaccionHandler(transaccionRepo *repositories.TransaccionRepository) *TransaccionHandler {
	return &TransaccionHandler{TransaccionRepo: transaccionRepo}
}

func (h *TransaccionHandler) GetTransactions(c *gin.Context) {
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

	transacciones, err := h.TransaccionRepo.GetTransactions(inicioMes, finMes, userClaims.FinanzaId)
}
