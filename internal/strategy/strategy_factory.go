package strategy

import (
	"go-trading-bot/config"
)

func CreateStrategy(tradingConfig *config.TradingConfig) TradingStrategy {
	switch strategy := tradingConfig.Strategy; strategy {
	case "moving-average-cross":
		return &MovingAverageCrossStrategy{tradingConfig.Strategy}
	default:
		return nil
	}
}
