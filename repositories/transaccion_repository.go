package repositories

import (
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

func (r *TransaccionRepository) GetTransactions(inicioMes, finMes time.Time, finanzaId uint) (gin.H, error) {

	return nil, nil
}
