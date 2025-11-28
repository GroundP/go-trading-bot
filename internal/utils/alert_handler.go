package utils

import (
	"fmt"
	"go-trading-bot/config"
	"go-trading-bot/internal/logger"
	"go-trading-bot/internal/model"
	"net/http"
	"net/url"

	LANG "golang.org/x/text/language"
	MSG "golang.org/x/text/message"
)

// SendTelegramAlert sends a trading signal alert to Telegram
func SendTelegramAlert(signal model.Signal) {
	cfg := config.GetConfig()
	if cfg.TelegramSend != "OK" {
		logger.Log.Debug("Telegram send is not enabled. Skipping alert.")
		return
	}

	token := cfg.TelegramBotToken
	chatID := cfg.TelegramChatID

	if token == "" || chatID == "" {
		logger.Log.Debug("Telegram configuration is missing. Skipping alert.")
		return
	}

	message := formatSignalMessage(signal, "")
	sendMessage(token, chatID, message)
}

func SendTelegramMultiAlert(actions []model.Action) {
	cfg := config.GetConfig()
	if cfg.TelegramSend != "OK" {
		logger.Log.Debug("Telegram send is not enabled. Skipping alert.")
		return
	}

	token := cfg.TelegramBotToken
	chatID := cfg.TelegramChatID

	if token == "" || chatID == "" {
		logger.Log.Debug("Telegram configuration is missing. Skipping alert.")
		return
	}

	var totalMessage string
	for _, action := range actions {
		message := formatActionMessage(action)
		totalMessage += message + "\n-----------------------------------------------------\n\n"
	}

	logger.Log.Debugf("Telegram alert message: %s", totalMessage)
	sendMessage(token, chatID, totalMessage)
}

func SendTelegramMessage(message string) {
	cfg := config.GetConfig()
	if cfg.TelegramSend != "OK" {
		logger.Log.Debug("Telegram send is not enabled. Skipping alert.")
		return
	}

	token := cfg.TelegramBotToken
	chatID := cfg.TelegramChatID

	if token == "" || chatID == "" {
		logger.Log.Debug("Telegram configuration is missing. Skipping alert.")
		return
	}
	sendMessage(token, chatID, message)
}

// formatSignalMessage formats the trading signal into a readable message
func formatSignalMessage(signal model.Signal, usdtPrice string) string {
	var emoji string
	var action string

	switch signal.Type {
	case model.BUY:
		emoji = "🟢"
		action = "매수 신호"
	case model.SELL:
		emoji = "🔴"
		action = "매도 신호"
	case model.HOLD:
		emoji = "⚪"
		action = "홀드 신호"
	default:
		emoji = "⚫"
		action = "알 수 없는 신호"
	}

	// 메시지 헤더
	message := fmt.Sprintf("<b>%s [%s] %s</b>\n\n", emoji, signal.Market, action)

	// 현재가 정보
	p := MSG.NewPrinter(LANG.Korean)
	message += p.Sprintf("💰 <b>현재가:</b> %.0f (%s)\n", signal.CurrentPrice, usdtPrice)

	// Stage 정보 (사이클 전략인 경우)
	if signal.Stage != nil {

		message += fmt.Sprintf("📊 사이클 단계: <b>%s</b>\n", signal.Stage.StageNumber)
		// 단계 방향 정보
		var dirIcon string
		var dirText string
		switch signal.Stage.StageDir {
		case model.STAGE_DIR_NORMAL:
			dirIcon = "➡️"
			dirText = "정상 진행"
		case model.STAGE_DIR_REVERSE:
			dirIcon = "🔙"
			dirText = "역방향 전환"
		case model.STAGE_DIR_MAINTAIN:
			dirIcon = "⏸️"
			dirText = "단계 유지"
		}

		message += fmt.Sprintf("✔ <i>%s, %s %s</i>\n", signal.Stage.Description, dirIcon, dirText)
		message += "\n"
	}

	// 설명
	if signal.Description != "" {
		message += fmt.Sprintf("📝 <b>상세:</b>\n%s\n\n", signal.Description)
	}

	// 전략 및 시각
	message += fmt.Sprintf("🎯 <b>전략:</b> %s\n", signal.StrategyName)
	message += fmt.Sprintf("🕐 <b>시각:</b> %s", signal.Timestamp)

	return message
}

func formatActionMessage(action model.Action) string {
	message := formatSignalMessage(action.Signal, action.USDTPrice)
	message += "\n\n"

	// INSERT_YOUR_CODE
	// 포지션 정보 출력
	message += "\n"
	message += "<b>📦 포지션 정보</b>\n"
	if action.Position.Status == model.POSITION_BUY {
		profit := (action.Signal.CurrentPrice - action.Position.EntryPrice) * action.Position.Quantity
		message += fmt.Sprintf(
			"상태: <b>보유중</b>\n수량: <b>%f</b>\n진입가: <b>%f</b>\n수익: <b>%f</b>\n",
			action.Position.Quantity,
			action.Position.EntryPrice,
			profit,
		)
	} else {
		message += "상태: <b>없음</b>\n"
	}
	return message
}

// sendMessage sends a message to Telegram using the Bot API
func sendMessage(token, chatID, message string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", message)
	data.Set("parse_mode", "HTML")

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		logger.Log.Errorf("텔레그램 메시지 전송 실패: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Log.Errorf("텔레그램 API 오류 (상태 코드: %d)", resp.StatusCode)
	} else {
		logger.Log.Debug("텔레그램 알림 전송 완료")
	}
}
