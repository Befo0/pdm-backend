package repositories

import (
	"pdm-backend/models"

	"gorm.io/gorm"
)

type IngresosRepository struct {
	DB *gorm.DB
}

func NewIngresosRepository(db *gorm.DB) *IngresosRepository {
	return &IngresosRepository{DB: db}
}

type IngresosOpciones struct {
	IdIngreso     uint   `json:"id_ingreso"`
	NombreIngreso string `json:"nombre_ingreso"`
}

func (r *IngresosRepository) GetIncomes(finanzaId uint) ([]IngresosOpciones, error) {

	var opciones []IngresosOpciones

	err := r.DB.Model(models.TipoIngresos{}).Where("finanzas_id = ?", finanzaId).
		Select("tipo_ingresos.id AS id_ingreso, tipo_ingresos.nombre_ingresos AS nombre_ingreso").
		Scan(&opciones).Error

	if err != nil {
		return nil, err
	}

	return opciones, nil
}

type IngresosLista struct {
	IdIngreso     uint    `json:"id_ingreso"`
	NombreIngreso string  `json:"nombre_ingreso"`
	MontoIngreso  float64 `json:"monto_ingreso"`
	NombreUsuario string  `json:"nombre_usuario"`
}

func (r *IngresosRepository) GetIncomesList(finanzaId uint) ([]IngresosLista, error) {

	var listaIngresos []IngresosLista

	err := r.DB.Model(models.TipoIngresos{}).Where("finanzas_id = ?", finanzaId).
		Select("tipo_ingresos.id AS id_ingreso, tipo_ingresos.nombre_ingresos AS nombre_ingreso, tipo_ingresos.monto_ingreso AS monto_ingreso, users.nombre AS nombre_usuario").
		Joins("LEFT JOIN users ON users.id = tipo_ingresos.user_id").
		Scan(&listaIngresos).Error

	if err != nil {
		return nil, err
	}

	return listaIngresos, nil
}

func (r *IngresosRepository) CreateIncome(ingreso *models.TipoIngresos) error {
	return r.DB.Create(&ingreso).Error
}
