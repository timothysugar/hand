package hand

import (
	"errors"
	"fmt"
	"sync"
)

type hand struct {
	players []*player
	m sync.RWMutex
	dealer *player
	nextToPlay *player
	cards []card
	stage stage
	pot	pot
}

func newHand(ps []*player, dealer *player, blinds ...int) (*hand, error) {
	if (len(ps) <= 1) { return nil, errors.New("hand requires at least 2 players") }

	var dIdx int
	for i, v := range(ps) {
		if (v == dealer) { dIdx = i}
	}

	sortedPs := append(ps[dIdx:], ps[:dIdx]...)

	state, err := initialGameState(sortedPs, blinds)
	if (err != nil) { return nil, err}
	return &hand{ players: sortedPs, pot: newPot(), dealer: dealer, stage: state, nextToPlay: dealer  }, nil
}

func (h *hand) activePlayers() []*player {
	h.m.Lock()
	defer h.m.Unlock()

	var active []*player
	for _, v := range(h.players) {
		if (!v.folded) { active = append(active, v) }
	}

	return active
}

func (h *hand) activePlayerAt(idx int) (*player, error) {
	if (idx > len(h.activePlayers())) { return nil, errors.New("index out of range of active players")}
	ps, err := h.activePlayersAt(idx, idx + 1)
	return ps[0], err
}

func (h *hand) activePlayersAt(startIdx int, endIdx int) ([]*player, error) {
	if (startIdx > endIdx) { return nil, errors.New("start index must preceed or equal end index")}
	h.m.Lock()
	defer h.m.Unlock()

	var dIdx int
	var active []*player

	for _, v := range(h.players) {
		if (v == h.dealer) { dIdx = len(active)}
		if (!v.folded) { active = append(active, v) }
	}

	if (len(active) < (endIdx - startIdx)) { return nil, errors.New("fewer active players than requested")}

	var i = (startIdx + dIdx) % len(active)
	var j = (endIdx + dIdx) % len(active)

	if (i <= j) { return active[i:j], nil }

	return append(active[i:], active[:j]...), nil
}

func initialGameState(ps []*player, blinds []int) (stage, error) {
	if (len(blinds) == 0) { return flop{}, nil}
	return newPreflop(ps, blinds)
}

type betTooLowError struct {
    player player
	betAmount int
    requiredAmount int
}

func (e betTooLowError) Error() string {
	return fmt.Sprintf("bet of %d played by %v is too low; %d required", e.betAmount, e.player, e.requiredAmount)
}

type outOfTurnError struct {
    attempted *player
    nextToPlay *player
}

func (e outOfTurnError) Error() string {
	return fmt.Sprintf("%v is next to play but %v attempted", e.nextToPlay, e.attempted )
}

func (h *hand) tableCard(num int) {
	cs := make([]card, num)
	h.cards = append(h.cards,  cs...)
}

func (h *hand) nextMove() {
	var playIdx int
	for i, v := range(h.players) {
		if (h.nextToPlay == v) {
			playIdx = i
		}
	}
	h.nextToPlay = h.players[(playIdx + 1) % len(h.players)]
}

func (h *hand) winner() *player {
	if (len(h.players) == 1) {
		return h.players[0]
	}
	return nil
}


func (h *hand) handleInput(p *player, inp input) error {
	if (p != h.nextToPlay) {
		return &outOfTurnError{ p, h.nextToPlay }
	}
	s, err := h.stage.handleInput(h, p, inp)
	if (err != nil) { return err }
	if (s != nil) {
		if (s.id() != h.stage.id()) {
			h.stage.exit(h)
			s.enter(h)
		} else {
			h.nextMove()
		}
		h.stage = s
		return nil
	}
	h.nextMove()
	return nil
}

func (h *hand) doFold(p *player) ([]*player, error) {
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

func (h *hand) fold(p *player) error {
	if (len(h.players) == 1) { return errors.New("final player cannot fold") }
	err := h.handleInput(p, input{ action: Fold, chips: 0 })
	if (err != nil) { return err }

	return nil
}

func (h *hand) blind(p *player) error {
	req := h.stage.requiredBet(h, p)
	return h.handleInput(p, input{ action: Blind, chips: req})
}