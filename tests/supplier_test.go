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

func TestCreateSupplier(t *testing.T) {
	s := &server.Server{}
	s.DB = &DBMock{}

	server := httptest.NewServer(http.HandlerFunc(s.AllSuppliers))
	defer server.Close()

	resp, err := http.Post(server.URL, "application/json", nil)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected := "{\"error\":\"EOF\"}"
	if expected != strings.TrimRight(string(body), "\n") {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}

	supplier := &types.Supplier{}
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(supplier)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status bad request; got %v", resp.Status)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected = "{\"error\":\"supplier_id cannot be empty\",\"field\":\"supplier_id\"}"
	if expected != strings.TrimRight(string(body), "\n") {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}

	supplier.SupplierId = "123"
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(supplier)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
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

	supplier.Name = "test"
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(supplier)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status created; got %v", resp.Status)
	}

	e := "test"
	supplier.ContactEmail = &e
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(supplier)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status bad request; got %v", resp.Status)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	expected = "{\"error\":\"invalid email\",\"field\":\"contact_email\"}"
	if expected != strings.TrimRight(string(body), "\n") {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}

	e = "test@test.com"
	supplier.ContactEmail = &e
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(supplier)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status created; got %v", resp.Status)
	}

}
