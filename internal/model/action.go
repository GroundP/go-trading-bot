package model

import "fmt"

type Action struct {
	Market   string
	Signal   Signal
	Position Position
}

func (a Action) String() string {
	return fmt.Sprintf("[%v] Signal: %v, Position: %v", a.Market, a.Signal, a.Position)
}
