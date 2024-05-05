package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
	"bca-go-final/internal/views/bca/transaction/partials/details"
)

func (s *Server) DetailsTable(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	parsedInvoiceId, _ := utils.ValidateUUID(mux.Vars(r)["invoiceId"], "factura")

	if r.Method == http.MethodPost {
		r.ParseForm()
		biId := r.Form.Get("item")
		parsedBudgetItemId, err := utils.ValidateUUID(biId, "partida")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			log.Println("Error parsing budgetItemId. Err: ", err)
			return
		}

		q := r.Form.Get("quantity")
		quantity, err := utils.ConvertFloat(q, "cantidad", true)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		c := r.Form.Get("cost")
		cost, err := utils.ConvertFloat(c, "costo", true)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		d := types.InvoiceDetailCreate{
			InvoiceId:    parsedInvoiceId,
			BudgetItemId: parsedBudgetItemId,
			CompanyId:    ctx.CompanyId,
			Quantity:     quantity,
			Cost:         cost,
			Total:        quantity * cost,
		}

		if err := s.DB.AddDetail(d); err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("Ya existe una partida con ese nombre en la factura"))
				return
			}
			if strings.Contains(err.Error(), "no rows") {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("No existe presupuesto para esa partida"))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error al insertar partida. Err: %s", err.Error())))
			log.Println(err)
			return
		}
	}

	det, err := s.DB.GetAllDetails(parsedInvoiceId, ctx.CompanyId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}

	component := details.InvoiceDetailsTable(det)
	component.Render(r.Context(), w)
}

func (s *Server) DetailsAdd(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	parsedInvoiceId, _ := utils.ValidateUUID(mux.Vars(r)["invoiceId"], "factura")

	budgetItems := s.returnAllSelects([]string{"budgetitems"}, ctx.CompanyId)["budgetitems"]

	component := details.EditDetails(budgetItems, parsedInvoiceId.String())
	component.Render(r.Context(), w)
}

func (s *Server) DetailsEdit(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	parsedInvoiceId, _ := utils.ValidateUUID(mux.Vars(r)["invoiceId"], "factura")
	parsedBudgetItemId, _ := utils.ValidateUUID(mux.Vars(r)["budgetItemId"], "partida")

	if err := s.DB.DeleteDetail(parsedInvoiceId, parsedBudgetItemId, ctx.CompanyId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	det, err := s.DB.GetAllDetails(parsedInvoiceId, ctx.CompanyId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	component := details.InvoiceDetailsTable(det)
	component.Render(r.Context(), w)
}
