package components

import (
	"context"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestInput(t *testing.T) {
	t.Run("Input with no value", func(t *testing.T) {
		inputType := "text"
		name := "name"
		value := ""
		id := "id"
		placeholder := "placeholder"
		r, w := io.Pipe()
		go func() {
			_ = Input(inputType, placeholder, id, name, value).Render(context.Background(), w)
			_ = w.Close()
		}()

		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			t.Fatal(err)
		}

		input := doc.Find(`[data-testid="input"]`)
		if input.Length() == 0 {
			t.Error("Input not found")
		}

		label := doc.Find(`[data-testid="label"]`)
		if label.Length() != 0 {
			t.Error("Label should not be found")
		}

		iType, _ := input.Attr("type")
		if iType != inputType {
			t.Errorf("Input type should be %s, but got %s", inputType, iType)
		}

		iName, _ := input.Attr("name")
		if iName != name {
			t.Errorf("Input name should be %s, but got %s", name, iName)
		}

		iId, _ := input.Attr("id")
		if iId != id {
			t.Errorf("Input id should be %s, but got %s", id, iId)
		}

		iPlaceholder, _ := input.Attr("placeholder")
		if iPlaceholder != placeholder {
			t.Errorf("Input placeholder should be %s, but got %s", placeholder, iPlaceholder)
		}

		iValue, _ := input.Attr("value")
		if iValue != value {
			t.Errorf("Input value should be %s, but got %s", value, iValue)
		}
	})

	t.Run("Input with value", func(t *testing.T) {
		inputType := "text"
		name := "name"
		value := "value"
		id := "id"
		placeholder := "placeholder"
		r, w := io.Pipe()
		go func() {
			_ = Input(inputType, placeholder, id, name, value).Render(context.Background(), w)
			_ = w.Close()
		}()

		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			t.Fatal(err)
		}

		input := doc.Find(`[data-testid="input"]`)
		if input.Length() == 0 {
			t.Error("Input not found")
		}

		label := doc.Find(`[data-testid="label"]`)
		if label.Length() == 0 {
			t.Error("Label not found")
		}

		iValue, _ := input.Attr("value")
		if iValue != value {
			t.Errorf("Input value should be %s, but got %s", value, iValue)
		}

		lText := label.Text()
		if lText != placeholder {
			t.Errorf("Label text should be %s, but got %s", placeholder, lText)
		}
	})
}
