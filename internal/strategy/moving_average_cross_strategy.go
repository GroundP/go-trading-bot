package strategy

import (
	"go-trading-bot/config"
	"go-trading-bot/internal/logger"
	"go-trading-bot/internal/model"
)

type MovingAverageCrossStrategy struct {
	name               string
	movingAverageCross config.MovingAverageCross
}

func (m *MovingAverageCrossStrategy) GetName() string {
	return m.name
}

func (m *MovingAverageCrossStrategy) Analyze(market string, candles []model.Candle) model.Signal {
	if len(candles) == 0 {
		logger.Log.Error("캔들이 비어있습니다. 분석할 수 없습니다. 🔴")
	}

	logger.Log.Info("캔들 분석을 시작합니다. 🔘")
	shortPeriod := m.movingAverageCross.ShortPeriod
	longPeriod := m.movingAverageCross.LongPeriod

	currentShortMA := m.calculateMA(candles, shortPeriod, 0)
	currentLongMA := m.calculateMA(candles, longPeriod, 0)

	previousShortMA := m.calculateMA(candles, shortPeriod, 1)
	previousLongMA := m.calculateMA(candles, longPeriod, 1)

	logger.Log.Infof("[%v] 현재 MA%v: %v, MA%v: %v", market, shortPeriod, currentShortMA, longPeriod, currentLongMA)
	logger.Log.Infof("[%v] 이전 MA%v: %v, MA%v: %v", market, shortPeriod, previousShortMA, longPeriod, previousLongMA)

	//currentCandle := candles[0]

	return model.Signal{}
}

func (m *MovingAverageCrossStrategy) GetRequiredCandleCount() int {
	return m.movingAverageCross.LongPeriod + 1
}

func (m *MovingAverageCrossStrategy) calculateMA(candles []model.Candle, period int, offset int) float64 {
	var sum float64
	for i := offset; i < offset+period; i++ {
		sum += candles[i].TradePrice
	}

	return sum / float64(period)
}
