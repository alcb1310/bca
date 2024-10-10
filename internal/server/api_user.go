package server

import (
	"encoding/json"
	"net/http"

	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	user, _ := s.DB.GetUser(ctx.Id, ctx.CompanyId)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (s *Server) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	users, _ := s.DB.GetAllUsers(ctx.CompanyId)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
