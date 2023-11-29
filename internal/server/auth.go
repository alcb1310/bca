package server

import (
	"bca-go-final/internal/types"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost, http.MethodOptions:
		resp := make(map[string]string)

		var c types.Company
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err == io.EOF {
				resp["error"] = err.Error()
			} else {
				resp["error"] = "employees must be a number"
			}
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Fatalf("error handling JSON marshal. Err: %v", err)
			}
			_, _ = w.Write(jsonResp)
			return
		}

		if c.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "name cannot be empty"
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Fatalf("error handling JSON marshal. Err: %v", err)
			}
			_, _ = w.Write(jsonResp)
			return
		}
		if c.Ruc == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "ruc cannot be empty"
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Fatalf("error handling JSON marshal. Err: %v", err)
			}
			_, _ = w.Write(jsonResp)
			return
		}
		if c.Employees <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "should pass at least one employee"
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Fatalf("error handling JSON marshal. Err: %v", err)
			}
			_, _ = w.Write(jsonResp)
			return
		}

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("error handling JSON marshal. Err: %v", err)
		}

		if err := s.db.CreateCompany(&c); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = err.Error()
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Fatalf("error handling JSON marshal. Err: %v", err)
			}
			_, _ = w.Write(jsonResp)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(jsonResp)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost, http.MethodOptions:
		resp := make(map[string]string)
		resp["message"] = "Login"

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("error handling JSON marshal. Err: %v", err)
		}

		_, _ = w.Write(jsonResp)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
