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

func TestCreateUser(t *testing.T) {
	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret", -5)
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

func TestUpdatePassword(t *testing.T) {
	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret", -5)
	token := createToken(s.TokenAuth)

	testData := []struct {
		name           string
		form           url.Values
		status         int
		body           []string
		updatePassword *mocks.Service_UpdatePassword_Call
	}{
		{
			name:           "should pass a form",
			form:           nil,
			status:         http.StatusBadRequest,
			body:           []string{},
			updatePassword: nil,
		},
		{
			name:           "should pass a password",
			form:           url.Values{},
			status:         http.StatusBadRequest,
			body:           []string{},
			updatePassword: nil,
		},
		{
			name:   "should update password",
			form:   url.Values{"password": {"test"}},
			status: http.StatusOK,
			updatePassword: db.EXPECT().UpdatePassword("test", uuid.UUID{}, uuid.UUID{}).Return(types.User{
				Id:        uuid.UUID{},
				Email:     "test@test.com",
				Name:      "test",
				CompanyId: uuid.UUID{},
				RoleId:    "a",
			}, nil),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			if tt.updatePassword != nil {
				tt.updatePassword.Times(1)
			}

			req, res := createRequest(token, http.MethodPut, "/bca/partials/users", strings.NewReader(tt.form.Encode()))
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

func TestUpdateUser(t *testing.T) {
	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret", -5)
	token := createToken(s.TokenAuth)
	id := uuid.New()

	testData := []struct {
		name            string
		id              uuid.UUID
		getUserMock     *mocks.Service_GetUser_Call
		updateUserMock  *mocks.Service_UpdateUser_Call
		getAllUsersMock *mocks.Service_GetAllUsers_Call
		form            url.Values
		status          int
	}{
		{
			name:           "should pass an id",
			id:             uuid.UUID{},
			getUserMock:    nil,
			updateUserMock: nil,
			form:           nil,
			status:         http.StatusForbidden,
		},
		{
			name: "should get user",
			id:   id,
			getUserMock: db.EXPECT().GetUser(id, uuid.UUID{}).Return(types.User{
				Id:        id,
				Email:     "test@test.com",
				Name:      "test",
				CompanyId: uuid.UUID{},
				RoleId:    "a",
			}, nil),
			updateUserMock: db.EXPECT().UpdateUser(types.User{
				Id:        id,
				Email:     "test@test.com",
				Name:      "test",
				CompanyId: uuid.UUID{},
				RoleId:    "a",
			}, id, uuid.UUID{}).Return(types.User{
				Id:        id,
				Email:     "test@test.com",
				Name:      "test",
				CompanyId: uuid.UUID{},
				RoleId:    "a",
			}, nil),
			getAllUsersMock: db.EXPECT().GetAllUsers(uuid.UUID{}).Return([]types.User{
				{
					Id:        id,
					Email:     "test@test.com",
					Name:      "test",
					CompanyId: uuid.UUID{},
					RoleId:    "a",
				},
			}, nil),
			form:   nil,
			status: http.StatusOK,
		},
		{
			name: "Invalid email",
			id:   id,
			form: url.Values{"email": {"test"}},
			getUserMock: db.EXPECT().GetUser(id, uuid.UUID{}).Return(types.User{
				Id:        id,
				Email:     "test@test.com",
				Name:      "test",
				CompanyId: uuid.UUID{},
				RoleId:    "a",
			}, nil),
			status: http.StatusBadRequest,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			if tt.getUserMock != nil {
				tt.getUserMock.Times(1)
			}

			if tt.updateUserMock != nil {
				tt.updateUserMock.Times(1)
			}

			if tt.getAllUsersMock != nil {
				tt.getAllUsersMock.Times(1)
			}
			req, res := createRequest(token, http.MethodPut, fmt.Sprintf("/bca/partials/users/%s", tt.id.String()), strings.NewReader(tt.form.Encode()))
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, tt.status, res.Code)
		})
	}
}
