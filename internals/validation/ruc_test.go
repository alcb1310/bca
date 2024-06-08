package validation_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alcb1310/bca/internals/validation"
)

var testData = []string{
	"0993366551001",
	"1793189551001",
	"1793189543001",
	"0993366570001",
	"1793189537001",
	"1793189538001",
	"1793189552001",
	"1793189544001",
	"1793189546001",
	"1793189549001",
	"0993366563001",
	"0691784262001",
	"0993366554001",
	"0993366556001",
	"0993366559001",
	"1291789518001",
	"0993366568001",
	"0993366553001",
	"1793189536001",
}

var validData = []string{
	"0195088608001",
	"1704749652001",
	"1791838300001",
}

func TestRucValidation(t *testing.T) {
	for i, ruc := range testData {
		expectedResponse := "El ruc no es válido, digito verificador incorrecto"
		res := validation.RucValidation(ruc, true)

		assert.Equal(t, expectedResponse, res)
		if res != expectedResponse {
			fmt.Println(fmt.Sprintf("testData number %d", i))
		}
	}

	for _, ruc := range validData {
		res := validation.RucValidation(ruc, true)
		assert.Empty(t, res)
	}
}
