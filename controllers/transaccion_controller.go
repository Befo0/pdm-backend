package controllers

import (
	"net/http"
	"pdm-backend/models"
	"pdm-backend/repositories"
	"pdm-backend/services"
	"time"

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

type TransactionRequest struct {
	TipoTransaccion uint    `json:"tipo_transaccion" binding:"required"`
	TipoMovimiento  uint    `json:"tipo_movimiento" binding:"required"`
	TipoCategoria   *uint   `json:"tipo_categoria,omitempty"`
	TipoGasto       *uint   `json:"tipo_gasto,omitempty"`
	Monto           float64 `json:"monto" binding:"required"`
	Descripcion     string  `json:"descripcion"`
}

func (h *TransaccionHandler) CreateTransaction(c *gin.Context) {

	var transaccionRequest TransactionRequest
	var transaccion models.Transacciones

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	if err := c.ShouldBindJSON(&transaccionRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato de la peticion esta incorrecto"})
		return
	}

	transaccion.FinanzasID = userClaims.FinanzaId
	transaccion.UserID = userClaims.UserId
	transaccion.TipoRegistroID = transaccionRequest.TipoTransaccion
	transaccion.Monto = transaccionRequest.Monto
	transaccion.Descripcion = &transaccionRequest.Descripcion
	transaccion.FechaRegistro = time.Now()

	switch transaccionRequest.TipoTransaccion {
	case 1:
		transaccion.TipoIngresosID = &transaccionRequest.TipoMovimiento

	case 2:
		transaccion.SubCategoriaEgresoID = &transaccionRequest.TipoMovimiento

		if transaccionRequest.TipoCategoria == nil || transaccionRequest.TipoGasto == nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Faltan datos obligatorios"})
			return
		}

		transaccion.CategoriaEgresoID = transaccionRequest.TipoCategoria
		transaccion.TipoPresupuestoID = transaccionRequest.TipoGasto

	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Tipo de transaccion invalida"})
		return
	}

	if err := h.TransaccionRepo.CreateTransaction(&transaccion); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al crear la transaccion"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "La transaccion fue creada corractamente"})
}
