package components

import (
	"context"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestPageTitle(t *testing.T) {
	t.Run("PageTitle", func(t *testing.T) {
		title := "test"
		r, w := io.Pipe()
		go func() {
			_ = PageTitle(title).Render(context.Background(), w)
			_ = w.Close()
		}()

		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			t.Fatal(err)
		}

		titleComponent := doc.Find(`[data-testid="page-title"]`)
		if titleComponent.Length() == 0 {
			t.Error("Page title not found")
		}

		if titleComponent.Text() != title {
			t.Errorf("Expected title to be %s, but got %s", title, titleComponent.Text())
		}

		if titleComponent.HasClass("text-xl") == false {
			t.Errorf("Page title should have text-xl class")
		}
	})
}
