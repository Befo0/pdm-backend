package repositories

import (
	"pdm-backend/models"
	"pdm-backend/services"
	"strings"
	"time"

	"gorm.io/gorm"
)

type InvitacionRepository struct {
	DB *gorm.DB
}

func NewInvitacionRepository(db *gorm.DB) *InvitacionRepository {
	return &InvitacionRepository{DB: db}
}

type Invitacion struct {
	Codigo string `json:"codigo_invitacion"`
}

func (r *InvitacionRepository) CreateInvite(finanzaId *uint) (*Invitacion, error) {

	var invitacionRespuesta Invitacion
	intentosMax := 5

	for i := 0; i < intentosMax; i++ {

		codigo, err := services.GenerateInvitacionCode(10)
		if err != nil {
			return nil, err
		}

		invitacion := models.Invitaciones{
			FinanzasID: *finanzaId,
			Codigo:     codigo,
			ExpiraEn:   time.Now().Add(time.Minute * 15),
		}

		err = r.DB.Create(&invitacion).Error
		if err == nil {
			invitacionRespuesta.Codigo = invitacion.Codigo
			break
		}

		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "duplicate key") {
			return nil, err
		}

		return nil, err
	}
	return &invitacionRespuesta, nil
}
