package main

import (
	"fmt"
	"runtime"

	"go-trading-bot/config"
)

func main() {
	fmt.Println("=== Go Trading Bot ===")
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println("Environment configured successfully!")

	c := config.GetConfig()
	fmt.Printf("config -> %v\n", c)

	config.ReadTradingConfig()
	t := config.GetTradingConfig()
	fmt.Printf("tradingConfig -> %+v\n", t)

}
