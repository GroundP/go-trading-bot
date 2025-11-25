package api

import (
	"go-trading-bot/internal/handler"
	"go-trading-bot/internal/service"

	"github.com/gin-gonic/gin"
)

// NewRouter는 TradingBot 인스턴스를 받아 라우터를 생성합니다
func NewRouter(tradingBot *service.TradingBot) *gin.Engine {
	router := gin.Default()

	tradingBotHandler := handler.NewHandler(tradingBot)

	v1Group := router.Group("/api/v1")
	{
		// 헬스체크
		v1Group.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "OK"})
		})

		// 트레이딩 신호 조회
		// GET /api/v1/signal?market=KRW-BTC (특정 마켓)
		// GET /api/v1/signal (모든 마켓)
		v1Group.GET("/signal", tradingBotHandler.GetSignal)
	}
	return router
}
