package tests

import (
	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginRouteValidation(t *testing.T) {
	s := &server.Server{}
	s.DB = &DBMock{}
	server := httptest.NewServer(http.HandlerFunc(s.Login))
	defer server.Close()

	resp, err := http.Post(server.URL, "application/json", nil)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}

	// Assertions
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status bad request; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	expected := "{\"error\":\"EOF\"}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}

	login := &types.Login{}

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(login)
	resp, err = http.Post(server.URL, "application/json", buf)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status bad request; got %v", resp.Status)
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected = "{\"error\":\"email cannot be empty\",\"field\":\"email\"}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}

	login.Email = "test"
	err = json.NewEncoder(buf).Encode(login)
	resp, err = http.Post(server.URL, "application/json", buf)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status bad request; got %v", resp.Status)
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected = "{\"error\":\"invalid email\",\"field\":\"email\"}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}

	login.Email = "test@test.com"
	err = json.NewEncoder(buf).Encode(login)
	resp, err = http.Post(server.URL, "application/json", buf)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status bad request; got %v", resp.Status)
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected = "{\"error\":\"password cannot be empty\",\"field\":\"password\"}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}

	login.Password = "test"
	err = json.NewEncoder(buf).Encode(login)
	resp, err = http.Post(server.URL, "application/json", buf)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status ok; got %v", resp.Status)
	}
}

func TestLoginSuccess(t *testing.T) {
	s := &server.Server{}
	s.DB = &DBMock{}
	server := httptest.NewServer(http.HandlerFunc(s.Login))
	buf := new(bytes.Buffer)

	defer server.Close()

	login := &types.Login{}
	login.Email = "test@test.com"
	login.Password = "test"

	_ = json.NewEncoder(buf).Encode(login)
	resp, _ := http.Post(server.URL, "application/json", buf)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status ok; got %v", resp.Status)
	}
}
