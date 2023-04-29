package hand

import "errors"

func (curr won) id() string {
	return "won"
}

func (curr won) requiredBet(h *hand, p *player) int {
	return h.pot.required(*p)
}

func (curr won) enter(h *hand) error {
	// distribute pot
	return nil
}

func (curr won) exit(h *hand) error {
	return nil
}

func (curr won) handleInput(h *hand, p *player, inp input) (stage, error) {
	return nil, errors.New("no action can be taken after winning")
}