package main

import (
	"log"
	"pdm-backend/models"
	"pdm-backend/repositories"

	"gorm.io/gorm"
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

	SeedData()

	log.Print("Migraciones realizadas con exito")
}

func SeedData() {

	db := repositories.GetDB()

	tipoFinanzas := []models.TipoFinanzas{
		{NombreTipo: "personal"},
		{NombreTipo: "conjunta"},
	}
	for _, tipo := range tipoFinanzas {
		var existing models.TipoFinanzas

		if err := db.Where("nombre_tipo = ?", tipo.NombreTipo).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&tipo)
			}
		}
	}

	rolFinanza := []models.RolFinanzaConjunto{
		{NombreRol: "admin"},
		{NombreRol: "colaborador"},
	}
	for _, rol := range rolFinanza {
		var existing models.RolFinanzaConjunto

		if err := db.Where("nombre_rol = ?", rol.NombreRol).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&rol)
			}
		}
	}

	tipoPresupuesto := []models.TipoPresupuesto{
		{NombreTipoPresupuesto: "Gastos variables"},
		{NombreTipoPresupuesto: "Gastos fijos"},
		{NombreTipoPresupuesto: "Gastos Provisionales"},
	}
	for _, tipoPresupuesto := range tipoPresupuesto {
		var existing models.TipoPresupuesto

		if err := db.Where("nombre_tipo_presupuesto = ?", tipoPresupuesto.NombreTipoPresupuesto).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&tipoPresupuesto)
			}
		}
	}

	tipoRegistroTransaccion := []models.TipoRegistro{
		{NombreTipoRegistro: "Ingreso"},
		{NombreTipoRegistro: "Egreso"},
	}
	for _, tipoRegistro := range tipoRegistroTransaccion {
		var existing models.TipoRegistro

		if err := db.Where("nombre_tipo_registro = ?", tipoRegistro.NombreTipoRegistro).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&tipoRegistro)
			}
		}
	}
}
