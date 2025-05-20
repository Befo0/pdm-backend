package routes

import (
	"pdm-backend/controllers"
	"pdm-backend/middlewares"
	"pdm-backend/repositories"

	"github.com/gin-gonic/gin"
)

func InvitacionRouter(r *gin.Engine) {

	invitacionRepository := repositories.NewInvitacionRepository(repositories.GetDB())
	handler := controllers.NewInvitacionHandler(invitacionRepository)

	invitacion := r.Group("/invitaciones")
	invitacion.Use(middlewares.AuthMiddleware())
	{
		invitacion.POST("/crear/:id", handler.CreateInvite)
	}
}
