package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost, http.MethodOptions:
		resp := make(map[string]string)
		resp["message"] = "Register"

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("error handling JSON marshal. Err: %v", err)
		}

		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(jsonResp)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	log.Printf("login request: %s %s\n", r.Method, r.URL)
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
