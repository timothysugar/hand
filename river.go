package hand

type river struct {
	bettingStage
}

func newRiverState(remaining []*Player) river {
	curr := func(bs bettingStage) stage {
		return river{bs}
	}
	next := func(remaining []*Player) stage {
		return newWon(remaining)
	}
	bs := newBettingStage(remaining, 5, curr, next)
	return river{bs}
}

func (curr river) id() string {
	return "river"
}
