package strategy

import "go-trading-bot/internal/model"

type TradingStrategy interface {
	GetName() string
	Analyze(market string, candles []model.Candle) model.Signal
}
