package server

import (
	"bca-go-final/internal/types"
	"bca-go-final/internal/views/bca/transaction/partials/details"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) DetailsTable(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)
	id := mux.Vars(r)["invoiceId"]
	parsedInvoiceId, _ := uuid.Parse(id)

	if r.Method == http.MethodPost {
		r.ParseForm()
		biId := r.Form.Get("item")
		if biId == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("budgetItemId is empty")
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
			log.Println("quantity is empty")
			return
		}
		if c == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("cost is empty")
			return
		}
		quantity, err := strconv.ParseFloat(q, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Error parsing quantity. Err: ", err)
			return
		}
		cost, err := strconv.ParseFloat(c, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Error parsing cost. Err: ", err)
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
			w.WriteHeader(http.StatusInternalServerError)
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
	ctx, _ := getMyPaload(r)

	id := mux.Vars(r)["invoiceId"]
	parsedInvoiceId, _ := uuid.Parse(id)

	budgetItems := []types.Select{}
	bi := s.DB.GetBudgetItemsByAccumulate(ctx.CompanyId, false)

	for _, v := range bi {
		x := types.Select{
			Key:   v.ID.String(),
			Value: v.Name,
		}
		budgetItems = append(budgetItems, x)
	}

	component := details.EditDetails(budgetItems, parsedInvoiceId.String())
	component.Render(r.Context(), w)
}
