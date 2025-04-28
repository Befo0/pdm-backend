package repositories

import (
	"gorm.io/gorm"
)

type SubCategoriaRepository struct {
	DB *gorm.DB
}

func NewSubCategoriaRepository(db *gorm.DB) *SubCategoriaRepository {
	return &SubCategoriaRepository{DB: db}
}
