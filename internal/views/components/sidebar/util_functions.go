package sidebar

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func hasActiveMenu(menu string, doc *goquery.Document) bool {
	testId := fmt.Sprintf(`[data-testid="menu-%s"]`, menu)
	return doc.Find(testId).HasClass("text-blue-gray-100")
}

func getMenu(active, menu string, t *testing.T) *goquery.Document {
	if menu != "parametros" && menu != "transacciones" && menu != "reportes" {
		t.Fatal("Menu not found")
	}
	r, w := io.Pipe()
	go func() {
		switch menu {
		case "parametros":
			_ = Parametros(active).Render(context.Background(), w)

		case "transacciones":
			_ = Transacciones(active).Render(context.Background(), w)

		case "reportes":
			_ = Reportes(active).Render(context.Background(), w)

		}
		_ = w.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatal(err)
	}
	return doc
}
