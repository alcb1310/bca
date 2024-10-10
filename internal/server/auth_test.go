package server_test

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/alcb1310/bca/internal/server"
	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/mocks"
)

func TestApiLogin(t *testing.T) {
	reqUrl := "/api/v1/login"
	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret")
	userId := uuid.New()

	testData := []struct {
		name   string
		form   map[string]interface{}
		status int
		body   map[string]string
		mock   *mocks.Service_Login_Call
	}{
		{
			name:   "should pass a form",
			form:   nil,
			status: http.StatusBadRequest,
			body:   map[string]string{"error": "credenciales inválidas"},
			mock:   nil,
		},
		{
			name:   "should pass an email",
			form:   map[string]interface{}{},
			status: http.StatusBadRequest,
			body:   map[string]string{"error": "credenciales inválidas"},
			mock:   nil,
		},
		{
			name:   "should pass a password",
			form:   map[string]interface{}{"email": "a@b.c"},
			status: http.StatusBadRequest,
			body:   map[string]string{"error": "credenciales inválidas"},
			mock:   nil,
		},
		{
			name: "should pass a valid email",
			form: map[string]interface{}{
				"email":    "abc",
				"password": "test",
			},
			status: http.StatusBadRequest,
			body:   map[string]string{"error": "credenciales inválidas"},
			mock:   nil,
		},
		{
			name: "should pass a valid password",
			form: map[string]interface{}{
				"email":    "a@b.c",
				"password": "test",
			},
			status: http.StatusOK,
			body: map[string]string{
				"Id":         userId.String(),
				"Email":      "a@b.c",
				"Name":       "test",
				"CompanyId": uuid.UUID{}.String(),
				"RoleId":    "a",
			},
			mock: db.EXPECT().Login(&types.Login{Email: "a@b.c", Password: "test"}).Return("", &types.User{
				Id:        userId,
				Email:     "a@b.c",
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

			var reader io.Reader
			if d.form == nil {
				reader = nil
			} else {
				jsonData, _ := json.Marshal(d.form)
				reader = strings.NewReader(string(jsonData))
			}

			req, res := createApiRequest("", http.MethodPost, reqUrl, reader)
			s.Router.ServeHTTP(res, req)

			assert.Equal(t, res.Header().Get("Content-Type"), "application/json")
			assert.Equal(t, d.status, res.Code)
			slog.Debug("TestApiLogin", "Auth Token", res.Header().Get("BCA-Auth-Token"))

			var jsonResp map[string]string
			err := json.NewDecoder(res.Body).Decode(&jsonResp)
			assert.NoError(t, err)
			assert.Equal(t, jsonResp, d.body)
		})
	}
}
