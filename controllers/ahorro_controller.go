package controllers

import (
	"net/http"
	"pdm-backend/models"
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

	ahorroData, err := h.AhorroRepo.GetSavingsData(userClaims.FinanzaId, anio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al conseguir los datos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": ahorroData})
}

type SavingRequest struct {
	Monto float64 `json:"monto"`
	Mes   int     `json:"mes"`
	Anio  int     `json:"anio"`
}

func (h *AhorroHandler) CreateSavingGoal(c *gin.Context) {

	var ahorroRequest SavingRequest
	var ahorro models.AhorroMensual

	if err := c.ShouldBindJSON(&ahorroRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato de la peticion esta incorrecto"})
		return
	}

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	if err := h.AhorroRepo.CreateOrUpdateSavingGoal(userClaims.FinanzaId, ahorroRequest.Monto, ahorroRequest.Mes, ahorro.Anio); err != nil {
	}

}
