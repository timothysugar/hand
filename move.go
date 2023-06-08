package hand

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

func NewBetRange(min, max int) RequiredBet {
	return RequiredBet{Minimum: min, Maximum: max}
}
func NewMove(a Action, b RequiredBet) Move { return Move{a, b} }
