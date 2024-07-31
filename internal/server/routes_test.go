package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
	"bca-go-final/mocks"
)

func TestHelloWorldHandler(t *testing.T) {
	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret")

	t.Run("should open the home page", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "Welcome")
	})
}
