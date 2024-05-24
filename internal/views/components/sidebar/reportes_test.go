package sidebar

import (
	"testing"
)

func TestReportes(t *testing.T) {
	t.Run("Reportes menu has to have all menu options", func(t *testing.T) {
		doc := getMenu("actual", "reportes", t)

		menuActual := doc.Find(`[data-testid="menu-actual"]`)
		if menuActual.Length() == 0 {
			t.Error("Actual menu not found")
		}
		textMenuActual := doc.Find(`[data-testid="text-actual"]`)
		if textMenuActual.Text() != "Actual" {
			t.Errorf("Expected menu text to be %s, but got %s", "Actual", textMenuActual.Text())
		}

		menuCuadre := doc.Find(`[data-testid="menu-balance"]`)
		if menuCuadre.Length() == 0 {
			t.Error("Cuadre menu not found")
		}
		textMenuCuadre := doc.Find(`[data-testid="text-balance"]`)
		if textMenuCuadre.Text() != "Cuadre" {
			t.Errorf("Expected menu text to be %s, but got %s", "Cuadre", textMenuCuadre.Text())
		}

		menuGastado := doc.Find(`[data-testid="menu-gastado"]`)
		if menuGastado.Length() == 0 {
			t.Error("Gastado menu not found")
		}
		textMenuGastado := doc.Find(`[data-testid="text-gastado"]`)
		if textMenuGastado.Text() != "Gastado por Partida" {
			t.Errorf("Expected menu text to be %s, but got %s", "Gastado", textMenuGastado.Text())
		}

		menuHistorico := doc.Find(`[data-testid="menu-historico"]`)
		if menuHistorico.Length() == 0 {
			t.Error("Historico menu not found")
		}
		textMenuHistorico := doc.Find(`[data-testid="text-historico"]`)
		if textMenuHistorico.Text() != "Histórico" {
			t.Errorf("Expected menu text to be %s, but got %s", "Historico", textMenuHistorico.Text())
		}
	})

	t.Run("Should highlight the selected menu", func(t *testing.T) {
		t.Run("Should highlight actual menu", func(t *testing.T) {
			doc := getMenu("actual", "reportes", t)
			if !hasActiveMenu("actual", doc) {
				t.Error("Actual menu not highlighted")
			}
		})

		t.Run("Should highlight balance menu", func(t *testing.T) {
			doc := getMenu("balance", "reportes", t)
			if !hasActiveMenu("balance", doc) {
				t.Error("Cuadre menu not highlighted")
			}
		})

		t.Run("Should highlight gastado menu", func(t *testing.T) {
			doc := getMenu("gastado", "reportes", t)
			if !hasActiveMenu("gastado", doc) {
				t.Error("Gastado menu not highlighted")
			}
		})

		t.Run("Should highlight historico menu", func(t *testing.T) {
			doc := getMenu("historico", "reportes", t)
			if !hasActiveMenu("historico", doc) {
				t.Error("Historico menu not highlighted")
			}
		})
	})
}
