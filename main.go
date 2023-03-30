package hand

import (
	"errors"

	"github.com/rs/xid"
)

type player struct {
	id xid.ID
}

func newPlayer() player {
	return player{
		id: xid.New(),
	}
}

type Hand struct {
	players []player
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

func (h *Hand) winner() *player {
	if (len(h.players) == 1) {
		return &h.players[0]
	}
	return nil
}

func newHand(ps []player) (*Hand, error) {
	if (len(ps) <= 1) { return nil, errors.New("hand requires at least 2 players") }
	return &Hand{ players:  ps }, nil
}