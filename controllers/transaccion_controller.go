package controllers

import (
	"errors"
	"net/http"
	"pdm-backend/models"
	"pdm-backend/repositories"
	"pdm-backend/services"
	"pdm-backend/websockets"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransaccionHandler struct {
	TransaccionRepo *repositories.TransaccionRepository
}

func NewTransaccionHandler(transaccionRepo *repositories.TransaccionRepository) *TransaccionHandler {
	return &TransaccionHandler{TransaccionRepo: transaccionRepo}
}

func (h *TransaccionHandler) GetTransactions(c *gin.Context) {

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

	mes := int(inicioMes.Month())
	year := inicioMes.Year()

	transacciones, err := h.TransaccionRepo.GetTransactions(inicioMes, finMes, finanzaId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Hubo un error al conseguir las transacciones"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "finanza_id": finanzaId, "mes": mes, "anio": year, "transacciones": transacciones})
}

func (h *TransaccionHandler) GetTransactionOptions(c *gin.Context) {

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

	opcionTransacciones, err := h.TransaccionRepo.GetOptions(finanzaId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Hubo un error al conseguir la transaccion"})
		return
	}

	c.JSON(http.StatusOK, opcionTransacciones)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "No se encontro la transacción"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Hubo un error al conseguir la transaccion"})
		return
	}

	c.JSON(http.StatusOK, transaccion)
}

type TransactionRequest struct {
	TipoTransaccion uint       `json:"tipo_transaccion" binding:"required"`
	TipoMovimiento  uint       `json:"tipo_movimiento" binding:"required"`
	Monto           float64    `json:"monto" binding:"required"`
	Descripcion     string     `json:"descripcion"`
	FechaRegistro   *time.Time `json:"fecha_registro" binding:"required"`
}

func (h *TransaccionHandler) CreateTransaction(c *gin.Context) {

	var finanzaId uint
	var ahorroId uint
	var transaccionRequest TransactionRequest
	var transaccion models.Transacciones

	ahora := time.Now()
	fechaMinima := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	if err := c.ShouldBindJSON(&transaccionRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato de la peticion esta incorrecto"})
		return
	}

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
	ahorroId = userClaims.AhorroId

	if id != 0 {
		finanzaId = id
		ahorroId, err = h.TransaccionRepo.GetSavingSubCategorie(finanzaId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al conseguir el id de la subcategoria"})
			return
		}
	}

	transaccion.FinanzasID = finanzaId
	transaccion.UserID = userClaims.UserId
	transaccion.TipoRegistroID = transaccionRequest.TipoTransaccion
	transaccion.Monto = transaccionRequest.Monto
	transaccion.Descripcion = &transaccionRequest.Descripcion

	fecha := transaccionRequest.FechaRegistro

	if transaccionRequest.FechaRegistro == nil {
		transaccionRequest.FechaRegistro = &ahora
	} else {
		if fecha.After(ahora) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "La fecha de registro no puede ser futura",
			})
			return
		}

		if fecha.Before(fechaMinima) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "La fecha de registro es demasiado antigua",
			})
			return
		}
		fechaAño, fechaMes, _ := fecha.Date()
		nowAño, nowMes, _ := ahora.Date()

		if fechaAño == nowAño && fechaMes == nowMes {
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Solo puedes registrar movimientos del mes actual",
			})
			return
		}
	}

	transaccion.FechaRegistro = *transaccionRequest.FechaRegistro

	switch transaccionRequest.TipoTransaccion {
	case 1:
		transaccion.TipoIngresosID = &transaccionRequest.TipoMovimiento

	case 2:
		transaccion.SubCategoriaEgresoID = &transaccionRequest.TipoMovimiento

		identificadores, err := h.TransaccionRepo.GetIds(*transaccion.SubCategoriaEgresoID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al conseguir los ids"})
			return
		}

		transaccion.CategoriaEgresoID = &identificadores.CategoriaId
		transaccion.TipoPresupuestoID = &identificadores.GastoId

	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Tipo de transaccion invalida"})
		return
	}

	if err := h.TransaccionRepo.CreateTransaction(&transaccion); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al crear la transaccion"})
		return
	}

	if transaccion.SubCategoriaEgresoID != nil && *transaccion.SubCategoriaEgresoID == ahorroId {
		if err := h.TransaccionRepo.CreateOrUpdateSaving(finanzaId, transaccion.Monto, transaccion.FechaRegistro); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al registrar el ahorro mensual"})
			return
		}
	}

	if id != 0 {
		webSocketEvent := h.TransaccionRepo.BuildWebSocketEvent(finanzaId, transaccion.SubCategoriaEgresoID, ahorroId)

		websockets.MensajeBroadcast <- *webSocketEvent
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "La transaccion fue creada correctamente"})
}
