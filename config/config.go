package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
)

var (
	instance      *Config
	once          sync.Once
	tradingConfig *TradingConfig
)

type Config struct {
	Env     string
	AppName string
	Port    int

	DBHost string
	DBPort int
	DBUser string
	DBPass string
	DBName string

	AccessKey string
	SecretKey string

	UpbitApiUrl string
}

type TradingConfig struct {
	Markets            []string           `json:"markets"`
	Strategy           string             `json:"strategy"`
	Candle             Candle             `json:"candle"`
	MovingAverageCross MovingAverageCross `json:"moving-average-cross"`
	AnalysisInterval   int                `json:"analysis-interval"`
	OrderAmount        float64            `json:"order-amount"`
}

type Candle struct {
	Category string `json:"category"`
	Unit     int    `json:"unit"`
}

type MovingAverageCross struct {
	ShortPeriod int `json:"short-period"`
	LongPeriod  int `json:"long-period"`
}

func GetConfig() *Config {
	// singleton
	once.Do(func() {
		instance = newConfig()
	})
	return instance
}

func newConfig() *Config {
	return &Config{
		Env:     getEnvStr("ENV", "development"),
		AppName: getEnvStr("APP_NAME", "go-trading-bot"),
		Port:    getEnvInt("PORT", 3000),

		DBHost: getEnvStr("DB_HOST", ""),
		DBPort: getEnvInt("DB_PORT", 6379),
		DBUser: getEnvStr("DB_USER", ""),
		DBPass: getEnvStr("DB_PASS", ""),
		DBName: getEnvStr("DB_NAME", ""),

		AccessKey: getEnvStr("ACCESS_KEY", ""),
		SecretKey: getEnvStr("SECRET_KEY", ""),

		UpbitApiUrl: getEnvStr("UPBIT_API_URL", "https://api.upbit.com/v1"),
	}
}

func getEnvStr(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}
		return parsed
	}
	return defaultValue
}

func ReadTradingConfig() bool {
	file, err := os.ReadFile("application.json")
	if err != nil {
		fmt.Println("Failed to open application.json")
		return false
	}

	err = json.Unmarshal(file, &tradingConfig)
	if err != nil {
		fmt.Println("Failed to parse application.json")
		return false
	}

	return true
}

func GetTradingConfig() *TradingConfig {
	return tradingConfig
}
