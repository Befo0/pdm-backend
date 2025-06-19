package repositories

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func GetDB() *gorm.DB {
	once.Do(func() {
		if os.Getenv("ENV") != "production" {
			err := godotenv.Load(".env")
			if err != nil {
				log.Println("No se pudo cargar .env (esto es normal en producción)")
			}
		}

		dsn := os.Getenv("POSTGRES_URL")
		if dsn == "" {
			log.Fatal("No se ha encontrado la url de la base de datos")
		}

		DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Error conectando a la base de datos:", err)
		}

		db = DB
	})

	return db
}
