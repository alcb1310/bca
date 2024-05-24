package server_test

import (
	"bytes"
	"encoding/json"
	"errors"
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

var validRegister = &types.CompanyCreate{
	Name:      "test",
	Ruc:       "123456789",
	Email:     "test@test.com",
	Password:  "test",
	Employees: 100,
	User:      "test",
}

func TestRegister(t *testing.T) {
	t.Run("Should return method not allowed on GET", func(t *testing.T) {
		db := mocks.NewServiceMock()
		_, srv := server.NewServer(db)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/bca/register", nil)

		srv.Register(response, request)

		assert.Equal(t, http.StatusMethodNotAllowed, response.Code)
	})

	t.Run("Should return 200 on Options", func(t *testing.T) {
		db := mocks.NewServiceMock()
		_, srv := server.NewServer(db)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodOptions, "/bca/register", nil)

		srv.Register(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("POST Method", func(t *testing.T) {
		t.Run("Valid Input", func(t *testing.T) {
			db := mocks.NewServiceMock()
			_, srv := server.NewServer(db)

			db.On("CreateCompany", validRegister).Return(nil)

			var buf bytes.Buffer
			json.NewEncoder(&buf).Encode(validRegister)

			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/bca/register", &buf)

			srv.Register(response, request)

			assert.Equal(t, http.StatusCreated, response.Code)
		})

		t.Run("Invalid Input", func(t *testing.T) {
			t.Run("Data validation", func(t *testing.T) {
				t.Run("Should return 400 when no data is sent", func(t *testing.T) {
					db := mocks.NewServiceMock()
					_, srv := server.NewServer(db)

					response := httptest.NewRecorder()
					request := httptest.NewRequest(http.MethodPost, "/bca/register", nil)

					srv.Register(response, request)

					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), `{"error":"EOF"}`)
				})

				t.Run("Name", func(t *testing.T) {
					t.Run("Empty name", func(t *testing.T) {
						db := mocks.NewServiceMock()
						_, srv := server.NewServer(db)
						data := make(map[string]interface{})
						data["name"] = ""

						var buf bytes.Buffer

						if err := json.NewEncoder(&buf).Encode(data); err != nil {
							t.Fatal(err)
						}

						response := httptest.NewRecorder()
						request := httptest.NewRequest(http.MethodPost, "/bca/register", &buf)

						srv.Register(response, request)

						assert.Equal(t, http.StatusBadRequest, response.Code)
						assert.Equal(t, response.Body.String(), `{"error":"name cannot be empty","field":"name"}`)
					})
				})

				t.Run("Ruc", func(t *testing.T) {
					t.Run("Empty ruc", func(t *testing.T) {
						db := mocks.NewServiceMock()
						_, srv := server.NewServer(db)
						data := make(map[string]interface{})
						data["name"] = "test"
						data["ruc"] = ""

						var buf bytes.Buffer

						if err := json.NewEncoder(&buf).Encode(data); err != nil {
							t.Fatal(err)
						}

						response := httptest.NewRecorder()
						request := httptest.NewRequest(http.MethodPost, "/bca/register", &buf)

						srv.Register(response, request)

						assert.Equal(t, http.StatusBadRequest, response.Code)
						assert.Equal(t, response.Body.String(), `{"error":"ruc cannot be empty","field":"ruc"}`)
					})
				})

				t.Run("Employees", func(t *testing.T) {
					t.Run("No employees", func(t *testing.T) {
						db := mocks.NewServiceMock()
						_, srv := server.NewServer(db)
						data := make(map[string]interface{})
						data["name"] = "test"
						data["ruc"] = "123456789"
						data["employees"] = 0

						var buf bytes.Buffer

						if err := json.NewEncoder(&buf).Encode(data); err != nil {
							t.Fatal(err)
						}

						response := httptest.NewRecorder()
						request := httptest.NewRequest(http.MethodPost, "/bca/register", &buf)

						srv.Register(response, request)

						assert.Equal(t, http.StatusBadRequest, response.Code)
						assert.Equal(t, response.Body.String(), `{"error":"should pass at least one employee","field":"employees"}`)
					})

					t.Run("Invalid employees", func(t *testing.T) {
						db := mocks.NewServiceMock()
						_, srv := server.NewServer(db)
						data := make(map[string]interface{})
						data["name"] = "test"
						data["ruc"] = "123456789"
						data["employees"] = "invalid"

						var buf bytes.Buffer

						if err := json.NewEncoder(&buf).Encode(data); err != nil {
							t.Fatal(err)
						}

						response := httptest.NewRecorder()
						request := httptest.NewRequest(http.MethodPost, "/bca/register", &buf)

						srv.Register(response, request)

						assert.Equal(t, http.StatusBadRequest, response.Code)
						assert.Equal(t, response.Body.String(), `{"error":"employees must be a valid positive number","field":"employees"}`)
					})

					t.Run("Negative employees", func(t *testing.T) {
						db := mocks.NewServiceMock()
						_, srv := server.NewServer(db)
						data := make(map[string]interface{})
						data["name"] = "test"
						data["ruc"] = "123456789"
						data["employees"] = -1

						var buf bytes.Buffer

						if err := json.NewEncoder(&buf).Encode(data); err != nil {
							t.Fatal(err)
						}

						response := httptest.NewRecorder()
						request := httptest.NewRequest(http.MethodPost, "/bca/register", &buf)

						srv.Register(response, request)

						assert.Equal(t, http.StatusBadRequest, response.Code)
						assert.Equal(t, response.Body.String(), `{"error":"employees must be a valid positive number","field":"employees"}`)
					})
				})

				t.Run("Email", func(t *testing.T) {
					t.Run("Empty email", func(t *testing.T) {
						db := mocks.NewServiceMock()
						_, srv := server.NewServer(db)
						data := make(map[string]interface{})
						data["name"] = "test"
						data["ruc"] = "123456789"
						data["employees"] = 1
						data["email"] = ""

						var buf bytes.Buffer

						if err := json.NewEncoder(&buf).Encode(data); err != nil {
							t.Fatal(err)
						}

						response := httptest.NewRecorder()
						request := httptest.NewRequest(http.MethodPost, "/bca/register", &buf)

						srv.Register(response, request)

						assert.Equal(t, http.StatusBadRequest, response.Code)
						assert.Equal(t, response.Body.String(), `{"error":"email cannot be empty","field":"email"}`)
					})

					t.Run("Invalid email", func(t *testing.T) {
						db := mocks.NewServiceMock()
						_, srv := server.NewServer(db)
						data := make(map[string]interface{})
						data["name"] = "test"
						data["ruc"] = "123456789"
						data["employees"] = 1
						data["email"] = "invalid"

						var buf bytes.Buffer

						if err := json.NewEncoder(&buf).Encode(data); err != nil {
							t.Fatal(err)
						}

						response := httptest.NewRecorder()
						request := httptest.NewRequest(http.MethodPost, "/bca/register", &buf)

						srv.Register(response, request)

						assert.Equal(t, http.StatusBadRequest, response.Code)
						assert.Equal(t, response.Body.String(), `{"error":"invalid email","field":"email"}`)
					})
				})

				t.Run("Password", func(t *testing.T) {
					t.Run("Empty password", func(t *testing.T) {
						db := mocks.NewServiceMock()
						_, srv := server.NewServer(db)
						data := make(map[string]interface{})
						data["name"] = "test"
						data["ruc"] = "123456789"
						data["employees"] = 1
						data["email"] = "valid@test.com"
						data["password"] = ""

						var buf bytes.Buffer

						if err := json.NewEncoder(&buf).Encode(data); err != nil {
							t.Fatal(err)
						}

						response := httptest.NewRecorder()
						request := httptest.NewRequest(http.MethodPost, "/bca/register", &buf)

						srv.Register(response, request)

						assert.Equal(t, http.StatusBadRequest, response.Code)
						assert.Equal(t, response.Body.String(), `{"error":"password cannot be empty","field":"password"}`)
					})
				})

				t.Run("User", func(t *testing.T) {
					t.Run("Empty user", func(t *testing.T) {
						db := mocks.NewServiceMock()
						_, srv := server.NewServer(db)
						data := make(map[string]interface{})
						data["name"] = "test"
						data["ruc"] = "123456789"
						data["employees"] = 1
						data["email"] = "valid@test.com"
						data["password"] = "valid"
						data["user"] = ""

						var buf bytes.Buffer

						if err := json.NewEncoder(&buf).Encode(data); err != nil {
							t.Fatal(err)
						}

						response := httptest.NewRecorder()
						request := httptest.NewRequest(http.MethodPost, "/bca/register", &buf)

						srv.Register(response, request)

						assert.Equal(t, http.StatusBadRequest, response.Code)
						assert.Equal(t, response.Body.String(), `{"error":"name of the user cannot be empty","field":"user"}`)
					})
				})
			})

			t.Run("Duplicate information", func(t *testing.T) {
				db := mocks.NewServiceMock()
				_, srv := server.NewServer(db)

				db.On("CreateCompany", validRegister).Return(errors.New("SQLSTATE 23505"))

				buf := new(bytes.Buffer)
				if err := json.NewEncoder(buf).Encode(validRegister); err != nil {
					t.Fatal(err)
				}

				response := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodPost, "/bca/register", buf)

				srv.Register(response, request)

				assert.Equal(t, http.StatusConflict, response.Code)
				assert.Equal(t, response.Body.String(), `{"error":"company already exists"}`)
			})

			t.Run("Other SQL errors", func(t *testing.T) {
				db := mocks.NewServiceMock()
				_, srv := server.NewServer(db)

				db.On("CreateCompany", validRegister).Return(errors.New("other error"))

				buf := new(bytes.Buffer)
				if err := json.NewEncoder(buf).Encode(validRegister); err != nil {
					t.Fatal(err)
				}

				response := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodPost, "/bca/register", buf)

				srv.Register(response, request)

				assert.Equal(t, http.StatusInternalServerError, response.Code)
				assert.Equal(t, response.Body.String(), `{"error":"other error"}`)
			})
		})
	})
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
