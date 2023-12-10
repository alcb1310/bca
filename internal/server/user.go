package server

import (
	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
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
		if r.Method != http.MethodGet {
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
