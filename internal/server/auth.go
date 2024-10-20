package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodPost:
		resp := make(map[string]string)

		c := &types.CompanyCreate{}

		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err == io.EOF {
				resp["error"] = err.Error()
			} else {
				resp["error"] = "employees must be a number"
				resp["field"] = "employees"
			}
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				slog.Error("error handling JSON marshal", "err", err)
			}
			_, _ = w.Write(jsonResp)
			return
		}

		if c.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "name cannot be empty"
			resp["field"] = "name"
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				slog.Error("error handling JSON marshal", "err", err)
			}
			_, _ = w.Write(jsonResp)
			return
		}
		if c.Ruc == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "ruc cannot be empty"
			resp["field"] = "ruc"
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				slog.Error("error handling JSON marshal", "err", err)
			}
			_, _ = w.Write(jsonResp)
			return
		}
		if c.Employees <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "should pass at least one employee"
			resp["field"] = "employees"
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				slog.Error("error handling JSON marshal", "err", err)
			}
			_, _ = w.Write(jsonResp)
			return
		}
		if c.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "email cannot be empty"
			resp["field"] = "email"
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				slog.Error("error handling JSON marshal", "err", err)
			}
			_, _ = w.Write(jsonResp)
			return
		}

		if !utils.IsValidEmail(c.Email) {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "invalid email"
			resp["field"] = "email"
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				slog.Error("error handling JSON marshal", "err", err)
			}
			_, _ = w.Write(jsonResp)
			return
		}

		if c.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "password cannot be empty"
			resp["field"] = "password"
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				slog.Error("error handling JSON marshal", "err", err)
			}
			_, _ = w.Write(jsonResp)
			return
		}
		if c.User == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "name of the user cannot be empty"
			resp["field"] = "user"
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				slog.Error("error handling JSON marshal", "err", err)
			}
			_, _ = w.Write(jsonResp)
			return
		}

		if err := s.DB.CreateCompany(c); err != nil {
			slog.Error("Error creating company: ", "err", err)
			if strings.Contains(err.Error(), "SQLSTATE 23505") {
				w.WriteHeader(http.StatusConflict)
				resp["error"] = "company already exists"
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				resp["error"] = err.Error()
			}

			jsonResp, err := json.Marshal(resp)
			if err != nil {
				slog.Error("error handling JSON marshal", "err", err)
			}
			_, _ = w.Write(jsonResp)
			return
		}

		w.WriteHeader(http.StatusCreated)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodPost:
		resp := make(map[string]string)

		l := &types.Login{}

		if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err == io.EOF {
				resp["error"] = err.Error()
			}
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				slog.Error("error handling JSON marshal", "err", err)
			}
			_, _ = w.Write(jsonResp)
			return
		}
		if l.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "email cannot be empty"
			resp["field"] = "email"
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				slog.Error("error handling JSON marshal", "err", err)
			}
			_, _ = w.Write(jsonResp)
			return
		}

		if !utils.IsValidEmail(l.Email) {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "invalid email"
			resp["field"] = "email"
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				slog.Error("error handling JSON marshal", "err", err)
			}
			_, _ = w.Write(jsonResp)
			return
		}

		if l.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "password cannot be empty"
			resp["field"] = "password"
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				slog.Error("error handling JSON marshal", "err", err)
			}
			_, _ = w.Write(jsonResp)
			return
		}

		token, _, err := s.DB.Login(l)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			resp["error"] = err.Error()
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				slog.Error("error handling JSON marshal", "err", err)
			}
			_, _ = w.Write(jsonResp)
			return
		}

		resp["token"] = token

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			slog.Error("error handling JSON marshal", "err", err)
		}

		_, _ = w.Write(jsonResp)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) ApiLogin(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody || r.Body == nil {
		errorResponse := make(map[string]string)
		errorResponse["error"] = "credenciales inválidas"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var login types.Login
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		errorResponse := make(map[string]string)
		w.WriteHeader(http.StatusBadRequest)
		if err == io.EOF {
			errorResponse["error"] = err.Error()
		}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if login.Email == "" || login.Password == "" || !utils.IsValidEmail(login.Email) {
		errorResponse := make(map[string]string)
		errorResponse["error"] = "credenciales inválidas"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	_, user, err := s.DB.Login(&login)
	if err != nil {
		errorResponse := make(map[string]string)
		w.WriteHeader(http.StatusUnauthorized)
		errorResponse["error"] = "credenciales inválidas"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	_, token, _ := s.TokenAuth.Encode(map[string]interface{}{"id": user.Id, "name": user.Name, "email": user.Email, "company_id": user.CompanyId, "role": user.RoleId})

	resp := make(map[string]interface{})
	resp["user"] = user
	resp["token"] = token

	json.NewEncoder(w).Encode(resp)
}
