// Package model
package model

type Candle struct {
	market       string  `json:"market"`
	openingPrice float64 `json:"opening_price"`
	highPrice    float64 `json:"high_price"`
	lowPrice     float64 `json:"low_price"`
	tradePrice   float64 `json:"trade_price"`
	timestamp    int64   `json:"timestamp"`
}
