package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bca-go-final/internal/database"
)

func TestProfile(t *testing.T) {
	db := database.ServiceMock{}
	_, router := NewServer(db)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bca/user/perfil", nil)

	router.Profile(response, request)

	got := response.Code
	want := http.StatusOK

	if got != want {
		t.Errorf("Got %d, want %d", got, want)
	}

	expected := "Mi Perfil"
	recieved := response.Body.String()

	if !strings.Contains(recieved, expected) {
		t.Errorf("Response does not contain %s, but is %s", expected, recieved)
	}
}

func TestAdmin(t *testing.T) {
	db := database.ServiceMock{}
	_, router := NewServer(db)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bca/user/admin", nil)

	router.Admin(response, request)

	got := response.Code
	want := http.StatusOK

	if got != want {
		t.Errorf("Got %d, want %d", got, want)
	}

	expected := "Administrar usuarios"
	recieved := response.Body.String()

	if !strings.Contains(recieved, expected) {
		t.Errorf("Response does not contain %s, but is %s", expected, recieved)
	}
}

// TEST: ChangePassword function
func TestChangePassword(t *testing.T) {
	t.Skip("Not implemented")
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
