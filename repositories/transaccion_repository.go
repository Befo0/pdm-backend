package repositories

import (
	"errors"
	"pdm-backend/models"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransaccionRepository struct {
	DB *gorm.DB
}

func NewTransaccionRepository(db *gorm.DB) *TransaccionRepository {
	return &TransaccionRepository{DB: db}
}

type ListaTransacciones struct {
	TransaccionId    uint    `json:"transaccion_id"`
	NombreCategoria  string  `json:"nombre_categoria"`
	Monto            float64 `json:"monto_transaccion"`
	TipoMovimientoId uint    `json:"tipo_movimiento_id"`
	TipoMovimiento   string  `json:"tipo_movimiento_nombre"`
	FechaTransaccion string  `json:"fecha_transaccion"`
	NombreUsuario    string  `json:"nombre_usuario"`
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
		Order("transacciones.fecha_registro DESC").
		Scan(&transacciones).Error

	if err != nil {
		return nil, err
	}

	return transacciones, nil
}

type OpcionesTransaccion struct {
	IdRegistro     uint          `json:"tipo_registro_id"`
	NombreRegistro string        `json:"tipo_registro_nombre"`
	Opciones       []interface{} `json:"opciones"`
}

func (r *TransaccionRepository) GetOptions(finanzaId uint) ([]OpcionesTransaccion, error) {

	opcionesTransaccion := []OpcionesTransaccion{}
	errCh := make(chan error, 2)

	ingresosRepo := NewIngresosRepository(r.DB)
	subCategoriaRepo := NewSubCategoriaRepository(r.DB)
	var wg sync.WaitGroup

	err := r.DB.Model(models.TipoRegistro{}).
		Select("tipo_registros.id AS id_registro, tipo_registros.nombre_tipo_registro AS nombre_registro").
		Scan(&opcionesTransaccion).Error
	if err != nil {
		return nil, err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		subCategoriaOpciones, err := subCategoriaRepo.GetSubCategories(finanzaId)
		if err != nil {
			errCh <- err
			return
		}
		for i := range opcionesTransaccion {
			if opcionesTransaccion[i].IdRegistro == 2 {
				for _, opcion := range subCategoriaOpciones {
					opcionesTransaccion[i].Opciones = append(opcionesTransaccion[i].Opciones, opcion)
				}
				break
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ingresosOpciones, err := ingresosRepo.GetIncomes(finanzaId)
		if err != nil {
			errCh <- err
			return
		}
		for i := range opcionesTransaccion {
			if opcionesTransaccion[i].IdRegistro == 1 {
				for _, opcion := range ingresosOpciones {
					opcionesTransaccion[i].Opciones = append(opcionesTransaccion[i].Opciones, opcion)
				}
				break
			}
		}
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	return opcionesTransaccion, nil
}

type Transaccion struct {
	TipoMovimientoID uint    `json:"tipo_movimiento_id" gorm:"column:tipo_movimiento_id"`
	TipoMovimiento   string  `json:"tipo_movimiento" gorm:"column:tipo_movimiento"`
	Movimiento       string  `json:"movimiento" gorm:"column:movimiento"`
	Categoria        string  `json:"categoria" gorm:"column:categoria"`
	TipoGasto        string  `json:"tipo_gasto" gorm:"column:tipo_gasto"`
	Presupuesto      float64 `json:"presupuesto" gorm:"column:presupuesto"`
	Monto            float64 `json:"monto" gorm:"column:monto"`
	DescripcionGasto string  `json:"descripcion_gasto" gorm:"column:descripcion_gasto"`
	NombreUsuario    string  `json:"nombre_usuario" gorm:"column:nombre_usuario"`
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
		meta_mensuals.monto_meta AS presupuesto,
		transacciones.monto AS monto,
		transacciones.descripcion AS descripcion_gasto,
		users.nombre AS nombre_usuario
	`).
		Joins("LEFT JOIN tipo_registros ON tipo_registros.id = transacciones.tipo_registro_id").
		Joins("LEFT JOIN tipo_ingresos ON tipo_ingresos.id = transacciones.tipo_ingresos_id").
		Joins("LEFT JOIN sub_categoria_egresos ON sub_categoria_egresos.id = transacciones.sub_categoria_egreso_id").
		Joins("LEFT JOIN categoria_egresos ON categoria_egresos.id = transacciones.categoria_egreso_id").
		Joins("LEFT JOIN tipo_presupuestos ON tipo_presupuestos.id = transacciones.tipo_presupuesto_id").
		Joins("LEFT JOIN meta_mensuals ON meta_mensuals.finanzas_id = transacciones.finanzas_id AND meta_mensuals.mes = EXTRACT(MONTH FROM transacciones.fecha_registro) AND meta_mensuals.anio = EXTRACT(YEAR FROM transacciones.fecha_registro)").
		Joins("LEFT JOIN users ON users.id = transacciones.user_id").
		Scan(&transaccion)

	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if err := tx.Error; err != nil {
		return nil, err
	}

	return &transaccion, nil
}

type IdSubCategorias struct {
	CategoriaId uint
	GastoId     uint
}

func (r *TransaccionRepository) GetIds(subCategoriaId uint) (*IdSubCategorias, error) {

	var identificadores IdSubCategorias

	err := r.DB.Model(models.SubCategoriaEgreso{}).Where("sub_categoria_egresos.id = ?", subCategoriaId).
		Select("sub_categoria_egresos.categoria_egreso_id AS categoria_id, sub_categoria_egresos.tipo_presupuesto_id AS gasto_id").
		Scan(&identificadores).Error
	if err != nil {
		return nil, err
	}

	return &identificadores, nil
}

func (r *TransaccionRepository) CreateTransaction(transaccion *models.Transacciones) error {
	return r.DB.Create(&transaccion).Error
}

func (r *TransaccionRepository) CreateOrUpdateSaving(finanzasId uint, monto float64, fecha time.Time) error {

	anio := fecha.Year()
	mes := int(fecha.Month())

	var ahorro models.AhorroMensual
	err := r.DB.Model(models.AhorroMensual{}).Where("finanzas_Id = ? AND anio = ? AND mes = ?", finanzasId, anio, mes).First(&ahorro).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			nuevoAhorro := models.AhorroMensual{
				FinanzasID: finanzasId,
				Monto:      monto,
				Mes:        mes,
				Anio:       anio,
			}
			return db.Create(&nuevoAhorro).Error
		}
		return err
	}

	ahorro.Monto += monto

	return r.DB.Save(&ahorro).Error
}

func (r *TransaccionRepository) GetSavingSubCategorie(finanzaId uint) (uint, error) {
	var subCategoriaId uint

	err := r.DB.Model(models.SubCategoriaEgreso{}).Where("finanzas_id = ? AND nombre_sub_categoria = ?", finanzaId, "Ahorro").
		Select("id").Scan(&subCategoriaId).Error
	if err != nil {
		return 0, err
	}

	return subCategoriaId, nil
}

type PayloadEvent struct {
	Event   string      `json:"event"`
	Payload interface{} `json:"payload"`
}

type BroadCastMessage struct {
	FinanzaId uint           `json:"finanza_id"`
	EventInfo []PayloadEvent `json:"event_info"`
}

func (r *TransaccionRepository) BuildWebSocketEvent(finanzaId uint, fecha time.Time, transactionSubCategorieId *uint, savingSubCategorieId uint) (*BroadCastMessage, error) {

	type result[T any] struct {
		data T
		err  error
	}

	var (
		eventInfo       []PayloadEvent
		wg              sync.WaitGroup
		chResumen       = make(chan result[gin.H])
		chDatos         = make(chan result[[]DashboardData])
		chTransacciones = make(chan result[[]ListaTransacciones])
		chAhorro        = make(chan result[[]AhorroResponse])
	)

	inicioMes := time.Date(fecha.Year(), fecha.Month(), 1, 0, 0, 0, 0, time.UTC)
	finMes := inicioMes.AddDate(0, 1, 0)

	finanzaRepo := NewFinanzaRepository(r.DB)
	ahorroRepo := NewAhorroRepository(r.DB)

	wg.Add(1)
	go func() {
		defer wg.Done()
		resumen, err := finanzaRepo.GetDashboardSummary(finanzaId, inicioMes, finMes)
		chResumen <- result[gin.H]{resumen, err}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		datos, err := finanzaRepo.GetDataSummary(inicioMes, finMes, finanzaId)
		chDatos <- result[[]DashboardData]{datos, err}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		lista, err := r.GetTransactions(inicioMes, finMes, finanzaId)
		chTransacciones <- result[[]ListaTransacciones]{lista, err}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if transactionSubCategorieId != nil && *transactionSubCategorieId == savingSubCategorieId {
			ahorro, err := ahorroRepo.GetSavingsData(finanzaId, fecha.Year())
			chAhorro <- result[[]AhorroResponse]{ahorro, err}
		} else {
			chAhorro <- result[[]AhorroResponse]{nil, nil}
		}
	}()

	go func() {
		wg.Wait()
		close(chResumen)
		close(chDatos)
		close(chTransacciones)
		close(chAhorro)
	}()

	resumenRes := <-chResumen
	if resumenRes.err != nil {
		return nil, resumenRes.err
	}

	datosRes := <-chDatos
	if datosRes.err != nil {
		return nil, datosRes.err
	}

	transaccionesRes := <-chTransacciones
	if transaccionesRes.err != nil {
		return nil, transaccionesRes.err
	}

	ahorroRes := <-chAhorro
	if ahorroRes.err != nil {
		return nil, ahorroRes.err
	}

	eventInfo = append(eventInfo, PayloadEvent{Event: "resumen_finanza", Payload: resumenRes.data})
	eventInfo = append(eventInfo, PayloadEvent{Event: "datos_finanza", Payload: datosRes.data})
	eventInfo = append(eventInfo, PayloadEvent{Event: "lista_transacciones", Payload: transaccionesRes.data})
	eventInfo = append(eventInfo, PayloadEvent{Event: "ahorro_finanza", Payload: ahorroRes.data})

	webSocketEvent := BroadCastMessage{
		FinanzaId: finanzaId,
		EventInfo: eventInfo,
	}

	return &webSocketEvent, nil
}
