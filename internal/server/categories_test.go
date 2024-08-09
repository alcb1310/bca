package server_test

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
	"bca-go-final/mocks"
)

func TestCreateCategory(t *testing.T) {
	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret")
	token := createToken(s.TokenAuth)

	testData := []struct {
		name             string
		form             url.Values
		status           int
		body             []string
		createCategory   *mocks.Service_CreateCategory_Call
		getAllCategories *mocks.Service_GetAllCategories_Call
	}{
		{
			name:             "should pass a form",
			form:             nil,
			status:           http.StatusBadRequest,
			body:             []string{},
			createCategory:   nil,
			getAllCategories: nil,
		},
		{
			name:             "should pass a name",
			form:             url.Values{},
			status:           http.StatusBadRequest,
			body:             []string{"Ingrese un nombre de categoría"},
			createCategory:   nil,
			getAllCategories: nil,
		},
		{
			name:             "should create a category",
			form:             url.Values{"name": {"test"}},
			status:           http.StatusOK,
			body:             []string{},
			createCategory:   db.EXPECT().CreateCategory(types.Category{Name: "test", CompanyId: uuid.UUID{}}).Return(nil),
			getAllCategories: db.EXPECT().GetAllCategories(uuid.UUID{}).Return([]types.Category{}, nil),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createCategory != nil {
				tt.createCategory.Times(1)
			}
			if tt.getAllCategories != nil {
				tt.getAllCategories.Times(1)
			}

			req, res := createRequest(token, http.MethodPost, "/bca/partials/categories", strings.NewReader(tt.form.Encode()))
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, tt.status, res.Code)

			for _, b := range tt.body {
				assert.Contains(t, res.Body.String(), b)
			}
		})
	}
}

func TestEditCategory(t *testing.T) {
	categoryId := uuid.New()
	testUrl := fmt.Sprintf("/bca/partials/categories/%s", categoryId.String())

	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret")
	token := createToken(s.TokenAuth)

	testData := []struct {
		name             string
		form             url.Values
		status           int
		body             []string
		updateCategory   *mocks.Service_UpdateCategory_Call
		getAllCategories *mocks.Service_GetAllCategories_Call
	}{
		{
			name:             "update a category",
			form:             url.Values{},
			status:           http.StatusOK,
			body:             []string{},
			updateCategory:   db.EXPECT().UpdateCategory(types.Category{Id: categoryId, Name: "test", CompanyId: uuid.UUID{}}).Return(nil),
			getAllCategories: db.EXPECT().GetAllCategories(uuid.UUID{}).Return([]types.Category{}, nil),
		},
		{
			name:             "update a category",
			form:             url.Values{"name": {"test 1"}},
			status:           http.StatusOK,
			body:             []string{},
			updateCategory:   db.EXPECT().UpdateCategory(types.Category{Id: categoryId, Name: "test 1", CompanyId: uuid.UUID{}}).Return(nil),
			getAllCategories: db.EXPECT().GetAllCategories(uuid.UUID{}).Return([]types.Category{}, nil),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			db.EXPECT().GetCategory(categoryId, uuid.UUID{}).Return(types.Category{
				Id:        categoryId,
				Name:      "test",
				CompanyId: uuid.UUID{},
			}, nil)

			if tt.updateCategory != nil {
				tt.updateCategory.Times(1)
			}
			if tt.getAllCategories != nil {
				tt.getAllCategories.Times(1)
			}

			req, res := createRequest(token, http.MethodPut, testUrl, strings.NewReader(tt.form.Encode()))
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, tt.status, res.Code)

			for _, b := range tt.body {
				assert.Contains(t, res.Body.String(), b)
			}
		})
	}
}
