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

func TestCreateProject(t *testing.T) {
	s := &server.Server{}
	s.DB = &DBMock{}

	server := httptest.NewServer(http.HandlerFunc(s.AllProjects))
	defer server.Close()

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

	p := &types.Project{}
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(p)
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

	p.Name = "test"
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(p)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status ok; got %v", resp.Status)
	}

}
