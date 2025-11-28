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

	"crypto/sha512"
	"encoding/hex"
	"sort"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (u *UpbitAPIClient) FetchBalance(accessKey, secretKey string) ([]model.Position, error) {
	baseURL := u.BaseURL + "/accounts"

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		logger.Log.Errorf("Failed to create request: %v", err)
		return nil, err
	}

	token, err := createJwt(accessKey, secretKey, nil)
	if err != nil {
		logger.Log.Errorf("Failed to create JWT: %v", err)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log.Errorf("Failed to fetch balance: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Errorf("Failed to read response body: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logger.Log.Errorf("Failed to fetch balance -> url: %v, statusCode: %v, msg: %v", baseURL, resp.StatusCode, string(body))
		return nil, errors.New("failed to fetch balance")
	}

	var positions []any
	if err := json.Unmarshal(body, &positions); err != nil {
		logger.Log.Errorf("Failed to convert data(positions): %v", err)
		return nil, err
	}

	var result model.Positions
	for _, position := range positions {
		positionMap := position.(map[string]any)
		balance, err := strconv.ParseFloat(positionMap["balance"].(string), 64)
		if err != nil {
			logger.Log.Errorf("Failed to parse balance: %v", err)
			continue
		}
		avgBuyPrice, err := strconv.ParseFloat(positionMap["avg_buy_price"].(string), 64)
		if err != nil {
			logger.Log.Errorf("Failed to parse avg_buy_price: %v", err)
			continue
		}

		result = append(result, model.Position{
			Status:     model.POSITION_BUY,
			Market:     positionMap["currency"].(string),
			Quantity:   balance,
			EntryPrice: avgBuyPrice,
			Profit:     0,
		})
	}

	return result, nil
}

func createJwt(accessKey, secretKey string, params map[string]string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["access_key"] = accessKey
	claims["nonce"] = uuid.NewString()
	if len(params) > 0 {
		claims["query_hash"] = generateQueryHash(params)
		claims["query_hash_alg"] = "SHA512"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(secretKey))
}

func generateQueryHash(params map[string]string) string {
	// Step 1: key 정렬
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Step 2: URL Query 생성
	values := url.Values{}
	for _, key := range keys {
		values.Add(key, params[key])
	}

	queryString := values.Encode()

	// Step 3: SHA-512 해싱
	hash := sha512.Sum512([]byte(queryString))
	return hex.EncodeToString(hash[:])
}
