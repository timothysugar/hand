package hand

import (
	"github.com/rs/xid"
)

type Player struct {
	Id     string
	Name   string
	Chips  int
	Cards  []Card
	Folded bool
}

func NewPlayer(name string, chips int) *Player {
	id := xid.New().String()
	return &Player{
		Id:    id,
		Name:  name,
		Chips: chips,
	}
}

func (p *Player) String() string {
	return p.Id
}

func (p *Player) bet(amount int) {
	p.Chips = p.Chips - amount
}
