package repositories

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetDB() *gorm.DB {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("La variable de entorno no ha sido cargada")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("No se ha encontrado la url de la base de datos")
	}

	db, err := sql.Open("libsql", dsn)
	if err != nil {
		log.Fatal("Error abriendo la conexión libsql:", err)
	}

	DB, err := gorm.Open(sqlite.Dialector{Conn: db}, &gorm.Config{})
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}

	return DB
}
