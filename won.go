package hand

import (
	"errors"
	"sort"
)

type won struct {
	ps []*player
}

func newWon(remaining []*player) won {
	return won{remaining}
}

func (curr won) id() string {
	return "won"
}

func (curr won) requiredBet(h *hand, p *player) int {
	return h.pot.required(*p)
}

type pHand struct {
	player *player
	cards  []card
}

type byHand []pHand

func (p byHand) Len() int      { return len(p) }
func (p byHand) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p byHand) Less(i, j int) bool { return false } // TODO: Sort by rank rather than deterministic sorting

func (curr won) enter(h *hand) error {
	// evaluate hands
	var pHands []pHand
	for _, v := range h.players {
		pH := pHand{v, append(h.cards, v.cards...)}
		pHands = append(pHands, pH)
	}
	sort.Sort(byHand(pHands))

	h.finish(finishedHand{pHands[0].player, h.pot.total()})
	return nil
}

func (curr won) exit(h *hand) error {
	return nil
}

func (curr won) handleInput(h *hand, p *player, inp input) (stage, error) {
	return nil, errors.New("no action can be taken after winning")
}
