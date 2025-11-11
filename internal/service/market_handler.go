package service

import (
	"go-trading-bot/config"
	"go-trading-bot/internal/client"
	"go-trading-bot/internal/logger"
	"go-trading-bot/internal/model"
)

type MarketHandler struct {
	upbitAPIClient *client.UpbitAPIClient
}

func (m *MarketHandler) validateAndFilterMarkets() (validMarkets []string) {
	tradingConfig := config.GetTradingConfig()

	if len(tradingConfig.Markets) == 0 {
		logger.Log.Error("설정된 마켓이 없습니다.")
		return []string{}
	}

	logger.Log.Info("마켓 검증 시작 🔘")
	logger.Log.Infof("설정된 마켓 수: %v", len(tradingConfig.Markets))
	logger.Log.Infof("설정된 마켓: %+v", tradingConfig.Markets)

	var userTargets []string
	for _, m := range tradingConfig.Markets {
		userTargets = append(userTargets, "KRW-"+m)
	}

	marketInfo, err := m.upbitAPIClient.GetAllMarkets()
	if err != nil {
		logger.Log.Errorf("Upbit 마켓 목록 조회 실패. 설정된 마켓을 그대로 사용합니다. %s 🔴", err.Error())
		return userTargets
	}

	if len(marketInfo) == 0 {
		logger.Log.Error("Upbit 마켓 목록이 비어있습니다. 설정된 마켓을 그대로 사용합니다. 🔴")
		return userTargets
	}

	logger.Log.Infof("업비트 지원 마켓 수: %v", len(marketInfo))

	for _, u := range userTargets {
		find := false
		for _, m := range marketInfo {
			if u == m.Market {
				logger.Log.Infof("[유효] %v - %v (%v)", m.Market, m.KoreanName, m.EnglishName)
				validMarkets = append(validMarkets, m.Market)
				find = true
			}
		}

		if !find {
			logger.Log.Warnf("[무효] 업비트에서 지원하지 않는 마켓입니다. 제외됩니다(%v) 🟠", u)
		}
	}

	logger.Log.Infof("유효한 마켓 수 : %v / %v", len(validMarkets), len(userTargets))
	logger.Log.Info("마켓 검증 완료 🟢")

	return validMarkets
}

func (m *MarketHandler) GetCandles(market string, requireCandleCount int) (candles []model.Candle) {
	candleConfig := config.GetTradingConfig().Candle
	path := candleConfig.BuildAPIPath()
	if len(path) == 0 {
		logger.Log.Errorf("Candle Path를 만드는데 실패했습니다. %+v", candleConfig)
		return candles
	}

	candles, err := m.upbitAPIClient.FetchCandles(market, path, requireCandleCount)
	if err != nil {
		logger.Log.Errorf("Failed to fetch Candles -> %s", err.Error())
		return candles
	}

	return candles
}
