package models

import (
	"pdm-backend/repositories"

	"gorm.io/gorm"
)

type TipoFinanzas struct {
	gorm.Model
	NombreTipo string     `json:"tipo_finanza"`
	Finanzas   []Finanzas `gorm:"foreignKey:TipoFinanzasID" json:"finanzas"`
}

type RolFinanzaConjunto struct {
	gorm.Model
	NombreRol string `json:"rol_usuario"`
}

type TipoPresupuesto struct {
	gorm.Model
	NombreTipoPresupuesto string          `json:"tipo_presupuesto"`
	Presupuestos          []Presupuesto   `json:"presupuestos"`
	Transacciones         []Transacciones `gorm:"foreignKey:TipoPresupuestoID" json:"transacciones"`
}

type TipoRegistro struct {
	gorm.Model
	NombreTipoRegistro string          `json:"tipo_registro"`
	Transacciones      []Transacciones `json:"transacciones"`
}

func SeedData() {

	db := repositories.GetDB()

	tipoFinanzas := []TipoFinanzas{
		{NombreTipo: "personal"},
		{NombreTipo: "conjunta"},
	}
	for _, tipo := range tipoFinanzas {
		var existing TipoFinanzas

		if err := db.Where("nombre_tipo = ?", tipo.NombreTipo).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&tipo)
			}
		}
	}

	rolFinanza := []RolFinanzaConjunto{
		{NombreRol: "admin"},
		{NombreRol: "colaborador"},
	}
	for _, rol := range rolFinanza {
		var existing RolFinanzaConjunto

		if err := db.Where("nombre_rol = ?", rol.NombreRol).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&rol)
			}
		}
	}

	tipoPresupuesto := []TipoPresupuesto{
		{NombreTipoPresupuesto: "Gastos variables"},
		{NombreTipoPresupuesto: "Gastos fijos"},
		{NombreTipoPresupuesto: "Gastos Provisionales"},
	}
	for _, tipoPresupuesto := range tipoPresupuesto {
		var existing TipoPresupuesto

		if err := db.Where("nombre_tipo_presupuesto = ?", tipoPresupuesto.NombreTipoPresupuesto).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&tipoPresupuesto)
			}
		}
	}

	tipoRegistroTransaccion := []TipoRegistro{
		{NombreTipoRegistro: "Ingreso"},
		{NombreTipoRegistro: "Egreso"},
	}
	for _, tipoRegistro := range tipoRegistroTransaccion {
		var existing TipoRegistro

		if err := db.Where("nombre_tipo_registro = ?", tipoRegistro.NombreTipoRegistro).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&tipoRegistro)
			}
		}
	}
}
