package sidebar

import "testing"

func TestCostosUnitarios(t *testing.T) {
	t.Run("Costos unitarios menu has to have all menu options", func(t *testing.T) {
		doc := getMenu("cantidad", "costosunitarios", t)

		menuCantidad := doc.Find(`[data-testid="menu-cantidad"]`)
		if menuCantidad.Length() == 0 {
			t.Error("Cantidad menu not found")
		}
		textMenuCantidad := doc.Find(`[data-testid="text-cantidad"]`)
		if textMenuCantidad.Text() != "Cantidades" {
			t.Errorf("Expected menu text to be %s, but got %s", "Cantidades", textMenuCantidad.Text())
		}

		menuAnalisis := doc.Find(`[data-testid="menu-analisis"]`)
		if menuAnalisis.Length() == 0 {
			t.Error("Analisis menu not found")
		}
		textMenuAnalisis := doc.Find(`[data-testid="text-analisis"]`)
		if textMenuAnalisis.Text() != "Analisis" {
			t.Errorf("Expected menu text to be %s, but got %s", "Analisis", textMenuAnalisis.Text())
		}
	})

	t.Run("Should highlight the selected menu", func(t *testing.T) {
		t.Run("Should highlight cantidad menu", func(t *testing.T) {
			doc := getMenu("cantidad", "costosunitarios", t)
			if !hasActiveMenu("cantidad", doc) {
				t.Error("Cantidad menu not highlighted")
			}
		})

		t.Run("Should highlight analisis menu", func(t *testing.T) {
			doc := getMenu("analisis", "costosunitarios", t)
			if !hasActiveMenu("analisis", doc) {
				t.Error("Analisis menu not highlighted")
			}
		})
	})
}
