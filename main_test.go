package hand

import (
	"testing"
)

func checkPlayers(t *testing.T, ps []*player, rem ...*player) {
	if (len(ps) != len(rem)) { t.Errorf("Player count should reduce by one; got: %v, want: %v", len(ps), len(rem))}
	for i, v := range(rem) {
		p := ps[i]
		if (p != v) { t.Errorf("Player %v should remain but players is %v", p, ps)}
	}
}

func TestPlayerOfManyFolds(t *testing.T) {
	p1 := newPlayer(2)
	p2 := newPlayer(2)
	p3 := newPlayer(2)
	players := []*player{ p1, p2, p3 }
	h, _ := newHand(players, p1)

	result, err := h.fold(p1)
	
	if (err != nil) { t.Error("Error should be nil")}

	checkPlayers(t, h.players, p2, p3)
	checkPlayers(t, result, p2, p3)
}

func TestPenultimatePlayerFolds(t *testing.T) {
	p1 := newPlayer(2)
	p2 := newPlayer(2)
	players := []*player{ p1, p2 }
	h, _ := newHand(players, p1, 1, 2)

	result, err := h.fold(p1)
	
	if (err != nil) { t.Error("Error should be nil" )}
	winner := h.winner()
	if (winner != p2) { t.Errorf("Hand should have been won by %v but wasn't", p1)}
	checkPlayers(t, h.players, p2)
	checkPlayers(t, result, p2)
}

func TestFinalPlayerCannotFold(t *testing.T) {
	p1 := newPlayer(2)
	p2 := newPlayer(2)
	players := []*player{ p1, p2 }
	h, _ := newHand(players, p1)

	h.fold(p1)
	_, err := h.fold(p2)

	if (err == nil) { t.Errorf("Last player folding should return error but did not; err: %v", err)}
}

func TestHandRequiresTwoPlayers(t *testing.T) {
	p1 := newPlayer(2)
	players := []*player{ p1 }

	_, err := newHand(players, p1)

	if (err == nil) { t.Error("New hand with less than one player should return error but did not")}
}

func TestPlayerPlaysOnlyBlind(t *testing.T) {
	blind := 1
	initial := blind
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{ p1, p2 }
	h, _ := newHand(players, p1, blind)

	if (p1.chips != initial) { t.Errorf("Player 1 should have %d chips before playing blind but has %d", initial, p1.chips)}
	if (p2.chips != initial) { t.Errorf("Player 2 should have %d chips before playing blind but has %d", initial, p2.chips)}
	err := h.blind(p1)
	if (err != nil) { t.Error() }
	if (p1.chips != 0) { t.Errorf("Player should have %d chip after playing small blind but has %d", initial - blind, p1.chips)}
}

func TestBothPlayersPlayBlinds(t *testing.T) {
	smallBlind := 1
	bigBlind := 2
	p1 := newPlayer(bigBlind)
	p2 := newPlayer(bigBlind)
	players := []*player{ p1, p2 }
	h, _ := newHand(players, p1, smallBlind, bigBlind)

	if (p1.chips != bigBlind) { t.Errorf("Player 1 should have 2 chips before playing blind but has %d", p1.chips)}
	if (p2.chips != bigBlind) { t.Errorf("Player 2 should have 2 chips before playing blind but has %d", p2.chips)}
	var err error
	err = h.blind(p1)
	if (err != nil) { t.Error() }
	if (p1.chips != 1) { t.Errorf("Player should have 1 chip after playing small blind but has %d", p1.chips)}
	err = h.blind(p2)
	if (err != nil) { t.Error() }
	if (p2.chips != 0) { t.Errorf("Player should have 0 chips after playing big blind but has %d", p2.chips)}
}


func TestPlayerCallsTheBlind(t *testing.T) {
	blind := 1
	p1 := newPlayer(blind)
	p2 := newPlayer(blind)
	players := []*player{ p1, p2 }
	h, _ := newHand(players, p1, blind)

	if (p1.chips != blind) { t.Errorf("Player 1 should have 2 chips before playing blind but has %d", p1.chips)}
	if (p2.chips != blind) { t.Errorf("Player 2 should have 2 chips before playing blind but has %d", p2.chips)}
	h.blind(p1)
	if (p1.chips != 0) { t.Errorf("Player 1 should have 0 chips after playing blind but has %d", p1.chips)}
	h.call(p2)
	if (p2.chips != 0) { t.Errorf("Player 2 should have 0 chips after calling but has %d", p2.chips)}
}

func TestBlindsPlayedFromDealer(t *testing.T) {
	blind := 1
	initial := blind * 2
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{ p1, p2 }
	// p1 is first from the dealer
	h1, _ := newHand(players, p1, blind)

	if (p1.chips != initial) { t.Errorf("Player 1 should have 2 chips before playing blind but has %d", p1.chips)}
	if (p2.chips != initial) { t.Errorf("Player 2 should have 2 chips before playing blind but has %d", p2.chips)}
	h1.blind(p1)
	if (p1.chips != initial - blind) { t.Errorf("Player has unexpected number of chips %d", p1.chips)}
	h1.fold(p2)
	if (p2.chips != initial) { t.Errorf("Player has unexpected number of chips %d", p2.chips)}

	// p2 is first from the dealer
	h2, _ := newHand(players, p2, blind)

	h2.blind(p2)
	if (p2.chips != initial - blind) { t.Errorf("Player has unexpected number of chips %d", p2.chips)}
}

func TestBlindsPlayedInWrongOrderReturnsError(t *testing.T) {
	blind := 1
	initial := blind * 2
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{ p1, p2 }
	h, _ := newHand(players, p1, blind)

	err := h.blind(p2)

	if (err == nil) { t.Error("Expected an error for out of order blind but none received")}
}

func TestCalledInWrongOrderReturnsError(t *testing.T) {
	blind := 1
	initial := blind * 2
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{ p1, p2 }
	h, _ := newHand(players, p1, blind)

	var err error
	err = h.blind(p1)
	if (err != nil) { t.Error() }

	err = h.call(p1)
	if (err == nil) { t.Error("Expected an error for out of order call but none received")}
}