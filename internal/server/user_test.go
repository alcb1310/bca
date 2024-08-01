package server_test

import (
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

func TestCreateUser(t *testing.T) {
	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret")
	token := createToken(s.TokenAuth)

	testData := []struct {
		name            string
		form            url.Values
		status          int
		body            []string
		createUserMock  *mocks.Service_CreateUser_Call
		getAllUsersMock *mocks.Service_GetAllUsers_Call
	}{
		{
			name:            "should pass a form",
			form:            nil,
			status:          http.StatusBadRequest,
			body:            []string{},
			createUserMock:  nil,
			getAllUsersMock: nil,
		},
		{
			name:            "should pass an email",
			form:            url.Values{},
			status:          http.StatusBadRequest,
			body:            []string{},
			createUserMock:  nil,
			getAllUsersMock: nil,
		},
		{
			name:            "should pass a valid email",
			form:            url.Values{"email": {"test"}},
			status:          http.StatusBadRequest,
			body:            []string{},
			createUserMock:  nil,
			getAllUsersMock: nil,
		},
		{
			name:            "should pass a password",
			form:            url.Values{"email": {"test@test.com"}},
			status:          http.StatusBadRequest,
			body:            []string{},
			createUserMock:  nil,
			getAllUsersMock: nil,
		},
		{
			name:            "should pass a name",
			form:            url.Values{"email": {"test@test.com"}, "password": {"test"}},
			status:          http.StatusBadRequest,
			body:            []string{},
			createUserMock:  nil,
			getAllUsersMock: nil,
		},
		{
			name:   "should create a new user",
			form:   url.Values{"email": {"test@test.com"}, "password": {"test"}, "name": {"test"}},
			status: http.StatusOK,
			body:   []string{"test", "test@test.com", "<table>"},
			createUserMock: db.EXPECT().CreateUser(&types.UserCreate{
				Email:     "test@test.com",
				Password:  "test",
				Name:      "test",
				CompanyId: uuid.UUID{},
				RoleId:    "a",
			}).Return(types.User{
				Id:        uuid.UUID{},
				Email:     "test@test.com",
				Name:      "test",
				CompanyId: uuid.UUID{},
				RoleId:    "a",
			}, nil),
			getAllUsersMock: db.EXPECT().GetAllUsers(uuid.UUID{}).Return([]types.User{
				{
					Id:        uuid.UUID{},
					Email:     "test@test.com",
					Name:      "test",
					CompanyId: uuid.UUID{},
					RoleId:    "a",
				},
			}, nil),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createUserMock != nil {
				tt.createUserMock.Times(1)
			}
			if tt.getAllUsersMock != nil {
				tt.getAllUsersMock.Times(1)
			}
			req, res := createRequest(token, http.MethodPost, "/bca/partials/users", strings.NewReader(tt.form.Encode()))
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, tt.status, res.Code)
			if len(tt.body) > 0 {
				for _, b := range tt.body {
					assert.Contains(t, res.Body.String(), b)
				}
			}
		})
	}
}
