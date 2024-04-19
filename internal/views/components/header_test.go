package components

import (
	"context"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestHeader(t *testing.T) {
	t.Run("ButtonGroup", func(t *testing.T) {
		r, w := io.Pipe()
		go func() {
			_ = Header().Render(context.Background(), w)
			_ = w.Close()
		}()

		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			t.Fatal(err)
		}

		headerTitle := doc.Find(`[data-testid="application-name"]`)
		if headerTitle.Length() == 0 {
			t.Error("Header title not found")
		}
		want := "Sistema Control Presupuestario"
		if headerTitle.Text() != want {
			t.Errorf("Expected button text to be %s, but got %s", want, headerTitle.Text())
		}

		userIcon := doc.Find(`[data-testid="user-icon"]`)
		if userIcon.Length() == 0 {
			t.Errorf("User icon not found")
		}
		if !userIcon.HasClass("text-blue-gray-100") {
			t.Errorf("User icon should have text-blue-gray-100 class")
		}
		aria, _ := userIcon.Attr("aria-label")
		if aria != "Usuario" {
			t.Errorf("User icon should have aria-label 'Usuario'")
		}

		helpIcon := doc.Find(`[data-testid="help-icon"]`)
		if helpIcon.Length() == 0 {
			t.Errorf("Help icon not found")
		}
		aria, _ = helpIcon.Attr("aria-label")
		if aria != "Ayuda" {
			t.Errorf("Help icon should have aria-label 'Ayuda'")
		}
		target, _ := helpIcon.Attr("target")
		if target != "_blank" {
			t.Errorf("Help icon should have target '_blank'")
		}

		logoutIcon := doc.Find(`[data-testid="logout-icon"]`)
		if logoutIcon.Length() == 0 {
			t.Errorf("Logout icon not found")
		}
		aria, _ = logoutIcon.Attr("aria-label")
		if aria != "Cerrar Sesión" {
			t.Errorf("Logout icon should have aria-label 'Cerrar Sesión'")
		}
		post, _ := logoutIcon.Attr("hx-post")
		if post != "/bca/logout" {
			t.Errorf("Logout icon should have hx-post attribute")
		}

		userContext := doc.Find(`[data-testid="user-context"]`)
		if userContext.Length() == 0 {
			t.Errorf("User context not found")
		}
		if !userContext.HasClass("hidden") {
			t.Errorf("User context should be hidden")
		}
	})
}
