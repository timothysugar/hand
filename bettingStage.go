package hand

import "errors"

type bettingStage struct {
	initial       []*Player
	plays         []Input
	numCards      int
	makeCurrStage func(bettingStage) stage
	makeNextStage func([]*Player) stage
}

func newBettingStage(
	activePlayers []*Player,
	numCards int,
	curr func(bettingStage) stage,
	nextStageFact func([]*Player) stage,
) bettingStage {
	plays := make([]Input, 0)
	return bettingStage{activePlayers, plays, numCards, curr, nextStageFact}
}

func (bs bettingStage) requiredBet(h *Hand, p *Player) int {
	return h.pot.required(*p)
}

func (bs bettingStage) enter(h *Hand) error {
	existing := len(h.cards)
	if existing < bs.numCards {
		h.tableCard(bs.numCards - existing)
	}
	return nil
}

func (bs bettingStage) exit(h *Hand) error {
	h.playFromDealer()
	return nil
}

func (bs bettingStage) handleInput(h *Hand, p *Player, inp Input) (stage, error) {
	var err error
	switch inp.Action {
	case Fold:
		var remaining []*Player
		remaining, err = h.fold(p)
		if len(remaining) == 1 {
			bs.exit(h)
			return newWon(remaining), nil
		}
	case Call:
		err = h.call(p)
	case Check:
		err = h.check(p)
	case Raise:
		err = h.raise(p, inp.Chips)
	default:
		return nil, errors.New("unsupported input")
	}

	if err != nil {
		return nil, err
	}

	bs.plays = append(bs.plays, inp)
	if bs.allPlayed(h.pot) {
		bs.exit(h)
		return bs.makeNextStage(h.activePlayers()), nil
	}

	return bs.makeCurrStage(bs), nil
}

func (bs bettingStage) allPlayed(pot pot) bool {
	return (len(bs.plays) >= len(bs.initial) && !pot.outstandingStake())
}
