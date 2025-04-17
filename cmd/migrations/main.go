package main

import (
	"log"
	"pdm-backend/models"
	"pdm-backend/repositories"
)

func main() {

	db := repositories.GetDB()

	err := db.AutoMigrate(
		&models.TipoFinanzas{},
		&models.TipoPresupuesto{},
		&models.TipoRegistro{},
		&models.TipoIngresos{},
		&models.RolFinanzaConjunto{},
		&models.User{},
		&models.Finanzas{},
		&models.FinanzasConjunto{},
		&models.CategoriaEgreso{},
		&models.SubCategoriaEgreso{},
		&models.Presupuesto{},
		&models.Transacciones{},
	)
	if err != nil {
		log.Fatal("Ocurrio un error al realizar las migraciones ", err)
	}

	models.SeedData()

	log.Print("Migraciones realizadas con exito")
}
