package tests

import (
	"bca-go-final/internal/server"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	assert := assert.New(t)
	s := &server.Server{}
	server := httptest.NewServer(http.HandlerFunc(s.HelloWorldHandler))
	defer server.Close()
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	// Assertions
	assert.Equal(http.StatusOK, resp.StatusCode, fmt.Sprintf("expected status OK; got %v", resp.Status))
	expected := "{\"message\":\"Hello World\"}"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	assert.Equal(expected, string(body), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))
}
