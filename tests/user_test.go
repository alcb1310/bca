package tests

import (
	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetAllUsers(t *testing.T) {
	s := &server.Server{}
	s.DB = &DBMock{}

	server := httptest.NewServer(http.HandlerFunc(s.GetAllUsers))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status ok; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	r := []types.User{}
	json.NewDecoder(strings.NewReader(string(body))).Decode(&r)

	if len(r) != 2 {
		t.Errorf("expected response body to be 2; got %v", len(body))
	}
}

func TestCreateUser(t *testing.T) {
	s := &server.Server{}
	s.DB = &DBMock{}

	server := httptest.NewServer(http.HandlerFunc(s.GetAllUsers))
	defer server.Close()

	c := &types.UserCreate{}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(c)
	resp, err := http.Post(server.URL, "application/json", nil)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status bad request; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected := "{\"error\":\"EOF\"}"
	if expected != strings.TrimRight(string(body), "\n") {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}

	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(c)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status bad request; got %v", resp.Status)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected = "{\"error\":\"email cannot be empty\",\"field\":\"email\"}"
	if expected != strings.TrimRight(string(body), "\n") {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}

	c.Email = "test"
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(c)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status bad request; got %v", resp.Status)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected = "{\"error\":\"invalid email\",\"field\":\"email\"}"
	if expected != strings.TrimRight(string(body), "\n") {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}

	c.Email = "test@test.com"
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(c)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status bad request; got %v", resp.Status)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected = "{\"error\":\"password cannot be empty\",\"field\":\"password\"}"
	if expected != strings.TrimRight(string(body), "\n") {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}

	c.Password = "test"
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(c)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status bad request; got %v", resp.Status)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected = "{\"error\":\"name cannot be empty\",\"field\":\"name\"}"
	if expected != strings.TrimRight(string(body), "\n") {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}

	c.Name = "test"
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(c)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status bad request; got %v", resp.Status)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected = "{\"error\":\"role cannot be empty\",\"field\":\"role\"}"
	if expected != strings.TrimRight(string(body), "\n") {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}

	c.RoleId = "a"
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(c)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status created; got %v", resp.Status)
	}
}
