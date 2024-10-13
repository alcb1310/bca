package server_test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

// creates a request and response recorder with a given JWT token
func createRequest(token, method, url string, body io.Reader) (*http.Request, *httptest.ResponseRecorder) {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if token != "" {
		req.AddCookie(&http.Cookie{
			Name:  "jwt",
			Value: token,
			Path:  "/",
		})
	}
	res := httptest.NewRecorder()
	return req, res
}

func createApiRequest(token, method, url string, body io.Reader) (*http.Request, *httptest.ResponseRecorder) {
	var req *http.Request
	if body == nil {
		req, _ = http.NewRequest(method, url, nil)
	} else {
		req, _ = http.NewRequest(method, url, body)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.AddCookie(&http.Cookie{
			Name:  "jwt",
			Value: token,
			Path:  "/",
		})
	}
	res := httptest.NewRecorder()
	return req, res
}

// creates a valid JWT token for testing purposes
func createToken(ja *jwtauth.JWTAuth) string {
	_, token, _ := ja.Encode(map[string]interface{}{
		"id":        uuid.UUID{},
		"email":     "test@test.com",
		"name":      "test",
		"companyId": uuid.UUID{},
		"roleId":    "a",
	})

	return token
}
