package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go-trading-bot/config"
	"go-trading-bot/internal/api"
	"go-trading-bot/internal/logger"
	"go-trading-bot/internal/service"
)

func main() {
	logger.Log.Infof("Go Trading Bot 🟢")
	logger.Log.Infof("Go Version: %s\n", runtime.Version())
	logger.Log.Infof("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	logger.Log.Infof("Environment configured successfully!")

	c := config.GetConfig()
	logger.Log.Infof("config -> %+v\n", c)

	config.ReadTradingConfig()
	t := config.GetTradingConfig()
	logger.Log.Infof("tradingConfig -> %+v\n", t)

	stopChan := make(chan struct{})
	tradingBot := &service.TradingBot{}
	tradingBot.Initialize()
	go tradingBot.RunTradingBot(stopChan)

	// Gin API 서버 시작 (goroutine으로)
	router := api.NewRouter()
	go func() {
		addr := fmt.Sprintf(":%d", c.Port)
		logger.Log.Infof("Starting Gin API server on %s 🌐", addr)
		if err := router.Run(addr); err != nil {
			logger.Log.Errorf("Failed to start API server: %v", err)
		}
	}()

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Graceful shutdown
	logger.Log.Info("Shutting down Trading Bot 🛑")
	close(stopChan)
}
