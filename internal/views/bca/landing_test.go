package bca

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestLanding(t *testing.T) {
	r, w := io.Pipe()
	go func() {
		_ = LandingPage("test").Render(context.Background(), w)
		_ = w.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatal(err)
	}

	if doc.Find(`[data-testid="welcome"]`).Text() != "test" {
		t.Fatal(fmt.Sprintf("Expected to welcome %s, and got %s", "test", doc.Find(`[data-testid="welcome"]`).Text()))
	}
}
