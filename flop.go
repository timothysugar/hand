package hand

type flop struct {
	bettingStage
}

func newFlopState(remaining []*player) flop {
	curr := func(bs bettingStage) stage {
		return flop{bs}
	}
	next := func(remaining []*player) stage {
		return newTurnState(remaining)
	}
	bs := newBettingStage(remaining, 3, curr, next)
	return flop{bs}
}

func (curr flop) id() string {
	return "flop"
}
