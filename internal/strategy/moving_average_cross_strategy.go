package strategy

import "go-trading-bot/internal/model"

type MovingAverageCrossStrategy struct {
	name string
}

func (o *MovingAverageCrossStrategy) GetName() string {
	return o.name
}

func (o *MovingAverageCrossStrategy) Analyze(market string, candles []model.Candle) model.Signal {
	return model.Signal{}
}
