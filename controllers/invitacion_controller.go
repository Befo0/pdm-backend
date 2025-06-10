package controllers

import (
	"net/http"
	"pdm-backend/repositories"
	"pdm-backend/services"

	"github.com/gin-gonic/gin"
)

type InvitacionHandler struct {
	InvitacionRepo *repositories.InvitacionRepository
}

func NewInvitacionHandler(invitacionRepo *repositories.InvitacionRepository) *InvitacionHandler {
	return &InvitacionHandler{InvitacionRepo: invitacionRepo}
}

func (h *InvitacionHandler) CreateInvite(c *gin.Context) {

	finanzaId, httpCode, jsonResponse := services.ParseUint(c)
	if finanzaId == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	invitacion, err := h.InvitacionRepo.CreateInvite(finanzaId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al crear la invitación"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "codigo_invitacion": invitacion.Codigo})
}
