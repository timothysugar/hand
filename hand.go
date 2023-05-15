package hand

import (
	"errors"
	"fmt"
	"sync"
)

type Hand struct {
	players    []*Player
	m          sync.RWMutex
	finished   chan FinishedHand
	dealer     *Player
	nextToPlay *Player
	cards      []card
	stage      stage
	pot        pot
}

type FinishedHand struct {
	winner *Player
	chips  int
}

func NewHand(ch chan FinishedHand, ps []*Player, dealer *Player, blinds ...int) (*Hand, error) {
	if len(ps) <= 1 {
		return nil, errors.New("hand requires at least 2 players")
	}

	var dIdx int
	for i, v := range ps {
		if v == dealer {
			dIdx = i
		}
	}

	sortedPs := append(ps[dIdx:], ps[:dIdx]...)

	state, err := initialGameState(sortedPs, blinds)
	if err != nil {
		return nil, err
	}
	return &Hand{players: sortedPs, pot: newPot(), dealer: dealer, stage: state, nextToPlay: dealer, finished: ch}, nil
}

func (h *Hand) finish(fh FinishedHand) {
	h.finished <- fh
	close(h.finished)
}

func (h *Hand) playFromDealer() {
	h.nextToPlay = h.dealer
}

func (h *Hand) activePlayers() []*Player {
	h.m.Lock()
	defer h.m.Unlock()

	var active []*Player
	for _, v := range h.players {
		if !v.folded {
			active = append(active, v)
		}
	}

	return active
}

func (h *Hand) activePlayerAt(idx int) (*Player, error) {
	if idx > len(h.activePlayers()) {
		return nil, errors.New("index out of range of active players")
	}
	ps, err := h.activePlayersAt(idx, idx+1)
	return ps[0], err
}

func (h *Hand) activePlayersAt(startIdx int, endIdx int) ([]*Player, error) {
	if startIdx > endIdx {
		return nil, errors.New("start index must preceed or equal end index")
	}
	h.m.Lock()
	defer h.m.Unlock()

	var dIdx int
	var active []*Player

	for _, v := range h.players {
		if v == h.dealer {
			dIdx = len(active)
		}
		if !v.folded {
			active = append(active, v)
		}
	}

	if len(active) < (endIdx - startIdx) {
		return nil, errors.New("fewer active players than requested")
	}

	var i = (startIdx + dIdx) % len(active)
	var j = (endIdx + dIdx) % len(active)

	if i <= j {
		return active[i:j], nil
	}

	return append(active[i:], active[:j]...), nil
}

func initialGameState(ps []*Player, blinds []int) (stage, error) {
	if len(blinds) == 0 {
		return newFlopState(ps), nil
	}
	return newPreflop(ps, blinds)
}

type betTooLowError struct {
	player         Player
	betAmount      int
	requiredAmount int
}

func (e betTooLowError) Error() string {
	return fmt.Sprintf("bet of %d played by %v is too low; %d required", e.betAmount, e.player, e.requiredAmount)
}

type unexpectedBetAmountError struct {
	player    Player
	betAmount int
}

func (e unexpectedBetAmountError) Error() string {
	return fmt.Sprintf("bet of %d played by %v is unexpected", e.betAmount, e.player)
}

type outOfTurnError struct {
	attempted  *Player
	nextToPlay *Player
}

func (e outOfTurnError) Error() string {
	return fmt.Sprintf("%v is next to play but %v attempted", e.nextToPlay, e.attempted)
}

func (h *Hand) tableCard(num int) {
	cs := make([]card, num)
	h.cards = append(h.cards, cs...)
}

func (h *Hand) nextMove() {
	var playIdx int
	for i, v := range h.players {
		if h.nextToPlay == v {
			playIdx = i
		}
	}
	h.nextToPlay = h.players[(playIdx+1)%len(h.players)]
}

func (h *Hand) HandleInput(p *Player, inp Input) error {
	if p != h.nextToPlay {
		return &outOfTurnError{p, h.nextToPlay}
	}
	s, err := h.stage.handleInput(h, p, inp)
	if err != nil {
		return err
	}
	if s != nil {
		if s.id() != h.stage.id() {
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

func (h *Hand) fold(p *Player) ([]*Player, error) {
	if len(h.players) == 1 {
		return nil, errors.New("final player cannot fold")
	}
	var idx int
	for i, v := range h.players {
		if v == p {
			idx = i
			break
		}
	}
	ret := make([]*Player, len(h.players)-1)
	copy(ret[:idx], h.players[:idx])
	copy(ret[idx:], h.players[idx+1:])
	h.players = ret

	return ret, nil
}

func (h *Hand) check(p *Player) error {
	req := h.pot.required(*p)
	if req != 0 {
		return errors.New("cannot check when required is not zero")
	}
	return nil
}

func (h *Hand) call(p *Player) error {
	req := h.pot.required(*p)
	h.pot.add(p, req)
	return nil
}

func (h *Hand) raise(p *Player, bet int) error {
	req := h.pot.required(*p)
	if bet < req {
		return betTooLowError{*p, bet, req}
	}
	if bet == req {
		return unexpectedBetAmountError{*p, bet}
	}
	h.pot.add(p, bet)
	return nil
}
