package sidebar

import (
	"context"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestTransacciones(t *testing.T) {
	t.Run("Transacciones menu has to have all menu options", func(t *testing.T) {
		doc := getTransaccionesMenu("presupuesto", t)

		// menuPresupuesto := doc.Find(`[data-testid="menu-presupuesto"]`)
		menuPresupuesto := doc.Find(`[data-testid="menu-presupuesto"]`)
		if menuPresupuesto.Length() == 0 {
			t.Error("Presupuesto menu not found")
		}
		textMenuPresupuesto := doc.Find(`[data-testid="text-presupuesto"]`)
		if textMenuPresupuesto.Text() != "Presupuesto" {
			t.Errorf("Expected menu text to be %s, but got %s", "Presupuesto", textMenuPresupuesto.Text())
		}

		menuFacturas := doc.Find(`[data-testid="menu-facturas"]`)
		if menuFacturas.Length() == 0 {
			t.Error("Facturas menu not found")
		}
		textMenuFacturas := doc.Find(`[data-testid="text-facturas"]`)
		if textMenuFacturas.Text() != "Facturas" {
			t.Errorf("Expected menu text to be %s, but got %s", "Facturas", textMenuFacturas.Text())
		}

		menuCierre := doc.Find(`[data-testid="menu-cierre"]`)
		if menuCierre.Length() == 0 {
			t.Error("Cierre menu not found")
		}
		textMenuCierre := doc.Find(`[data-testid="text-cierre"]`)
		if textMenuCierre.Text() != "Cierre Mensual" {
			t.Errorf("Expected menu text to be %s, but got %s", "Cierre Mensual", textMenuCierre.Text())
		}
	})

	t.Run("Should highlight the selected menu", func(t *testing.T) {
		t.Run("Should highlight presupuesto menu", func(t *testing.T) {
			doc := getTransaccionesMenu("presupuesto", t)
			if !hasActiveMenu("presupuesto", doc) {
				t.Error("Presupuesto menu not highlighted")
			}
		})

		t.Run("Should highlight facturas menu", func(t *testing.T) {
			doc := getTransaccionesMenu("facturas", t)
			if !hasActiveMenu("facturas", doc) {
				t.Error("Facturas menu not highlighted")
			}
		})

		t.Run("Should highlight cierre menu", func(t *testing.T) {
			doc := getTransaccionesMenu("cierre", t)
			if !hasActiveMenu("cierre", doc) {
				t.Error("Cierre menu not highlighted")
			}
		})
	})
}

func getTransaccionesMenu(active string, t *testing.T) *goquery.Document {
	r, w := io.Pipe()
	go func() {
		_ = Transacciones(active).Render(context.Background(), w)
		_ = w.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatal(err)
	}
	return doc
}
