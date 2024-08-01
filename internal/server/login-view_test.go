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

func TestLoginView(t *testing.T) {
	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret")

	t.Run("should validate login form", func(t *testing.T) {
		t.Run("must pass a form", func(t *testing.T) {
      req, res := createRequest("", "POST", "/login", nil)
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, http.StatusBadRequest, res.Code)
			assert.Contains(t, res.Body.String(), "missing form body")
		})

		t.Run("must pass an email", func(t *testing.T) {
			form := url.Values{}
      req, res := createRequest("", "POST", "/login", strings.NewReader(form.Encode()))
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, http.StatusBadRequest, res.Code)
			assert.Contains(t, res.Body.String(), "credenciales inválidas")
		})

		t.Run("must pass a valid email", func(t *testing.T) {
			form := url.Values{}
			form.Add("email", "test")
      req, res := createRequest("", "POST", "/login", strings.NewReader(form.Encode()))
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, http.StatusBadRequest, res.Code)
			assert.Contains(t, res.Body.String(), "credenciales inválidas")
		})

		t.Run("must pass a password", func(t *testing.T) {
			form := url.Values{}
			form.Add("email", "test@test.com")
      req, res := createRequest("", "POST", "/login", strings.NewReader(form.Encode()))
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, http.StatusBadRequest, res.Code)
			assert.Contains(t, res.Body.String(), "credenciales inválidas")
		})

		t.Run("should not login on invalid credentials", func(t *testing.T) {
			db.EXPECT().Login(&types.Login{Email: "test@test.com", Password: "test"}).Return("", &types.User{
				Id:        uuid.UUID{},
				Email:     "test@test.com",
				Name:      "test",
				CompanyId: uuid.UUID{},
				RoleId:    "a",
			}, nil)
			form := url.Values{}
			form.Add("email", "test@test.com")
			form.Add("password", "test")
      req, res := createRequest("", "POST", "/login", strings.NewReader(form.Encode()))
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, http.StatusSeeOther, res.Code)
		})
	})
}
