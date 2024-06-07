package server

import "net/http"

func (s *Service) Home(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World"))
	return nil
}
