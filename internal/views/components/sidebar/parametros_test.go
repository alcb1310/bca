package sidebar

import (
	"testing"
)

func TestParametros(t *testing.T) {
	t.Run("Parametros menu has to have all menu options", func(t *testing.T) {
		doc := getMenu("partida", "parametros", t)

		menuPartida := doc.Find(`[data-testid="menu-partida"]`)
		if menuPartida.Length() == 0 {
			t.Error("Partidas menu not found")
		}
		textMenuPartida := doc.Find(`[data-testid="text-partida"]`)
		if textMenuPartida.Text() != "Partidas" {
			t.Errorf("Expected menu text to be %s, but got %s", "Partidas", textMenuPartida.Text())
		}

		menuCategorias := doc.Find(`[data-testid="menu-categorias"]`)
		if menuCategorias.Length() == 0 {
			t.Error("Categorias menu not found")
		}
		textMenuCategorias := doc.Find(`[data-testid="text-categorias"]`)
		if textMenuCategorias.Text() != "Categorias" {
			t.Errorf("Expected menu text to be %s, but got %s", "Categorias", textMenuCategorias.Text())
		}

		menuMateriales := doc.Find(`[data-testid="menu-materiales"]`)
		if menuMateriales.Length() == 0 {
			t.Error("Materiales menu not found")
		}
		textMenuMateriales := doc.Find(`[data-testid="text-materiales"]`)
		if textMenuMateriales.Text() != "Materiales" {
			t.Errorf("Expected menu text to be %s, but got %s", "Materiales", textMenuMateriales.Text())
		}

		menuProyecto := doc.Find(`[data-testid="menu-proyecto"]`)
		if menuProyecto.Length() == 0 {
			t.Error("Proyectos menu not found")
		}
		textMenuProyecto := doc.Find(`[data-testid="text-proyecto"]`)
		if textMenuProyecto.Text() != "Proyectos" {
			t.Errorf("Expected menu text to be %s, but got %s", "Proyectos", textMenuProyecto.Text())
		}

		menuProveedor := doc.Find(`[data-testid="menu-proveedor"]`)
		if menuProveedor.Length() == 0 {
			t.Error("Proveedores menu not found")
		}
		textMenuProveedor := doc.Find(`[data-testid="text-proveedor"]`)
		if textMenuProveedor.Text() != "Proveedores" {
			t.Errorf("Expected menu text to be %s, but got %s", "Proveedores", textMenuProveedor.Text())
		}

		menuRubros := doc.Find(`[data-testid="menu-rubros"]`)
		if menuRubros.Length() == 0 {
			t.Error("Rubros menu not found")
		}
		textMenuRubros := doc.Find(`[data-testid="text-rubros"]`)
		if textMenuRubros.Text() != "Rubros" {
			t.Errorf("Expected menu text to be %s, but got %s", "Rubros", textMenuRubros.Text())
		}
	})

	t.Run("Should highlight active menu", func(t *testing.T) {
		t.Run("Should hightlight partida menu", func(t *testing.T) {
			active := "partida"
			doc := getMenu(active, "parametros", t)
			if !hasActiveMenu(active, doc) {
				t.Error("Partida menu not highlighted")
			}
		})

		t.Run("Should hightlight categorias menu", func(t *testing.T) {
			active := "categorias"
			doc := getMenu(active, "parametros", t)
			if !hasActiveMenu(active, doc) {
				t.Error("Partida menu not highlighted")
			}

		})

		t.Run("Should hightlight materiales menu", func(t *testing.T) {
			active := "materiales"
			doc := getMenu(active, "parametros", t)
			if !hasActiveMenu(active, doc) {
				t.Error("Partida menu not highlighted")
			}

		})

		t.Run("Should hightlight proyecto menu", func(t *testing.T) {
			active := "proyecto"
			doc := getMenu(active, "parametros", t)
			if !hasActiveMenu(active, doc) {
				t.Error("Partida menu not highlighted")
			}

		})

		t.Run("Should hightlight proveedor menu", func(t *testing.T) {
			active := "proveedor"
			doc := getMenu(active, "parametros", t)
			if !hasActiveMenu(active, doc) {
				t.Error("Partida menu not highlighted")
			}

		})

		t.Run("Should hightlight rubros menu", func(t *testing.T) {
			active := "rubros"
			doc := getMenu(active, "parametros", t)
			if !hasActiveMenu(active, doc) {
				t.Error("Partida menu not highlighted")
			}

		})
	})
}
