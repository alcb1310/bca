package database

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s ServiceMock) GetAllCategories(companyId uuid.UUID) ([]types.Category, error) {
	return []types.Category{}, nil
}

func (s ServiceMock) CreateCategory(category types.Category) error {
	return nil
}

func (s ServiceMock) GetCategory(id, companyId uuid.UUID) (types.Category, error) {
	return types.Category{}, nil
}

func (s ServiceMock) UpdateCategory(category types.Category) error {
	return nil
}
