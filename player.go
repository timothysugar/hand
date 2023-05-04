package hand

import "github.com/rs/xid"

type player struct {
	id     string
	chips  int
	cards  []card
	folded bool
}

func newPlayer(chips int) *player {
	return &player{
		id:    xid.New().String(),
		chips: chips,
	}
}

func (p *player) String() string {
	return p.id
}
