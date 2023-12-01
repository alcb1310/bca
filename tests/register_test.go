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

func TestRegistrationRouteValidation(t *testing.T) {
	s := &server.Server{}
	server := httptest.NewServer(http.HandlerFunc(s.Register))
	defer server.Close()

	resp, err := http.Post(server.URL, "application/json", nil)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

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

	c := &types.CompanyCreate{}

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(c)
	resp, err = http.Post(server.URL, "application/json", buf)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status bad request; got %v", resp.Status)
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected = "{\"error\":\"name cannot be empty\",\"field\":\"name\"}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}

	c.Name = "test"
	err = json.NewEncoder(buf).Encode(c)
	resp, err = http.Post(server.URL, "application/json", buf)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status bad request; got %v", resp.Status)
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected = "{\"error\":\"ruc cannot be empty\",\"field\":\"ruc\"}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}

	c.Ruc = "123"
	err = json.NewEncoder(buf).Encode(c)
	resp, err = http.Post(server.URL, "application/json", buf)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status bad request; got %v", resp.Status)
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected = "{\"error\":\"should pass at least one employee\",\"field\":\"employees\"}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}

	type BadRequest struct {
		Name      string `json:"name"`
		Ruc       string `json:"ruc"`
		Employees string `json:"employees"`
	}

	b := &BadRequest{
		Name:      "test",
		Ruc:       "123",
		Employees: "bad request",
	}
	err = json.NewEncoder(buf).Encode(b)
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
	expected = "{\"error\":\"employees must be a number\",\"field\":\"employees\"}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}

	c.Employees = 1
	if err := json.NewEncoder(buf).Encode(&c); err != nil {
		t.Fatalf("error marshaling. Err: %v", err)
	}

	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err.Error())
	}
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

	c.Email = "test"
	if err := json.NewEncoder(buf).Encode(&c); err != nil {
		t.Fatalf("error marshaling. Err: %v", err)
	}
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

	c.Email = "LJGQ6@example.com"
	if err := json.NewEncoder(buf).Encode(&c); err != nil {
		t.Fatalf("error marshaling. Err: %v", err)
	}
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err.Error())
	}
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

	c.Password = "test"
	if err := json.NewEncoder(buf).Encode(&c); err != nil {
		t.Fatalf("error marshaling. Err: %v", err)
	}
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err.Error())
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status bad request; got %v", resp.Status)
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected = "{\"error\":\"name of the user cannot be empty\",\"field\":\"user\"}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}
