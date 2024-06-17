package validation

import (
	"strconv"
)

var coef = [9]int{4, 3, 2, 7, 6, 5, 4, 3, 2}

func IdValidation(ruc string, required bool) string {
	if _, err := strconv.Atoi(ruc); err != nil {
		return "El ruc o cédula no es válido, debe tener solo caracteres numéricos"
	}

	prov, _ := strconv.Atoi(ruc[:2])
	if prov > 24 || prov < 1 {
		return "El ruc o cédula no es válido"
	}

	switch len(ruc) {
	case 10:
		return CedulaValidation(ruc, required)
	case 13:
		return RucValidation(ruc, required)
	default:
		return "El ruc o cédula no es válido"
	}
}

func RucValidation(ruc string, required bool) string {
	if !required && ruc == "" {
		return ""
	}

	if required && ruc == "" {
		return "El ruc es requerido"
	}

	if ruc[10:] != "001" {
		return "El ruc debe terminar en 001"
	}

	if ruc[2] != '9' {
		return CedulaValidation(ruc[:10], required)
	}

	sum := 0
	for i, val := range coef {
		x, _ := strconv.Atoi(string(ruc[i]))
		sum += x * val
	}

	ver := 11 - (sum % 11)
	if ver == 11 {
		ver = 0
	}
	rucVer, _ := strconv.Atoi(string(ruc[9]))

	if rucVer != ver {
		return "El ruc no es válido, digito verificador incorrecto"
	}

	return ""
}

func CedulaValidation(cedula string, required bool) string {
	if !required && cedula == "" {
		return ""
	}

	if required && cedula == "" {
		return "La cédula es requerida"
	}

	sum := 0

	for i := 1; i < 10; i++ {
		coef := 1
		if i%2 != 0 {
			coef = 2
		}

		x, _ := strconv.Atoi(string(cedula[i-1]))

		res := x * coef
		if res > 9 {
			res = res%10 + 1
		}

		sum += res
	}

	ver := sum % 10
	if ver != 0 {
		ver = 10 - ver
	}
	cedulaVer, _ := strconv.Atoi(string(cedula[9]))

	if ver != cedulaVer {
		return "La cédula no es válida"
	}

	return ""
}
