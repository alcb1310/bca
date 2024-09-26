package server_test

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/alcb1310/bca/internal/server"
	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/mocks"
)

func TestLoginView(t *testing.T) {
	reqUrl := "/login"
	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret")
	testData := []struct {
		name   string
		form   io.Reader
		status int
		body   []string
		mock   *mocks.Service_Login_Call
	}{
		{
			name:   "should pass a form",
			form:   nil,
			status: http.StatusBadRequest,
			body:   []string{"missing form body"},
			mock:   nil,
		},
		{
			name:   "should pass an email",
			form:   strings.NewReader(url.Values{}.Encode()),
			status: http.StatusBadRequest,
			body:   []string{"credenciales inválidas"},
			mock:   nil,
		},
		{
			name:   "should pass a valid email",
			form:   strings.NewReader(url.Values{"email": {"test"}}.Encode()),
			status: http.StatusBadRequest,
			body:   []string{"credenciales inválidas"},
			mock:   nil,
		},
		{
			name:   "should pass a valid password",
			form:   strings.NewReader(url.Values{"email": {"test@test.com"}}.Encode()),
			status: http.StatusBadRequest,
			body:   []string{"credenciales inválidas"},
			mock:   nil,
		},
		{
			name:   "should login on valid credentials",
			form:   strings.NewReader(url.Values{"email": {"test@test.com"}, "password": {"test"}}.Encode()),
			status: http.StatusSeeOther,
			body:   []string{""},
			mock: db.EXPECT().Login(&types.Login{Email: "test@test.com", Password: "test"}).Return("", &types.User{
				Id:        uuid.UUID{},
				Email:     "test@test.com",
				Name:      "test",
				CompanyId: uuid.UUID{},
				RoleId:    "a",
			}, nil),
		},
	}

	for _, d := range testData {
		t.Run(d.name, func(t *testing.T) {
			if d.mock != nil {
				d.mock.Times(1)
			}
			req, res := createRequest("", http.MethodPost, reqUrl, d.form)
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, d.status, res.Code)
			for _, b := range d.body {
				assert.Contains(t, res.Body.String(), b)
			}
		})
	}
}
