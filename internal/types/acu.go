package types

import "github.com/google/uuid"

type Quantity struct {
	Id        uuid.UUID
	Project   Project
	Rubro     Rubro
	Quantity  float64
	CompanyId uuid.UUID
}

type AnalysisReport struct {
	ProjectName  string
	CategoryName string
	MaterialName string
	Quantity     float64
}
