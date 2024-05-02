package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
	"bca-go-final/mocks"
)

var validCredentials = &types.Login{
	Email:    "test@test.com",
	Password: "test",
}

func TestLogin(t *testing.T) {
	t.Run("Should return method not allowed on GET", func(t *testing.T) {
		db := mocks.NewServiceMock()
		_, srv := server.NewServer(db)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/bca/login", nil)

		srv.Login(response, request)

		assert.Equal(t, http.StatusMethodNotAllowed, response.Code)
	})

	t.Run("Should return 200 on Options", func(t *testing.T) {
		db := mocks.NewServiceMock()
		_, srv := server.NewServer(db)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodOptions, "/bca/login", nil)

		srv.Login(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("POST Method", func(t *testing.T) {
		t.Run("Valid Credentials", func(t *testing.T) {
			db := mocks.NewServiceMock()
			_, srv := server.NewServer(db)
			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(validCredentials); err != nil {
				t.Fatal(err)
			}

			db.On("Login", validCredentials).Return("token", nil)

			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/bca/login", &buf)

			srv.Login(response, request)

			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, `{"token":"token"}`, response.Body.String())
		})

		t.Run("Invalid Credentials", func(t *testing.T) {
			t.Run("Data validation", func(t *testing.T) {
				t.Run("Should return 400 when no credentials provided", func(t *testing.T) {
					db := mocks.NewServiceMock()
					_, srv := server.NewServer(db)

					response := httptest.NewRecorder()
					request := httptest.NewRequest(http.MethodPost, "/bca/login", nil)

					srv.Login(response, request)

					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, `{"error":"EOF"}`, response.Body.String())
				})

				t.Run("Invalid email", func(t *testing.T) {
					t.Run("Empty email", func(t *testing.T) {
						db := mocks.NewServiceMock()
						_, srv := server.NewServer(db)

						data := make(map[string]interface{})
						data["email"] = ""

						var buf bytes.Buffer
						if err := json.NewEncoder(&buf).Encode(data); err != nil {
							t.Fatal(err)
						}

						response := httptest.NewRecorder()
						request := httptest.NewRequest(http.MethodPost, "/bca/login", &buf)

						srv.Login(response, request)

						assert.Equal(t, http.StatusBadRequest, response.Code)
						assert.Equal(t, `{"error":"email cannot be empty","field":"email"}`, response.Body.String())
					})

					t.Run("Invalid email", func(t *testing.T) {
						db := mocks.NewServiceMock()
						_, srv := server.NewServer(db)

						data := make(map[string]interface{})
						data["email"] = "test"

						var buf bytes.Buffer
						if err := json.NewEncoder(&buf).Encode(data); err != nil {
							t.Fatal(err)
						}

						response := httptest.NewRecorder()
						request := httptest.NewRequest(http.MethodPost, "/bca/login", &buf)

						srv.Login(response, request)

						assert.Equal(t, http.StatusBadRequest, response.Code)
						assert.Equal(t, `{"error":"invalid email","field":"email"}`, response.Body.String())

					})
				})

				t.Run("Invalid password", func(t *testing.T) {
					t.Run("Empty password", func(t *testing.T) {
						db := mocks.NewServiceMock()
						_, srv := server.NewServer(db)

						data := make(map[string]interface{})
						data["email"] = "test@email.com"
						data["password"] = ""

						var buf bytes.Buffer
						if err := json.NewEncoder(&buf).Encode(data); err != nil {
							t.Fatal(err)
						}

						response := httptest.NewRecorder()
						request := httptest.NewRequest(http.MethodPost, "/bca/login", &buf)

						srv.Login(response, request)

						assert.Equal(t, http.StatusBadRequest, response.Code)
						assert.Equal(t, `{"error":"password cannot be empty","field":"password"}`, response.Body.String())
					})
				})
			})

			t.Run("Incorrect credentials", func(t *testing.T) {
				db := mocks.NewServiceMock()
				_, srv := server.NewServer(db)

				var buf bytes.Buffer
				if err := json.NewEncoder(&buf).Encode(validCredentials); err != nil {
					t.Fatal(err)
				}

				db.On("Login", validCredentials).Return("", fmt.Errorf("invalid credentials"))

				response := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodPost, "/bca/login", &buf)

				srv.Login(response, request)

				assert.Equal(t, http.StatusUnauthorized, response.Code)
				assert.Equal(t, `{"error":"invalid credentials"}`, response.Body.String())
			})
		})
	})
}
