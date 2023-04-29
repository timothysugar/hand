package hand

import "errors"

func (p *player) bet(amount int) {
	p.chips = p.chips - amount
}

func (h *hand) doCheck(p *player) error {
	req := h.pot.required(*p)
	if (req != 0) { return errors.New("cannot check when required is not zero") }
	return nil
}


func (h *hand) doCall(p *player) error {
	req := h.pot.required(*p)
	h.pot.add(p, req)
	return nil
}

func (h *hand) call(p *player) error {
	req := h.stage.requiredBet(h, p)
	return h.handleInput(p, input{ action: Call, chips: req})
}