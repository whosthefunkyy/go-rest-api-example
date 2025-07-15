package handlers

import (
	"github.com/whosthefunkyy/go-rest-api-example/models"
	"gorm.io/gorm"
	"errors"
)

// mockUserRepository for UserRepository 
type mockUserRepository struct{}

func (m *mockUserRepository) GetAll() ([]models.User, error) {
	return []models.User{
		{ID: 1, Name: "Artem", Age: 28},
		{ID: 2, Name: "Daniel", Age: 29},
	}, nil
}

func (m *mockUserRepository) GetByID(id int) (*models.User, error) {
	if id == 1 {
		return &models.User{ID: 1, Name: "Artem", Age: 28}, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepository) Create(user *models.User) error {
	return nil
}

func (m *mockUserRepository) Update(user *models.User) error {
	return nil
}

func (m *mockUserRepository) Delete(id int) error {
	return nil
}
// Failed
type mockFailingUserRepository struct{}

func (m *mockFailingUserRepository) GetAll() ([]models.User, error) {
	return nil, errors.New("database failure")
}

func (m *mockFailingUserRepository) GetByID(id int) (*models.User, error)   { return nil, nil }
func (m *mockFailingUserRepository) Create(user *models.User) error         { return nil }
func (m *mockFailingUserRepository) Update(user *models.User) error         { return nil }
func (m *mockFailingUserRepository) Delete(id int) error                    { return nil }

// ERROR 
type mockErrorUserRepository struct{}

func (m *mockErrorUserRepository) GetAll() ([]models.User, error) {
	return nil, nil
}

func (m *mockErrorUserRepository) GetByID(id int) (*models.User, error) {
	return nil, errors.New("unexpected database error")
}

func (m *mockErrorUserRepository) Create(user *models.User) error { return nil }
func (m *mockErrorUserRepository) Update(user *models.User) error { return nil }
func (m *mockErrorUserRepository) Delete(id int) error            { return nil }