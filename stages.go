//go:generate stringer -type=Action

package hand

type stage interface {
	enter(h *Hand) error
	handleInput(h *Hand, p *Player, inp Input) (stage, error)
	validMoves(h *Hand) map[string][]Move
	requiredBet(h *Hand, p *Player) int // TODO: remove this - only used in tests
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
