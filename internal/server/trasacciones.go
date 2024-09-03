package server

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
	"github.com/alcb1310/bca/internal/views/bca/transaction"
)

func (s *Server) Budget(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	p := s.DB.GetActiveProjects(ctx.CompanyId, true)
	projects := []types.Select{}
	for _, v := range p {
		x := types.Select{
			Key:   v.ID.String(),
			Value: v.Name,
		}
		projects = append(projects, x)
	}

	component := transaction.BudgetView(projects)
	component.Render(r.Context(), w)
}

func (s *Server) Invoice(w http.ResponseWriter, r *http.Request) {
	component := transaction.InvoiceView()
	component.Render(r.Context(), w)
}

func (s *Server) Closure(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	r.ParseForm()
	pId := r.Form.Get("proyecto")
	parsedProjectId, err := uuid.Parse(pId)
	success := "true"
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if strings.Contains(err.Error(), "length: 0") {
			w.Write([]byte("Seleccione un proyecto"))
			return
		}
		w.Write([]byte("Proyecto inválido"))
		return
	}

	d := r.Form.Get("date")
	if d == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ingrese una fecha"))
		return
	}
	date, err := time.Parse("2006-01-02", d)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ingrese una fecha válida"))
		return
	}

	if err := s.DB.CreateClosure(ctx.CompanyId, parsedProjectId, date); err != nil {
		slog.Error("No se pudo cerrar el proyecto", "err", err, "project", parsedProjectId, "date", utils.ConvertDate(date))
		success = "false"
	}

	w.Header().Set("HX-Redirect", "/bca/transacciones/cierre?success="+success)
	http.Redirect(w, r, "/bca/transacciones/cierre?success=true", http.StatusOK)
}

func (s *Server) ClosureForm(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	p := s.DB.GetActiveProjects(ctx.CompanyId, true)
	projects := []types.Select{}
	for _, v := range p {
		x := types.Select{
			Key:   v.ID.String(),
			Value: v.Name,
		}
		projects = append(projects, x)
	}

	component := transaction.ClosureView(projects)
	component.Render(r.Context(), w)
}
