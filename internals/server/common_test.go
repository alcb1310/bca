package server_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alcb1310/bca/internals/server"
	"github.com/alcb1310/bca/mocks"
)

// mount a test server
func mount() *server.Service {
	db := &mocks.DatabaseService{}
	return server.New(&slog.Logger{}, db)
}

// executeRequest executes a request against a test server
func executeRequest(t *testing.T, s *server.Service, method, url string, body io.Reader) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	assert.NoError(t, err)
	s.Router.ServeHTTP(rr, req)

	return rr
}
