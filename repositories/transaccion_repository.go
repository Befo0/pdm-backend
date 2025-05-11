package repositories

import (
	"pdm-backend/models"
	"time"

	"gorm.io/gorm"
)

type TransaccionRepository struct {
	DB *gorm.DB
}

func NewTransaccionRepository(db *gorm.DB) *TransaccionRepository {
	return &TransaccionRepository{DB: db}
}

type ListaTransacciones struct {
	TransaccionId    uint
	NombreCategoria  string
	Monto            float64
	TipoMovimientoId uint
	TipoMovimiento   string
	FechaTransaccion string
	NombreUsuario    string
}

func (r *TransaccionRepository) GetTransactions(inicioMes, finMes time.Time, finanzaId uint) ([]ListaTransacciones, error) {

	transacciones := []ListaTransacciones{}

	err := r.DB.Model(models.Transacciones{}).Where("transacciones.finanzas_id = ? AND transacciones.fecha_registro >= ? AND transacciones.fecha_registro < ?", finanzaId, inicioMes, finMes).
		Select(`
		transacciones.id AS transaccion_id, 
		CASE 
			WHEN transacciones.tipo_registro_id = 1 THEN tipo_ingresos.nombre_ingresos
			WHEN transacciones.tipo_registro_id = 2 THEN categoria_egresos.nombre_categoria 
			ELSE ''
		END AS nombre_categoria, 
		transacciones.monto AS monto, 
		transacciones.tipo_registro_id AS tipo_movimiento_id, 
		tipo_registros.nombre_tipo_registro AS tipo_movimiento, 
		transacciones.fecha_registro AS fecha_transaccion, 
		users.nombre AS nombre_usuario`).
		Joins("LEFT JOIN categoria_egresos ON categoria_egresos.id = transacciones.categoria_egreso_id").
		Joins("LEFT JOIN tipo_ingresos ON tipo_ingresos.id = transacciones.tipo_ingresos_id").
		Joins("LEFT JOIN tipo_registros ON tipo_registros.id = transacciones.tipo_registro_id").
		Joins("LEFT JOIN users ON users.id = transacciones.user_id").
		Scan(&transacciones).Error

	if err != nil {
		return nil, err
	}

	return transacciones, nil
}

type Transaccion struct {
	TipoMovimientoID uint    `gorm:"column:tipo_movimiento_id"`
	TipoMovimiento   string  `gorm:"column:tipo_movimiento"`
	Movimiento       string  `gorm:"column:movimiento"`
	Categoria        string  `gorm:"column:categoria"`
	TipoGasto        string  `gorm:"column:tipo_gasto"`
	Presupuesto      float64 `gorm:"column:presupuesto"`
	Monto            float64 `gorm:"column:monto"`
	DescripcionGasto string  `gorm:"column:descripcion_gasto"`
	NombreUsuario    string  `gorm:"column:nombre_usuario"`
}

func (r *TransaccionRepository) GetTransactionById(transaccionId *uint) (*Transaccion, error) {

	var transaccion Transaccion

	tx := r.DB.Model(models.Transacciones{}).Where("transacciones.id = ?", transaccionId).
		Select(`
		transacciones.tipo_registro_id AS tipo_movimiento_id,
		tipo_registros.nombre_tipo_registro AS tipo_movimiento, 
		CASE
			WHEN transacciones.tipo_registro_id = 1 THEN tipo_ingresos.nombre_ingresos
			WHEN transacciones.tipo_registro_id = 2 THEN sub_categoria_egresos.nombre_sub_categoria
			ELSE ''
		END AS movimiento,
		CASE
			WHEN transacciones.tipo_registro_id = 2 THEN categoria_egresos.nombre_categoria 
			ELSE ''
		END AS categoria, 
		CASE 
			WHEN transacciones.tipo_registro_id = 2 THEN tipo_presupuestos.nombre_tipo_presupuesto 
			ELSE ''
		END AS tipo_gasto, 
		sub_categoria_egresos.presupuesto_mensual AS presupuesto, 
		transacciones.monto AS monto, 
		transacciones.descripcion AS descripcion_gasto, 
		users.nombre AS nombre_usuario`).
		Joins("LEFT JOIN tipo_registros ON tipo_registros.id = transacciones.tipo_registro_id").
		Joins("LEFT JOIN tipo_ingresos ON tipo_ingresos.id = transacciones.tipo_ingresos_id").
		Joins("LEFT JOIN sub_categoria_egresos ON sub_categoria_egresos.id = transacciones.sub_categoria_egreso_id").
		Joins("LEFT JOIN categoria_egresos ON categoria_egresos.id = transacciones.categoria_egreso_id").
		Joins("LEFT JOIN tipo_presupuestos ON tipo_presupuestos.id = transacciones.tipo_presupuesto_id").
		Joins("LEFT JOIN users ON users.id = transacciones.user_id").Scan(&transaccion)

	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &transaccion, nil
}

func (r *TransaccionRepository) CreateTransaction(transaccion *models.Transacciones) error {
	return r.DB.Create(&transaccion).Error
}
