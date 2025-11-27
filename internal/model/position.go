package model

import "fmt"

type PositionStatus string

const (
	POSITION_NONE PositionStatus = "NONE" // 포지션 없음
	POSITION_BUY  PositionStatus = "BUY"  // 매수 상태
)

type Position struct {
	Status     PositionStatus
	Market     string
	Quantity   float64
	EntryPrice float64
	Profit     float64
}

func (p *Position) String() string {
	return fmt.Sprintf("Market: %v, Status: %v, Quantity: %.2f, EntryPrice: %.2f, Profit: %.2f", p.Market, p.Status, p.Quantity, p.EntryPrice, p.Profit)
}
