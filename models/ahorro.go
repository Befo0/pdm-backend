package models

import "gorm.io/gorm"

type MetaMensual struct {
	gorm.Model
	FinanzasID uint     `gorm:"index;not null" json:"finanza_id"`
	Anio       int      `gorm:"not null" json:"anio"`
	Mes        int      `gorm:"not null" json:"mes"`
	MontoMeta  float64  `gorm:"not null" json:"monto"`
	Finanzas   Finanzas `gorm:"foreignKey:FinanzasID" json:"finanza"`
}

type AhorroMensual struct {
	gorm.Model
	FinanzasID uint     `gorm:"index;not null"`
	Anio       int      `gorm:"not null"`
	Mes        int      `gorm:"not null"`
	Monto      float64  `gorm:"not null"`
	Finanzas   Finanzas `gorm:"foreignKey:FinanzasID" json:"finanza"`
}
