package hand

import "errors"

type river struct {
	initial []*player
	plays   []input
}

func newRiverState(initial []*player) river {
	plays := make([]input, 0)
	return river{initial: initial, plays: plays}
}

func (curr river) id() string {
	return "river"
}

func (curr river) requiredBet(h *hand, p *player) int {
	return h.pot.required(*p)
}

func (curr river) enter(h *hand) error {
	cs := len(h.cards)
	expected := 5
	if cs < expected {
		h.tableCard(expected - cs)
	}
	return nil
}

func (curr river) exit(h *hand) error {
	h.playFromDealer()
	return nil
}

func (curr river) handleInput(h *hand, p *player, inp input) (stage, error) {
	var err error
	switch inp.action {
	case Fold:
		var remaining []*player
		remaining, err = h.fold(p)
		if len(remaining) == 1 {
			curr.exit(h)
			return won{}, nil
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

	curr.plays = append(curr.plays, inp)
	if curr.allPlayed(h.pot) {
		curr.exit(h)
		return newWon(h.activePlayers()), nil
	}

	return curr, nil
}

func (curr river) allPlayed(pot pot) bool {
	return (len(curr.plays) >= len(curr.initial) && !pot.outstandingStake())
}
