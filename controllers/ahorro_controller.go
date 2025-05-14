package controllers

import (
	"errors"
	"net/http"
	"pdm-backend/repositories"
	"pdm-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AhorroHandler struct {
	AhorroRepo *repositories.AhorroRepository
}

func NewAhorroHandler(ahorroRepo *repositories.AhorroRepository) *AhorroHandler {
	return &AhorroHandler{AhorroRepo: ahorroRepo}
}

func (h *AhorroHandler) GetSavingsData(c *gin.Context) {

	var metaMensual float64

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	ahorro, err := h.AhorroRepo.GetSaving(userClaims.FinanzaId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "No se encontro la cuenta de ahorro"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Hubo un error al conseguir el ahorro"})
		return
	}

	metaAhorro = ahorro.Monto / 12
}

type AhorroRequest struct {
	MetaAhorro float64 `json:"meta_ahorro"`
}

func (h *AhorroHandler) UpdateSaving(c *gin.Context) {

	var ahorroRequest AhorroRequest

	if err := c.ShouldBindJSON(&ahorroRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato de la petición es incorrecto"})
		return
	}

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	ahorro, err := h.AhorroRepo.GetSaving(userClaims.FinanzaId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "No se encontro la cuenta de ahorro"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Hubo un error al conseguir el ahorro"})
		return
	}

	ahorro.Monto = ahorroRequest.MetaAhorro

	if err := h.AhorroRepo.UpdateSaving(ahorro); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al actualizar el ahorro"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "El ahorro fue actualizado correctamente"})
}
