package hand

type stage interface {
	id() string
	enter(h *Hand) error
	handleInput(h *Hand, p *Player, inp Input) (stage, error)
	requiredBet(h *Hand, p *Player) int
	exit(h *Hand) error
}

type Input struct {
	Action Action
	Chips  int
}

type Action int

const (
	Undefined Action = iota
	Blind
	Check
	Fold
	Call
	Raise
)
