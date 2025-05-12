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

type SubCategoriaHandler struct {
	SubCategoriaRepo *repositories.SubCategoriaRepository
}

func NewSubCategoriaHandler(subCategoriaRepo *repositories.SubCategoriaRepository) *SubCategoriaHandler {
	return &SubCategoriaHandler{SubCategoriaRepo: subCategoriaRepo}
}

func (h *SubCategoriaHandler) GetSubCategories(c *gin.Context) {

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	subCategorias, err := h.SubCategoriaRepo.GetSubCategories(userClaims.FinanzaId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al traer las opciones de subCategorias"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sub_categorias": subCategorias})
}

func (h *SubCategoriaHandler) GetSubCategoriesExpensesType(c *gin.Context) {

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	tipoGastos, err := h.SubCategoriaRepo.GetSubCategoriesExpensesType()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al conseguir las opciones de gastos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"opciones": tipoGastos})
}

func (h *SubCategoriaHandler) GetSubCategoriesList(c *gin.Context) {

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	subCategoriasList, err := h.SubCategoriaRepo.GetSubCategoriesList(userClaims.FinanzaId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al conseguir la lista de sub categorias"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "lista_sub_categorias": subCategoriasList})
}

type SubCategoriaRequest struct {
	CategoriaId        uint    `json:"categoria_id" binding:"required,gt=0"`
	NombreSubCategoria string  `json:"nombre_sub_categoria" binding:"required"`
	TipoGastoId        uint    `json:"tipo_gasto_id" binding:"required,gt=0"`
	Presupuesto        float64 `json:"presupuesto" binding:"required,gte=0"`
}

func (h *SubCategoriaHandler) CreateSubCategoria(c *gin.Context) {
	var subCategoriaRequest SubCategoriaRequest
	var subCategoria models.SubCategoriaEgreso

	if err := c.ShouldBindJSON(&subCategoriaRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato de la peticion esta incorrecto"})
		return
	}

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	subCategoria.FinanzasID = userClaims.FinanzaId
	subCategoria.UserID = userClaims.UserId
	subCategoria.CategoriaEgresoID = subCategoriaRequest.CategoriaId
	subCategoria.NombreSubCategoria = subCategoriaRequest.NombreSubCategoria
	subCategoria.PresupuestoMensual = subCategoriaRequest.Presupuesto
	subCategoria.TipoPresupuestoID = subCategoriaRequest.TipoGastoId

	if err := h.SubCategoriaRepo.CreateSubCategory(&subCategoria); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al crear la sub categoria"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "La sub categoria fue creada con exito"})
}

func (h *SubCategoriaHandler) UpdateSubCategoria(c *gin.Context) {

	var updateRequest SubCategoriaRequest

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato de la peticion esta incorrecto"})
		return
	}

	idSubCategoria, httpCode, jsonResponse := services.ParseUint(c)
	if idSubCategoria == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	subCategoria, err := h.SubCategoriaRepo.GetSubCategoryById(idSubCategoria)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "No se encontro la sub categoria"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Hubo un error al conseguir la sub categoria"})
		return
	}

	subCategoria.CategoriaEgresoID = updateRequest.CategoriaId
	subCategoria.NombreSubCategoria = updateRequest.NombreSubCategoria
	subCategoria.PresupuestoMensual = updateRequest.Presupuesto
	subCategoria.TipoPresupuestoID = updateRequest.TipoGastoId

	if err := h.SubCategoriaRepo.DB.Save(&subCategoria).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al modificar la sub categoria"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "La sub categoria fue modificada correctamente"})

}
