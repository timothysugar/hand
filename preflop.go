package hand

import "errors"

type preflop struct {
	blinds map[*player]blind
}

func newPreflop(remaining []*player, blinds []int) (preflop, error) {
	bs := make(map[*player]blind)

	for i, v := range remaining {
		if i >= len(blinds) {
			break
		}
		bs[v] = newBlind(blinds[i])
	}
	return preflop{blinds: bs}, nil
}

func (curr preflop) requiredBet(h *hand, p *player) int {
	blind := curr.blinds[p]
	return blind.required
}

func (curr preflop) enter(h *hand) error {
	return nil
}

func (curr preflop) exit(h *hand) error {
	next, err := h.activePlayerAt(len(curr.blinds))
	if err != nil {
		return err
	}
	h.nextToPlay = next
	return nil
}

func (curr preflop) id() string {
	return "preflop"
}

func (curr preflop) handleInput(h *hand, p *player, inp input) (stage, error) {
	switch inp.action {
	case Blind:
		blinds := curr.blinds
		blind := blinds[p]
		b, err := blind.play(inp.chips)
		if err != nil {
			return nil, err
		}
		h.pot.add(p, blind.required)
		blinds[p] = *b
		for _, v := range blinds {
			if v.required != 0 && !v.played() {
				return curr, nil
			}
		}
		curr.exit(h)
		return newFlopState(h.activePlayers()), nil
	default:
		return nil, errors.New("unsupported action in preflop")
	}
}
