package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alcb1310/bca/internal/server"
	"github.com/alcb1310/bca/mocks"
)

func TestHelloWorldHandler(t *testing.T) {
	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret", -5)

	t.Run("should open the home page", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "Welcome")
	})
}
