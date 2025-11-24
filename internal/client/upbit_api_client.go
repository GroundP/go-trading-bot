package client

import (
	"encoding/json"
	"errors"
	"go-trading-bot/internal/logger"
	"go-trading-bot/internal/model"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type UpbitAPIClient struct {
	BaseURL string
}

func (u *UpbitAPIClient) GetAllMarkets() ([]model.MarketInfo, error) {
	url := u.BaseURL + "/market/all"

	resp, err := http.Get(url)
	if err != nil {
		logger.Log.Errorf("Failed to get all markets: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Errorf("Failed to parse all markets: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logger.Log.Errorf("Failed to fetch candles -> statusCode: %v, msg: %v", resp.StatusCode, string(body))
		return nil, errors.New("failed to fetch /market/all")
	}

	var marketInfo []model.MarketInfo
	if err := json.Unmarshal(body, &marketInfo); err != nil {
		logger.Log.Errorf("Failed to convert data(all markets): %v", err)
		return nil, err
	}

	return marketInfo, nil
}

func (u *UpbitAPIClient) FetchCandles(market string, path string, requireCandleCount int) ([]model.Candle, error) {
	baseURL := u.BaseURL + path

	params := url.Values{}
	params.Add("market", market)
	params.Add("count", strconv.Itoa(requireCandleCount))

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		logger.Log.Errorf("Failed to create request: %v", err)
		return nil, err
	}

	req.URL.RawQuery = params.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log.Errorf("Failed to fetch candles: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Errorf("Failed to read response body: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logger.Log.Errorf("Failed to fetch candles -> url: %v, statusCode: %v, msg: %v", baseURL, resp.StatusCode, string(body))
		return nil, errors.New("failed to fetch candles")
	}

	var candles []model.Candle
	if err := json.Unmarshal(body, &candles); err != nil {
		logger.Log.Errorf("Failed to convert data(candles): %v", err)
		return nil, err
	}

	return candles, nil
}
