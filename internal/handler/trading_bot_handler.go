// Package handler
package handler

import (
	"fmt"
	"go-trading-bot/internal/service"

	"github.com/gin-gonic/gin"
)

// Handler는 API 핸들러에 필요한 의존성을 관리합니다
type TradingBotHandler struct {
	TradingBot *service.TradingBot
}

// NewHandler는 새로운 Handler 인스턴스를 생성합니다
func NewHandler(tradingBot *service.TradingBot) *TradingBotHandler {
	return &TradingBotHandler{
		TradingBot: tradingBot,
	}
}

// GetSignal은 특정 마켓의 최신 트레이딩 신호를 반환합니다
func (h *TradingBotHandler) GetSignal(c *gin.Context) {
	market := c.Query("market")

	// market 파라미터가 없으면 모든 신호 반환
	if market == "" {
		signals := h.TradingBot.GetAllLatestSignals()
		c.JSON(200, gin.H{
			"success": true,
			"data":    signals,
			"count":   len(signals),
		})
		return
	}

	// 특정 마켓의 신호 반환
	signal := h.TradingBot.GetLatestSignal(market)

	if signal.Market == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": fmt.Sprintf("Market %s not found", market),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    signal,
	})
}
