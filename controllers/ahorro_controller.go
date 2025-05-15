package controllers

import (
	"net/http"
	"pdm-backend/models"
	"pdm-backend/repositories"
	"pdm-backend/services"

	"github.com/gin-gonic/gin"
)

type AhorroHandler struct {
	AhorroRepo *repositories.AhorroRepository
}

func NewAhorroHandler(ahorroRepo *repositories.AhorroRepository) *AhorroHandler {
	return &AhorroHandler{AhorroRepo: ahorroRepo}
}

func (h *AhorroHandler) GetSavingsData(c *gin.Context) {

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

type AhorroRequest struct {
	MetaAhorro float64 `json:"meta_ahorro"`
}

func (h *AhorroHandler) UpdateSaving(c *gin.Context) {

}
