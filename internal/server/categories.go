package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
	"bca-go-final/internal/views/bca/settings/partials"
)

func (s *Server) CategoriesTable(w http.ResponseWriter, r *http.Request) {
  var err error
  ctxPayload, _ := utils.GetMyPaload(r)

  if r.Method == http.MethodPost {
    r.ParseForm()
    c := types.Category{
      Name: r.Form.Get("name"),
      CompanyId: ctxPayload.CompanyId,
    }
    err = s.DB.CreateCategory(c)
    if err != nil {
      if strings.Contains(err.Error(), "duplicate") {
        w.WriteHeader(http.StatusConflict)
        w.Write([]byte(fmt.Sprintf("La categoria %s ya existe", c.Name)))
        return
      }
      w.WriteHeader(http.StatusInternalServerError)
      log.Println(err)
      return
    }
  }

  categories, _ := s.DB.GetAllCategories(ctxPayload.CompanyId)
  component := partials.CategoriesTable(categories)
  
  component.Render(r.Context(), w)
}

func (s *Server) CategoryAdd(w http.ResponseWriter, r *http.Request) {
  component := partials.EditCategory(nil)
  component.Render(r.Context(), w)
}
