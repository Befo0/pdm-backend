package controllers

import (
	"errors"
	"net/http"
	"pdm-backend/models"
	"pdm-backend/repositories"
	"pdm-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IngresosHandler struct {
	IngresosRepo *repositories.IngresosRepository
}

func NewIngresosHandler(ingresosRepo *repositories.IngresosRepository) *IngresosHandler {
	return &IngresosHandler{IngresosRepo: ingresosRepo}
}

func (h *IngresosHandler) GetIncomesList(c *gin.Context) {

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

	listaIngreso, err := h.IngresosRepo.GetIncomesList(finanzaId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al conseguir la lista de ingresos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "ingresos": listaIngreso})
}

func (h *IngresosHandler) GetIncomeById(c *gin.Context) {

	idIngreso, httpCode, jsonResponse := services.ParseUint(c)
	if idIngreso == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	ingreso, err := h.IngresosRepo.GetIncome(idIngreso)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Ese ingreso no existe"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al conseguir e ingreso"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "ingreso": ingreso})
}

type IncomeRequest struct {
	NombreIngreso      string  `json:"nombre_ingreso"`
	DescripcionIngreso string  `json:"descripcion_ingreso"`
	MontoIngreso       float64 `json:"monto_ingreso"`
}

func (h *IngresosHandler) CreateIncome(c *gin.Context) {

	var ingresoRequest IncomeRequest
	var ingreso models.TipoIngresos
	var finanzaId uint

	if err := c.ShouldBindJSON(&ingresoRequest); err != nil {
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

	ingreso.FinanzasID = finanzaId
	ingreso.UserID = userClaims.UserId
	ingreso.NombreIngresos = ingresoRequest.NombreIngreso
	ingreso.Descripcion = ingresoRequest.DescripcionIngreso
	ingreso.MontoIngreso = ingresoRequest.MontoIngreso

	if err := h.IngresosRepo.CreateIncome(&ingreso); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al crear el ingreso"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Ingreso creado con exito"})
}

func (h *IngresosHandler) UpdateIncome(c *gin.Context) {

	var ingresoRequest IncomeRequest

	idIngreso, httpCode, jsonResponse := services.ParseUint(c)
	if idIngreso == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	if err := c.ShouldBindJSON(&ingresoRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato de la peticion esta incorrecto"})
		return
	}

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	ingreso, err := h.IngresosRepo.GetIncomeById(idIngreso)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Ese ingreso no existe"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al conseguir e ingreso"})
		return
	}

	ingreso.NombreIngresos = ingresoRequest.NombreIngreso
	ingreso.Descripcion = ingresoRequest.DescripcionIngreso
	ingreso.MontoIngreso = ingresoRequest.MontoIngreso

	if err := h.IngresosRepo.UpdateIncome(ingreso); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al actualizar el ingreso"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "El ingreso se ha actualizado con exito"})
}
