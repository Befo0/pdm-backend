package routes

import "github.com/gin-gonic/gin"

func UserRoutes(r *gin.Engine) {
	user := r.Group("/user")
	user.POST("/login")
	user.POST("/register")
}
