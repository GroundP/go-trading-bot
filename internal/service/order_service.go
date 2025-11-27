package service

import (
	"go-trading-bot/config"
	"go-trading-bot/internal/logger"
	"go-trading-bot/internal/model"
)

type OrderService struct {
	positions map[string]model.Position
}

func (o *OrderService) GetPosition(market string) *model.Position {
	if position, exists := o.positions[market]; exists {
		return &position
	}
	return nil
}

func (o *OrderService) SetPosition(market string, position *model.Position) {
	o.positions[market] = *position
}

func (o *OrderService) RemovePosition(market string) {
	delete(o.positions, market)
}

func (o *OrderService) PlaceOrder(market string, signalType model.SignalType, currentPrice float64) {
	if signalType == model.BUY {
		orderAmount := config.GetTradingConfig().OrderAmount
		quantity := float64(int((orderAmount/currentPrice)*10000)) / 10000
		position := &model.Position{
			Market:     market,
			Status:     model.POSITION_BUY,
			Quantity:   quantity,
			EntryPrice: currentPrice,
			Profit:     0,
		}
		logger.Log.Infof("[%v] 매수 주문을 실행합니다. 주문 금액: %v, 주문 수량: %v", market, orderAmount, quantity)
		logger.Log.Infof("[%v] 포지션 정보: %v", market, position)
		o.SetPosition(market, position)
	} else if signalType == model.SELL {
		position := o.GetPosition(market)
		if position == nil {
			logger.Log.Infof("[%v] 포지션이 없습니다.", market)
			return
		}

		profit := (currentPrice - position.EntryPrice) * position.Quantity
		position.Profit = profit
		position.Status = model.POSITION_NONE

		logger.Log.Infof("[%v] 매도 주문을 실행합니다. 포지션 수량: %v, 수익: %v", market, position.Quantity, profit)
		logger.Log.Infof("[%v] 포지션 정보: %v", market, position)
		o.RemovePosition(market)
	}
}
