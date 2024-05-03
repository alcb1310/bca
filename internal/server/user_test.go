package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
	"bca-go-final/mocks"
)

func TestProfile(t *testing.T) {
	db := mocks.NewServiceMock()
	_, srv := server.NewServer(db)

	db.On("GetUser", uuid.UUID{}, uuid.UUID{}).Return(types.User{
		Id:        uuid.New(),
		Name:      "Test",
		Email:     "test@b.com",
		RoleId:    "a",
		CompanyId: uuid.New(),
	}, nil)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bca/user/perfil", nil)

	srv.Profile(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Perfil")
	assert.Contains(t, response.Body.String(), "test@b.com")
	assert.Contains(t, response.Body.String(), "Test")
}

func TestAdmin(t *testing.T) {
	db := mocks.NewServiceMock()
	_, srv := server.NewServer(db)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bca/user/admin", nil)

	srv.Admin(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Admin")
}
