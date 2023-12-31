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

func TestGetAllUsers(t *testing.T) {
	assert := assert.New(t)
	s := &server.Server{}
	s.DB = &DBMock{}

	server := httptest.NewServer(http.HandlerFunc(s.AllUsers))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(http.StatusOK, resp.StatusCode, fmt.Sprintf("expected status ok; got %v", resp.Status))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	r := []types.User{}
	json.NewDecoder(strings.NewReader(string(body))).Decode(&r)
	assert.Equal(len(r), 2, fmt.Sprintf("expected response body to be 2; got %v", len(body)))
}

func TestCreateUser(t *testing.T) {
	assert := assert.New(t)
	s := &server.Server{}
	s.DB = &DBMock{}

	server := httptest.NewServer(http.HandlerFunc(s.AllUsers))
	defer server.Close()

	c := &types.UserCreate{}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(c)
	resp, err := http.Post(server.URL, "application/json", nil)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected := "{\"error\":\"EOF\"}"
	assert.Equal(expected, strings.Trim(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(c)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected = "{\"error\":\"email cannot be empty\",\"field\":\"email\"}"
	assert.Equal(expected, strings.Trim(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	c.Email = "test"
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(c)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected = "{\"error\":\"invalid email\",\"field\":\"email\"}"
	assert.Equal(expected, strings.Trim(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	c.Email = "test@test.com"
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(c)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected = "{\"error\":\"password cannot be empty\",\"field\":\"password\"}"
	assert.Equal(expected, strings.Trim(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	c.Password = "test"
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(c)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected = "{\"error\":\"name cannot be empty\",\"field\":\"name\"}"
	assert.Equal(expected, strings.Trim(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	c.Name = "test"
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(c)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected = "{\"error\":\"role cannot be empty\",\"field\":\"role\"}"
	assert.Equal(expected, strings.Trim(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	c.RoleId = "a"
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(c)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(http.StatusCreated, resp.StatusCode, fmt.Sprintf("expected status created; got %v", resp.Status))
}
