package hand

type turn struct {
	bettingStage
}

func newTurnState(remaining []*player) turn {
	curr := func(bs bettingStage) stage {
		return turn{bs}
	}
	next := func(remaining []*player) stage {
		return newRiverState(remaining)
	}
	bs := newBettingStage(remaining, 4, curr, next)
	return turn{bs}
}

func (curr turn) id() string {
	return "turn"
}
