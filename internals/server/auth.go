package server

import (
	"log/slog"
	"net/http"

	"github.com/alcb1310/bca/externals/views/register"
	"github.com/alcb1310/bca/internals/types"
	"github.com/alcb1310/bca/internals/utils"
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

	if len(formErrors) > 0 {
		slog.Debug("RegisterForm: if", "formErrors", formErrors)
		return renderPage(w, r, register.Register(fields, formErrors))
	}

	if err := s.DB.CreateCompany(company, user); err != nil {
		renderPage(w, r, register.Register(fields, formErrors))
		return err
	}

	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

func (s *Service) Login(w http.ResponseWriter, r *http.Request) error {
	user := types.User{}
	var err error

	r.ParseForm()
	fields := make(map[string]string)
	fields["email"] = r.Form.Get("email")
	fields["password"] = r.Form.Get("password")

	formErrors := validation.ValidateLogin(fields)
	if len(formErrors) > 0 {
		slog.Debug("Login Process", "formErrors", formErrors)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Credenciales inválidas"))
		return nil
	}

	if user, err = s.DB.Login(fields["email"], fields["password"]); err != nil {
		slog.Debug("Login", "error", "invalid credentials")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Credenciales inválidas"))
		return nil
	}

	jwtToken, err := utils.GenerateJWT(user)
	if err != nil {
		slog.Error("Login: Unable to generate JWT Token", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	http.SetCookie(w, utils.GenerateCookie(jwtToken))
	http.Redirect(w, r, "/bca", http.StatusFound)
	return nil
}
