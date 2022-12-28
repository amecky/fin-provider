package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
	"workhorse/fin-provider/model"

	"github.com/amecky/fin-math/math"
)

type Plus500Provider struct {
}

func NewPlus500Provider() *Plus500Provider {
	return &Plus500Provider{}
}

func (p *Plus500Provider) LoadCandles(id string) (*math.Matrix, error) {
	html, err := LoadHtmlContent(fmt.Sprintf("https://www.plus500.com/de/api/LiveData/ChartUpdate?InstrumentId=%s&ResolutionLevel=2&SeriesType=2&Precision=2", id))
	if err != nil {
		return nil, err
	}
	lines := strings.Split(html, "],")
	m := math.NewMatrix(6)
	for _, l := range lines {
		entries := strings.Split(l, ",")
		date := entries[0]
		date = date[1:]
		// FIXME: very first entry is "[[
		t, _ := strconv.Atoi(date)
		tm := time.Unix(int64(t/1000), 0)
		row := m.AddRow(tm.Format("2006-01-02 15:04"))
		for i := 0; i < 4; i++ {
			cur := entries[i+1]
			if strings.Index(cur, "]") != -1 {
				cur = cur[0:strings.Index(cur, "]")]
			}
			v, _ := strconv.ParseFloat(cur, 32)
			row.Set(i, v)
		}
	}
	return m, nil
}

func buildPlus500URL(tickers []string) string {
	ids := strings.Join(tickers, ",")
	u := "https://www.plus500.com/de/api/LiveData/FeedsUpdates"
	v := url.Values{
		"instrumentIds": {ids},
	}
	return fmt.Sprintf("%s?%s", u, v.Encode())
}

func (p *Plus500Provider) LoadSummary(tickers []string) ([]model.PriceSummary, error) {
	var ret = make([]model.PriceSummary, 0)
	html, err := LoadHtmlContent(buildPlus500URL(tickers))
	if err != nil {
		return nil, err
	}
	type InstrumentResponse []struct {
		InstrumentId      string  `json:"InstrumentId"`
		BuyPrice          float64 `json:"BuyPrice"`
		SellPrice         float64 `json:"SellPrice"`
		High              float64 `json:"HighPrice"`
		Low               float64 `json:"LowPrice"`
		ChangePercent     float64 `json:"ChangePercent"`
		ChangePercentText string  `json:"ChangePercentText"`
	}
	var resp InstrumentResponse
	err = json.Unmarshal([]byte(html), &resp)
	if err != nil {
		return nil, err
	}
	for _, i := range resp {
		ret = append(ret, model.PriceSummary{
			Id:                i.InstrumentId,
			Buy:               i.BuyPrice,
			Sell:              i.SellPrice,
			High:              i.High,
			Low:               i.Low,
			ChangePercent:     i.ChangePercent,
			ChangePercentText: i.ChangePercentText,
		})
	}
	return ret, nil
}
