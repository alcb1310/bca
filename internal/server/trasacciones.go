package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"bca-go-final/internal/utils"
	"bca-go-final/internal/views/bca/transaction"
)

func (s *Server) Budget(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	projects := s.returnAllSelects([]string{"projects"}, ctx.CompanyId)["projects"]

	component := transaction.BudgetView(projects)
	component.Render(r.Context(), w)
}

func (s *Server) Invoice(w http.ResponseWriter, r *http.Request) {
	component := transaction.InvoiceView()
	component.Render(r.Context(), w)
}

func (s *Server) Closure(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	projects := s.returnAllSelects([]string{"projects"}, ctx.CompanyId)["projects"]

	if r.Method == http.MethodPost {
		r.ParseForm()
		success := "true"

		pId := r.Form.Get("proyecto")
		parsedProjectId, err := utils.ValidateUUID(pId, "proyecto")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		d := r.Form.Get("date")
		if d == "" {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Seleccione una fecha")
			return
		}
		date, err := time.Parse("2006-01-02", d)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := s.DB.CreateClosure(ctx.CompanyId, parsedProjectId, date); err != nil {
			log.Println(err)
			log.Println(fmt.Sprintf("No se pudo cerrar el proyecto: %s para la fecha: %s", parsedProjectId, utils.ConvertDate(date)))
			success = "false"
		}
		// w.WriteHeader(http.StatusOK)
		w.Header().Set("HX-Redirect", "/bca/transacciones/cierre?success="+success)
		http.Redirect(w, r, "/bca/transacciones/cierre?success=true", http.StatusOK)
		return
	}

	component := transaction.ClosureView(projects)
	component.Render(r.Context(), w)
}
