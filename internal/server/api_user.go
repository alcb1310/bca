package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	user, _ := s.DB.GetUser(ctx.Id, ctx.CompanyId)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

func (s *Server) ApiGetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	users, _ := s.DB.GetAllUsers(ctx.CompanyId)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}

func (s *Server) ApiCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody || r.Body == nil {
		errorResponse := make(map[string]string)
		errorResponse["error"] = "Invalid request body"

		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	ctx, _ := utils.GetMyPaload(r)

	user := types.UserCreate{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		errorResponse := make(map[string]string)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
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
		_ = json.NewEncoder(w).Encode(userErrors)
		return
	}

	createdUser, err := s.DB.CreateUser(&user)
	if err != nil {
		var e *pgconn.PgError

		if errors.As(err, &e) && e.Code == "23505" {
			errorResponse := make(map[string]string)
			errorResponse["error"] = "Ya existe un usuario con ese correo"
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(errorResponse)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(createdUser)
}

func (s *Server) ApiDeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, _ := uuid.Parse(id)

	if ctx.Id == parsedId {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err := s.DB.DeleteUser(parsedId, ctx.CompanyId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) ApiUpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, _ := uuid.Parse(id)
	var user, userInfo types.User
	var err error

	if user, err = s.DB.GetUser(parsedId, ctx.CompanyId); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	slog.Info("ApiUpdateUser", "user", user)
	_ = json.NewDecoder(r.Body).Decode(&userInfo)

	if userInfo.Email != "" {
		user.Email = userInfo.Email
	}

	if userInfo.Name != "" {
		user.Name = userInfo.Name
	}

	if userInfo.RoleId != "" {
		user.RoleId = userInfo.RoleId
	}

	slog.Info("ApiUpdateUser", "user", user)

	if _, err = s.DB.UpdateUser(user, parsedId, ctx.CompanyId); err != nil {
		var e *pgconn.PgError

		if errors.As(err, &e) && e.Code == "23505" {
			errorResponse := make(map[string]string)
			errorResponse["error"] = "Ya existe un usuario con ese correo"
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(errorResponse)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}
