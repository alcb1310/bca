package server_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
	"bca-go-final/mocks"
)

func TestLoginView(t *testing.T) {
	t.Run("GET Method", func(t *testing.T) {
		db := mocks.NewServiceMock()
		_, srv := server.NewServer(db)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/login", nil)

		srv.LoginView(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Contains(t, response.Body.String(), "Login")
	})

	t.Run("POST Method", func(t *testing.T) {
		t.Run("valid credentials", func(t *testing.T) {
			db := mocks.NewServiceMock()
			_, srv := server.NewServer(db)

			form := url.Values{}
			form.Add("email", validCredentials.Email)
			form.Add("password", validCredentials.Password)
			buf := strings.NewReader(form.Encode())

			db.On("Login", validCredentials).Return("token", nil)

			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/login", buf)
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			srv.LoginView(response, request)

			assert.Equal(t, http.StatusSeeOther, response.Code)
		})

		t.Run("invalid credentials", func(t *testing.T) {
			t.Run("email", func(t *testing.T) {
				t.Run("empty email", func(t *testing.T) {
					db := mocks.NewServiceMock()
					_, srv := server.NewServer(db)

					form := url.Values{}
					buf := strings.NewReader(form.Encode())

					response := httptest.NewRecorder()
					request := httptest.NewRequest(http.MethodPost, "/login", buf)

					srv.LoginView(response, request)

					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Contains(t, response.Body.String(), "credenciales inválidas")
				})

				t.Run("invalid email", func(t *testing.T) {
					db := mocks.NewServiceMock()
					_, srv := server.NewServer(db)

					form := url.Values{}
					form.Add("email", "invalid")
					buf := strings.NewReader(form.Encode())

					response := httptest.NewRecorder()
					request := httptest.NewRequest(http.MethodPost, "/login", buf)

					srv.LoginView(response, request)

					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Contains(t, response.Body.String(), "credenciales inválidas")
				})
			})

			t.Run("password", func(t *testing.T) {
				t.Run("empty password", func(t *testing.T) {
					db := mocks.NewServiceMock()
					_, srv := server.NewServer(db)

					form := url.Values{}
					form.Add("email", validCredentials.Email)
					buf := strings.NewReader(form.Encode())

					response := httptest.NewRecorder()
					request := httptest.NewRequest(http.MethodPost, "/login", buf)

					srv.LoginView(response, request)

					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Contains(t, response.Body.String(), "credenciales inválidas")
				})
			})

			t.Run("invalid credentials", func(t *testing.T) {
				db := mocks.NewServiceMock()
				_, srv := server.NewServer(db)

				form := url.Values{}
				form.Add("email", validCredentials.Email)
				form.Add("password", validCredentials.Password)
				buf := strings.NewReader(form.Encode())

				db.On("Login", validCredentials).Return("", errors.New("unknown error"))

				response := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodPost, "/login", buf)
				request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

				srv.LoginView(response, request)

				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "credenciales inválidas")
			})
		})
	})

	t.Run("unimplemented method", func(t *testing.T) {
		db := mocks.NewServiceMock()
		_, srv := server.NewServer(db)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPut, "/login", nil)

		srv.LoginView(response, request)

		assert.Equal(t, http.StatusMethodNotAllowed, response.Code)
	})
}

func TestBcaView(t *testing.T) {
	db := mocks.NewServiceMock()
	_, srv := server.NewServer(db)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bca", nil)

	srv.BcaView(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Bienvenido")
}

func TestLogout(t *testing.T) {
	db := mocks.NewServiceMock()
	_, srv := server.NewServer(db)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/logout", nil)

	srv.Logout(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
}
