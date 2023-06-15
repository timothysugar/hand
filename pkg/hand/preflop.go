package hand

import "errors"

type preflop struct {
	blinds map[*Player]blind
}

func newPreflop(remaining []*Player, blinds []int) (preflop, error) {
	bs := make(map[*Player]blind)

	for i, v := range remaining {
		if i >= len(blinds) {
			break
		}
		bs[v] = newBlind(blinds[i])
	}
	return preflop{blinds: bs}, nil
}

func (curr preflop) requiredBet(h *Hand, p *Player) int {
	blind := curr.blinds[p]
	return blind.required
}

func (curr preflop) enter(h *Hand) error {
	return nil
}

func (curr preflop) exit(h *Hand) error {
	next, err := h.activePlayerAt(len(curr.blinds))
	if err != nil {
		return err
	}
	h.nextToPlay = next
	return nil
}

func (curr preflop) handleInput(h *Hand, p *Player, inp Input) (stage, error) {
	switch inp.Action {
	case Blind:
		blinds := curr.blinds
		blind := blinds[p]
		b, err := blind.play(inp.Chips)
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

func (curr preflop) validMoves(h *Hand) map[string][]Move {
	mvs := make(map[string][]Move)
	plyr := h.nextToPlay
	req := curr.requiredBet(h, plyr)
	mvs[plyr.Id] = []Move{NewMove(Blind, NewExactBet(req))}
	return mvs
}
