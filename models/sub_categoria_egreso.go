package models

import "gorm.io/gorm"

type SubCategoriaEgreso struct {
	gorm.Model
	FinanzasID         uint              `json:"finanza_id" gorm:"index;not null"`
	Finanzas           Finanzas          `gorm:"foreignKey:FinanzasID"`
	NombreSubCategoria string            `json:"nombre_sub_categoria" gorm:"not null"`
	EsCompartida       bool              `json:"es_conjunta" gorm:"not null"`
	FinanzasConjuntoID *uint             `json:"finanza_conjunto_id" gorm:"index"`
	FinanzasConjunto   *FinanzasConjunto `gorm:"foreignKey:FinanzasConjuntoID"`
	UserID             *uint             `json:"id_usuario_registro" gorm:"index"`
	User               User              `gorm:"foreignKey:UserID"`
	Transacciones      []Transacciones   `json:"transacciones"`
	Presupuesto        []Presupuesto     `json:"presupuesto"`
}
