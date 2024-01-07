package server

import (
	"bca-go-final/internal/views"
	"bca-go-final/internal/views/base"
	"bca-go-final/internal/views/derrors"
	"net/http"
)

func (s *Server) LoginView(w http.ResponseWriter, r *http.Request) {
	err := make(map[string]string)
	switch r.Method {
	case http.MethodGet:
		component := views.LoginView(err)
		base := base.Layout(component)
		base.Render(r.Context(), w)
	default:
		err["method"] = r.Method
		err["url"] = r.RequestURI
		component := derrors.NotImplemented(err)
		base := base.Layout(component)
		base.Render(r.Context(), w)
	}
}
