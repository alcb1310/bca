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

func TestCreateSupplier(t *testing.T) {
	assert := assert.New(t)
	s := &server.Server{}
	s.DB = &DBMock{}

	server := httptest.NewServer(http.HandlerFunc(s.AllSuppliers))
	defer server.Close()

	resp, err := http.Post(server.URL, "application/json", nil)
	assert.Equal(err, nil, fmt.Sprintf("error making request to server. Err: %v", err))
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected := "{\"error\":\"EOF\"}"
	assert.Equal(expected, strings.TrimRight(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	supplier := &types.Supplier{}
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(supplier)
	resp, err = http.Post(server.URL, "application/json", buf)
	assert.Equal(err, nil, fmt.Sprintf("error making request to server. Err: %v", err))
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected = "{\"error\":\"supplier_id cannot be empty\",\"field\":\"supplier_id\"}"
	assert.Equal(expected, strings.TrimRight(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	supplier.SupplierId = "123"
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(supplier)
	resp, err = http.Post(server.URL, "application/json", buf)
	assert.Equal(err, nil, fmt.Sprintf("error making request to server. Err: %v", err))
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected = "{\"error\":\"name cannot be empty\",\"field\":\"name\"}"
	assert.Equal(expected, strings.TrimRight(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	supplier.Name = "test"
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(supplier)
	resp, err = http.Post(server.URL, "application/json", buf)
	assert.Equal(err, nil, fmt.Sprintf("error making request to server. Err: %v", err))
	assert.Equal(http.StatusCreated, resp.StatusCode, fmt.Sprintf("expected status created; got %v", resp.Status))

	e := "test"
	supplier.ContactEmail = &e
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(supplier)
	resp, err = http.Post(server.URL, "application/json", buf)
	assert.Equal(err, nil, fmt.Sprintf("error making request to server. Err: %v", err))
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected = "{\"error\":\"invalid email\",\"field\":\"contact_email\"}"
	assert.Equal(expected, strings.TrimRight(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	e = "test@test.com"
	supplier.ContactEmail = &e
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(supplier)
	resp, err = http.Post(server.URL, "application/json", buf)
	assert.Equal(err, nil, fmt.Sprintf("error making request to server. Err: %v", err))
	assert.Equal(http.StatusCreated, resp.StatusCode, fmt.Sprintf("expected status created; got %v", resp.Status))
}
