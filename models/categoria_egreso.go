package models

import "gorm.io/gorm"

type CategoriaEgreso struct {
	gorm.Model
	FinanzasID         uint                 `json:"finanza_id" gorm:"index;not null"`
	Finanzas           Finanzas             `gorm:"foreignKey:FinanzasID"`
	FinanzasConjuntoID *uint                `json:"finanza_conjunto_id" gorm:"index"`
	FinanzasConjunto   *FinanzasConjunto    `gorm:"foreignKey:FinanzasConjuntoID"`
	NombreCategoria    string               `json:"nombre_categoria" gorm:"not null"`
	UserID             *uint                `json:"id_usuario_registro" gorm:"index"`
	User               *User                `gorm:"foreignKey:UserID"`
	Transacciones      []Transacciones      `json:"transacciones"`
	SubCategoriaEgreso []SubCategoriaEgreso `json:"sub_categorias"`
}
