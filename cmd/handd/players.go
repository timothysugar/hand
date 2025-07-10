package main

import "errors"

type Players struct {
	players []string
	limit   int
}

func NewPlayers(max int) Players {
	return Players{
		players: make([]string, 0),
		limit:   max,
	}
}

func (ps *Players) Add(name string) error {
	if len(ps.players) >= ps.limit {
		return errors.New("maximum number of players reached")
	}
	ps.players = append(ps.players, name)
	return nil
}
