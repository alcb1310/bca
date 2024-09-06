package server_test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/alcb1310/bca/internal/server"
	"github.com/alcb1310/bca/mocks"
)

func TestPostClosure(t *testing.T) {
	testUrl := "/bca/transacciones/cierre"

	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret")
	token := createToken(s.TokenAuth)
	projectId := uuid.New()

	testData := []struct {
		name    string
		form    url.Values
		status  int
		body    []string
		closure *mocks.Service_CreateClosure_Call
	}{
		{
			name:    "should pass a form",
			form:    nil,
			status:  http.StatusBadRequest,
			body:    []string{""},
			closure: nil,
		},
		{
			name:    "should pass a project",
			form:    url.Values{},
			status:  http.StatusBadRequest,
			body:    []string{"Seleccione un proyecto"},
			closure: nil,
		},
		{
			name:    "should pass a valid project",
			form:    url.Values{"proyecto": []string{"invalid"}},
			status:  http.StatusBadRequest,
			body:    []string{"Proyecto inválido"},
			closure: nil,
		},
		{
			name:    "should pass a date",
			form:    url.Values{"proyecto": []string{uuid.New().String()}},
			status:  http.StatusBadRequest,
			body:    []string{"Ingrese una fecha"},
			closure: nil,
		},
		{
			name:    "should pass a valid date",
			form:    url.Values{"proyecto": []string{uuid.New().String()}, "date": []string{"invalid"}},
			status:  http.StatusBadRequest,
			body:    []string{"Ingrese una fecha válida"},
			closure: nil,
		},
		{
			name:    "should create a closure",
			form:    url.Values{"proyecto": []string{projectId.String()}, "date": []string{"2024-01-01"}},
			status:  http.StatusOK,
			body:    []string{""},
			closure: db.EXPECT().CreateClosure(uuid.UUID{}, projectId, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)).Return(nil),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			if tt.closure != nil {
				tt.closure.Times(1)
			}

			req, res := createRequest(token, http.MethodPost, testUrl, strings.NewReader(tt.form.Encode()))
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, tt.status, res.Code)

			for _, b := range tt.body {
				assert.Contains(t, res.Body.String(), b)
			}
		})
	}
}
