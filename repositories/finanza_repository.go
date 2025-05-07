package repositories

import (
	"pdm-backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FinanzaRepository struct {
	DB *gorm.DB
}

func NewFinanzaRepository(db *gorm.DB) *FinanzaRepository {
	return &FinanzaRepository{DB: db}
}

func SumarMonto(db *gorm.DB, modelo interface{}, finanzaId uint, tipo int, inicio, final time.Time) (float64, error) {
	var total float64

	err := db.Model(modelo).
		Where("finanzas_id = ? AND tipo_registro_id = ? AND fecha_registro >= ? AND fecha_registro < ?", finanzaId, tipo, inicio, final).
		Select("COALESCE(SUM(monto), 0)").Scan(&total).Error

	return total, err
}

func (r *FinanzaRepository) GetFinanceSummary(finanzaId uint, inicio, final time.Time) (gin.H, error) {

	ingresosTotales, err := SumarMonto(r.DB, models.Transacciones{}, finanzaId, 1, inicio, final)
	if err != nil {
		return nil, err
	}

	egresosTotales, err := SumarMonto(r.DB, models.Transacciones{}, finanzaId, 2, inicio, final)
	if err != nil {
		return nil, err
	}

	diferencia := ingresosTotales - egresosTotales

	resumenJSON := gin.H{
		"ingresos_totales": ingresosTotales,
		"egresos_totales":  egresosTotales,
		"diferencia":       diferencia,
	}

	return resumenJSON, nil
}

func (r *FinanzaRepository) GetEgresoSummary(finanzaId uint, inicio, final time.Time) (gin.H, error) {
	var presupuestoMensual float64

	err := r.DB.Model(models.SubCategoriaEgreso{}).
		Where("finanzas_id = ?", finanzaId).
		Select("COALESCE(SUM(presupuesto_mensual), 0)").
		Scan(&presupuestoMensual).Error

	if err != nil {
		return nil, err
	}

	egresosTotales, err := SumarMonto(r.DB, models.Transacciones{}, finanzaId, 2, inicio, final)
	if err != nil {
		return nil, err
	}

	variacion := presupuestoMensual - egresosTotales

	registroJSON := gin.H{
		"presupuesto_mensual": presupuestoMensual,
		"consumo_mensual":     egresosTotales,
		"variacion_mensual":   variacion,
	}

	return registroJSON, nil
}

func (r *FinanzaRepository) GetSavingsSummary(finanzaId uint, inicio, final time.Time) (gin.H, error) {
	var metaAhorro float64
	var ahorroGuardado float64
	var subCategoriaId uint
	var porcentajeAhorro float64

	err := r.DB.Model(models.Ahorro{}).Where("finanzas_id = ?", finanzaId).Select("monto").Scan(&metaAhorro).Error
	if err != nil {
		return nil, err
	}

	err = r.DB.Model(models.SubCategoriaEgreso{}).Where("finanzas_id  = ? AND nombre_sub_categoria = ?", finanzaId, "Ahorro").
		Select("id").Scan(&subCategoriaId).Error

	if err != nil {
		return nil, err
	}

	err = r.DB.Model(models.Transacciones{}).
		Where("finanzas_id = ? AND tipo_registro_id = ? AND fecha_registro >= ? AND fecha_registro < ? AND sub_categoria_egreso_id = ?", finanzaId, 2, inicio, final, subCategoriaId).
		Select("COALESCE(SUM(monto), 0)").Scan(&ahorroGuardado).Error

	if err != nil {
		return nil, err
	}

	if metaAhorro != 0 {
		porcentajeAhorro = (ahorroGuardado * 100) / metaAhorro
	} else {
		porcentajeAhorro = 0
	}

	ahorroJSON := gin.H{
		"meta":                metaAhorro,
		"acumulado":           ahorroGuardado,
		"progreso_porcentaje": porcentajeAhorro,
	}

	return ahorroJSON, nil
}

type DashboardData struct {
	CategoriaId      uint `json:"-"`
	CategoriaNombre  string
	TotalPresupuesto float64
	Gasto            float64
	Diferencia       float64
}

func (r *FinanzaRepository) GetDataSummary(inicioMes, finMes time.Time, finanzaId uint) (*[]DashboardData, error) {

	var resultados []DashboardData

	err := r.DB.Model(models.CategoriaEgreso{}).Select("categoria_egresos.id AS categoria_id, categoria_egresos.nombre_categoria AS categoria_nombre, COALESCE(SUM(sub_categoria_egresos.presupuesto_mensual), 0) AS total_presupuesto").
		Joins("LEFT JOIN sub_categoria_egresos ON sub_categoria_egresos.categoria_egreso_id = categoria_egresos.id").
		Group("categoria_egresos.nombre_categoria").
		Order("categoria_egresos.nombre_categoria").
		Scan(&resultados).Error
	if err != nil {
		return nil, err
	}

	for index := range resultados {
		var totalGasto float64

		err := r.DB.Model(models.Transacciones{}).Where("finanzas_id = ? AND tipo_registro_id = ? AND fecha_registro >= ? AND fecha_registro < ? AND categoria_egreso_id = ?", finanzaId, 2, inicioMes, finMes, resultados[index].CategoriaId).Select("COALESCE(SUM(monto), 0)").Scan(&totalGasto).Error

		if err != nil {
			return nil, err
		}

		diferencia := resultados[index].TotalPresupuesto - totalGasto

		resultados[index].Gasto = totalGasto
		resultados[index].Diferencia = diferencia

	}

	return &resultados, nil
}
