package main

import (
	"log"
	"pdm-backend/models"
	"pdm-backend/repositories"
)

func main() {

	db := repositories.GetDB()

	err := db.Migrator().DropTable(
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
		&models.Transacciones{},
		&models.Ahorro{},
		&models.Invitaciones{},
	)
	if err != nil {
		log.Fatal("Error al reiniciar la base de datos")
	}

}
