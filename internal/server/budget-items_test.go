package server_test

import (
	"database/sql"
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

func TestCreateBudgetItem(t *testing.T) {
	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret")
	token := createToken(s.TokenAuth)

	testData := []struct {
		name             string
		form             url.Values
		status           int
		body             []string
		createBudgetItem *mocks.Service_CreateBudgetItem_Call
		getBudgetItems   *mocks.Service_GetBudgetItems_Call
	}{
		{
			name:             "should pass a form",
			form:             nil,
			status:           http.StatusBadRequest,
			body:             []string{},
			createBudgetItem: nil,
			getBudgetItems:   nil,
		},
		{
			name:             "should provide a budget item code",
			form:             url.Values{},
			status:           http.StatusBadRequest,
			body:             []string{"Debe proporcionar un código de la partida"},
			createBudgetItem: nil,
			getBudgetItems:   nil,
		},
		{
			name:             "should provide a budget item name",
			form:             url.Values{"code": {"1234"}},
			status:           http.StatusBadRequest,
			body:             []string{"Debe proporcionar un nombre de la partida"},
			createBudgetItem: nil,
			getBudgetItems:   nil,
		},
		{
			name:             "should provide a valid id for parent",
			form:             url.Values{"parent": {"1234"}},
			status:           http.StatusBadRequest,
			body:             []string{"Código de la partida padre es inválido"},
			createBudgetItem: nil,
			getBudgetItems:   nil,
		},
		{
			name:   "should create a budget item",
			form:   url.Values{"code": {"1234"}, "name": {"test"}, "accumulate": {"false"}},
			status: http.StatusOK,
			body:   []string{},
			createBudgetItem: db.EXPECT().CreateBudgetItem(&types.BudgetItem{
				Code:       "1234",
				Name:       "test",
				Accumulate: sql.NullBool{Bool: false, Valid: true},
			}).Return(nil),
			getBudgetItems: db.EXPECT().GetBudgetItems(uuid.UUID{}, "").Return([]types.BudgetItemResponse{
				{
					ID:         uuid.New(),
					Code:       "1234",
					Name:       "test",
					Accumulate: sql.NullBool{Bool: false, Valid: true},
					CompanyId:  uuid.UUID{},
					ParentId:   uuid.NullUUID{},
				},
			}, nil),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createBudgetItem != nil {
				tt.createBudgetItem.Times(1)
			}

			if tt.getBudgetItems != nil {
				tt.getBudgetItems.Times(1)
			}

			req, res := createRequest(
				token,
				http.MethodPost,
				"/bca/partials/budget-item",
				strings.NewReader(tt.form.Encode()),
			)
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, tt.status, res.Code)

			for _, b := range tt.body {
				assert.Contains(t, res.Body.String(), b)
			}
		})
	}
}

func TestUpdateBudgetItem(t *testing.T) {
	id := uuid.New()
	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret")
	token := createToken(s.TokenAuth)

	testData := []struct {
		name             string
		form             url.Values
		status           int
		body             []string
		updateBudgetItem *mocks.Service_UpdateBudgetItem_Call
		getBudgetItems   *mocks.Service_GetBudgetItems_Call
	}{
		{
			name:             "should provide a valid id for parent",
			form:             url.Values{"parent": {"1234"}},
			status:           http.StatusBadRequest,
			body:             []string{"Código de la partida padre es inválido"},
			updateBudgetItem: nil,
			getBudgetItems:   nil,
		},
		{
			name:   "should update a budget item",
			form:   url.Values{},
			status: http.StatusOK,
			body:   []string{},
			updateBudgetItem: db.EXPECT().UpdateBudgetItem(&types.BudgetItem{
				ID:         id,
				Code:       "1234",
				Name:       "test",
				Accumulate: sql.NullBool{Bool: false, Valid: true},
			}).Return(nil),
			getBudgetItems: db.EXPECT().GetBudgetItems(uuid.UUID{}, "").Return([]types.BudgetItemResponse{
				{
					ID:         id,
					Code:       "1234",
					Name:       "test",
					Accumulate: sql.NullBool{Bool: false, Valid: true},
					CompanyId:  uuid.UUID{},
					ParentId:   uuid.NullUUID{},
				},
			}, nil),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			db.EXPECT().GetOneBudgetItem(id, uuid.UUID{}).Return(&types.BudgetItem{
				ID:         id,
				Code:       "1234",
				Name:       "test",
				Accumulate: sql.NullBool{Bool: false, Valid: true},
				CompanyId:  uuid.UUID{},
				ParentId:   nil,
			}, nil).Times(1)

			if tt.updateBudgetItem != nil {
				tt.updateBudgetItem.Times(1)
			}

			if tt.getBudgetItems != nil {
				tt.getBudgetItems.Times(1)
			}

			req, res := createRequest(
				token,
				http.MethodPut,
				fmt.Sprintf("/bca/partials/budget-item/%s", id),
				strings.NewReader(tt.form.Encode()),
			)
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, tt.status, res.Code)
			for _, b := range tt.body {
				assert.Contains(t, res.Body.String(), b)
			}
		})
	}
}
