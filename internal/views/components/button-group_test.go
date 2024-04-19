package components

import (
	"context"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestButtonGroup(t *testing.T) {
	t.Run("ButtonGroup", func(t *testing.T) {
		r, w := io.Pipe()
		go func() {
			_ = ButtonGroup().Render(context.Background(), w)
			_ = w.Close()
		}()

		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			t.Fatal(err)
		}

		// check if save button exists
		saveButton := doc.Find(`[data-testid="save-button"]`)
		if saveButton.Length() == 0 {
			t.Error("Save button not found")
		}
		expectedText := "Grabar"
		if saveButton.Text() != expectedText {
			t.Errorf("Expected button text to be %s, but got %s", expectedText, saveButton.Text())
		}

		// check if cancel button exists
		cancelButton := doc.Find(`[data-testid="cancel-button"]`)
		if cancelButton.Length() == 0 {
			t.Error("Cancel button not found")
		}
		expectedText = "Cancelar"

		if cancelButton.Text() != expectedText {
			t.Errorf("Expected button text to be %s, but got %s", expectedText, cancelButton.Text())
		}
	})

}
