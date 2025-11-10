package main

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go-trading-bot/config"
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

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// alert on stop
	logger.Log.Info("stopped Trading Bot 🛑")
}
