package hand

import "errors"

type flop struct {
	initial []*player
	plays []input
	// players []*player
}

func newFlopState(initial []*player) flop {
	plays := make([]input, 0)
	return flop{ initial: initial, plays: plays }
}

func (curr flop) id() string {
	return "flop"
}

func (curr flop) requiredBet(h *hand, p *player) int {
	return h.pot.required(*p)
}

func (curr flop) enter(h *hand) error {
	cs := len(h.cards)
	expected := 3
	if (cs < 3) { h.tableCard(expected - cs)}
	return nil
}

func (curr flop) exit(h *hand) error {
	h.nextToPlay = h.dealer
	return nil
}

func (curr flop) handleInput(h *hand, p *player, inp input) (stage, error) {
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
		return turn{}, nil
	}

	return curr, nil
}

func (curr flop) allPlayed(pot pot) bool {
	return (len(curr.plays) >= len(curr.initial) && !pot.outstandingStake())
}