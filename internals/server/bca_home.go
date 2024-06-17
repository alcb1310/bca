package server

import (
	"net/http"

	"github.com/alcb1310/bca/externals/views/bca"
)

func (s *BCAService) BCAHome(w http.ResponseWriter, r *http.Request) error {
	user, err := getUserFromContext(r)
	if err != nil {
		return err
	}

	return renderPage(w, r, bca.Index(user))
}
