package server

import (
	"net/http"

	"github.com/alcb1310/bca/externals/views/bca"
)

func (s *BCAService) BCAHome(w http.ResponseWriter, r *http.Request) error {
	return renderPage(w, r, bca.Index())
}
