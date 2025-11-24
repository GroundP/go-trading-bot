package strategy

import (
	"fmt"
	"go-trading-bot/config"
	"go-trading-bot/internal/logger"
	"go-trading-bot/internal/model"
	"time"
)

type MovingAverageCycleStrategy struct {
	name               string
	movingAverageCycle config.MovingAverageCycle
	stage              model.Stage
}

func (m *MovingAverageCycleStrategy) GetName() string {
	return m.name
}

func (m *MovingAverageCycleStrategy) Analyze(market string, candles []model.Candle) model.Signal {
	if len(candles) == 0 {
		logger.Log.Error("캔들이 비어있습니다. 분석할 수 없습니다. 🔴")
	}

	logger.Log.Info("캔들 분석을 시작합니다. 🔘")
	periods := [3]int{m.movingAverageCycle.ShortPeriod, m.movingAverageCycle.MediumPeriod, m.movingAverageCycle.LongPeriod}
	maCurrent := [3]float64{}
	maPrevious := [3]float64{}

	for i, period := range periods {
		maCurrent[i] = m.calculateMA(candles, period, 0)
		maPrevious[i] = m.calculateMA(candles, period, 1)
	}

	logger.Log.Infof("[%v] 이전 MA%v: %.2f, MA%v: %.2f, MA%v: %.2f", market, periods[0], maPrevious[0], periods[1], maPrevious[1], periods[2], maPrevious[2])
	logger.Log.Infof("[%v] 현재 MA%v: %.2f, MA%v: %.2f, MA%v: %.2f", market, periods[0], maCurrent[0], periods[1], maCurrent[1], periods[2], maCurrent[2])

	currentCandle := candles[0]
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Stage를 분석하고 Signal 생성
	signal := m.calculateSignal(market, currentCandle.TradePrice, currentTime, periods, maCurrent, maPrevious)
	return signal
}

func (m *MovingAverageCycleStrategy) GetRequiredCandleCount() int {
	return m.movingAverageCycle.LongPeriod + 1
}

func (m *MovingAverageCycleStrategy) calculateMA(candles []model.Candle, period int, offset int) float64 {
	var sum float64
	for i := offset; i < offset+period; i++ {
		sum += candles[i].TradePrice
	}

	return sum / float64(period)
}

func (m *MovingAverageCycleStrategy) calculateSignal(market string, currentPrice float64, currentTime string, periods [3]int, maCurrent [3]float64, maPrevious [3]float64) model.Signal {
	currentShortMA := maCurrent[0]
	currentMediumMA := maCurrent[1]
	currentLongMA := maCurrent[2]
	previousShortMA := maPrevious[0]
	previousMediumMA := maPrevious[1]
	previousLongMA := maPrevious[2]

	// Stage 초기화
	stageNumber := model.STAGE_0
	stageDescription := "알 수 없는 단계"
	stageDir := model.STAGE_DIR_NONE
	signalType := model.HOLD

	// Stage 분석 및 신호 결정
	if currentShortMA > currentMediumMA && currentShortMA > currentLongMA {
		if currentMediumMA > currentLongMA {
			// STAGE_1: 안정 상승기, 단/중/장 배치
			stageNumber = model.STAGE_1
			stageDescription = "안정 상승기, 단/중/장 배치"
			if currentShortMA > previousShortMA && currentMediumMA > previousMediumMA && currentLongMA > previousLongMA {
				signalType = model.BUY // 모두 우상향 중인 경우 매수
			}
		} else {
			// STAGE_6: 상승 추세 시작, 단/장/중 배치
			stageNumber = model.STAGE_6
			stageDescription = "상승 추세 시작, 단/장/중 배치"
		}
	} else if currentMediumMA > currentLongMA && currentMediumMA > currentShortMA {
		if currentShortMA > currentLongMA {
			// STAGE_2: 상승 추세 끝, 중/단/장 배치
			stageNumber = model.STAGE_2
			stageDescription = "상승 추세 끝, 중/단/장 배치"
		} else {
			// STAGE_3: 하락 추세 시작, 중/장/단 배치
			stageNumber = model.STAGE_3
			stageDescription = "하락 추세 시작, 중/장/단 배치"
		}
	} else if currentLongMA > currentMediumMA && currentLongMA > currentShortMA {
		if currentMediumMA > currentShortMA {
			// STAGE_4: 안정 하락기, 장/중/단 배치
			stageNumber = model.STAGE_4
			stageDescription = "안정 하락기, 장/중/단 배치"
			if currentShortMA < previousShortMA && currentMediumMA < previousMediumMA && currentLongMA < previousLongMA {
				signalType = model.SELL // 모두 우하향 중인 경우 매도
			}
		} else {
			// STAGE_5: 하락 추세 끝, 장/단/중 배치
			stageNumber = model.STAGE_5
			stageDescription = "하락 추세 끝, 장/단/중 배치"
		}
	}

	if m.stage.StageNumber != model.STAGE_0 {
		if m.stage.StageNumber == stageNumber {
			stageDir = model.STAGE_DIR_MAINTAIN
		} else if m.stage.StageNumber > stageNumber {
			stageDir = model.STAGE_DIR_REVERSE
		} else if m.stage.StageNumber < stageNumber {
			stageDir = model.STAGE_DIR_NORMAL
		}
	}

	m.stage = model.Stage{
		StageNumber: stageNumber,
		StageDir:    stageDir,
		Description: stageDescription,
	}

	// Stage 정보를 포함한 상세 Description 생성
	description := m.generateDescription(signalType, periods, maCurrent)

	// Stage 복사 (Signal에 포함시키기 위해)
	stageCopy := m.stage

	// Signal 생성
	signal := model.Signal{
		Type:         signalType,
		Market:       market,
		CurrentPrice: currentPrice,
		Timestamp:    currentTime,
		Description:  description,
		StrategyName: m.GetName(),
		Stage:        &stageCopy,
	}

	return signal
}

// generateDescription은 신호 타입과 MA 정보를 기반으로 설명을 생성합니다
func (m *MovingAverageCycleStrategy) generateDescription(signalType model.SignalType, periods [3]int, maCurrent [3]float64) string {
	stageDesc := m.stage.Description
	maInfo := ""

	switch signalType {
	case model.BUY:
		maInfo = "📈 매수 신호 - "
	case model.SELL:
		maInfo = "📉 매도 신호 - "
	case model.HOLD:
		maInfo = "⏸️ 관망 - "
	}

	maInfo += stageDesc
	maInfo += " | "
	maInfo += fmt.Sprintf("MA%d: %.2f, MA%d: %.2f, MA%d: %.2f",
		periods[0], maCurrent[0],
		periods[1], maCurrent[1],
		periods[2], maCurrent[2])

	return maInfo
}
