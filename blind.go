package hand

import "errors"

type blind struct {
	required    int
	contributed int
}

func newBlind(required int) blind {
	return blind{
		required: required,
	}
}

func (b blind) played() bool {
	return b.contributed >= b.required
}

func (b blind) play(value int) (*blind, error) {
	if value != b.required {
		return nil, errors.New("blind value played does not match required")
	}
	ret := newBlind(b.required)
	ret.contributed = value
	return &ret, nil
}
