package strategy

import (
	"go-trading-bot/config"
	"go-trading-bot/internal/model"
)

func CreateStrategy(tradingConfig *config.TradingConfig) TradingStrategy {
	switch strategy := tradingConfig.Strategy; strategy {
	case "moving-average-cross":
		return &MovingAverageCrossStrategy{tradingConfig.Strategy, tradingConfig.MovingAverageCross}
	case "moving-average-cycle":
		return &MovingAverageCycleStrategy{tradingConfig.Strategy, tradingConfig.MovingAverageCycle, model.Stage{}}
	default:
		return nil
	}
}
