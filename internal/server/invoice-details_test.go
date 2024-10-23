package server_test

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/alcb1310/bca/internal/server"
	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/mocks"
)

func TestCreateInvoiceDetails(t *testing.T) {
	invoiceId := uuid.New()
	testUrl := fmt.Sprintf("/bca/partials/invoices/%s/details", invoiceId.String())
	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret", -5)
	token := createToken(s.TokenAuth)
	itemId := uuid.New()

	testData := []struct {
		name          string
		form          url.Values
		status        int
		body          []string
		addDetail     *mocks.Service_AddDetail_Call
		getAllDetails *mocks.Service_GetAllDetails_Call
	}{
		{
			name:          "should pass a form",
			form:          nil,
			status:        http.StatusBadRequest,
			body:          []string{},
			addDetail:     nil,
			getAllDetails: nil,
		},
		{
			name:          "should pass an item id",
			form:          url.Values{},
			status:        http.StatusBadRequest,
			body:          []string{"Ingrese una partida"},
			addDetail:     nil,
			getAllDetails: nil,
		},
		{
			name: "should pass a valid item id",
			form: url.Values{
				"item": {"test"},
			},
			status:        http.StatusBadRequest,
			body:          []string{"Ingrese una partida válida"},
			addDetail:     nil,
			getAllDetails: nil,
		},
		{
			name: "should pass a quantity",
			form: url.Values{
				"item": {itemId.String()},
			},
			status:        http.StatusBadRequest,
			body:          []string{"Ingrese una cantidad"},
			addDetail:     nil,
			getAllDetails: nil,
		},
		{
			name: "should pass a valid quantity",
			form: url.Values{
				"item":     {itemId.String()},
				"quantity": {"test"},
			},
			status:        http.StatusBadRequest,
			body:          []string{"Cantidad debe ser un número válido"},
			addDetail:     nil,
			getAllDetails: nil,
		},
		{
			name: "should pass a cost",
			form: url.Values{
				"item":     {itemId.String()},
				"quantity": {"1"},
			},
			status:        http.StatusBadRequest,
			body:          []string{"Ingrese un costo"},
			addDetail:     nil,
			getAllDetails: nil,
		},
		{
			name: "should pass a valid cost",
			form: url.Values{
				"item":     {itemId.String()},
				"quantity": {"1"},
				"cost":     {"test"},
			},
			status:        http.StatusBadRequest,
			body:          []string{"Costo debe ser un número válido"},
			addDetail:     nil,
			getAllDetails: nil,
		},
		{
			name: "should create a detail",
			form: url.Values{
				"item":     {itemId.String()},
				"quantity": {"1"},
				"cost":     {"1"},
			},
			status: http.StatusOK,
			body:   []string{},
			addDetail: db.EXPECT().AddDetail(types.InvoiceDetailCreate{
				InvoiceId:    invoiceId,
				BudgetItemId: itemId,
				CompanyId:    uuid.UUID{},
				Quantity:     1,
				Cost:         1,
				Total:        1,
			}).Return(nil),
			getAllDetails: db.EXPECT().GetAllDetails(invoiceId, uuid.UUID{}).Return([]types.InvoiceDetailsResponse{}, nil),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			if tt.addDetail != nil {
				tt.addDetail.Times(1)
			}

			if tt.getAllDetails != nil {
				tt.getAllDetails.Times(1)
			}

			req, res := createRequest(token, http.MethodPost, testUrl, strings.NewReader(tt.form.Encode()))
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, tt.status, res.Code)

			for _, b := range tt.body {
				assert.Contains(t, res.Body.String(), b)
			}
		})
	}
}
