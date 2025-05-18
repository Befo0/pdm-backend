package routes

import (
	"pdm-backend/controllers"
	"pdm-backend/middlewares"
	"pdm-backend/repositories"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.Engine) {

	userRepo := repositories.NewUserRepository(repositories.GetDB())
	financeRepo := repositories.NewFinanzaRepository(repositories.GetDB())
	handler := controllers.NewUserHandler(userRepo, financeRepo)

	user := r.Group("/usuario")
	user.POST("/login", handler.Login)
	user.POST("/registro", handler.Register)

	user.Use(middlewares.AuthMiddleware())
	{
		user.PATCH("/cambiar-perfil", handler.UpdateProfile)
		user.PATCH("/cambiar-contrasena", handler.UpdatePassword)
	}
}
