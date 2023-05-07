package hand

import "errors"

type bettingStage struct {
	initial []*player
	plays   []input
	numCards int
	nextStageFact func([]*player) stage
}

func newBettingStage() bettingStage {
	plays := make([]input, 0)
	return bettingStage{plays} 

}

func (bs bettingStage) requiredBet(h *hand, p *player) int {
	return h.pot.required(*p)
}

func (bs bettingStage) enter(h *hand) error {
	existing := len(h.cards)
	if existing < bs.numCards {
		h.tableCard(bs.numCards - existing)
	}
	return nil
}

func (bs bettingStage) exit(h *hand) error {
	h.playFromDealer()
	return nil
}

func (bs bettingStage) handleInput(curr func(bettingStage) stage, h *hand, p *player, inp input) (stage, error) {
	var err error
	switch inp.action {
	case Fold:
		var remaining []*player
		remaining, err = h.fold(p)
		if len(remaining) == 1 {
			bs.exit(h)
			return newWon(remaining), nil
		}
	case Call:
		err = h.doCall(p)
	case Check:
		err = h.doCheck(p)
	case Raise:
		// todo
	default:
		return nil, errors.New("unsupported input")
	}

	if err != nil {
		return nil, err
	}

	bs.plays = append(bs.plays, inp)
	if bs.allPlayed(h.pot) {
		bs.exit(h)
		return bs.nextStageFact(h.activePlayers()), nil
	}

	return curr(bs), nil
}

func (bs bettingStage) allPlayed(pot pot) bool {
	return (len(bs.plays) >= len(bs.initial) && !pot.outstandingStake())
}
