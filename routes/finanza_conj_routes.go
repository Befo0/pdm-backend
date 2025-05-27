package routes

import (
	"pdm-backend/controllers"
	"pdm-backend/middlewares"
	"pdm-backend/repositories"

	"github.com/gin-gonic/gin"
)

func FinanzaConjuntoRouter(r *gin.Engine) {

	finanzaRepo := repositories.NewFinanzaConjRepository(repositories.GetDB())
	handler := controllers.NewFinanzaConjHandler(finanzaRepo)

	finanzaConj := r.Group("/finanza-conjunta")
	finanzaConj.Use(middlewares.AuthMiddleware())
	{
		finanzaConh.GET("/lista", handler.GetConjFinances)
		finanzaConj.POST("/crear", handler.CreateConjFinance)
	}
}
