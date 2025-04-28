package models

import "gorm.io/gorm"

type Finanzas struct {
	gorm.Model
	UserID            uint                 `json:"user_id" gorm:"index;not null"`
	User              User                 `gorm:"foreignKey:UserID"`
	TipoFinanzasID    uint                 `json:"tipo_finanza_id" gorm:"index;not null"`
	TipoFinanzas      TipoFinanzas         `gorm:"foreignKey:TipoFinanzasID"`
	Titulo            *string              `json:"titulo" gorm:"size:255"`
	Descripcion       *string              `json:"descripcion" gorm:"size:500"`
	Transacciones     []Transacciones      `json:"transacciones"`
	SubCategorias     []SubCategoriaEgreso `json:"sub_categorias"`
	Categorias        []CategoriaEgreso    `json:"categorias"`
	TipoIngresos      []TipoIngresos       `json:"tipo_ingresos"`
	FinanzasConjuntos []FinanzasConjunto   `gorm:"foreignKey:FinanzasID" json:"finanzas_conjuntos"`
}
