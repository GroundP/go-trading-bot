package client

import (
	"encoding/json"
	"go-trading-bot/config"
	"go-trading-bot/internal/logger"
	"go-trading-bot/internal/model"
	"io"
	"net/http"
)

type UpbitApiClient struct {
}

func (o *UpbitApiClient) GetAllMarkets() ([]model.MarketInfo, error) {
	url := config.GetConfig().UpbitApiUrl + "/market/all"

	resp, err := http.Get(url)
	if err != nil {
		logger.Log.Errorf("Failed to get all markets: %w", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Errorf("Failed to parse all markets: %w", err)
		return nil, err
	}

	var marketInfo []model.MarketInfo
	if err := json.Unmarshal(body, &marketInfo); err != nil {
		logger.Log.Errorf("Failed to convert data(all markets): %w", err)
		return nil, err
	}

	return marketInfo, nil
}
