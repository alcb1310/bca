package database

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s ServiceMock) GetAllUsers(companyId uuid.UUID) ([]types.User, error) {
	return []types.User{}, nil
}

func (s ServiceMock) CreateUser(u *types.UserCreate) (types.User, error) {
	return types.User{}, nil
}

func (s ServiceMock) GetUser(id, companyId uuid.UUID) (types.User, error) {
	return types.User{}, nil
}

func (s ServiceMock) UpdateUser(u types.User, id, companyId uuid.UUID) (types.User, error) {
	return types.User{}, nil
}

func (s ServiceMock) UpdatePassword(pass string, id, companyId uuid.UUID) (types.User, error) {
	return types.User{}, nil
}

func (s ServiceMock) DeleteUser(id, companyId uuid.UUID) error {
	return nil
}
