package model

type PriceSummary struct {
	Id                string
	Buy               float64
	Sell              float64
	High              float64
	Low               float64
	ChangePercent     float64
	ChangePercentText string
}
