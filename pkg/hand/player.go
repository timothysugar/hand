package hand

import "github.com/rs/xid"

type Player struct {
	id     string
	chips  int
	cards  []card
	folded bool
}

func NewPlayer(chips int) *Player {
	return &Player{
		id:    xid.New().String(),
		chips: chips,
	}
}

func (p *Player) String() string {
	return p.id
}

func (p *Player) bet(amount int) {
	p.chips = p.chips - amount
}
