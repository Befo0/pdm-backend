package routes

import (
	"pdm-backend/controllers"
	"pdm-backend/middlewares"
	"pdm-backend/repositories"

	"github.com/gin-gonic/gin"
)

func CategoriaRouter(r *gin.Engine) {

	categoriaRepo := repositories.NewCategoriaRepository(repositories.GetDB())
	handler := controllers.NewCategoriaHandler(categoriaRepo)

	categoria := r.Group("/categoria")
	categoria.Use(middlewares.AuthMiddleware())
	{
		categoria.GET("/lista-categorias", handler.GetCategories)
		categoria.GET("/datos-categorias/:id", handler.GetCategoriesData)
		categoria.POST("/crear-categoria", handler.CreateCategoria)
		categoria.PATCH("/actualizar-categoria/:id", handler.UpdateCategoria)
	}
}
