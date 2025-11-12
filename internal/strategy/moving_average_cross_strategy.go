package strategy

import (
	"fmt"
	"go-trading-bot/config"
	"go-trading-bot/internal/logger"
	"go-trading-bot/internal/model"
	"time"
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

	logger.Log.Infof("[%v] 이전 MA%v: %.2f, MA%v: %.2f", market, shortPeriod, previousShortMA, longPeriod, previousLongMA)
	logger.Log.Infof("[%v] 현재 MA%v: %.2f, MA%v: %.2f", market, shortPeriod, currentShortMA, longPeriod, currentLongMA)

	currentCandle := candles[0]
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	var signal model.Signal
	if previousShortMA < previousLongMA && currentShortMA > currentLongMA {
		description := fmt.Sprintf("▲ 골든 크로스 발생 -> MA%d(%.2f)이 MA%d(%.2f)를 상향 돌파", shortPeriod, currentShortMA, longPeriod, currentLongMA)
		signal = model.Signal{Type: model.BUY, Market: market, CurrentPrice: currentCandle.TradePrice, Timestamp: currentTime, Description: description, StrategyName: m.GetName()}
	}

	if previousShortMA > previousLongMA && currentShortMA < currentLongMA {
		description := fmt.Sprintf("▼ 데드 크로스 발생 -> MA%d(%.2f)이 MA%d(%.2f)를 하향 돌파", shortPeriod, currentShortMA, longPeriod, currentLongMA)
		signal = model.Signal{Type: model.BUY, Market: market, CurrentPrice: currentCandle.TradePrice, Timestamp: currentTime, Description: description, StrategyName: m.GetName()}
	}

	description := "이동평균선 교차 없음 - 관망"
	signal = model.Signal{Type: model.HOLD, Market: market, CurrentPrice: currentCandle.TradePrice, Timestamp: currentTime, Description: description, StrategyName: m.GetName()}
	m.printSignal(&signal)
	return signal
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

func (m *MovingAverageCrossStrategy) printSignal(signal *model.Signal) {
	logger.Log.Infof("마켓: %v", signal.Market)
	logger.Log.Infof("신호: %v", signal.Type)
	logger.Log.Infof("현재가: %.2f", signal.CurrentPrice)
	logger.Log.Infof("전략: %v", signal.StrategyName)
	logger.Log.Infof("설명: %v", signal.Description)
	logger.Log.Infof("시각: %v", signal.Timestamp)
	logger.Log.Info("캔들 분석 완료 🟢")
}
