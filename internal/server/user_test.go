package server_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
)

func TestProfile(t *testing.T) {
	srv, db := server.MakeServer()

	db.On("GetUser", uuid.UUID{}, uuid.UUID{}).Return(types.User{
		Id:        uuid.New(),
		Name:      "Test",
		Email:     "test@b.com",
		RoleId:    "a",
		CompanyId: uuid.New(),
	}, nil)

	request, response := server.MakeRequest(http.MethodGet, "/bca/partials/users", nil)

	srv.Profile(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Perfil")
	assert.Contains(t, response.Body.String(), "test@b.com")
	assert.Contains(t, response.Body.String(), "Test")
}

func TestAdmin(t *testing.T) {
	srv, _ := server.MakeServer()
	request, response := server.MakeRequest(http.MethodGet, "/bca/partials/users", nil)

	srv.Admin(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Admin")
}

func TestChangePassword(t *testing.T) {
	srv, _ := server.MakeServer()
	request, response := server.MakeRequest(http.MethodGet, "/bca/partials/users", nil)

	srv.ChangePassword(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Cambiar Contraseña")
}

func TestUsersTable(t *testing.T) {
	t.Run("GET request", func(t *testing.T) {
		srv, db := server.MakeServer()

		db.On("GetAllUsers", uuid.UUID{}).Return([]types.User{
			{
				Id:        uuid.New(),
				Name:      "Test",
				Email:     "test@b.com",
				RoleId:    "a",
				CompanyId: uuid.New(),
			},
			{
				Id:        uuid.New(),
				Name:      "Test2",
				Email:     "test2@b.com",
				RoleId:    "a",
				CompanyId: uuid.New(),
			},
		}, nil)

		request, response := server.MakeRequest(http.MethodGet, "/bca/partials/users", nil)

		srv.UsersTable(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Contains(t, response.Body.String(), "<td>test@b.com</td>")
		assert.Contains(t, response.Body.String(), "<td>test2@b.com</td>")
	})

	t.Run("POST request", func(t *testing.T) {
		t.Run("Invalid data", func(t *testing.T) {
			t.Run("Email", func(t *testing.T) {
				t.Run("Empty", func(t *testing.T) {
					srv, _ := server.MakeServer()

					form := url.Values{}
					form.Add("email", "")
					buf := bytes.NewBufferString(form.Encode())

					request, response := server.MakeRequest(http.MethodPost, "/bca/partials/users", buf)

					srv.UsersTable(response, request)

					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "invalid email")
				})

				t.Run("Invalid", func(t *testing.T) {
					srv, _ := server.MakeServer()

					form := url.Values{}
					form.Add("email", "test")
					buf := bytes.NewBufferString(form.Encode())

					request, response := server.MakeRequest(http.MethodPost, "/bca/partials/users", buf)

					srv.UsersTable(response, request)

					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "invalid email")
				})
			})

			t.Run("Password", func(t *testing.T) {
				srv, _ := server.MakeServer()

				form := url.Values{}
				form.Add("email", "test@b.com")
				form.Add("password", "")
				buf := bytes.NewBufferString(form.Encode())

				request, response := server.MakeRequest(http.MethodPost, "/bca/partials/users", buf)

				srv.UsersTable(response, request)

				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Equal(t, response.Body.String(), "invalid password")
			})

			t.Run("Name", func(t *testing.T) {
				srv, _ := server.MakeServer()

				form := url.Values{}
				form.Add("email", "test@b.com")
				form.Add("password", "test")
				form.Add("name", "")
				buf := bytes.NewBufferString(form.Encode())

				request, response := server.MakeRequest(http.MethodPost, "/bca/partials/users", buf)

				srv.UsersTable(response, request)

				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Equal(t, response.Body.String(), "invalid name")
			})
		})

		t.Run("Valid data", func(t *testing.T) {
			user := &types.UserCreate{
				Id:        uuid.UUID{},
				Name:      "Test",
				Email:     "test@b.com",
				Password:  "test",
				RoleId:    "a",
				CompanyId: uuid.UUID{},
			}

			t.Run("Successfull request", func(t *testing.T) {
				srv, db := server.MakeServer()

				db.On("CreateUser", user).Return(types.User{
					Id:        uuid.New(),
					Name:      user.Name,
					Email:     user.Email,
					RoleId:    user.RoleId,
					CompanyId: user.CompanyId,
				}, nil)

				db.On("GetAllUsers", uuid.UUID{}).Return([]types.User{
					{
						Id:        uuid.New(),
						Name:      "Test",
						Email:     "test@b.com",
						RoleId:    "a",
						CompanyId: uuid.New(),
					},
					{
						Id:        uuid.New(),
						Name:      "Test2",
						Email:     "test2@b.com",
						RoleId:    "a",
						CompanyId: uuid.New(),
					},
				}, nil)

				form := url.Values{}
				form.Add("email", user.Email)
				form.Add("password", user.Password)
				form.Add("name", user.Name)
				form.Add("role", user.RoleId)
				buf := bytes.NewBufferString(form.Encode())

				request, response := server.MakeRequest(http.MethodPost, "/bca/partials/users", buf)

				srv.UsersTable(response, request)

				assert.Equal(t, http.StatusOK, response.Code)
			})

			t.Run("Failed request", func(t *testing.T) {
				t.Run("Duplicate", func(t *testing.T) {
					srv, db := server.MakeServer()

					db.On("CreateUser", user).Return(types.User{}, errors.New("duplicate"))
					form := url.Values{}
					form.Add("email", user.Email)
					form.Add("password", user.Password)
					form.Add("name", user.Name)
					form.Add("role", user.RoleId)
					buf := bytes.NewBufferString(form.Encode())

					request, response := server.MakeRequest(http.MethodPost, "/bca/partials/users", buf)

					srv.UsersTable(response, request)

					assert.Equal(t, http.StatusConflict, response.Code)
					assert.Equal(t, response.Body.String(), "Usuario ya existe")
				})

				t.Run("Unknown", func(t *testing.T) {
					srv, db := server.MakeServer()

					db.On("CreateUser", user).Return(types.User{}, errors.New("unknown"))
					form := url.Values{}
					form.Add("email", user.Email)
					form.Add("password", user.Password)
					form.Add("name", user.Name)
					form.Add("role", user.RoleId)
					buf := bytes.NewBufferString(form.Encode())

					request, response := server.MakeRequest(http.MethodPost, "/bca/partials/users", buf)

					srv.UsersTable(response, request)

					assert.Equal(t, http.StatusInternalServerError, response.Code)
					assert.Equal(t, response.Body.String(), "Error al crear usuario")
				})
			})
		})
	})

	t.Run("PUT request", func(t *testing.T) {
		t.Run("Empty password", func(t *testing.T) {
			srv, _ := server.MakeServer()

			form := url.Values{}
			form.Add("password", "")
			buf := bytes.NewBufferString(form.Encode())

			request, response := server.MakeRequest(http.MethodPut, "/bca/partials/users", buf)

			srv.UsersTable(response, request)

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.Equal(t, response.Body.String(), "contraseña inválida")
		})

		t.Run("Successfull request", func(t *testing.T) {
			srv, db := server.MakeServer()

			db.On("UpdatePassword", "test", uuid.UUID{}, uuid.UUID{}).Return(types.User{}, nil)

			form := url.Values{}
			form.Add("password", "test")
			buf := bytes.NewBufferString(form.Encode())

			request, response := server.MakeRequest(http.MethodPut, "/bca/partials/users", buf)

			srv.UsersTable(response, request)

			assert.Equal(t, http.StatusOK, response.Code)
		})

		t.Run("Failed request", func(t *testing.T) {
			srv, db := server.MakeServer()

			db.On("UpdatePassword", "test", uuid.UUID{}, uuid.UUID{}).Return(types.User{}, errors.New("unknown"))

			form := url.Values{}
			form.Add("password", "test")
			buf := bytes.NewBufferString(form.Encode())

			request, response := server.MakeRequest(http.MethodPut, "/bca/partials/users", buf)

			srv.UsersTable(response, request)

			assert.Equal(t, http.StatusInternalServerError, response.Code)
		})
	})
}

func TestUserAdd(t *testing.T) {
	srv, _ := server.MakeServer()

	request, response := server.MakeRequest(http.MethodGet, "/bca/partials/users/add", nil)

	srv.UserAdd(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Agregar usuario")
}
