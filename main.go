package hand

import (
	"errors"

	"github.com/rs/xid"
)

type player struct {
	id xid.ID
	chips int
}


func newPlayer(chips int) player {
	return player{
		id: xid.New(),
		chips: chips,
	}
}

func (p *player) bet(amount int) {
	p.chips = p.chips - amount
}

type pot struct {
	contribs map[player]int
}

func newPot() pot {
	return pot{
		contribs: make(map[player]int),
	}
}

func (p pot) add(pl *player, amount int) {
	// if (p.contribs[pl] == nil) {}
	pl.bet(amount)
	p.contribs[*pl] += amount
}

func (p pot) required(pl player) int {
	curr := p.contribs[pl]
	max := 0
	for _, v := range(p.contribs) {
		if (v > max) { max = v }
	}
	return max - curr
}

type Hand struct {
	players []player
	pot	pot
	smallBlind int
	bigBlind int
}

func (h *Hand) fold(p player) ([]player, error) {
	if (len(h.players) == 1) { return nil, errors.New("final player cannot fold") }

	var idx int
	for i, v := range(h.players) {
		if (v == p) {
			idx = i
			break
		}
	}
	ret := make([]player, len(h.players) - 1)
	copy(ret[:idx], h.players[:idx])
	copy(ret[idx:], h.players[idx+1:])
	h.players = ret

	return ret, nil
}

func (h *Hand) blind(p *player, amount int) {
	h.pot.add(p, amount)
}

func (h *Hand) call(p *player) {
	req := h.pot.required(*p)
	h.pot.add(p, req)
}

func (h *Hand) winner() *player {
	if (len(h.players) == 1) {
		return &h.players[0]
	}
	return nil
}

func newHand(ps []player, smallBlind int, bigBlind int) (*Hand, error) {
	if (len(ps) <= 1) { return nil, errors.New("hand requires at least 2 players") }
	return &Hand{ players:  ps, pot: newPot(), smallBlind: smallBlind, bigBlind: bigBlind }, nil
}