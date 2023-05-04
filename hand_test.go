package hand

import "testing"

const (
	initial    = 2
	bigBlind   = 2
	smallBlind = 1
)

func TestPenultimatePlayerFolds(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{p1, p2}

	done := make(chan finishedHand)
	go func() {
		h, _ := newHand(done, players, p1)

		err := h.fold(p1)

		if err != nil {
			t.Error("Error should be nil")
		}
		winner := h.winner()
		if winner != p2 {
			t.Errorf("Hand should have been won by %v but wasn't", p1)
		}
		checkPlayers(t, h.players, p2)
	}()
	v := <-done
	want := finishedHand{}
	if v != want {
		t.Errorf("expected %v but got %v", want, v)
	}
	_, ok := <-done
	if ok {
		t.Error("expected done channel to be closed")
	}
}

func TestPenultimatePlayerFoldsFromBlind(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{p1, p2}
	done := make(chan finishedHand)
	h, _ := newHand(done, players, p1, smallBlind)

	// preflop
	go func() {

	var err error
	err = h.blind(p1)
	if err != nil {
		t.Error(err)
	}

	var winner *player
	winner = h.winner()
	if winner != nil {
		t.Error("Hand should not yet have been won")
	}

	err = h.fold(p2)
	if err != nil {
		t.Error(err)
	}
	winner = h.winner()
	if winner != p1 {
		t.Errorf("Hand should have been won by %v but wasn't", p1)
	}
	}()
	<-done
	checkPlayers(t, h.players, p1)
}

func TestFinalPlayerCannotFold(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{p1, p2}

	done := make(chan finishedHand)
	h, _ := newHand(done, players, p1)

	go func() {

	var err error
	err = h.fold(p1)
	if err != nil {
		t.Error()
	}

	err = h.fold(p2)
	if err == nil {
		t.Error("Last player folding should return error but did not")
	}
	}()
	v := <-done
	want := finishedHand{}
	if v != want {
		t.Errorf("expected %v but got %v", want, v)
	}
	_, ok := <-done
	if ok {
		t.Error("expected done channel to be closed")
	}
}

func TestBlindPlayerCannotFold(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{p1, p2}

	done := make(chan finishedHand)
	h, _ := newHand(done, players, p1, 1)

	err := h.fold(p1)

	if err == nil {
		t.Error("Player playing blind cannot fold")
	}
}

func TestBlindPlayerCannotCall(t *testing.T) {
	p1 := newPlayer(1)
	p2 := newPlayer(1)
	players := []*player{p1, p2}

	done := make(chan finishedHand)
	h, _ := newHand(done, players, p1, 1)

	err := h.call(p1)

	if err == nil {
		t.Error("Player playing blind cannot call")
	}
}

func TestHandRequiresTwoPlayers(t *testing.T) {
	p1 := newPlayer(initial)
	players := []*player{p1}

	done := make(chan finishedHand)
	_, err := newHand(done, players, p1)

	if err == nil {
		t.Error("New hand with less than one player should return error but did not")
	}
}

func TestAllPlayersCallTheBlind(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	p3 := newPlayer(initial)
	players := []*player{p1, p2, p3}

	done := make(chan finishedHand)
	h, _ := newHand(done, players, p1, smallBlind, bigBlind)

	// preflop
	if p1.chips != initial {
		t.Errorf("Player 1 should have %d chips before playing blind but has %d", initial, p1.chips)
	}
	if p2.chips != initial {
		t.Errorf("Player 2 should have %d chips before playing blind but has %d", initial, p2.chips)
	}
	if p3.chips != initial {
		t.Errorf("Player 3 should have %d chips before playing blind but has %d", initial, p3.chips)
	}
	var err error
	err = h.blind(p1)
	if err != nil {
		t.Error(err)
	}
	if p1.chips != (initial - smallBlind) {
		t.Errorf("Player 1 should have %d chips after playing blind but has %d", initial-smallBlind, p1.chips)
	}
	err = h.blind(p2)
	if err != nil {
		t.Error(err)
	}
	if p2.chips != (initial - bigBlind) {
		t.Errorf("Player 2 should have %d chips after playing blind but has %d", initial-bigBlind, p2.chips)
	}

	// flop
	err = h.call(p3)
	if err != nil {
		t.Error(err)
	}
	if p3.chips != (initial - bigBlind) {
		t.Errorf("Player 3 should have %d chips after calling but has %d", initial-bigBlind, p3.chips)
	}
	err = h.call(p1)
	if err != nil {
		t.Error(err)
	}
	if p1.chips != (initial - bigBlind) {
		t.Errorf("Player 1 should have %d chips after calling but has %d", initial-bigBlind, p1.chips)
	}
	err = h.check(p2)
	if err != nil {
		t.Error(err)
	}
	if p2.chips != (initial - bigBlind) {
		t.Errorf("Player 2 should have %d chips after calling but has %d", initial-bigBlind, p2.chips)
	}
}

func TestOneFoldsAndOneCallsBlind(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	p3 := newPlayer(initial)
	players := []*player{p1, p2, p3}

	done := make(chan finishedHand)
	h, _ := newHand(done, players, p1, smallBlind)

	if p1.chips != initial {
		t.Errorf("Player 1 should have %d chips before playing blind but has %d", initial, p1.chips)
	}
	if p2.chips != initial {
		t.Errorf("Player 2 should have %d chips before playing blind but has %d", initial, p2.chips)
	}
	if p3.chips != initial {
		t.Errorf("Player 3 should have %d chips before playing blind but has %d", initial, p3.chips)
	}
	var err error
	err = h.blind(p1)
	if err != nil {
		t.Error(err)
	}
	if p1.chips != (initial - smallBlind) {
		t.Errorf("Player 1 should have %d chips after playing blind but has %d", initial-smallBlind, p1.chips)
	}
	err = h.fold(p2)
	if err != nil {
		t.Error(err)
	}
	err = h.call(p3)
	if err != nil {
		t.Error(err)
	}
	if p3.chips != (initial - smallBlind) {
		t.Errorf("Player 3 should have %d chips after calling but has %d", initial-smallBlind, p3.chips)
	}
	err = h.call(p1)
	if err != nil {
		t.Error(err)
	}
	if p1.chips != (initial - smallBlind) {
		t.Errorf("Player 1 should have %d chips after calling but has %d", initial-smallBlind, p1.chips)
	}
}

func TestOneFoldsAndOneChecksBlind(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	p3 := newPlayer(initial)
	players := []*player{p1, p2, p3}

	done := make(chan finishedHand)
	h, _ := newHand(done, players, p1, smallBlind)

	if p1.chips != initial {
		t.Errorf("Player 1 should have %d chips before playing blind but has %d", initial, p1.chips)
	}
	if p2.chips != initial {
		t.Errorf("Player 2 should have %d chips before playing blind but has %d", initial, p2.chips)
	}
	if p3.chips != initial {
		t.Errorf("Player 3 should have %d chips before playing blind but has %d", initial, p3.chips)
	}
	var err error
	err = h.blind(p1)
	if err != nil {
		t.Error(err)
	}
	if p1.chips != (initial - smallBlind) {
		t.Errorf("Player 1 should have %d chips after playing blind but has %d", initial-smallBlind, p1.chips)
	}
	err = h.fold(p2)
	if err != nil {
		t.Error(err)
	}
	err = h.call(p3)
	if err != nil {
		t.Error(err)
	}
	if p3.chips != (initial - smallBlind) {
		t.Errorf("Player 3 should have %d chips after calling but has %d", initial-smallBlind, p3.chips)
	}
	err = h.check(p1)
	if err != nil {
		t.Error(err)
	}
	if p1.chips != (initial - smallBlind) {
		t.Errorf("Player 1 should have %d chips after checking but has %d", initial-smallBlind, p1.chips)
	}
}

func TestCheckWhenBlindDueReturnsError(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{p1, p2}

	done := make(chan finishedHand)
	h, _ := newHand(done, players, p1, smallBlind)

	err := h.check(p1)
	if err == nil {
		t.Error()
	}
}

func TestCheckWhenBetDueReturnsError(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{p1, p2}

	done := make(chan finishedHand)
	h, _ := newHand(done, players, p1, smallBlind)

	var err error
	err = h.blind(p1)
	if err != nil {
		t.Error()
	}
	err = h.check(p2)
	if err == nil {
		t.Error()
	}
}

func TestBlindsPlayedFromDealer(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{p1, p2}

	var done chan finishedHand
	done = make(chan finishedHand)
	// p1 is first from the dealer
	h1, _ := newHand(done, players, p1, smallBlind)

	go func() {

	if p1.chips != initial {
		t.Errorf("Player 1 should have %d chips before playing blind but has %d", initial, p1.chips)
	}
	if p2.chips != initial {
		t.Errorf("Player 2 should have %d chips before playing blind but has %d", initial, p2.chips)
	}
	var err error
	err = h1.blind(p1)
	if err != nil {
		t.Error(err)
	}
	if p1.chips != initial-smallBlind {
		t.Errorf("Player has unexpected number of chips %d", p1.chips)
	}
	err = h1.fold(p2)
	if err != nil {
		t.Error(err)
	}
	if p2.chips != initial {
		t.Errorf("Player has unexpected number of chips %d", p2.chips)
	}
	}()
	v := <-done
	want := finishedHand{}
	if v != want {
		t.Errorf("expected %v but got %v", want, v)
	}
	_, ok := <-done
	if ok {
		t.Error("expected done channel to be closed")
	}

	// p2 is next from the dealer
	done = make(chan finishedHand)
	h2, err := newHand(done, players, p2, smallBlind)
	if err != nil {
		t.Error(err)
	}

	err = h2.blind(p2)
	if err != nil {
		t.Error(err)
	}
	if p2.chips != initial-smallBlind {
		t.Errorf("Player has unexpected number of chips %d", p2.chips)
	}
}

func TestBlindsPlayedInWrongOrderReturnsError(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{p1, p2}

	done := make(chan finishedHand)
	h, _ := newHand(done, players, p1, smallBlind)

	err := h.blind(p2)

	if err == nil {
		t.Error("Expected an error for out of order blind but none received")
	}
}

func TestSamePlayerCallsImmediatelyAfterBlindReturnsError(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{p1, p2}

	done := make(chan finishedHand)
	h, _ := newHand(done, players, p1, smallBlind)

	var err error
	err = h.blind(p1)
	if err != nil {
		t.Error()
	}

	err = h.call(p1)
	if err == nil {
		t.Error("Expected an error for out of order call but none received")
	}
}

func TestSecondCallInWrongOrderReturnsError(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	p3 := newPlayer(initial)
	players := []*player{p1, p2, p3}

	done := make(chan finishedHand)
	h, _ := newHand(done, players, p1, smallBlind)

	var err error
	err = h.blind(p1)
	if err != nil {
		t.Error()
	}

	err = h.call(p2)
	if err != nil {
		t.Error(err)
	}

	err = h.call(p1)
	if err == nil {
		t.Error("Expected an error for out of order call but none received")
	}
}

func TestCallBlindThenChecksToTheRiver(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{p1, p2}

	done := make(chan finishedHand)
	h, _ := newHand(done, players, p1, smallBlind)

	// preflop
	var err error
	err = h.blind(p1)
	if err != nil {
		t.Error(err)
	}

	// flop
	if len(h.cards) != 3 {
		t.Errorf("unexpected number of cards, got %d	", len(h.cards))
	}
	err = h.call(p2)
	if err != nil {
		t.Error(err)
	}
	err = h.check(p1)
	if err != nil {
		t.Error(err)
	}

	// turn
	if len(h.cards) != 4 {
		t.Errorf("unexpected number of cards, got %d	", len(h.cards))
	}
	err = h.check(p1)
	if err != nil {
		t.Error(err)
	}
	err = h.check(p2)
	if err != nil {
		t.Error(err)
	}

	// river
	if len(h.cards) != 5 {
		t.Errorf("unexpected number of cards, got %d	", len(h.cards))
	}
	// err = h.check(p1)
}

func (h *hand) check(p *player) error {
	return h.handleInput(p, input{action: Check, chips: 0})
}

func checkPlayers(t *testing.T, ps []*player, rem ...*player) {
	if len(ps) != len(rem) {
		t.Errorf("Player count should reduce by one; got: %v, want: %v", len(ps), len(rem))
	}
	for i, v := range rem {
		p := ps[i]
		if p != v {
			t.Errorf("Player %v should remain but players is %v", p, ps)
		}
	}
}
