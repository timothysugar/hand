package hand

import "github.com/rs/xid"

type player struct {
	id xid.ID
	chips int
	folded bool
}

func newPlayer(chips int) *player {
	return &player{
		id: xid.New(),
		chips: chips,
	}
}

func (p *player) String() string {
	return p.id.String()
}