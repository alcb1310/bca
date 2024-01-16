package types

type BalanceResponse struct {
	Invoices []InvoiceResponse
	Total    float64
}
