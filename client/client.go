package client

import (
	"workhorse/fin-provider/model"

	"github.com/amecky/fin-math/math"
)

type PriceProvider interface {
	LoadSummary(tickers []string) ([]model.PriceSummary, error)
	LoadCandles(id string) (*math.Matrix, error)
}
