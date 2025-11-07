package service

import (
	"fmt"
	"go-trading-bot/config"
	"time"
)

func RunTradingBot() {
	tradingConfig := config.GetTradingConfig()
	ticker := time.NewTicker(time.Duration(tradingConfig.AnalysisInterval) * time.Second)

	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:

		}
	}

}

func runTask() {

}
