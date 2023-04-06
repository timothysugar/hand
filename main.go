package hand

import (
	"errors"
	"fmt"

	"github.com/rs/xid"
)

type player struct {
	id xid.ID
	chips int
}


func newPlayer(chips int) *player {
	return &player{
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
	players []*player
	nextToPlay *player
	dealer *player
	pot	pot
	blinds map[*player]int
}

func (h *Hand) fold(p *player) ([]*player, error) {
	if (len(h.players) == 1) { return nil, errors.New("final player cannot fold") }

	var idx int
	for i, v := range(h.players) {
		if (v == p) {
			idx = i
			break
		}
	}
	ret := make([]*player, len(h.players) - 1)
	copy(ret[:idx], h.players[:idx])
	copy(ret[idx:], h.players[idx+1:])
	h.players = ret

	return ret, nil
}

type outOfTurnError struct {
    attempted player
    nextToPlay player
}

func (e outOfTurnError) Error() string {
	return fmt.Sprintf("%v is next to play but %v attempted", e.nextToPlay, e.attempted )
}

func (h *Hand) blind(p *player) error {
	if (p != h.nextToPlay) {
		return &outOfTurnError{ *p, *h.nextToPlay }
	}
	req := h.blinds[p]
	h.pot.add(p, req)
	h.nextMove()
	return nil
}

func (h *Hand) nextMove() {
	var playIdx int
	for i, v := range(h.players) {
		if (h.nextToPlay == v) {
			playIdx = i
		}
	}
	h.nextToPlay = h.players[(playIdx + 1) % len(h.players)]
}

func (h *Hand) call(p *player) error {
	if (p != h.nextToPlay) {
		return nil
		// return &outOfTurnError{ *p, *h.nextToPlay }
	}

	req := h.pot.required(*p)
	h.pot.add(p, req)
	return nil
}

func (h *Hand) winner() *player {
	if (len(h.players) == 1) {
		return h.players[0]
	}
	return nil
}

func newHand(ps []*player, dealer *player, blinds ...int) (*Hand, error) {
	if (len(ps) <= 1) { return nil, errors.New("hand requires at least 2 players") }
	bs := make(map[*player]int)
	// Find index of dealer in players
	var dIdx int
	for i, v := range(ps) {
		if (v == dealer) { dIdx = i}
	}
	// Set blinds from the dealer position. Assumes blinds passed from small and increasing
	var idx int
	for i, v := range(blinds) {
		idx = (dIdx + i) % len(ps)
		bs[ps[idx]] = v
	}
	return &Hand{ players:  ps, pot: newPot(), nextToPlay: dealer, blinds: bs }, nil
}