package repositories

import (
	"errors"
	"pdm-backend/models"

	"gorm.io/gorm"
)

type AhorroRepository struct {
	DB *gorm.DB
}

func NewAhorroRepository(db *gorm.DB) *AhorroRepository {
	return &AhorroRepository{DB: db}
}

var meses = []string{"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"}

type AhorroResponse struct {
	NombreMes              string  `json:"nombre_mes"`
	Anio                   int     `json:"anio"`
	MetaAhorro             float64 `json:"meta_ahorro"`
	MontoAhorrado          float64 `json:"monto_ahorrado"`
	PorcentajeCumplimiento float64 `json:"porcentaje_cumplimiento"`
}

func (r *AhorroRepository) GetSavingsData(finanzaId uint, anio int) ([]AhorroResponse, error) {

	var ahorroResponse []AhorroResponse

	for index, mes := range meses {
		var metaAhorro models.MetaMensual
		var ahorroMensual models.AhorroMensual

		err := r.DB.Model(models.MetaMensual{}).Where("finanzas_id = ?, anio = ?, mes = ?", finanzaId, anio, index+1).
			First(&metaAhorro).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return nil, err
		}
		err = r.DB.Model(models.AhorroMensual{}).Where("finanzas_id = ?, anio = ?, mes = ?", finanzaId, anio, index+1).
			First(&ahorroMensual).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		var porcentaje float64
		if metaAhorro.MontoMeta > 0 {
			porcentaje = (ahorroMensual.Monto / metaAhorro.MontoMeta) * 100

		}

		ahorroResponse = append(ahorroResponse, AhorroResponse{
			NombreMes:              mes,
			Anio:                   metaAhorro.Anio,
			MetaAhorro:             metaAhorro.MontoMeta,
			MontoAhorrado:          ahorroMensual.Monto,
			PorcentajeCumplimiento: porcentaje,
		})

	}

	return ahorroResponse, nil
}

func (r *AhorroRepository) CreateOrUpdateSavingGoal(finanzaId uint, monto float64, mes, anio int) error {

	var ahorro models.MetaMensual
	err := r.DB.Model(models.MetaMensual{}).Where("finanzasId = ? AND anio = ? AND mes = ?", finanzaId, anio, mes).
		First(&ahorro).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			nuevaMeta := models.MetaMensual{
				FinanzasID: finanzaId,
				MontoMeta:  monto,
				Mes:        mes,
				Anio:       anio,
			}

			return r.DB.Create(&nuevaMeta).Error
		}
		return err
	}

	ahorro.MontoMeta += monto

	return r.DB.Save(&ahorro).Error
}
