package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
	"bca-go-final/internal/views/bca/users"
	"bca-go-final/internal/views/partials"
)

func (s *Server) Profile(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	user, _ := s.DB.GetUser(ctx.Id, ctx.CompanyId)

	component := users.ProfileView(user)
	component.Render(r.Context(), w)
}

func (s *Server) Admin(w http.ResponseWriter, r *http.Request) {
	component := users.AdminView()
	component.Render(r.Context(), w)
}

func (s *Server) ChangePassword(w http.ResponseWriter, r *http.Request) {
	component := users.ChangePasswordView()
	component.Render(r.Context(), w)
}

func (s *Server) SingleUser(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, _ := uuid.Parse(id)

	if ctx.Id == parsedId {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		if err := s.DB.DeleteUser(parsedId, ctx.CompanyId); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.MethodPut:
		u, _ := s.DB.GetUser(parsedId, ctx.CompanyId)

		r.ParseForm()
		if r.Form.Get("name") != "" {
			u.Name = r.Form.Get("name")
		}
		if r.Form.Get("email") != "" {
			u.Email = r.Form.Get("email")
		}
		if _, err := s.DB.UpdateUser(u, parsedId, ctx.CompanyId); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	users, _ := s.DB.GetAllUsers(ctx.CompanyId)
	component := partials.UsersTable(users, ctx.Id)
	component.Render(r.Context(), w)
}

func (s *Server) UsersTable(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	switch r.Method {
	case http.MethodPut:
		r.ParseForm()
		pass := r.Form.Get("password")

		if pass == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, err := s.DB.UpdatePassword(pass, ctx.Id, ctx.CompanyId); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return

	case http.MethodPost:
		u := &types.UserCreate{}
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		u.Name = r.Form.Get("name")
		u.Email = r.Form.Get("email")
		u.Password = r.Form.Get("password")
		u.RoleId = r.Form.Get("role")
		u.CompanyId = ctx.CompanyId

		if u.Email != "" && !utils.IsValidEmail(u.Email) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if u.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if u.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if u.RoleId == "" {
			u.RoleId = "a"
		}

		_, err = s.DB.CreateUser(u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	users, _ := s.DB.GetAllUsers(ctx.CompanyId)
	component := partials.UsersTable(users, ctx.Id)
	component.Render(r.Context(), w)
}

func (s *Server) UserAdd(w http.ResponseWriter, r *http.Request) {
	component := partials.EditUser(nil)
	component.Render(r.Context(), w)
}

func (s *Server) UserEdit(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, _ := uuid.Parse(id)
	u, _ := s.DB.GetUser(parsedId, ctx.CompanyId)

	component := partials.EditUser(&u)
	component.Render(r.Context(), w)
}
