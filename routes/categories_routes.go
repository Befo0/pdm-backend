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
		categoria.GET("/opciones", handler.GetCategories)
		categoria.GET("/lista", handler.GetCategoriesList)
		categoria.GET("/datos/:id", handler.GetCategoriesData)
		categoria.POST("/crear", handler.CreateCategoria)
		categoria.PATCH("/actualizar/:id", handler.UpdateCategoria)
	}
}
