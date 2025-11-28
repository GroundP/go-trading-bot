package service

import (
	"fmt"
	"go-trading-bot/config"
	"go-trading-bot/internal/client"
	"go-trading-bot/internal/logger"
	"go-trading-bot/internal/model"
	"go-trading-bot/internal/strategy"
	"go-trading-bot/internal/utils"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type TradingBot struct {
	strategy        strategy.TradingStrategy
	marketHandler   *MarketHandler
	validateMarkets []string
	latestSignal    map[string]model.Signal
	orderService    *OrderService
}

func (t *TradingBot) Initialize() {
	t.marketHandler = &MarketHandler{upbitAPIClient: &client.UpbitAPIClient{BaseURL: config.GetConfig().UpbitAPIUrl}, binanceAPIClient: &client.BinanceAPIClient{}}
	t.validateMarkets = t.marketHandler.validateAndFilterMarkets()
	t.latestSignal = make(map[string]model.Signal)
	t.orderService = &OrderService{positions: make(map[string]model.Position)}
}

func (t *TradingBot) RunTradingBot(stopChan <-chan struct{}) {
	if len(t.validateMarkets) == 0 {
		logger.Log.Errorf("유효한 마켓이 없습니다. 봇을 시작할 수 없습니다. 🔴")
		return
	}

	go t.runTask()
	tradingConfig := config.GetTradingConfig()
	t.strategy = strategy.CreateStrategy(tradingConfig)

	ticker := time.NewTicker(time.Duration(tradingConfig.AnalysisInterval) * time.Minute)
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

func (t *TradingBot) GetLatestSignal(market string) model.Signal {
	if signal, exists := t.latestSignal[market]; exists {
		return signal
	}
	return model.Signal{}
}

func (t *TradingBot) GetAllLatestSignals() []model.Signal {
	signals := make([]model.Signal, 0, len(t.latestSignal))
	for _, signal := range t.latestSignal {
		signals = append(signals, signal)
	}
	return signals
}

func (t *TradingBot) runTask() {
	logger.Log.Info("=========runTask===========")

	requireCandleCount := t.strategy.GetRequiredCandleCount()

	for _, m := range t.validateMarkets {
		candles := t.marketHandler.GetCandles(m, requireCandleCount)
		signal := t.strategy.Analyze(m, candles)
		t.handleSignal(signal)
	}

	signals := t.GetAllLatestSignals()
	positions := t.marketHandler.GetPositions()
	actions := t.createActions(signals, positions)
	utils.SendTelegramMultiAlert(actions)
}

func (t *TradingBot) handleSignal(signal model.Signal) {
	t.latestSignal[signal.Market] = signal
	//t.printSignal(&signal)
	logger.Log.Infof("SIGNAL INFO:\n%v", t.createSignalInfo(&signal))
	//utils.SendTelegramAlert(signal)

	switch signal.Type {
	case model.BUY:
		logger.Log.Infof("[%v] 매수 신호 -> BUY 주문을 실행합니다.", signal.Market)
		t.orderService.PlaceOrder(signal.Market, model.BUY, signal.CurrentPrice)
	case model.SELL:
		logger.Log.Infof("[%v] 매도 신호 -> SELL 주문을 실행합니다.", signal.Market)
		t.orderService.PlaceOrder(signal.Market, model.SELL, signal.CurrentPrice)
	case model.HOLD:
		logger.Log.Infof("[%v] HOLD 신호 -> 매매 없음, 포지션 상태: %v", signal.Market, "")
	}
}

// shouldSendHoldAlert는 HOLD 신호에서도 알림을 보낼지 결정합니다
func (t *TradingBot) shouldSendHoldAlert(signal *model.Signal) bool {
	// Stage 정보가 있고, 단계가 변경된 경우에만 알림 전송
	if signal.Stage != nil {
		// 정상 진행이나 역방향 전환 시 알림 전송
		if signal.Stage.StageDir == model.STAGE_DIR_NORMAL || signal.Stage.StageDir == model.STAGE_DIR_REVERSE {
			return true
		}
	}
	return false
}

func (t *TradingBot) printSignal(signal *model.Signal) {
	logger.Log.Infof("마켓: %v", signal.Market)
	logger.Log.Infof("신호: %v", signal.Type)
	logger.Log.Infof("현재가: %.2f", signal.CurrentPrice)
	logger.Log.Infof("전략: %v", signal.StrategyName)

	// Stage 정보가 있으면 출력
	if signal.Stage != nil {
		logger.Log.Infof("사이클 단계: %v (%v) (%v)", signal.Stage.StageNumber, signal.Stage.StageDir, signal.Stage.Description)
	}

	logger.Log.Infof("설명: %v", signal.Description)
	logger.Log.Infof("시각: %v", signal.Timestamp)
	logger.Log.Info("캔들 분석 완료 🟢")
}

func (t *TradingBot) createSignalInfo(signal *model.Signal) string {
	p := message.NewPrinter(language.Korean)

	info := fmt.Sprintf("✔ 마켓: %v\n", signal.Market)
	info += fmt.Sprintf("✔ 신호: %v\n", signal.Type)
	info += p.Sprintf("✔ 현재가: %.0f\n", signal.CurrentPrice)
	info += fmt.Sprintf("✔ 전략: %v\n", signal.StrategyName)

	if signal.Stage != nil {
		info += fmt.Sprintf("✔ Stage: %v (%v) (%v)\n", signal.Stage.StageNumber, signal.Stage.StageDir, signal.Stage.Description)
	}

	info += fmt.Sprintf("✔ 설명: %v\n", signal.Description)
	info += fmt.Sprintf("✔ 시각: %v\n", signal.Timestamp)
	return info
}

func (t *TradingBot) createActions(signals []model.Signal, positions model.Positions) []model.Action {
	actions := make([]model.Action, 0, len(signals))
	for _, signal := range signals {
		asset := strings.Split(signal.Market, "-")[1]
		var position model.Position
		for _, p := range positions {
			if asset == p.Market {
				position = p
				break
			}
		}

		var usdtPrice string
		binancePrices := t.marketHandler.GetBinancePrices()
		for _, price := range binancePrices {
			if price.Asset == asset {
				usdtPrice = price.Price
				break
			}
		}

		action := model.Action{
			Market:   signal.Market,
			Signal:   signal,
			Position: position,
			USDTPrice: usdtPrice,
		}
		actions = append(actions, action)
	}

	for _, action := range actions {
		logger.Log.Infof("ACTION: %#v", action)
	}

	return actions
}
