package hand

import (
	"errors"
	"sort"
)

type won struct {
	ps []*player
}

func (curr won) newWon(remaining []*player) won {
	return won{remaining}
}

func (curr won) id() string {
	return "won"
}

func (curr won) requiredBet(h *hand, p *player) int {
	return h.pot.required(*p)
}

type pHand struct {
	cards []card
}

func (pHand) rank() int {
	return 0
}

type byHand []pHand

func (p byHand) Len() int      { return len(p) }
func (p byHand) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// func (p ByHand) Less(i, j int) bool { return p[i].rank() < p[j].rank() }
func (p byHand) Less(i, j int) bool { return false } // TODO: Sort by rank rather than deterministic sorting

func (curr won) enter(h *hand) error {
	var pHands []pHand
	for _, v := range h.players {
		pH := pHand{append(h.cards, v.cards...)}
		pHands = append(pHands, pH)
	}
	sort.Sort(byHand(pHands))
	// evaluate hands
	// distribute pot
	h.finish()
	return nil
}

func (curr won) exit(h *hand) error {
	return nil
}

func (curr won) handleInput(h *hand, p *player, inp input) (stage, error) {
	return nil, errors.New("no action can be taken after winning")
}
