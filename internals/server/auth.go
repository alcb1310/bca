package server

import (
	"log/slog"
	"net/http"

	"github.com/alcb1310/bca/externals/views/register"
	"github.com/alcb1310/bca/internals/types"
	"github.com/alcb1310/bca/internals/validation"
)

func (s *Service) Register(w http.ResponseWriter, r *http.Request) error {
	m := make(map[string]string)
	return renderPage(w, r, register.Register(m, m))
}

func (s *Service) RegisterForm(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	fields := make(map[string]string)

	fields["ruc"] = r.Form.Get("ruc")
	fields["name"] = r.Form.Get("name")
	fields["employees"] = r.Form.Get("employees")
	fields["email"] = r.Form.Get("email")
	fields["password"] = r.Form.Get("password")
	fields["username"] = r.Form.Get("username")

	company := &types.Company{
		IsActive: true,
	}
	user := &types.CreateUser{}
	formErrors := validation.ValidateRegisterForm(fields, company, user)

	slog.Debug("Register", "company", company)
	slog.Debug("Register", "user", user)

	if len(formErrors) > 0 {
		slog.Debug("RegisterForm: if", "formErrors", formErrors)
		return renderPage(w, r, register.Register(fields, formErrors))
	}

	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}
