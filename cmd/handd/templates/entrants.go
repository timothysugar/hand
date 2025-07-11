package templates

import "github.com/timothysugar/hand/pkg/hand"

type PlayerViewModel struct {
	Id      string
	TableId string
	HandId  string
	Entrant
	Cards []struct {
		Card  hand.Card
		Class string
	}
	Moves []hand.Move
}

type OpponentViewModel struct {
	Entrant
	FaceDownCards []struct{}
}

type Entrant struct {
	Name   string
	Chips  int
	Bet    int
	Active bool
	Folded bool
}
