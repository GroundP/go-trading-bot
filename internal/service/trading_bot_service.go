package service

import (
	"go-trading-bot/config"
	"go-trading-bot/internal/client"
	"go-trading-bot/internal/logger"
	"go-trading-bot/internal/strategy"
	"time"
)

type TradingBot struct {
	strategy        strategy.TradingStrategy
	marketHandler   *MarketHandler
	validateMarkets []string
}

func (t *TradingBot) Initialize() {
	t.marketHandler = &MarketHandler{upbitApiClient: &client.UpbitAPIClient{BaseURL: config.GetConfig().UpbitAPIUrl}}
	t.validateMarkets = t.marketHandler.validateAndFilterMarkets()
}

func (t *TradingBot) RunTradingBot(stopChan <-chan struct{}) {
	if len(t.validateMarkets) == 0 {
		logger.Log.Errorf("유효한 마켓이 없습니다. 봇을 시작할 수 없습니다. 🔴")
		return
	}

	go t.runTask()
	tradingConfig := config.GetTradingConfig()
	t.strategy = strategy.CreateStrategy(tradingConfig)

	ticker := time.NewTicker(time.Duration(tradingConfig.AnalysisInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go t.runTask()
		case <-stopChan:
			logger.Log.Infof("스케줄러 종료 요청")
			return
		}
	}
}

func (t *TradingBot) runTask() {
	logger.Log.Info("=========runTask===========")

	requireCandleCount := t.strategy.GetRequiredCandleCount()

	for _, m := range t.validateMarkets {
		candles := t.marketHandler.GetCandles(m, requireCandleCount)
		signals := t.strategy.Analyze(m, candles)
		//handleSignals(signals)
	}

}
