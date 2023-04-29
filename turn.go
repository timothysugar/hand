package hand

import "errors"

type turn struct {
	initial []*player
	plays []input
}

func newTurnState(initial []*player) turn {
	plays := make([]input, 0)
	return turn{ initial: initial, plays: plays }
}

func (curr turn) id() string {
	return "turn"
}

func (curr turn) requiredBet(h *hand, p *player) int {
	return h.pot.required(*p)
}

func (curr turn) enter(h *hand) error {
	cs := len(h.cards)
	expected := 4
	if (cs < expected) { h.tableCard(expected - cs)}
	return nil
}

func (curr turn) exit(h *hand) error {
	h.playFromDealer()
	return nil
}

func (curr turn) handleInput(h *hand, p *player, inp input) (stage, error) {
	var err error
	switch inp.action {
	case Fold:
		var remaining []*player
		remaining, err = h.doFold(p)
		if (len(remaining) == 1) { 
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

	if (err != nil) { return nil, err }

	curr.plays = append(curr.plays, inp)
	if (curr.allPlayed(h.pot)) { 
		curr.exit(h)
		return river{}, nil
	}

	return curr, nil
}

func (curr turn) allPlayed(pot pot) bool {
	return (len(curr.plays) >= len(curr.initial) && !pot.outstandingStake())
}