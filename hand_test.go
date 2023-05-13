package hand

import (
	"testing"
)

const (
	initial    = 10
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

		if err := playFold(h, p1); err != nil {
			t.Error("Error should be nil")
		}
		checkPlayers(t, h.players, p2)
	}()
	v := <-done
	want := finishedHand{winner: p2, chips: 0}
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

	go func() {
		var err error
		err = playBlind(h, p1)
		if err != nil {
			t.Error(err)
		}

		err = playFold(h, p2)
		if err != nil {
			t.Error(err)
		}
	}()

	v := <-done
	want := finishedHand{winner: p1, chips: 1}
	if v != want {
		t.Errorf("expected %v but got %v", want, v)
	}
}

func TestFinalPlayerCannotFold(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{p1, p2}

	done := make(chan finishedHand)
	h, _ := newHand(done, players, p1)

	go func() {
		var err error
		err = playFold(h, p1)
		if err != nil {
			t.Error()
		}

		err = playFold(h, p2)
		if err == nil {
			t.Error("last player folding should return error")
		}
	}()
	v := <-done
	want := finishedHand{winner: p2}
	if v != want {
		t.Errorf("expected %v but got %v", want, v)
	}
}

func TestBlindPlayerCannotFold(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{p1, p2}

	done := make(chan finishedHand)
	h, _ := newHand(done, players, p1, 1)

	err := playFold(h, p1)

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

	err := playCall(h, p1)

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
	err = playBlind(h, p1)
	if err != nil {
		t.Error(err)
	}
	if p1.chips != (initial - smallBlind) {
		t.Errorf("Player 1 should have %d chips after playing blind but has %d", initial-smallBlind, p1.chips)
	}
	err = playBlind(h, p2)
	if err != nil {
		t.Error(err)
	}
	if p2.chips != (initial - bigBlind) {
		t.Errorf("Player 2 should have %d chips after playing blind but has %d", initial-bigBlind, p2.chips)
	}

	// flop
	err = playCall(h, p3)
	if err != nil {
		t.Error(err)
	}
	if p3.chips != (initial - bigBlind) {
		t.Errorf("Player 3 should have %d chips after calling but has %d", initial-bigBlind, p3.chips)
	}
	err = playCall(h, p1)
	if err != nil {
		t.Error(err)
	}
	if p1.chips != (initial - bigBlind) {
		t.Errorf("Player 1 should have %d chips after calling but has %d", initial-bigBlind, p1.chips)
	}
	err = playCheck(h, p2)
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
	err = playBlind(h, p1)
	if err != nil {
		t.Error(err)
	}
	if p1.chips != (initial - smallBlind) {
		t.Errorf("Player 1 should have %d chips after playing blind but has %d", initial-smallBlind, p1.chips)
	}
	err = playFold(h, p2)
	if err != nil {
		t.Error(err)
	}
	err = playCall(h, p3)
	if err != nil {
		t.Error(err)
	}
	if p3.chips != (initial - smallBlind) {
		t.Errorf("Player 3 should have %d chips after calling but has %d", initial-smallBlind, p3.chips)
	}
	err = playCall(h, p1)
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
	err = playBlind(h, p1)
	if err != nil {
		t.Error(err)
	}
	if p1.chips != (initial - smallBlind) {
		t.Errorf("Player 1 should have %d chips after playing blind but has %d", initial-smallBlind, p1.chips)
	}
	err = playFold(h, p2)
	if err != nil {
		t.Error(err)
	}
	err = playCall(h, p3)
	if err != nil {
		t.Error(err)
	}
	if p3.chips != (initial - smallBlind) {
		t.Errorf("Player 3 should have %d chips after calling but has %d", initial-smallBlind, p3.chips)
	}
	err = playCheck(h, p1)
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

	err := playCheck(h, p1)
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
	err = playBlind(h, p1)
	if err != nil {
		t.Error()
	}
	err = playCheck(h, p2)
	if err == nil {
		t.Error()
	}
}

func TestBlindsPlayedFromDealer(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{p1, p2} // p2 is not the first in order

	done := make(chan finishedHand)
	h, err := newHand(done, players, p2, smallBlind)
	if err != nil {
		t.Error(err)
	}

	err = playBlind(h, p2)
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

	err := playBlind(h, p2)

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
	err = playBlind(h, p1)
	if err != nil {
		t.Error()
	}

	err = playCall(h, p1)
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
	err = playBlind(h, p1)
	if err != nil {
		t.Error()
	}

	err = playCall(h, p2)
	if err != nil {
		t.Error(err)
	}

	err = playCall(h, p1)
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

	go func() {
		// preflop
		var err error
		err = playBlind(h, p1)
		if err != nil {
			t.Error(err)
		}

		// flop
		if len(h.cards) != 3 {
			t.Errorf("unexpected number of cards, got %d	", len(h.cards))
		}
		err = playCall(h, p2)
		if err != nil {
			t.Error(err)
		}
		err = playCheck(h, p1)
		if err != nil {
			t.Error(err)
		}

		// turn
		if len(h.cards) != 4 {
			t.Errorf("unexpected number of cards, got %d	", len(h.cards))
		}
		err = playCheck(h, p1)
		if err != nil {
			t.Error(err)
		}
		err = playCheck(h, p2)
		if err != nil {
			t.Error(err)
		}

		// river
		if len(h.cards) != 5 {
			t.Errorf("unexpected number of cards, got %d	", len(h.cards))
		}
		err = playCheck(h, p1)
		if err != nil {
			t.Error(err)
		}
		err = playCheck(h, p2)
		if err != nil {
			t.Error(err)
		}
		// players hands evaluated
	}()
	v := <-done
	want := finishedHand{winner: p1, chips: smallBlind * len(players)}
	if v != want {
		t.Errorf("expected %v but got %v", want, v)
	}
}

func TestPlayerFoldsAfterReraise(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{p1, p2}

	done := make(chan finishedHand)
	go func() {
		h, _ := newHand(done, players, p1, smallBlind)

		// preflop
		if len(h.cards) != 0 {
			t.Errorf("unexpected number of cards, got %d", len(h.cards))
		}

		if err := playBlind(h, p1); err != nil {
			t.Error("Error should be nil")
		}

		// flop
		if len(h.cards) != 3 {
			t.Errorf("unexpected number of cards, got %d", len(h.cards))
		}

		if err := playRaise(h, p2, 1); err != nil {
			t.Error(err)
		}
		if err := playRaise(h, p1, 2); err != nil {
			t.Error(err)
		}
		// outstanding action in flop
		if len(h.cards) != 3 {
			t.Errorf("unexpected number of cards, got %d", len(h.cards))
		}
		if err := playFold(h, p2); err != nil {
			t.Error(err)
		}
	}()
	v := <-done
	want := finishedHand{winner: p1, chips: smallBlind + 3}
	if v != want {
		t.Errorf("expected %v but got %v", want, v)
	}
}

func TestPlayerFoldsAfterRaise(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{p1, p2}

	done := make(chan finishedHand)
	go func() {
		h, _ := newHand(done, players, p1, smallBlind)

		// preflop
		if len(h.cards) != 0 {
			t.Errorf("unexpected number of cards, got %d", len(h.cards))
		}

		if err := playBlind(h, p1); err != nil {
			t.Error("Error should be nil")
		}

		// flop
		if len(h.cards) != 3 {
			t.Errorf("unexpected number of cards, got %d", len(h.cards))
		}

		if err := playRaise(h, p2, smallBlind+1); err != nil {
			t.Error(err)
		}
		if err := playFold(h, p1); err != nil {
			t.Error(err)
		}
	}()
	v := <-done
	want := finishedHand{winner: p2, chips: 3}
	if v != want {
		t.Errorf("expected %v but got %v", want, v)
	}
}

func TestRaiseByLessThanRequiredBetDueReturnsError(t *testing.T) {
	p1 := newPlayer(initial)
	p2 := newPlayer(initial)
	players := []*player{p1, p2}

	done := make(chan finishedHand)
	h, _ := newHand(done, players, p1)

	if err := playRaise(h, p1, 2); err != nil {
		t.Error(err)
	}
	if err := playRaise(h, p2, 1); err == nil {
		t.Error()
	}
}

func playBlind(h *hand, p *player) error {
	req := h.stage.requiredBet(h, p)
	return h.handleInput(p, input{action: Blind, chips: req})
}

func playFold(h *hand, p *player) error {
	return h.handleInput(p, input{action: Fold, chips: 0})
}

func playCheck(h *hand, p *player) error {
	return h.handleInput(p, input{action: Check, chips: 0})
}

func playCall(h *hand, p *player) error {
	req := h.stage.requiredBet(h, p)
	return h.handleInput(p, input{action: Call, chips: req})
}

func playRaise(h *hand, p *player, amount int) error {
	return h.handleInput(p, input{action: Raise, chips: amount})
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
