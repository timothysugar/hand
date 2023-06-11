package hand

import "math"

type Move struct {
	Action Action
	Bet    RequiredBet
}

type RequiredBet struct {
	Minimum int
	Maximum int
}

func NewExactBet(i int) RequiredBet {
	return RequiredBet{Minimum: i, Maximum: i}
}

func NewMinumumBet(min int) RequiredBet {
	return RequiredBet{Minimum: min, Maximum: math.MaxInt32}
}

func NewMove(a Action, b RequiredBet) Move { return Move{a, b} }
