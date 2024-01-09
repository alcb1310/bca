package server

import (
	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
	"bca-go-final/internal/views/bca/users"
	"bca-go-final/internal/views/partials"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) AllUsers(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)

	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodPatch:
		resp := map[string]string{}
		type expect struct {
			Password string `json:"password"`
		}
		var u = &expect{}

		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		if u.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "password cannot be empty"
			resp["field"] = "password"
			json.NewEncoder(w).Encode(resp)
			return
		}

		user, err := s.DB.UpdatePassword(u.Password, ctxPayload.Id, ctxPayload.CompanyId)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				resp["error"] = fmt.Sprintf("User with ID: `%s` not found", ctxPayload.Id)
				json.NewEncoder(w).Encode(resp)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	case http.MethodGet:
		users, err := s.DB.GetAllUsers(ctxPayload.CompanyId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
	case http.MethodPost:
		resp := map[string]string{}
		var u = &types.UserCreate{}

		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		if u.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "email cannot be empty"
			resp["field"] = "email"
			json.NewEncoder(w).Encode(resp)
			return
		}
		if !utils.IsValidEmail(u.Email) {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "invalid email"
			resp["field"] = "email"
			json.NewEncoder(w).Encode(resp)
			return
		}
		if u.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "password cannot be empty"
			resp["field"] = "password"
			json.NewEncoder(w).Encode(resp)
			return
		}
		if u.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "name cannot be empty"
			resp["field"] = "name"
			json.NewEncoder(w).Encode(resp)
			return
		}
		if u.RoleId == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "role cannot be empty"
			resp["field"] = "role"
			json.NewEncoder(w).Encode(resp)
			return
		}
		u.CompanyId = ctxPayload.CompanyId

		ux, err := s.DB.CreateUser(u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ux)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) OneUser(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	ctxPayload, _ := getMyPaload(r)
	id := mux.Vars(r)["id"]
	if strings.ToLower(id) == "me" {
		if r.Method != http.MethodGet && r.Method != http.MethodOptions {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		id = ctxPayload.Id.String()
	}

	parsedId, err := uuid.Parse(id)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	user, err := s.DB.GetUser(parsedId, ctxPayload.CompanyId)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			resp["error"] = fmt.Sprintf("User with ID: `%s` not found", id)
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		resp["error"] = err.Error()
		json.NewEncoder(w).Encode(resp)
		return
	}

	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)

	case http.MethodDelete:
		if parsedId == ctxPayload.Id {
			w.WriteHeader(http.StatusForbidden)
			resp["error"] = "cannot delete yourself"
			json.NewEncoder(w).Encode(resp)
			return
		}
		if err := s.DB.DeleteUser(parsedId, ctxPayload.CompanyId); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	case http.MethodPut:
		var u = &types.User{}

		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		if u.Email != "" && !utils.IsValidEmail(u.Email) {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "invalid email"
			resp["field"] = "email"
			json.NewEncoder(w).Encode(resp)
			return
		}

		if u.Email != "" {
			user.Email = u.Email
		}
		if u.Name != "" {
			user.Name = u.Name
		}
		if u.RoleId != "" {
			user.RoleId = u.RoleId
		}

		ux, err := s.DB.UpdateUser(user, parsedId, ctxPayload.CompanyId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ux)
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) Profile(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)

	user, _ := s.DB.GetUser(ctx.Id, ctx.CompanyId)

	component := users.ProfileView(user)
	component.Render(r.Context(), w)
}

func (s *Server) Admin(w http.ResponseWriter, r *http.Request) {
	component := users.AdminView()
	component.Render(r.Context(), w)
}

func (s *Server) ChangePassword(w http.ResponseWriter, r *http.Request) {
	component := users.ChangePasswordView()
	component.Render(r.Context(), w)
}

func (s *Server) SingleUser(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)
	id := mux.Vars(r)["id"]
	parsedId, _ := uuid.Parse(id)

	if ctx.Id == parsedId {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method == http.MethodDelete {
		if err := s.DB.DeleteUser(parsedId, ctx.CompanyId); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	users, _ := s.DB.GetAllUsers(ctx.CompanyId)
	component := partials.UsersTable(users, ctx.Id)
	component.Render(r.Context(), w)
}

func (s *Server) UsersTable(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)

	if r.Method == http.MethodPost {
		u := &types.UserCreate{}
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		u.Name = r.Form.Get("name")
		u.Email = r.Form.Get("email")
		u.Password = r.Form.Get("password")
		u.RoleId = r.Form.Get("role")
		u.CompanyId = ctx.CompanyId

		if u.Email != "" && !utils.IsValidEmail(u.Email) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if u.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if u.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if u.RoleId == "" {
			u.RoleId = "a"
		}

		_, err = s.DB.CreateUser(u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	users, _ := s.DB.GetAllUsers(ctx.CompanyId)
	component := partials.UsersTable(users, ctx.Id)
	component.Render(r.Context(), w)
}
