package utils

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func PrintFloat(value float64) string {
	p := message.NewPrinter(language.LatinAmericanSpanish)
	return p.Sprintf("%.2f", value)
}
