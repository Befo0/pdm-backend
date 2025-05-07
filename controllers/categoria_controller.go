package controllers

import (
	"net/http"
	"pdm-backend/repositories"
	"pdm-backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoriaHandler struct {
	CategoriaRepo *repositories.CategoriaRepository
}

func NewCategoriaHandler(categoriaRepo *repositories.CategoriaRepository) *CategoriaHandler {
	return &CategoriaHandler{CategoriaRepo: categoriaRepo}
}

func (h *CategoriaHandler) GetCategories(c *gin.Context) {

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	categorias, err := h.CategoriaRepo.GetCategories(userClaims.FinanzaId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al conseguir las categorias de la finanza"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categorias": categorias})
}

func (h *CategoriaHandler) GetCategoriesData(c *gin.Context) {

	idParam := c.Param("id")

	idUint, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El id no es un numero valido"})
		return
	}
	idCategoria := uint(idUint)

	userClaims, httpCode, jsonResponse := services.GetClaims(c)
	if userClaims == nil {
		c.JSON(httpCode, jsonResponse)
		return
	}

	datosFiltro, err := h.CategoriaRepo.GetCategoriesData(userClaims.FinanzaId, idCategoria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al conseguir los datos de la categoria"})
		return
	}

	c.JSON(http.StatusOK, datosFiltro)
}
