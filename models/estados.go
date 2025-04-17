package models

import (
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
