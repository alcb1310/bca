package tests

import (
	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
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
