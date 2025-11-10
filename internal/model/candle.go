// Package model
package model

type Candle struct {
	Market       string  `json:"market"`
	OpeningPrice float64 `json:"opening_price"`
	HighPrice    float64 `json:"high_price"`
	LowPrice     float64 `json:"low_price"`
	TradePrice   float64 `json:"trade_price"`
	Timestamp    int64   `json:"timestamp"`
}
