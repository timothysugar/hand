package hand

import (
	"errors"
	"sort"
)

type won struct {
	ps []*Player
}

func newWon(remaining []*Player) won {
	return won{remaining}
}

func (curr won) requiredBet(h *Hand, p *Player) int {
	return h.pot.required(*p)
}

type pHand struct {
	player *Player
	cards  []Card
}

type byHand []pHand

func (p byHand) Len() int      { return len(p) }
func (p byHand) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p byHand) Less(i, j int) bool { return false } // TODO: Sort by rank rather than deterministic sorting

func (curr won) enter(h *Hand) error {
	// evaluate hands
	var pHands []pHand
	for _, v := range h.players {
		pH := pHand{v, append(h.Cards, v.Cards...)}
		pHands = append(pHands, pH)
	}
	sort.Sort(byHand(pHands))

	h.finish(FinishedHand{pHands[0].player, h.pot.total()})
	return nil
}

func (curr won) exit(h *Hand) error {
	return nil
}

func (curr won) handleInput(h *Hand, p *Player, inp Input) (stage, error) {
	return nil, errors.New("no action can be taken after winning")
}

func (curr won) validMoves(h *Hand) map[string][]Move {
	return make(map[string][]Move)
}
