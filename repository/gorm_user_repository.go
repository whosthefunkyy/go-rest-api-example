package repository

import (
	"gorm.io/gorm"
	"github.com/whosthefunkyy/go-rest-api-example/models"
)

type GormUserRepository struct {
	DB *gorm.DB
}

func (r *GormUserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *GormUserRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, id).Error
	return &user, err
}

func (r *GormUserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *GormUserRepository) Update(user *models.User) error {
	return r.DB.Save(user).Error
}

func (r *GormUserRepository) Delete(id int) error {
	return r.DB.Delete(&models.User{}, id).Error
}
