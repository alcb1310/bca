package server_test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
	"bca-go-final/mocks"
)

func TestCreateUser(t *testing.T) {
	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret")
	token := createToken(s.TokenAuth)

	t.Run("should pass a form", func(t *testing.T) {
		req, res := createRequest(token, http.MethodPost, "/bca/partials/users", nil)
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("should pass an email", func(t *testing.T) {
		form := url.Values{}
		req, res := createRequest(token, http.MethodPost, "/bca/partials/users", strings.NewReader(form.Encode()))
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("should pass a valid email", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "test")
		req, res := createRequest(token, http.MethodPost, "/bca/partials/users", strings.NewReader(form.Encode()))
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("should pass a password", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "test@test.com")
		req, res := createRequest(token, http.MethodPost, "/bca/partials/users", strings.NewReader(form.Encode()))
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("should pass a name", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "test@test.com")
		form.Add("password", "test")
		req, res := createRequest(token, http.MethodPost, "/bca/partials/users", strings.NewReader(form.Encode()))
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("should create a new user", func(t *testing.T) {
		db.EXPECT().CreateUser(&types.UserCreate{
			Email:     "test@test.com",
			Password:  "test",
			Name:      "test",
			CompanyId: uuid.UUID{},
			RoleId:    "a",
		}).Return(types.User{
			Id:        uuid.UUID{},
			Email:     "test@test.com",
			Name:      "test",
			CompanyId: uuid.UUID{},
			RoleId:    "a",
		}, nil)
		db.EXPECT().GetAllUsers(uuid.UUID{}).Return([]types.User{
			{
				Id:        uuid.UUID{},
				Email:     "test@test.com",
				Name:      "test",
				CompanyId: uuid.UUID{},
				RoleId:    "a",
			},
		}, nil)
		form := url.Values{}
		form.Add("email", "test@test.com")
		form.Add("password", "test")
		form.Add("name", "test")
		req, res := createRequest(token, http.MethodPost, "/bca/partials/users", strings.NewReader(form.Encode()))
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "test")
		assert.Contains(t, res.Body.String(), "test@test.com")
		assert.Contains(t, res.Body.String(), "<table>")
	})
}
