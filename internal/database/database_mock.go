package database

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

type ServiceMock struct {
	Service
}

func (s ServiceMock) Health() map[string]string {
	return map[string]string{
		"message": "It's healthy",
	}
}

func (s ServiceMock) Levels(companyId uuid.UUID) []types.Select {
	return []types.Select{
		{Key: "level 1", Value: "level 1"},
		{Key: "level 2", Value: "level 2"},
		{Key: "level 3", Value: "level 3"},
	}
}
