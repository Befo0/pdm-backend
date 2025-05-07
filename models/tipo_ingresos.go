package models

import "gorm.io/gorm"

type TipoIngresos struct {
	gorm.Model
	FinanzasID         uint              `json:"finanza_id" gorm:"index;not null"`
	Finanzas           Finanzas          `gorm:"foreignKey:FinanzasID"`
	FinanzasConjuntoID *uint             `json:"finanza_conjunto_id" gorm:"index"`
	FinanzasConjunto   *FinanzasConjunto `gorm:"foreignKey:FinanzasConjuntoID"`
	NombreIngresos     string            `json:"nombre_ingresos" gorm:"not null"`
	MontoIngreso       float64           `json:"monto_ingreso" gorm:"not null"`
	UserID             uint              `json:"id_usuario_registro" gorm:"index"`
	User               User              `gorm:"foreignKey:UserID"`
	Transacciones      []Transacciones
}
