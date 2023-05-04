package hand

type pot struct {
	contribs map[string]int
}

func newPot() pot {
	return pot{
		contribs: make(map[string]int),
	}
}

func (p pot) add(pl *player, amount int) {
	pl.bet(amount)
	p.contribs[pl.id] += amount
}

func (p pot) maxStake() int {
	max := 0
	for _, v := range p.contribs {
		if v > max {
			max = v
		}
	}
	return max
}

func (p pot) outstandingStake() bool {
	var anyStake int
	for _, v := range p.contribs {
		anyStake = v
		break
	}
	for _, v := range p.contribs {
		if anyStake != v {
			return true
		}
	}
	return false
}

func (p pot) required(pl player) int {
	curr := p.contribs[pl.id]
	max := p.maxStake()
	return max - curr
}
