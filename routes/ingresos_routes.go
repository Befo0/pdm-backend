package routes

import (
	"pdm-backend/controllers"
	"pdm-backend/middlewares"
	"pdm-backend/repositories"

	"github.com/gin-gonic/gin"
)

func IngresosRouter(r *gin.Engine) {

	ingresosRepo := repositories.NewIngresosRepository(repositories.GetDB())
	handler := controllers.NewIngresosHandler(ingresosRepo)

	ingresos := r.Group("/ingresos")
	ingresos.Use(middlewares.AuthMiddleware())
	{
		ingresos.GET("/lista", handler.GetIncomesList)
		ingresos.GET("/ingreso/:id", handler.GetIncomeById)
		ingresos.POST("/crear", handler.CreateIncome)
		ingresos.PUT("/actualizar/:id", handler.UpdateIncome)
	}
}
