package components

import (
	"context"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestDrawerTitle(t *testing.T) {
	t.Run("DrawerTitle", func(t *testing.T) {
		title := "test"
		r, w := io.Pipe()
		go func() {
			_ = DrawerTitle(title).Render(context.Background(), w)
			_ = w.Close()
		}()

		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			t.Fatal(err)
		}

		titleComponent := doc.Find(`[data-testid="drawer-title"]`)
		if titleComponent.Length() == 0 {
			t.Error("Drawer title not found")
		}

		if titleComponent.Text() != title {
			t.Errorf("Expected title to be %s, but got %s", title, titleComponent.Text())
		}
	})
}
