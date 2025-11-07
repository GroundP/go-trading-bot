package strategy

type TradingStrategy interface {
	GetName()
	Analyze(market string, candles []Candle)
}

type MovingAverageCrossStrategy struct {
	name string
}
