package controllers

import (
	"net/http"
	"pdm-backend/repositories"
	"pdm-backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AhorroHandler struct {
	AhorroRepo *repositories.AhorroRepository
}

func NewAhorroHandler(ahorroRepo *repositories.AhorroRepository) *AhorroHandler {
	return &AhorroHandler{AhorroRepo: ahorroRepo}
}

func (h *AhorroHandler) GetSavingsData(c *gin.Context) {

	var finanzaId uint
	anioString := c.Query("anio")

	anio, err := strconv.Atoi(anioString)
	if err != nil || anio < 2000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Año inválido"})
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

	if id != 0 {
		finanzaId = id
	}

	ahorroData, err := h.AhorroRepo.GetSavingsData(finanzaId, anio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al conseguir los datos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "finanza_id": finanzaId, "data": ahorroData})
}

type SavingRequest struct {
	Monto float64 `json:"monto"`
	Mes   int     `json:"mes"`
	Anio  int     `json:"anio"`
}

func (h *AhorroHandler) CreateSavingGoal(c *gin.Context) {

	var finanzaId uint
	var ahorroRequest SavingRequest

	if err := c.ShouldBindJSON(&ahorroRequest); err != nil {
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

	if id != 0 {
		finanzaId = id
	}

	if err := h.AhorroRepo.CreateOrUpdateSavingGoal(finanzaId, ahorroRequest.Monto, ahorroRequest.Mes, ahorroRequest.Anio); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Hubo un error al crear o actualizar la meta mensual"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "La meta fue creada/actualizada correctamente"})
}
