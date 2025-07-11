package controllers

import (
	"errors"
	"net/http"
	"pdm-backend/models"
	"pdm-backend/repositories"
	"pdm-backend/services"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoriaHandler struct {
	CategoriaRepo *repositories.CategoriaRepository
}

func NewCategoriaHandler(categoriaRepo *repositories.CategoriaRepository) *CategoriaHandler {
	return &CategoriaHandler{CategoriaRepo: categoriaRepo}
}

func (h *CategoriaHandler) GetCategories(c *gin.Context) {

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

	categorias, err := h.CategoriaRepo.GetCategories(finanzaId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al conseguir las categorias de la finanza"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categorias": categorias})
}

func (h *CategoriaHandler) GetCategoriesData(c *gin.Context) {

	var finanzaId uint

	idCategoria, httpCode, jsonResponse := services.ParseUint(c)
	if idCategoria == nil {
		c.JSON(httpCode, jsonResponse)
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

	datosFiltro, err := h.CategoriaRepo.GetCategoriesData(finanzaId, idCategoria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al conseguir los datos de la categoria"})
		return
	}

	c.JSON(http.StatusOK, datosFiltro)
}

type CategoriaRequest struct {
	NombreCategoria string `json:"nombre_categoria"`
}

func (h *CategoriaHandler) CreateCategoria(c *gin.Context) {

	var finanzaId uint
	var categoriaRequest CategoriaRequest
	var categoria models.CategoriaEgreso

	if err := c.ShouldBindJSON(&categoriaRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato de la peticion esta incorrecto"})
		return
	}

	nombre := strings.TrimSpace(strings.ToLower(categoriaRequest.NombreCategoria))
	if nombre == "ahorro" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "No puedes crear otra categoria llamada Ahorro"})
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

	categoria.NombreCategoria = categoriaRequest.NombreCategoria
	categoria.FinanzasID = finanzaId
	categoria.UserID = userClaims.UserId

	if err := h.CategoriaRepo.CreateCategory(&categoria); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Hubo un error al crear la categoria"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "La categoria fue creada con exito"})
}

func (h *CategoriaHandler) UpdateCategoria(c *gin.Context) {

	var updateRequest CategoriaRequest

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato de la peticion esta incorrecto"})
		return
	}

	idCategoria, httpCode, jsonResponse := services.ParseUint(c)
	if idCategoria == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	categoria, err := h.CategoriaRepo.GetCategoryById(idCategoria)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "No se encontro la categoria"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Hubo un error al conseguir la categoria"})
		return
	}

	categoria.NombreCategoria = updateRequest.NombreCategoria

	if err := h.CategoriaRepo.UpdateCategory(categoria); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al modificar la categoria"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "La categoria fue modificada correctamente"})
}

func (h *CategoriaHandler) GetCategoriesList(c *gin.Context) {

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

	listaCategorias, err := h.CategoriaRepo.GetCategoriesList(finanzaId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al conseguir las categorias de la finanza"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"lista_categorias": listaCategorias, "finanza_id": finanzaId})
}
