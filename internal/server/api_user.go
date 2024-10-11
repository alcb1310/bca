package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	user, _ := s.DB.GetUser(ctx.Id, ctx.CompanyId)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (s *Server) ApiGetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	users, _ := s.DB.GetAllUsers(ctx.CompanyId)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (s *Server) ApiCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody || r.Body == nil {
		errorResponse := make(map[string]string)
		errorResponse["error"] = "Invalid request body"

		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	ctx, _ := utils.GetMyPaload(r)

	user := types.UserCreate{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		errorResponse := make(map[string]string)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse["error"] = err.Error()
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	user.RoleId = "a"
	user.CompanyId = ctx.CompanyId

	userErrors := make(map[string]string)
	if user.Email == "" || !utils.IsValidEmail(user.Email) {
		userErrors["email"] = "Ingrese un email valido"
	}

	if user.Name == "" {
		userErrors["name"] = "Ingrese un nombre"
	}

	if user.Password == "" {
		userErrors["password"] = "Ingrese una contraseña"
	}

	if len(userErrors) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(userErrors)
		return
	}

	createdUser, err := s.DB.CreateUser(&user)
	if err != nil {
		var e *pgconn.PgError

		if errors.As(err, &e) && e.Code == "23505" {
			errorResponse := make(map[string]string)
			errorResponse["error"] = "Ya existe un usuario con ese correo"
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}
