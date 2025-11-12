// Package model
package model

type SignalType int

const (
	BUY  SignalType = iota // 매수 신호
	SELL                   // 매도 신호
	HOLD                   // 관망
)

func (s SignalType) String() string {
	return [...]string{"BUY", "SELL", "HOLD"}[s]
}

type Signal struct {
	Type   SignalType
	Market string

	CurrentPrice float64
	Timestamp    string

	Description  string
	StrategyName string
}
