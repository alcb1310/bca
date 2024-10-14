package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiGetAllSuppliers(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	queryParams := r.URL.Query()
	search := queryParams.Get("query")

	suppliers, _ := s.DB.GetAllSuppliers(ctx.CompanyId, search)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(suppliers)
}

func (s *Server) ApiCreateSupplier(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody || r.Body == nil {
		errorResponse := make(map[string]string)
		errorResponse["error"] = "Invalid request body"

		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	ctx, _ := utils.GetMyPaload(r)

	var supplierCreate types.SupplierCreate
	if err := json.NewDecoder(r.Body).Decode(&supplierCreate); err != nil {
		errorResponse := make(map[string]string)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	contactName := sql.NullString{
		String: "",
		Valid:  false,
	}
	contactEmail := sql.NullString{
		String: "",
		Valid:  false,
	}
	contactPhone := sql.NullString{
		String: "",
		Valid:  false,
	}

	if supplierCreate.ContactName != "" {
		contactName = sql.NullString{
			String: supplierCreate.ContactName,
			Valid:  true,
		}
	}
	if supplierCreate.ContactEmail != "" {
		contactEmail = sql.NullString{
			String: supplierCreate.ContactEmail,
			Valid:  true,
		}
	}
	if supplierCreate.ContactPhone != "" {
		contactPhone = sql.NullString{
			String: supplierCreate.ContactPhone,
			Valid:  true,
		}
	}

	supplier := types.Supplier{
		SupplierId:   supplierCreate.SupplierId,
		Name:         supplierCreate.Name,
		ContactName:  contactName,
		ContactEmail: contactEmail,
		ContactPhone: contactPhone,
		CompanyId:    ctx.CompanyId,
	}

	if err := s.DB.CreateSupplier(&supplier); err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == "23505" {
			errorResponse := make(map[string]string)
			errorResponse["error"] = "Proveedor con ese ruc y/o nombre ya existe"
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(errorResponse)
			return
		}

		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(supplier)
}
