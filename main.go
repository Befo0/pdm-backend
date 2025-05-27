package main

import (
	"log"
	"os"
	"pdm-backend/routes"
	"pdm-backend/websockets"

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

	go websockets.HandleBroadCast()

	routes.UserRouter(r)
	routes.FinanzaRouter(r)
	routes.CategoriaRouter(r)
	routes.TransaccionRouter(r)
	routes.SubCategoriaRouter(r)
	routes.IngresosRouter(r)
	routes.AhorroRouter(r)
	websockets.WebSocketRouter(r)

	PORT := os.Getenv("PORT")
	r.Run(PORT)
}
