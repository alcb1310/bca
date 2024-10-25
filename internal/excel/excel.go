package excel

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"

	"github.com/alcb1310/bca/internal/database"
	"github.com/alcb1310/bca/internal/types"
)

func Balance(companyId, projectId uuid.UUID, date time.Time, db database.Service) *excelize.File {
	f := excelize.NewFile()
	id := uuid.New()
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("Failed to balance", "err", err)
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet("cuadre")
	if err != nil {
		fmt.Println(err)
		return f
	}
	f.SetActiveSheet(index)
	_ = f.DeleteSheet("Sheet1")

	_ = f.SetCellValue("cuadre", "A1", "Fecha")
	_ = f.SetCellValue("cuadre", "B1", "Proveedor")
	_ = f.SetCellValue("cuadre", "C1", "Factura")
	_ = f.SetCellValue("cuadre", "D1", "Valor")

	balance := db.GetBalance(companyId, projectId, date)

	row := 2
	for _, invoice := range balance.Invoices {
		cell := fmt.Sprintf("A%d", row)
		year := invoice.InvoiceDate.Year()
		month := int(invoice.InvoiceDate.Month())
		day := invoice.InvoiceDate.Day()
		// _=f.SetCellValue("cuadre", cell, invoice.InvoiceDate.Format("2006-01-02"))
		_ = f.SetCellFormula("cuadre", cell, fmt.Sprintf("=DATE(%d,%d,%d)", year, month, day))
		cell = fmt.Sprintf("B%d", row)
		_ = f.SetCellValue("cuadre", cell, invoice.Supplier.Name)
		cell = fmt.Sprintf("C%d", row)
		_ = f.SetCellValue("cuadre", cell, invoice.InvoiceNumber)
		cell = fmt.Sprintf("D%d", row)
		_ = f.SetCellFloat("cuadre", cell, invoice.InvoiceTotal, 2, 64)

		row++
	}

	_ = f.MergeCell("cuadre", fmt.Sprintf("A%d", row), fmt.Sprintf("C%d", row))
	_ = f.SetCellValue("cuadre", fmt.Sprintf("A%d", row), "TOTAL")

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

	_ = f.SetColWidth("cuadre", "B", "B", 33)
	_ = f.SetColWidth("cuadre", "C", "C", 20)

	_ = f.SetCellStyle("cuadre", "A1", "D1", titleStyle)
	_ = f.SetCellStyle("cuadre", "D2", fmt.Sprintf("D%d", row), numberStyle)
	_ = f.SetCellStyle("cuadre", "A2", fmt.Sprintf("A%d", row), dateStyle)

	_ = f.SetCellFormula("cuadre", fmt.Sprintf("D%d", row), fmt.Sprintf("=SUM(D2:D%d)", row-1))

	if err := f.SaveAs("./public/" + id.String() + ".xlsx"); err != nil {
		slog.Error("Failed to export file", "err", err)
	}

	return f
}

func Actual(companyId, projectId uuid.UUID, budgets []types.GetBudget, date *time.Time, db database.Service) *excelize.File {
	var d time.Time
	f := excelize.NewFile()
	id := uuid.New()
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("Failed to actual", "err", err)
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

	_ = f.DeleteSheet("Sheet1")

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
	areaStyle, _ := f.NewStyle(&excelize.Style{
		NumFmt: 4,
	})

	p, _ := db.GetProject(projectId, companyId)

	_ = f.SetCellValue("actual", "A1", "CONTROL PRESUPUESTARIO")
	_ = f.MergeCell("actual", "A1", "K1")
	_ = f.SetCellStyle("actual", "A1", "K1", pageTitleStyle)

	exp := "yyyy-mm-dd"
	dateStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
		},
		CustomNumFmt: &exp,
	})

	_ = f.SetColWidth("actual", "A", "A", 11.50)
	_ = f.SetColWidth("actual", "B", "B", 40)
	_ = f.SetColWidth("actual", "C", "K", 13.50)

	_ = f.SetCellValue("actual", "A3", "Fecha")
	_ = f.SetCellFormula("actual", "B3", fmt.Sprintf("=DATE(%d,%d,%d)", d.Year(), int(d.Month()), d.Day()))
	_ = f.SetCellStyle("actual", "B3", "B3", dateStyle)
	_ = f.SetCellValue("actual", "A4", "Proyecto")
	_ = f.SetCellValue("actual", "B4", p.Name)
	_ = f.SetCellStyle("actual", "A3", "A6", descTitleStyle)
	_ = f.SetCellValue("actual", "A5", "Area Bruta")
	_ = f.SetCellFloat("actual", "B5", p.GrossArea, 2, 64)
	_ = f.SetCellValue("actual", "A6", "Area Util")
	_ = f.SetCellFloat("actual", "B6", p.NetArea, 2, 64)
	_ = f.SetCellStyle("actual", "B5", "B6", areaStyle)

	_ = f.SetCellValue("actual", "A8", "Código")
	_ = f.MergeCell("actual", "A8", "A9")

	_ = f.SetCellValue("actual", "B8", "Partida")
	_ = f.MergeCell("actual", "B8", "B9")

	_ = f.SetCellValue("actual", "C8", "Inicial")
	_ = f.MergeCell("actual", "C8", "E8")
	_ = f.SetCellValue("actual", "C9", "Cantidad")
	_ = f.SetCellValue("actual", "D9", "Costo")
	_ = f.SetCellValue("actual", "E9", "Total")

	_ = f.SetCellValue("actual", "F8", "Rendido")
	_ = f.MergeCell("actual", "F8", "G8")
	_ = f.SetCellValue("actual", "F9", "Cantidad")
	_ = f.SetCellValue("actual", "G9", "Total")

	_ = f.SetCellValue("actual", "H8", "Por Gastar")
	_ = f.MergeCell("actual", "H8", "J8")
	_ = f.SetCellValue("actual", "H9", "Cantidad")
	_ = f.SetCellValue("actual", "I9", "Costo")
	_ = f.SetCellValue("actual", "J9", "Total")

	_ = f.SetCellValue("actual", "K8", "Actualizado")
	_ = f.MergeCell("actual", "K8", "K9")

	_ = f.SetCellStyle("actual", "A8", "K9", colTitleStyle)

	row := 10

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
			_ = f.SetCellStyle("actual", fmt.Sprintf("A%d", row), fmt.Sprintf("K%d", row), level1Style)
		case 2:
			row++
			_ = f.SetCellStyle("actual", fmt.Sprintf("A%d", row), fmt.Sprintf("K%d", row), level2Style)
		case 3:
			row++
			_ = f.SetCellStyle("actual", fmt.Sprintf("A%d", row), fmt.Sprintf("K%d", row), level3Style)
		default:
			_ = f.SetCellStyle("actual", fmt.Sprintf("A%d", row), fmt.Sprintf("K%d", row), defaultStyle)
		}

		_ = f.SetCellValue("actual", fmt.Sprintf("A%d", row), budget.BudgetItem.Code)
		_ = f.SetCellValue("actual", fmt.Sprintf("B%d", row), budget.BudgetItem.Name)

		if budget.InitialQuantity.Valid {
			_ = f.SetCellFloat("actual", fmt.Sprintf("C%d", row), budget.InitialQuantity.Float64, 2, 64)
			_ = f.SetCellFloat("actual", fmt.Sprintf("D%d", row), budget.InitialCost.Float64, 2, 64)
		}
		_ = f.SetCellFloat("actual", fmt.Sprintf("E%d", row), budget.InitialTotal, 2, 64)

		if budget.SpentQuantity.Valid {
			_ = f.SetCellFloat("actual", fmt.Sprintf("F%d", row), budget.SpentQuantity.Float64, 2, 64)
		}
		_ = f.SetCellFloat("actual", fmt.Sprintf("G%d", row), budget.SpentTotal, 2, 64)

		if budget.RemainingQuantity.Valid {
			_ = f.SetCellFloat("actual", fmt.Sprintf("H%d", row), budget.RemainingQuantity.Float64, 2, 64)
			_ = f.SetCellFloat("actual", fmt.Sprintf("I%d", row), budget.RemainingCost.Float64, 2, 64)
		}
		_ = f.SetCellFloat("actual", fmt.Sprintf("J%d", row), budget.RemainingTotal, 2, 64)
		_ = f.SetCellFloat("actual", fmt.Sprintf("K%d", row), budget.UpdatedBudget, 2, 64)

		row++
	}

	if err := f.SaveAs("./public/" + id.String() + ".xlsx"); err != nil {
		slog.Error("Failed to export", "err", err)
	}
	return f
}

func Spent(project types.Project, data []types.Spent, date time.Time) *excelize.File {
	f := excelize.NewFile()
	id := uuid.New()
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("Failed to close file", "err", err)
		}
	}()

	index, err := f.NewSheet("gastado")
	if err != nil {
		fmt.Println(err)
		return f
	}

	_ = f.SetColWidth("gastado", "A", "A", 11.50)
	_ = f.SetColWidth("gastado", "B", "B", 40)
	_ = f.SetColWidth("gastado", "C", "C", 13.50)

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
			Horizontal: "left",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})
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

	exp := "yyyy-mm-dd"
	dateStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
		},
		CustomNumFmt: &exp,
	})

	f.SetActiveSheet(index)
	_ = f.DeleteSheet("Sheet1")

	_ = f.SetCellValue("gastado", "A1", "Gastado Por Partida")
	_ = f.MergeCell("gastado", "A1", "C1")
	_ = f.SetCellStyle("gastado", "A1", "C1", pageTitleStyle)

	_ = f.SetCellValue("gastado", "A3", "Fecha de corte")
	_ = f.SetCellFormula("gastado", "B3", fmt.Sprintf("=DATE(%d,%d,%d)", date.Year(), int(date.Month()), date.Day()))
	_ = f.SetCellStyle("gastado", "B3", "B3", dateStyle)
	_ = f.SetCellValue("gastado", "A4", "Proyecto")
	_ = f.SetCellValue("gastado", "B4", project.Name)
	_ = f.SetCellStyle("gastado", "A3", "A4", colTitleStyle)

	_ = f.SetCellValue("gastado", "A6", "Código")
	_ = f.SetCellValue("gastado", "B6", "Partida")
	_ = f.SetCellValue("gastado", "C6", "Total")
	_ = f.SetCellStyle("gastado", "A6", "C6", titleStyle)

	defaultStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: false,
		},
		NumFmt: 4,
	})

	row := 7
	for _, d := range data {
		_ = f.SetCellValue("gastado", fmt.Sprintf("A%d", row), d.BudgetItem.Code)
		_ = f.SetCellValue("gastado", fmt.Sprintf("B%d", row), d.BudgetItem.Name)
		_ = f.SetCellFloat("gastado", fmt.Sprintf("C%d", row), d.Spent, 2, 64)

		_ = f.SetCellStyle("gastado", fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), defaultStyle)
		row++
	}

	if err := f.SaveAs("./public/" + id.String() + ".xlsx"); err != nil {
		slog.Error("Failed to export", "err", err)
	}

	return f
}
