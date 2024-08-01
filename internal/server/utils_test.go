package server_test

import (
	"io"
	"net/http"
	"net/http/httptest"
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
