package tests

import (
	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateProject(t *testing.T) {
	assert := assert.New(t)
	s := &server.Server{}
	s.DB = &DBMock{}

	server := httptest.NewServer(http.HandlerFunc(s.AllProjects))
	defer server.Close()

	resp, err := http.Post(server.URL, "application/json", nil)
	assert.Equal(err, nil, fmt.Sprintf("error making request to server. Err: %v", err))
	defer resp.Body.Close()
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err := io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected := "{\"error\":\"EOF\"}"
	assert.Equal(expected, strings.TrimRight(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	p := &types.Project{}
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(p)
	resp, err = http.Post(server.URL, "application/json", buf)
	assert.Equal(err, nil, fmt.Sprintf("error making request to server. Err: %v", err))
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected = "{\"error\":\"name cannot be empty\",\"field\":\"name\"}"
	assert.Equal(expected, strings.TrimRight(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	p.Name = "test"
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(p)
	resp, err = http.Post(server.URL, "application/json", buf)
	assert.Equal(err, nil, fmt.Sprintf("error making request to server. Err: %v", err))
	assert.Equal(http.StatusCreated, resp.StatusCode, fmt.Sprintf("expected status created; got %v", resp.Status))
}
