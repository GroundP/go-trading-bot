package client

import (
	"encoding/json"
	"go-trading-bot/internal/model"
	"net/http"
)

type BinanceAPIClient struct {
}

func (b *BinanceAPIClient) GetPrices() ([]model.PriceTicker, error) {
	resp, err := http.Get("https://api.binance.com/api/v3/ticker/price")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    var tickers []model.PriceTicker
    if err := json.NewDecoder(resp.Body).Decode(&tickers); err != nil {
        panic(err)
    }

	return tickers, nil
}