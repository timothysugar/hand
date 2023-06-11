package hand

type flop struct {
	bettingStage
}

func newFlopState(remaining []*Player) flop {
	curr := func(bs bettingStage) stage {
		return flop{bs}
	}
	next := func(remaining []*Player) stage {
		return newTurnState(remaining)
	}
	bs := newBettingStage(remaining, 3, curr, next)
	return flop{bs}
}
