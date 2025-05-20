package models

import (
	"time"

	"gorm.io/gorm"
)

type Transacciones struct {
	gorm.Model
	FinanzasID           uint                `json:"finanza_id" gorm:"index;not null"`
	Finanzas             Finanzas            `gorm:"foreignKey:FinanzasID"`
	Descripcion          *string             `json:"descripcion" gorm:"size:500"`
	TipoRegistroID       uint                `json:"tipo_registro_id" gorm:"index;not null"`
	TipoRegistro         TipoRegistro        `gorm:"foreignKey:TipoRegistroID"`
	TipoIngresosID       *uint               `json:"tipo_ingreso_id" gorm:"index"`
	TipoIngresos         *TipoIngresos       `gorm:"foreignKey:TipoIngresosID"`
	CategoriaEgresoID    *uint               `json:"categoria_egreso_id" gorm:"index"`
	CategoriaEgreso      *CategoriaEgreso    `gorm:"foreignKey:CategoriaEgresoID"`
	SubCategoriaEgresoID *uint               `json:"sub_categoria_egreso_id" gorm:"index"`
	SubCategoriaEgreso   *SubCategoriaEgreso `gorm:"foreignKey:SubCategoriaEgresoID"`
	TipoPresupuestoID    *uint               `json:"tipo_gasto_id" gorm:"index"`
	TipoPresupuesto      *TipoPresupuesto    `gorm:"foreignKey:TipoPresupuestoID"`
	UserID               uint                `json:"id_usuario_registro" gorm:"index"`
	User                 User                `gorm:"foreignKey:UserID"`
	FechaRegistro        time.Time           `json:"fecha_registro" gorm:"not null"`
	Monto                float64             `json:"monto" gorm:"not null"`
}
