package hand

import "testing"

func TestCallBlindThenChecksToTheRiver(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{ p1, p2 }
	h, _ := newHand(players, p1, smallBlind)
	
	// preflop
	var err error
	err = h.blind(p1)
	if (err != nil) { t.Error(err) }

	// flop
	if (len(h.cards) != 3) { t.Errorf("unexpected number of cards, got %d	", len(h.cards)) }
	err = h.call(p2)
	if (err != nil) { t.Error(err)}
	err = h.check(p1)
	if (err != nil) { t.Error(err) }

	// turn
	if (len(h.cards) != 4) { t.Errorf("unexpected number of cards, got %d	", len(h.cards)) }
	err = h.check(p1)
	if (err != nil) { t.Error(err) }

}

func (h *hand) check(p *player) error {
	return h.handleInput(p, input{ action: Check, chips: 0})
}
