package server_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHome(t *testing.T) {
	s := mount()

	rr := executeRequest(t, s, http.MethodGet, "/", nil)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Home Page")
}
