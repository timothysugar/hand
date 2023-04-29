package hand

func (curr turn) id() string {
	return "turn"
}

func (curr turn) requiredBet(h *hand, p *player) int {
	return h.pot.required(*p)
}

func (curr turn) enter(h *hand) error {
	cs := len(h.cards)
	expected := 4
	if (cs < 4) { h.tableCard(expected - cs)}
	return nil
}

func (curr turn) exit(h *hand) error {
	h.nextToPlay = h.dealer
	return nil
}

func (curr turn) handleInput(h *hand, p *player, inp input) (stage, error) {
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