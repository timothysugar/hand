package hand

type river struct {
	plays []input
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
	if (cs < expected) { h.tableCard(expected - cs)}
	return nil
}

func (curr river) exit(h *hand) error {
	h.playFromDealer()
	return nil
}

func (curr river) handleInput(h *hand, p *player, inp input) (stage, error) {
	if (len(curr.plays) >= len(h.players) && !h.pot.outstandingStake()) { return turn{}, nil}
	switch inp.action {
	case Fold:
		_, err := h.doFold(p);
		if (err != nil) { return nil, err }
		if (len(h.players) == 1) { 
			curr.exit(h)
			return won{}, nil 
		}
	case Call:
		err := h.call(p)
		if (err != nil) { return nil, err }
		plays := append(curr.plays, inp)
		curr.exit(h)
		return flop{ plays: plays }, nil // TODO river
	case Raise:
		// TODO
	}
	return curr, nil
}