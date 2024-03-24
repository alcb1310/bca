package server

import (
	"fmt"
	"net/http"

	"bca-go-final/internal/utils"
	"bca-go-final/internal/views/bca/settings/partials"
)

func (s *Server) RubrosTable(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	rubros, _ := s.DB.GetAllRubros(ctxPayload.CompanyId)

	fmt.Println(rubros)
	component := partials.RubrosTable(rubros)
	component.Render(r.Context(), w)
}
