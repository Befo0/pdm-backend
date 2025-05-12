package routes

import (
	"pdm-backend/controllers"
	"pdm-backend/middlewares"
	"pdm-backend/repositories"

	"github.com/gin-gonic/gin"
)

func SubCategoriaRouter(r *gin.Engine) {

	subCategoriaRepo := repositories.NewSubCategoriaRepository(repositories.GetDB())
	handler := controllers.NewSubCategoriaHandler(subCategoriaRepo)

	subCategoria := r.Group("/sub-categoria")
	subCategoria.Use(middlewares.AuthMiddleware())
	{
		subCategoria.GET("/opciones", handler.GetSubCategories)
		subCategoria.GET("/lista", handler.GetSubCategoriesList)
		subCategoria.GET("/opciones-gasto", handler.GetSubCategoriesExpensesType)
		subCategoria.POST("/crear", handler.CreateSubCategoria)
		subCategoria.PUT("/actualizar/:id", handler.UpdateSubCategoria)
	}
}
