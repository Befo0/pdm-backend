package controllers

import (
	"net/http"
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
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Hubo un error al conseguir las transacciones"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "transacciones": transacciones})
}

func (h *TransaccionHandler) GetTransactionById(c *gin.Context) {

	idTransaccion, httpCode, jsonResponse := services.ParseUint(c)
	if idTransaccion == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	transaccion, err := h.TransaccionRepo.GetTransactionById(idTransaccion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Hubo un error al conseguir la transaccion"})
		return
	}

	c.JSON(http.StatusOK, transaccion)
}

func (h *TransaccionHandler) CreateTransaction(c *gin.Context) {
}
