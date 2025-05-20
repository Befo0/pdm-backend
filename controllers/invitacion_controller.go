package controllers

import (
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

	err := h.InvitacionRepo.CreateInvite(finanzaId)
	if err != nil {
	}
}
