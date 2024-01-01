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
	supplierId uuid.UUID = uuid.New()
	projectId  uuid.UUID = uuid.New()
)

func TestCreateInvoice(t *testing.T) {
	assert := assert.New(t)
	s := &server.Server{}
	s.DB = &DBMock{}

	server := httptest.NewServer(http.HandlerFunc(s.AllInvoices))
	defer server.Close()

	resp, err := http.Post(server.URL, "application/json", nil)
	assert.Equal(err, nil, fmt.Sprintf("error making request to server. Err: %v", err))
	defer resp.Body.Close()
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err := io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected := "{\"error\":\"EOF\"}"
	assert.Equal(expected, strings.TrimRight(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	invoice := &types.InvoiceCreate{}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(invoice)
	resp, err = http.Post(server.URL, "application/json", buf)
	assert.Equal(err, nil, fmt.Sprintf("error making request to server. Err: %v", err))
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected = "{\"error\":\"supplier_id cannot be empty\",\"field\":\"supplier_id\"}"
	assert.Equal(expected, strings.TrimRight(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	x := make(map[string]string)
	x["supplier_id"] = "test"
	buf = new(bytes.Buffer)
	json.NewEncoder(buf).Encode(x)
	resp, err = http.Post(server.URL, "application/json", buf)
	assert.Equal(err, nil, fmt.Sprintf("error making request to server. Err: %v", err))
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected = "invalid UUID"
	assert.Contains(string(body), expected, fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	invoice.SupplierId = &supplierId
	buf = new(bytes.Buffer)
	json.NewEncoder(buf).Encode(invoice)
	resp, err = http.Post(server.URL, "application/json", buf)
	assert.Equal(err, nil, fmt.Sprintf("error making request to server. Err: %v", err))
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected = "{\"error\":\"project_id cannot be empty\",\"field\":\"project_id\"}"
	assert.Equal(expected, strings.TrimRight(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	x["supplier_id"] = supplierId.String()
	x["project_id"] = "test"
	buf = new(bytes.Buffer)
	json.NewEncoder(buf).Encode(x)
	resp, err = http.Post(server.URL, "application/json", buf)
	assert.Equal(err, nil, fmt.Sprintf("error making request to server. Err: %v", err))
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected = "invalid UUID"
	assert.Contains(string(body), expected, fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	invoice.ProjectId = &projectId
	buf = new(bytes.Buffer)
	json.NewEncoder(buf).Encode(invoice)
	resp, err = http.Post(server.URL, "application/json", buf)
	assert.Equal(err, nil, fmt.Sprintf("error making request to server. Err: %v", err))
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected = "{\"error\":\"invoice_number cannot be empty\",\"field\":\"invoice_number\"}"
	assert.Equal(expected, strings.TrimRight(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	in := ""
	invoice.InvoiceNumber = &in
	buf = new(bytes.Buffer)
	json.NewEncoder(buf).Encode(invoice)
	resp, err = http.Post(server.URL, "application/json", buf)
	assert.Equal(err, nil, fmt.Sprintf("error making request to server. Err: %v", err))
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected = "{\"error\":\"invoice_number cannot be empty\",\"field\":\"invoice_number\"}"
	assert.Equal(expected, strings.TrimRight(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	in = "test"
	invoice.InvoiceNumber = &in
	buf = new(bytes.Buffer)
	json.NewEncoder(buf).Encode(invoice)
	resp, err = http.Post(server.URL, "application/json", buf)
	assert.Equal(err, nil, fmt.Sprintf("error making request to server. Err: %v", err))
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected = "{\"error\":\"invoice_date cannot be empty\",\"field\":\"invoice_date\"}"
	assert.Equal(expected, strings.TrimRight(string(body), "\n"), fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	x["project_id"] = projectId.String()
	x["invoice_number"] = "test"
	x["invoice_date"] = "test"
	buf = new(bytes.Buffer)
	json.NewEncoder(buf).Encode(x)
	resp, err = http.Post(server.URL, "application/json", buf)
	assert.Equal(err, nil, fmt.Sprintf("error making request to server. Err: %v", err))
	assert.Equal(http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected status bad request; got %v", resp.Status))
	body, err = io.ReadAll(resp.Body)
	assert.Equal(err, nil, fmt.Sprintf("error reading response body. Err: %v", err))
	expected = "{\"error\":\"parsing time"
	assert.Contains(string(body), expected, fmt.Sprintf("expected response body to be %v; got %v", expected, string(body)))

	x["invoice_date"] = "2023-12-31T23:59:59Z"
	buf = new(bytes.Buffer)
	json.NewEncoder(buf).Encode(x)
	resp, err = http.Post(server.URL, "application/json", buf)
	assert.Equal(err, nil, fmt.Sprintf("error making request to server. Err: %v", err))
	assert.Equal(http.StatusCreated, resp.StatusCode, fmt.Sprintf("expected status created; got %v", resp.Status))
}
