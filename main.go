package main

import (
	"log"
	"os"
	"pdm-backend/routes"
	"pdm-backend/websockets"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("La variable de entorno no ha sido cargada")
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			return origin == "" || origin == "null"
		},
	}))
	go websockets.HandleBroadCast()

	routes.UserRouter(r)
	routes.FinanzaRouter(r)
	routes.CategoriaRouter(r)
	routes.TransaccionRouter(r)
	routes.SubCategoriaRouter(r)
	routes.IngresosRouter(r)
	routes.AhorroRouter(r)
	routes.InvitacionRouter(r)
	routes.FinanzaConjuntoRouter(r)
	websockets.WebSocketRouter(r)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = ":8080"
	}
	r.Run(PORT)
}
