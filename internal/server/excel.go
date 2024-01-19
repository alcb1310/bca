package server

import (
	"bca-go-final/internal/excel"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (s *Server) BalanceExcel(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)
	id := r.URL.Query().Get("project")
	parsedProjectId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error parsing projectId. Err: ", err)
		return
	}
	d := r.URL.Query().Get("date")
	dateVal, _ := time.Parse("2006-01-02", d)

	f := excel.Balance(ctx.CompanyId, parsedProjectId, dateVal, s.DB)
	fName := "/" + strings.Trim(f.Path, "./public")

	w.Header().Set("HX-Redirect", fName)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=cuadre.xlsx"))
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	if err := f.Write(w); err != nil {
		log.Println(err.Error())
	}

	go func() {
		time.Sleep(1 * time.Second)
		if err := os.Remove(f.Path); err != nil {
			log.Println(err.Error())
		}
	}()

}

func (s *Server) ActualExcel(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)
	parsedProjectId, err := uuid.Parse(r.URL.Query().Get("proyecto"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error parsing projectId. Err: ", err)
		return
	}

	l, err := strconv.ParseUint(r.URL.Query().Get("nivel"), 10, 8)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error parsing level. Err: ", err)
		return
	}
	level := uint8(l)
	budgets, _ := s.DB.GetBudgetsByProjectId(ctx.CompanyId, parsedProjectId, &level)

	f := excel.Actual(ctx.CompanyId, parsedProjectId, budgets, nil, s.DB)
	fName := "/" + strings.Trim(f.Path, "./public")

	w.Header().Set("HX-Redirect", fName)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=actual.xlsx"))
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	if err := f.Write(w); err != nil {
		log.Println(err.Error())
	}

	go func() {
		time.Sleep(1 * time.Second)
		if err := os.Remove(f.Path); err != nil {
			log.Println(err.Error())
		}
	}()
}

func (s *Server) HistoricExcel(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)
	parsedProjectId, err := uuid.Parse(r.URL.Query().Get("proyecto"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error parsing projectId. Err: ", err)
		return
	}

	l, err := strconv.ParseUint(r.URL.Query().Get("nivel"), 10, 8)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error parsing level. Err: ", err)
		return
	}
	level := uint8(l)
	d := r.URL.Query().Get("fecha")
	dateVal, _ := time.Parse("2006-01-02", d)

	budgets := s.DB.GetHistoricByProject(ctx.CompanyId, parsedProjectId, dateVal, level)

	f := excel.Actual(ctx.CompanyId, parsedProjectId, budgets, &dateVal, s.DB)
	fName := "/" + strings.Trim(f.Path, "./public")

	w.Header().Set("HX-Redirect", fName)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=actual.xlsx"))
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	if err := f.Write(w); err != nil {
		log.Println(err.Error())
	}

	go func() {
		time.Sleep(1 * time.Second)
		if err := os.Remove(f.Path); err != nil {
			log.Println(err.Error())
		}
	}()
}
