package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
	"bca-go-final/internal/views/bca/transaction/partials/details"
)

func (s *Server) DetailsTable(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	id := mux.Vars(r)["invoiceId"]
	parsedInvoiceId, _ := uuid.Parse(id)

	if r.Method == http.MethodPost {
		r.ParseForm()
		biId := r.Form.Get("item")
		if biId == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Ingrese una partida"))
			return
		}
		parsedBudgetItemId, err := uuid.Parse(biId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Error parsing budgetItemId. Err: ", err)
			return
		}
		q := r.Form.Get("quantity")
		c := r.Form.Get("cost")

		if q == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Ingrese una cantidad"))
			return
		}
		if c == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Ingrese un costo"))
			return
		}
		quantity, err := strconv.ParseFloat(q, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Cantidad debe ser un número válido"))
			return
		}
		cost, err := strconv.ParseFloat(c, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Costo debe ser un número válido"))
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
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Ya existe una partida con ese nombre en la factura"))
				return
			}
			if strings.Contains(err.Error(), "no rows") {
				w.WriteHeader(http.StatusBadRequest)
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
		return
	}

	component := details.InvoiceDetailsTable(det)
	component.Render(r.Context(), w)
}

func (s *Server) DetailsAdd(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	id := mux.Vars(r)["invoiceId"]
	parsedInvoiceId, _ := uuid.Parse(id)

	budgetItems := s.returnAllSelects([]string{"budgetitems"}, ctx.CompanyId)["budgetitems"]

	component := details.EditDetails(budgetItems, parsedInvoiceId.String())
	component.Render(r.Context(), w)
}

func (s *Server) DetailsEdit(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	iId := mux.Vars(r)["invoiceId"]
	bId := mux.Vars(r)["budgetItemId"]
	parsedInvoiceId, _ := uuid.Parse(iId)
	parsedBudgetItemId, _ := uuid.Parse(bId)
	_ = parsedBudgetItemId

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
