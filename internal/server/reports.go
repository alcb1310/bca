package server

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
	"bca-go-final/internal/views/bca/reports"
	"bca-go-final/internal/views/bca/reports/partials"
)

func (s *Server) Actual(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	data := []string{"projects", "levels"}
	results := s.returnAllSelects(data, ctxPayload.CompanyId)
	projects := results["projects"]
	levels := results["levels"]

	component := reports.ActualView(projects, levels)
	component.Render(r.Context(), w)
}

func (s *Server) Balance(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		pId := r.Form.Get("project")
		parsedProjectId, _ := utils.ValidateUUID(pId, "proyecto")

		d := r.Form.Get("date")
		date, _ := time.Parse("2006-01-02", d)

		invoices := s.DB.GetBalance(ctxPayload.CompanyId, parsedProjectId, date)

		component := partials.BalanceView(invoices)
		component.Render(r.Context(), w)

	case http.MethodGet:
		projects := s.returnAllSelects([]string{"projects"}, ctxPayload.CompanyId)["projects"]

		component := reports.BalanceView(projects)
		component.Render(r.Context(), w)
	}

}

func (s *Server) Historic(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	data := []string{"projects", "levels"}
	results := s.returnAllSelects(data, ctxPayload.CompanyId)
	projects := results["projects"]
	levels := results["levels"]

	if r.URL.Query().Get("proyecto") != "" && r.URL.Query().Get("fecha") != "" && r.URL.Query().Get("nivel") != "" {
		pId := r.URL.Query().Get("proyecto")
		parsedProjectId, _ := utils.ValidateUUID(pId, "proyecto")

		d := r.URL.Query().Get("fecha")
		date, _ := time.Parse("2006-01-02", d)
		l, _ := strconv.ParseUint(r.URL.Query().Get("nivel"), 10, 64)
		nivel := uint8(l)

		budgets := s.DB.GetHistoricByProject(ctxPayload.CompanyId, parsedProjectId, date, nivel)
		component := partials.BudgetView(budgets)
		component.Render(r.Context(), w)
		return
	}

	component := reports.HistoricView(projects, levels)
	component.Render(r.Context(), w)
}

func (s *Server) Spent(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	data := []string{"projects", "levels"}
	results := s.returnAllSelects(data, ctxPayload.CompanyId)
	projects := results["projects"]
	levels := results["levels"]

	if r.URL.Query().Get("proyecto") != "" && r.URL.Query().Get("fecha") != "" && r.URL.Query().Get("nivel") != "" {
		pId := r.URL.Query().Get("proyecto")
		parsedProjectId, _ := utils.ValidateUUID(pId, "proyecto")

		d := r.URL.Query().Get("fecha")
		date, _ := time.Parse("2006-01-02", d)
		l, _ := strconv.ParseUint(r.URL.Query().Get("nivel"), 10, 64)
		nivel := uint8(l)

		budgetItems := s.DB.GetBudgetItemsByLevel(ctxPayload.CompanyId, nivel)
		reportData := []types.Spent{}
		var grandTotal float64 = 0

		for _, bi := range budgetItems {
			x := []types.BudgetItem{bi}
			res := []uuid.UUID{}
			res = s.DB.GetNonAccumulateChildren(&ctxPayload.CompanyId, &parsedProjectId, x, res)

			total := s.DB.GetSpentByBudgetItem(ctxPayload.CompanyId, parsedProjectId, bi.ID, date, res)
			grandTotal += total
			if total > 0 {
				reportData = append(reportData, types.Spent{
					Spent:      total,
					BudgetItem: bi,
				})
			}
		}

		nCtx := context.WithValue(r.Context(), "date", utils.ConvertDate(date))

		component := partials.SpentView(types.SpentResponse{
			Spent:   reportData,
			Total:   grandTotal,
			Project: parsedProjectId,
		})
		component.Render(nCtx, w)
		return
	}

	component := reports.SpentView(projects, levels)
	component.Render(r.Context(), w)
}

func (s *Server) ActualGenerate(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)
	r.ParseForm()
	p := r.Form.Get("proyecto")
	projectId, _ := utils.ValidateUUID(p, "proyecto")
	z := r.Form.Get("nivel")
	var l uint64
	var err error
	if z == "" {
		l = 0
	} else {
		l, err = strconv.ParseUint(z, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
	level := uint8(l)

	budgets, err := s.DB.GetBudgetsByProjectId(ctxPayload.CompanyId, projectId, &level)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	component := partials.BudgetView(budgets)
	component.Render(r.Context(), w)
}

func (s *Server) SpentByBudgetItem(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)
	budgetItemId, err := utils.ValidateUUID(mux.Vars(r)["budgetItemId"], "partida")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	parsedProjectId, err := utils.ValidateUUID(mux.Vars(r)["projectId"], "proyecto")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	d := mux.Vars(r)["date"]
	date, _ := time.Parse("2006-01-02", d)

	budgetItem, _ := s.DB.GetOneBudgetItem(budgetItemId, ctxPayload.CompanyId)

	x := []types.BudgetItem{*budgetItem}
	var res []uuid.UUID
	res = s.DB.GetNonAccumulateChildren(&ctxPayload.CompanyId, &parsedProjectId, x, res)

	spent := s.DB.GetDetailsByBudgetItem(ctxPayload.CompanyId, parsedProjectId, budgetItemId, date, res)

	component := partials.SpentDetails(spent, *budgetItem)
	component.Render(r.Context(), w)
}
