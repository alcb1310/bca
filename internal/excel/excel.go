package excel

import (
	"bca-go-final/internal/database"
	"bca-go-final/internal/types"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

func Balance(companyId, projectId uuid.UUID, date time.Time, db database.Service) *excelize.File {
	f := excelize.NewFile()
	id := uuid.New()
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet("cuadre")
	if err != nil {
		fmt.Println(err)
		return f
	}
	f.SetActiveSheet(index)

	f.DeleteSheet("Sheet1")

	f.SetCellValue("cuadre", "A1", "Fecha")
	f.SetCellValue("cuadre", "B1", "Proveedor")
	f.SetCellValue("cuadre", "C1", "Factura")
	f.SetCellValue("cuadre", "D1", "Valor")

	balance := db.GetBalance(companyId, projectId, date)

	row := 2
	for _, invoice := range balance.Invoices {
		cell := fmt.Sprintf("A%d", row)
		year := invoice.InvoiceDate.Year()
		month := int(invoice.InvoiceDate.Month())
		day := invoice.InvoiceDate.Day()
		// f.SetCellValue("cuadre", cell, invoice.InvoiceDate.Format("2006-01-02"))
		f.SetCellFormula("cuadre", cell, fmt.Sprintf("=DATE(%d,%d,%d)", year, month, day))
		cell = fmt.Sprintf("B%d", row)
		f.SetCellValue("cuadre", cell, invoice.Supplier.Name)
		cell = fmt.Sprintf("C%d", row)
		f.SetCellValue("cuadre", cell, invoice.InvoiceNumber)
		cell = fmt.Sprintf("D%d", row)
		f.SetCellFloat("cuadre", cell, invoice.InvoiceTotal, 2, 64)

		row++
	}

	f.MergeCell("cuadre", fmt.Sprintf("A%d", row), fmt.Sprintf("C%d", row))
	f.SetCellValue("cuadre", fmt.Sprintf("A%d", row), "TOTAL")

	// styling the data
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#D9D9D9"},
		},
	})

	numberStyle, _ := f.NewStyle(&excelize.Style{
		NumFmt: 4,
	})

	exp := "yyyy-mm-dd"
	dateStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
		},
		CustomNumFmt: &exp,
	})

	f.SetColWidth("cuadre", "B", "B", 33)
	f.SetColWidth("cuadre", "C", "C", 20)

	f.SetCellStyle("cuadre", "A1", "D1", titleStyle)
	f.SetCellStyle("cuadre", "D2", fmt.Sprintf("D%d", row), numberStyle)
	f.SetCellStyle("cuadre", "A2", fmt.Sprintf("A%d", row), dateStyle)

	f.SetCellFormula("cuadre", fmt.Sprintf("D%d", row), fmt.Sprintf("=SUM(D2:D%d)", row-1))

	if err := f.SaveAs("./public/" + id.String() + ".xlsx"); err != nil {
		log.Println(err)
	}

	return f
}

func Actual(companyId, projectId uuid.UUID, budgets []types.GetBudget, date *time.Time, db database.Service) *excelize.File {
	var d time.Time
	f := excelize.NewFile()
	id := uuid.New()
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()

	if date == nil {
		d = time.Now()
	} else {
		d = *date
	}

	// Create a new sheet.
	index, err := f.NewSheet("actual")
	if err != nil {
		fmt.Println(err)
		return f
	}
	f.SetActiveSheet(index)

	f.DeleteSheet("Sheet1")

	pageTitleStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold: true,
			Size: 18,
		},
	})
	colTitleStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})
	descTitleStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})

	p, _ := db.GetProject(projectId, companyId)

	f.SetCellValue("actual", "A1", "CONTROL PRESUPUESTARIO")
	f.MergeCell("actual", "A1", "K1")
	f.SetCellStyle("actual", "A1", "K1", pageTitleStyle)

	exp := "yyyy-mm-dd"
	dateStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
		},
		CustomNumFmt: &exp,
	})

	f.SetColWidth("actual", "A", "A", 11.50)
	f.SetColWidth("actual", "B", "B", 40)
	f.SetColWidth("actual", "C", "K", 13.50)

	f.SetCellValue("actual", "A3", "Fecha")
	f.SetCellFormula("actual", "B3", fmt.Sprintf("=DATE(%d,%d,%d)", d.Year(), int(d.Month()), d.Day()))
	f.SetCellStyle("actual", "B3", "B3", dateStyle)
	f.SetCellValue("actual", "A4", "Proyecto")
	f.SetCellValue("actual", "B4", p.Name)
	f.SetCellStyle("actual", "A3", "A4", descTitleStyle)

	f.SetCellValue("actual", "A6", "Código")
	f.MergeCell("actual", "A6", "A7")

	f.SetCellValue("actual", "B6", "Partida")
	f.MergeCell("actual", "B6", "B7")

	f.SetCellValue("actual", "C6", "Inicial")
	f.MergeCell("actual", "C6", "E6")
	f.SetCellValue("actual", "C7", "Cantidad")
	f.SetCellValue("actual", "D7", "Costo")
	f.SetCellValue("actual", "E7", "Total")

	f.SetCellValue("actual", "F6", "Rendido")
	f.MergeCell("actual", "F6", "G6")
	f.SetCellValue("actual", "F7", "Cantidad")
	f.SetCellValue("actual", "G7", "Total")

	f.SetCellValue("actual", "H6", "Por Gastar")
	f.MergeCell("actual", "H6", "J6")
	f.SetCellValue("actual", "H7", "Cantidad")
	f.SetCellValue("actual", "I7", "Costo")
	f.SetCellValue("actual", "J7", "Total")

	f.SetCellValue("actual", "K6", "Actualizado")
	f.MergeCell("actual", "K6", "K7")

	f.SetCellStyle("actual", "A6", "K7", colTitleStyle)

	row := 8

	level1Style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "#000000", Style: 6},
			{Type: "bottom", Color: "#000000", Style: 6},
		},
		NumFmt: 4,
	})

	level2Style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "#000000", Style: 2},
			{Type: "bottom", Color: "#000000", Style: 6},
		},
		NumFmt: 4,
	})

	level3Style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
		NumFmt: 4,
	})

	defaultStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: false,
		},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#000000", Style: 4},
		},
		NumFmt: 4,
	})

	for _, budget := range budgets {

		switch budget.BudgetItem.Level {
		case 1:
			row++
			f.SetCellStyle("actual", fmt.Sprintf("A%d", row), fmt.Sprintf("K%d", row), level1Style)
		case 2:
			row++
			f.SetCellStyle("actual", fmt.Sprintf("A%d", row), fmt.Sprintf("K%d", row), level2Style)
		case 3:
			row++
			f.SetCellStyle("actual", fmt.Sprintf("A%d", row), fmt.Sprintf("K%d", row), level3Style)
		default:
			f.SetCellStyle("actual", fmt.Sprintf("A%d", row), fmt.Sprintf("K%d", row), defaultStyle)
		}

		f.SetCellValue("actual", fmt.Sprintf("A%d", row), budget.BudgetItem.Code)
		f.SetCellValue("actual", fmt.Sprintf("B%d", row), budget.BudgetItem.Name)

		if budget.InitialQuantity != nil {
			f.SetCellFloat("actual", fmt.Sprintf("C%d", row), *budget.InitialQuantity, 2, 64)
			f.SetCellFloat("actual", fmt.Sprintf("D%d", row), *budget.InitialCost, 2, 64)
		}
		f.SetCellFloat("actual", fmt.Sprintf("E%d", row), budget.InitialTotal, 2, 64)

		if budget.SpentQuantity != nil {
			f.SetCellFloat("actual", fmt.Sprintf("F%d", row), *budget.SpentQuantity, 2, 64)
		}
		f.SetCellFloat("actual", fmt.Sprintf("G%d", row), budget.SpentTotal, 2, 64)

		if budget.RemainingQuantity != nil {
			f.SetCellFloat("actual", fmt.Sprintf("H%d", row), *budget.RemainingQuantity, 2, 64)
			f.SetCellFloat("actual", fmt.Sprintf("I%d", row), *budget.RemainingCost, 2, 64)
		}
		f.SetCellFloat("actual", fmt.Sprintf("J%d", row), budget.RemainingTotal, 2, 64)
		f.SetCellFloat("actual", fmt.Sprintf("K%d", row), budget.UpdatedBudget, 2, 64)

		row++
	}

	if err := f.SaveAs("./public/" + id.String() + ".xlsx"); err != nil {
		log.Println(err)
	}
	return f
}