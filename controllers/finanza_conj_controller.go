package controllers

import (
	"errors"
	"net/http"
	"pdm-backend/repositories"
	"pdm-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FinanzaConjHandler struct {
	FinanceConjRepo *repositories.FinanzaConjRepository
}

func NewFinanzaConjHandler(financeRepo *repositories.FinanzaConjRepository) *FinanzaConjHandler {
	return &FinanzaConjHandler{FinanceConjRepo: financeRepo}
}

type CreateFinanceRequest struct {
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
}

func (h *FinanzaConjHandler) CreateConjFinance(c *gin.Context) {

	var createRequest CreateFinanceRequest

	if err := c.ShouldBindJSON(&createRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato de la petición esta incorrecto"})
		return
	}

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	err := h.FinanceConjRepo.CreateConjFinance(userClaims.UserId, createRequest.Titulo, createRequest.Descripcion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al crear la finanza"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "La finanza se ha creado correctamente"})
}

type JoinRequest struct {
	Codigo string `json:"codigo"`
}

func (h *FinanzaConjHandler) JoinUser(c *gin.Context) {
	var joinRequest JoinRequest

	if err := c.ShouldBindJSON(&joinRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato de la petición esta incorrecto"})
		return
	}

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	err := h.FinanceConjRepo.JoinUser(userClaims.UserId, joinRequest.Codigo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "El codigo de la invitacion no existe"})
			return
		}

		if err.Error() == "Ya perteneces a esta finanza conjunta" {
			c.JSON(http.StatusConflict, gin.H{"success": false, "message": err.Error()})
			return
		}

		if err.Error() == "El codigo ya ha expirado" {
			c.JSON(http.StatusConflict, gin.H{"success": false, "message": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al unir el usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "El usuario se ha unido con exito"})
}

func (h *FinanzaConjHandler) GetConjFinances(c *gin.Context) {

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	finanzasConjuntas, err := h.FinanceConjRepo.GetConjFinances(userClaims.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al traer las finanzas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "finanzas": finanzasConjuntas})
}

func (h *FinanzaConjHandler) GetConjFinancesDetails(c *gin.Context) {

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	finanzaId, err := services.GetFinanceId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato del query es incorrecto"})
		return
	}

	financeDetails, err := h.FinanceConjRepo.GetConjFinancesDetails(finanzaId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al conseguir los detalles de la finanza"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "detalles_finanza": financeDetails})
}

func (h *FinanzaConjHandler) DeleteUserFromFinance(c *gin.Context) {

	userId, httpCode, jsonResponse := services.ParseUint(c)
	if userId == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	finanzaId, err := services.GetFinanceId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato del query es incorrecto"})
		return
	}

	err = h.FinanceConjRepo.LeaveConjFinance(*userId, finanzaId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al eliminar el usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "El usuario ha sido eliminado correctamente"})
}

func (h *FinanzaConjHandler) LeaveConjFinance(c *gin.Context) {

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	finanzaId, err := services.GetFinanceId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato del query es incorrecto"})
		return
	}

	err = h.FinanceConjRepo.LeaveConjFinance(userClaims.UserId, finanzaId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al eliminar el usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "El usuario ha sido eliminado correctamente"})
}
