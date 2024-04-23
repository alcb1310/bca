package utils

import (
	"errors"
	"fmt"
	"strconv"
)

func ConvertFloat(value, field string, required bool) (val float64, err error) {
	if required && value == "" {
		return 0, errors.New(fmt.Sprintf("%s es requerido", field))
	} else if value == "" {
		return 0, nil
	}

	val, err = strconv.ParseFloat(value, 64)

	if err != nil {
		return val, errors.New(fmt.Sprintf("%s debe ser un número válido", field))
	}

	if val < 0 {
		return val, errors.New(fmt.Sprintf("%s debe ser un número positivo", field))
	}

	return val, nil
}
