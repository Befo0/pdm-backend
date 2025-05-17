package repositories

import (
	"pdm-backend/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(&user).Error
}

func (r *UserRepository) CreateUserAndFinance(user *models.User) error {

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		finanza := models.Finanzas{
			UserID:         user.ID,
			TipoFinanzasID: 1,
		}
		if err := tx.Create(&finanza).Error; err != nil {
			return err
		}

		categoria := models.CategoriaEgreso{
			FinanzasID:      finanza.ID,
			NombreCategoria: "Ahorro",
			UserID:          user.ID,
		}
		if err := tx.Create(&categoria).Error; err != nil {
			return err
		}

		subCategoria := models.SubCategoriaEgreso{
			FinanzasID:         finanza.ID,
			NombreSubCategoria: "Ahorro",
			TipoPresupuestoID:  3,
			CategoriaEgresoID:  categoria.ID,
			PresupuestoMensual: 0.00,
			UserID:             user.ID,
		}
		if err := tx.Create(&subCategoria).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	if err := r.DB.Where("correo = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserById(id uint) (*models.User, error) {
	var user models.User

	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

type Identificadores struct {
	FinanzaId uint
	AhorroId  uint
}

func (r *UserRepository) GetFinanceAndSavingSubCategoryByUserId(userId uint) (*Identificadores, error) {
	var identificadores Identificadores

	err := r.DB.Model(&models.Finanzas{}).
		Select("finanzas.id AS finanza_id, sub_categoria_egresos.id AS ahorro_id").
		Joins("JOIN sub_categoria_egresos ON finanzas.id = sub_categoria_egresos.finanzas_id").
		Where("finanzas.user_id = ? AND finanzas.tipo_finanzas_id = ? AND sub_categoria_egresos.nombre_sub_categoria = ?", userId, 1, "Ahorro").
		Scan(&identificadores).Error

	if err != nil {
		return nil, err
	}

	return &identificadores, nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	return r.DB.Save(&user).Error
}
