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
		finanzaConj.GET("/lista", handler.GetConjFinances)
		finanzaConj.GET("/detalles/:id", handler.GetConjFinancesDetails)
		finanzaConj.POST("/unirse", handler.JoinUser)
		finanzaConj.POST("/crear", handler.CreateConjFinance)
		finanzaConj.DELETE("/salir/:id", handler.LeaveConjFinance)
		finanzaConj.DELETE("/user/:id", handler.DeleteUserFromFinance)
	}
}
