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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	project  uuid.UUID = uuid.New()
	budget   uuid.UUID = uuid.New()
	quantity float64   = 100
	cost     float64   = 500
)

func TestBudget(t *testing.T) {
	assert := assert.New(t)
	s := &server.Server{}
	s.DB = &DBMock{}

	server := httptest.NewServer(http.HandlerFunc(s.AllBudgets))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(http.StatusOK, resp.StatusCode, fmt.Sprintf("expected status ok; got %v", resp.Status))
}

func TestCreateBudget(t *testing.T) {
	assert := assert.New(t)
	s := &server.Server{}
	s.DB = &DBMock{}

	server := httptest.NewServer(http.HandlerFunc(s.AllBudgets))
	defer server.Close()

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

	b := &types.CreateBudget{}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(b)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected = "{\"error\":\"project_id cannot be empty\",\"field\":\"project_id\"}"
	assert.Equal(expected, strings.TrimRight(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	x := make(map[string]string)
	x["project_id"] = "test"
	json.NewEncoder(buf).Encode(x)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected = "invalid UUID"
	assert.Contains(string(body), expected, fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	b.ProjectId = project
	json.NewEncoder(buf).Encode(b)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected = "{\"error\":\"budget_item_id cannot be empty\",\"field\":\"budget_item_id\"}"
	assert.Equal(expected, strings.TrimRight(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	x["project_id"] = project.String()
	x["budget_item_id"] = "test"
	json.NewEncoder(buf).Encode(x)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected = "invalid UUID"
	assert.Contains(string(body), expected, fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	b.BudgetItemId = budget
	json.NewEncoder(buf).Encode(b)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected = "{\"error\":\"quantity cannot be empty\",\"field\":\"quantity\"}"
	assert.Equal(expected, strings.TrimRight(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	x["budget_item_id"] = budget.String()
	x["quantity"] = "test"
	json.NewEncoder(buf).Encode(x)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected = "cannot unmarshal"
	assert.Contains(string(body), expected, fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	b.Quantity = &quantity
	json.NewEncoder(buf).Encode(b)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected = "{\"error\":\"cost cannot be empty\",\"field\":\"cost\"}"
	assert.Equal(expected, strings.TrimRight(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	b.Cost = &cost
	json.NewEncoder(buf).Encode(b)
	resp, err = http.Post(server.URL, "application/json", buf)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	assert.Equal(http.StatusCreated, resp.StatusCode, fmt.Sprintf("expected status created; got %v", resp.Status))
}
