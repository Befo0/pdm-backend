package main

import (
	"log"
	"pdm-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("La variable de entorno no ha sido cargada")
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET,POST,PUT,PATCH,DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	routes.UserRouter(r)
	routes.FinanzaRouter(r)

	r.Run(":8000")
}
