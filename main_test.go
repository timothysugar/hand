package hand

import (
	"testing"
)

const (
	initial = 2
	bigBlind = 2
	smallBlind = 1
)

func checkPlayers(t *testing.T, ps []*player, rem ...*player) {
	if (len(ps) != len(rem)) { t.Errorf("Player count should reduce by one; got: %v, want: %v", len(ps), len(rem))}
	for i, v := range(rem) {
		p := ps[i]
		if (p != v) { t.Errorf("Player %v should remain but players is %v", p, ps)}
	}
}

func TestPlayerOfManyFolds(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	p3 := newPlayer(initial)
	players := []*player{ p1, p2, p3 }
	h, _ := newHand(players, p1)

	result, err := h.fold(p1)
	
	if (err != nil) { t.Error() }

	checkPlayers(t, h.players, p2, p3)
	checkPlayers(t, result, p2, p3)
}

func TestPenultimatePlayerFolds(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{ p1, p2 }
	h, _ := newHand(players, p1)

	result, err := h.fold(p1)
	
	if (err != nil) { t.Error("Error should be nil" )}
	winner := h.winner()
	if (winner != p2) { t.Errorf("Hand should have been won by %v but wasn't", p1)}
	checkPlayers(t, h.players, p2)
	checkPlayers(t, result, p2)
}

func TestPenultimatePlayerFoldsFromBlind(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{ p1, p2 }
	h, _ := newHand(players, p1, smallBlind)

	var err error
	err = h.blind(p1)
	
	if (err != nil) { t.Error("Error should be nil" )}
	var winner *player
	winner = h.winner()
	if (winner != nil) { t.Errorf("Hand should not yet have been won")}
	_, err = h.fold(p2)
	if (err != nil) { t.Error("Error should be nil" )}
	winner = h.winner()
	if (winner != p1) { t.Errorf("Hand should have been won by %v but wasn't", p1)}
	checkPlayers(t, h.players, p1)
}

func TestFinalPlayerCannotFold(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{ p1, p2 }
	h, _ := newHand(players, p1)

	var err error
	_, err = h.fold(p1)
	if (err != nil) { t.Error() }

	_, err = h.fold(p2)
	if (err == nil) { t.Error("Last player folding should return error but did not") }
}

func TestBlindPlayerCannotFold(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{ p1, p2 }
	h, _ := newHand(players, p1, 1)

	_, err := h.fold(p1)

	if (err == nil) { t.Error("Player playing blind cannot fold")}
}

func TestBlindPlayerCannotCall(t *testing.T) {
	p1 := newPlayer(1)
	p2 := newPlayer(1)
	players := []*player{ p1, p2 }
	h, _ := newHand(players, p1, 1)

	err := h.call(p1)

	if (err == nil) { t.Error("Player playing blind cannot call")}
}

func TestHandRequiresTwoPlayers(t *testing.T) {
	p1 := newPlayer(initial)
	players := []*player{ p1 }

	_, err := newHand(players, p1)

	if (err == nil) { t.Error("New hand with less than one player should return error but did not")}
}

func TestAllPlayersCallTheBlind(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	p3 := newPlayer(initial)
	players := []*player{ p1, p2, p3 }
	h, _ := newHand(players, p1, smallBlind, bigBlind)

	if (p1.chips != initial) { t.Errorf("Player 1 should have %d chips before playing blind but has %d", initial, p1.chips)}
	if (p2.chips != initial) { t.Errorf("Player 2 should have %d chips before playing blind but has %d", initial, p2.chips)}
	if (p3.chips != initial) { t.Errorf("Player 3 should have %d chips before playing blind but has %d", initial, p3.chips)}
	var err error
	err = h.blind(p1)
	if (err != nil) { t.Error()}
	if (p1.chips != (initial - smallBlind)) { t.Errorf("Player 1 should have %d chips after playing blind but has %d", initial - smallBlind, p1.chips)}
	err = h.blind(p2)
	if (err != nil) { t.Error()}
	if (p2.chips != (initial - bigBlind)) { t.Errorf("Player 2 should have %d chips after playing blind but has %d", initial - bigBlind, p2.chips)}
	err = h.call(p3)
	if (err != nil) { t.Error()}
	if (p3.chips != (initial - bigBlind)) { t.Errorf("Player 3 should have %d chips after calling but has %d", initial - bigBlind, p3.chips)}
	err = h.call(p1)
	if (err != nil) { t.Error()}
	if (p1.chips != (initial - bigBlind)) { t.Errorf("Player 1 should have %d chips after calling but has %d", initial - bigBlind, p1.chips)}
	err = h.call(p2)
	if (err != nil) { t.Error()}
	if (p2.chips != (initial - bigBlind)) { t.Errorf("Player 2 should have %d chips after calling but has %d", initial - bigBlind, p2.chips)}
}

func TestOneFoldsAndOneCallsBlind(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	p3 := newPlayer(initial)
	players := []*player{ p1, p2, p3 }
	h, _ := newHand(players, p1, smallBlind)

	if (p1.chips != initial) { t.Errorf("Player 1 should have %d chips before playing blind but has %d", initial, p1.chips)}
	if (p2.chips != initial) { t.Errorf("Player 2 should have %d chips before playing blind but has %d", initial, p2.chips)}
	if (p3.chips != initial) { t.Errorf("Player 3 should have %d chips before playing blind but has %d", initial, p3.chips)}
	var err error
	err = h.blind(p1)
	if (err != nil) { t.Error()}
	if (p1.chips != (initial - smallBlind)) { t.Errorf("Player 1 should have %d chips after playing blind but has %d", initial - smallBlind, p1.chips)}
	_, err = h.fold(p2)
	if (err != nil) { t.Error()}
	err = h.call(p3)
	if (err != nil) { t.Error()}
	if (p3.chips != (initial - smallBlind)) { t.Errorf("Player 3 should have %d chips after calling but has %d", initial - smallBlind, p3.chips)}
	err = h.call(p1)
	if (err != nil) { t.Error()}
	if (p1.chips != (initial - smallBlind)) { t.Errorf("Player 1 should have %d chips after calling but has %d", initial - smallBlind, p1.chips)}
}

func TestOneFoldsAndOneChecksBlind(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	p3 := newPlayer(initial)
	players := []*player{ p1, p2, p3 }
	h, _ := newHand(players, p1, smallBlind)

	if (p1.chips != initial) { t.Errorf("Player 1 should have %d chips before playing blind but has %d", initial, p1.chips)}
	if (p2.chips != initial) { t.Errorf("Player 2 should have %d chips before playing blind but has %d", initial, p2.chips)}
	if (p3.chips != initial) { t.Errorf("Player 3 should have %d chips before playing blind but has %d", initial, p3.chips)}
	var err error
	err = h.blind(p1)
	if (err != nil) { t.Error()}
	if (p1.chips != (initial - smallBlind)) { t.Errorf("Player 1 should have %d chips after playing blind but has %d", initial - smallBlind, p1.chips)}
	_, err = h.fold(p2)
	if (err != nil) { t.Error()}
	err = h.call(p3)
	if (err != nil) { t.Error()}
	if (p3.chips != (initial - smallBlind)) { t.Errorf("Player 3 should have %d chips after calling but has %d", initial - smallBlind, p3.chips)}
	err = h.check(p1)
	if (err != nil) { t.Error()}
	if (p1.chips != (initial - smallBlind)) { t.Errorf("Player 1 should have %d chips after calling but has %d", initial - smallBlind, p1.chips)}
}

func TestCheckWhenBlindDueReturnsError(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{ p1, p2 }
	h, _ := newHand(players, p1, smallBlind)

	err := h.check(p1)
	if (err == nil) { t.Error() }
}

func TestCheckWhenBetDueReturnsError(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{ p1, p2 }
	h, _ := newHand(players, p1, smallBlind)

	var err error
	err = h.blind(p1)
	if (err != nil) { t.Error() }
	err = h.check(p2)
	if (err == nil) { t.Error() }
}

func TestBlindsPlayedFromDealer(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{ p1, p2 }
	// p1 is first from the dealer
	h1, _ := newHand(players, p1, smallBlind)

	if (p1.chips != initial) { t.Errorf("Player 1 should have %d chips before playing blind but has %d", initial, p1.chips)}
	if (p2.chips != initial) { t.Errorf("Player 2 should have %d chips before playing blind but has %d", initial, p2.chips)}
	h1.blind(p1)
	if (p1.chips != initial - smallBlind) { t.Errorf("Player has unexpected number of chips %d", p1.chips)}
	h1.fold(p2)
	if (p2.chips != initial) { t.Errorf("Player has unexpected number of chips %d", p2.chips)}

	// p2 is next from the dealer
	h2, _ := newHand(players, p2, smallBlind)

	h2.blind(p2)
	if (p2.chips != initial - smallBlind) { t.Errorf("Player has unexpected number of chips %d", p2.chips)}
}

func TestBlindsPlayedInWrongOrderReturnsError(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{ p1, p2 }
	h, _ := newHand(players, p1, smallBlind)

	err := h.blind(p2)

	if (err == nil) { t.Error("Expected an error for out of order blind but none received")}
}

func TestSamePlayerCallsImmediatelyAfterBlindReturnsError(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{ p1, p2 }
	h, _ := newHand(players, p1, smallBlind)

	var err error
	err = h.blind(p1)
	if (err != nil) { t.Error() }

	err = h.call(p1)
	if (err == nil) { t.Error("Expected an error for out of order call but none received")}
}

func TestSecondCallInWrongOrderReturnsError(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	p3 := newPlayer(initial)
	players := []*player{ p1, p2, p3 }
	h, _ := newHand(players, p1, smallBlind)

	var err error
	err = h.blind(p1)
	if (err != nil) { t.Error() }

	err = h.call(p2)
	if (err != nil) { t.Error()}

	err = h.call(p1)
	if (err == nil) { t.Error("Expected an error for out of order call but none received")}
}