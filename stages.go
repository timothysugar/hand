package hand

type stage interface {
	id() string
	enter (h *hand) error
	handleInput(h *hand, p *player, inp input) (stage, error)
	requiredBet(h *hand, p *player) int
	exit (h *hand) error
}

type input struct {
	action action
	chips int
}

type action int

const (
	Undefined action = iota
	Blind
	Check
	Fold
	Call
	Raise
)

type preflop struct {
	blinds map[*player]blind
}

type turn struct {
	plays []input
}

type won struct {
}

