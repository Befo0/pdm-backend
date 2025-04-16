package repositories

import (
	"database/sql"
	"log"
	"os"
	"pdm-backend/models"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() {

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

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Error en las migraciones:", err)
	}

	log.Println("Conexión y Migraciones realizadas con exito")
}

func GetDB() *gorm.DB {

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
