package components

import (
	"context"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestSidebar(t *testing.T) {
	t.Run("Sidebar", func(t *testing.T) {
		r, w := io.Pipe()
		go func() {
			_ = SidebarComponent("test").Render(context.Background(), w)
			_ = w.Close()
		}()

		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			t.Fatal(err)
		}

		sidebar := doc.Find(`[data-testid="sidebar"]`)
		if sidebar.Length() == 0 {
			t.Error("Sidebar not found")
		}
	})
}
