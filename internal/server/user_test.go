package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"bca-go-final/mocks"
)

func TestProfile(t *testing.T) {
	db := mocks.NewServiceMock()
	_, router := NewServer(db)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bca/user/perfil", nil)

	router.Profile(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Mi Perfil")
}

func TestAdmin(t *testing.T) {
	db := mocks.NewServiceMock()
	_, router := NewServer(db)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bca/user/admin", nil)

	router.Admin(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Administrar usuarios")
}

func TestChangePassword(t *testing.T) {
	db := mocks.NewServiceMock()
	_, router := NewServer(db)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bca/user/cambio", nil)

	router.ChangePassword(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Cambiar Contraseña")
}

// TEST: SingleUser function
func TestSingleUser(t *testing.T) {
	t.Skip("Not implemented")
}

// TEST: UsersTable function
func TestUsersTable(t *testing.T) {
	t.Skip("Not implemented")
}

// TEST: UserAdd function
func TestUserAdd(t *testing.T) {
	t.Skip("Not implemented")
}

// TEST: Profile function
func TestUserEdit(t *testing.T) {
	t.Skip("Not implemented")
}
