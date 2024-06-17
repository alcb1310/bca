package server

import (
	"net/http"

	"github.com/alcb1310/bca/externals/views/home"
)

func (s *Service) Home(w http.ResponseWriter, r *http.Request) error {
	return renderPage(w, r, home.HomeIndex())
}
