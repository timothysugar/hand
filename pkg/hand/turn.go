package hand

type turn struct {
	bettingStage
}

func newTurnState(remaining []*Player) turn {
	curr := func(bs bettingStage) stage {
		return turn{bs}
	}
	next := func(remaining []*Player) stage {
		return newRiverState(remaining)
	}
	bs := newBettingStage(remaining, 4, curr, next)
	return turn{bs}
}
