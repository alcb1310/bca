package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
	"github.com/alcb1310/bca/internal/views/bca/users"
	"github.com/alcb1310/bca/internal/views/partials"
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

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, _ := uuid.Parse(id)

	if ctx.Id == parsedId {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	u, _ := s.DB.GetUser(parsedId, ctx.CompanyId)

	r.ParseForm()
	if r.Form.Get("name") != "" {
		u.Name = r.Form.Get("name")
	}
	if r.Form.Get("email") != "" {
		if !utils.IsValidEmail(r.Form.Get("email")) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		u.Email = r.Form.Get("email")
	}
	if _, err := s.DB.UpdateUser(u, parsedId, ctx.CompanyId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	users, _ := s.DB.GetAllUsers(ctx.CompanyId)
	component := partials.UsersTable(users, ctx.Id)
	component.Render(r.Context(), w)
}

func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, _ := uuid.Parse(id)

	if ctx.Id == parsedId {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err := s.DB.DeleteUser(parsedId, ctx.CompanyId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	users, _ := s.DB.GetAllUsers(ctx.CompanyId)
	component := partials.UsersTable(users, ctx.Id)
	component.Render(r.Context(), w)
}

func (s *Server) SingleUserGet(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, _ := uuid.Parse(id)

	if ctx.Id == parsedId {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	users, _ := s.DB.GetAllUsers(ctx.CompanyId)
	component := partials.UsersTable(users, ctx.Id)
	component.Render(r.Context(), w)
}

func (s *Server) UsersChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

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
}

func (s *Server) UsersCreate(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

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

	if u.Email == "" || !utils.IsValidEmail(u.Email) {
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

	users, _ := s.DB.GetAllUsers(ctx.CompanyId)
	component := partials.UsersTable(users, ctx.Id)
	component.Render(r.Context(), w)
}

func (s *Server) UsersTableDisplay(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

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
