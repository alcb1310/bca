package server

import (
	"bca-go-final/internal/utils"
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) FlowDetails(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	r.ParseForm()
	p := r.Form.Get("project")
	if p == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	projectId, err := uuid.Parse(p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Código de proyecto es inválido"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		_, _ = s.DB.GetProject(projectId, ctx.CompanyId)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}
