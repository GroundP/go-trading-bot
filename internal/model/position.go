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

func (p Position) String() string {
	return fmt.Sprintf("[%v] Quantity: %f, EntryPrice: %f, Profit: %f", p.Market, p.Quantity, p.EntryPrice, p.Profit)
}

type Positions []Position

func (p Positions) String() string {
	result := ""
	for _, position := range p {
		result += position.String() + "\n"
	}
	return result
}
