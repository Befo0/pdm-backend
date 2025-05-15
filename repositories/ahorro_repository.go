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

func (r *AhorroRepository) CreateOrUpdateSavingGoal(finanzaId uint, monto float64, mes, anio int) error {

	var ahorro models.MetaMensual
	err := r.DB.Model(models.AhorroMensual{}).Where("finanzasId = ? AND anio = ? AND mes = ?", finanzaId, anio, mes).
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
