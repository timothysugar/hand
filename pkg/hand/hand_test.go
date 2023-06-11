package hand

import (
	"reflect"
	"testing"
)

const (
	initial    = 10
	bigBlind   = 2
	smallBlind = 1
)

func TestPenultimatePlayerFolds(t *testing.T) {
	th := createMinimalHand(t)

	if err := playFold(th.h, th.p1); err != nil {
		t.Error("Error should be nil")
	}

	checkPlayers(t, th.h.players, th.p2)
	want := FinishedHand{winner: th.p2, chips: 0}
	fin := <-th.fin
	if fin != want {
		t.Errorf("expected %v but got %v", want, fin)
	}
	_, ok := <-th.fin
	if ok {
		t.Error("expected done channel to be closed")
	}
}

func TestPenultimatePlayerFoldsFromBlind(t *testing.T) {
	th := createMinimalHandWithBlind(t)
	if err := playBlind(th.h, th.p1); err != nil {
		t.Error(err)
	}
	if err := playFold(th.h, th.p2); err != nil {
		t.Error(err)
	}

	want := FinishedHand{winner: th.p1, chips: 1}
	fin := <-th.fin
	if fin != want {
		t.Errorf("expected %v but got %v", want, fin)
	}
}

func TestFinalPlayerCannotFold(t *testing.T) {
	th := createMinimalHand(t)
	if err := playFold(th.h, th.p1); err != nil {
		t.Error(err)
	}
	if err := playFold(th.h, th.p2); err == nil {
		t.Error("last player folding should return error")
	}

	want := FinishedHand{winner: th.p2}
	fin := <-th.fin
	if fin != want {
		t.Errorf("expected %v but got %v", want, fin)
	}
}

func TestBlindPlayerCannotFold(t *testing.T) {
	th := createMinimalHandWithBlind(t)

	err := playFold(th.h, th.p1)

	if err == nil {
		t.Error("Player playing blind cannot fold")
	}
}

func TestBlindPlayerCannotCall(t *testing.T) {
	th := createMinimalHandWithBlind(t)

	err := playCall(th.h, th.p1)

	if err == nil {
		t.Error("Player playing blind cannot call")
	}
}

func TestHandRequiresTwoPlayers(t *testing.T) {
	p1 := NewPlayer(initial)
	players := []*Player{p1}

	_, err := NewHand(players, p1)

	if err == nil {
		t.Error("New hand with less than one player should return error but did not")
	}
}

func TestAllPlayersCallTheBlind(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	p3 := NewPlayer(initial)
	players := []*Player{p1, p2, p3}

	var err error
	h, err := NewHand(players, p1, smallBlind, bigBlind)
	if err != nil {
		t.Error(err)
	}
	h.Begin()

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
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	p3 := NewPlayer(initial)
	players := []*Player{p1, p2, p3}
	var err error
	h, err := NewHand(players, p1, smallBlind)
	if err != nil {
		t.Error(err)
	}
	h.Begin()

	if p1.chips != initial {
		t.Errorf("Player 1 should have %d chips before playing blind but has %d", initial, p1.chips)
	}
	if p2.chips != initial {
		t.Errorf("Player 2 should have %d chips before playing blind but has %d", initial, p2.chips)
	}
	if p3.chips != initial {
		t.Errorf("Player 3 should have %d chips before playing blind but has %d", initial, p3.chips)
	}
	if err = playBlind(h, p1); err != nil {
		t.Error(err)
	}
	if p1.chips != (initial - smallBlind) {
		t.Errorf("Player 1 should have %d chips after playing blind but has %d", initial-smallBlind, p1.chips)
	}
	if err = playFold(h, p2); err != nil {
		t.Error(err)
	}
	if err = playCall(h, p3); err != nil {
		t.Error(err)
	}
	if p3.chips != (initial - smallBlind) {
		t.Errorf("Player 3 should have %d chips after calling but has %d", initial-smallBlind, p3.chips)
	}
	if err = playCall(h, p1); err != nil {
		t.Error(err)
	}
	if p1.chips != (initial - smallBlind) {
		t.Errorf("Player 1 should have %d chips after calling but has %d", initial-smallBlind, p1.chips)
	}
}

func TestOneFoldsAndOneChecksBlind(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	p3 := NewPlayer(initial)
	players := []*Player{p1, p2, p3}

	var err error
	h, err := NewHand(players, p1, smallBlind)
	if err != nil {
		t.Error(err)
	}
	h.Begin()

	if p1.chips != initial {
		t.Errorf("Player 1 should have %d chips before playing blind but has %d", initial, p1.chips)
	}
	if p2.chips != initial {
		t.Errorf("Player 2 should have %d chips before playing blind but has %d", initial, p2.chips)
	}
	if p3.chips != initial {
		t.Errorf("Player 3 should have %d chips before playing blind but has %d", initial, p3.chips)
	}
	if err = playBlind(h, p1); err != nil {
		t.Error(err)
	}
	if p1.chips != (initial - smallBlind) {
		t.Errorf("Player 1 should have %d chips after playing blind but has %d", initial-smallBlind, p1.chips)
	}
	if err = playFold(h, p2); err != nil {
		t.Error(err)
	}
	if err = playCall(h, p3); err != nil {
		t.Error(err)
	}
	if p3.chips != (initial - smallBlind) {
		t.Errorf("Player 3 should have %d chips after calling but has %d", initial-smallBlind, p3.chips)
	}
	if err = playCheck(h, p1); err != nil {
		t.Error(err)
	}
	if p1.chips != (initial - smallBlind) {
		t.Errorf("Player 1 should have %d chips after checking but has %d", initial-smallBlind, p1.chips)
	}
}

func TestCheckWhenBlindDueReturnsError(t *testing.T) {
	th := createMinimalHandWithBlind(t)

	if err := playCheck(th.h, th.p1); err == nil {
		t.Error("Expected an error for out of order check but none received")
	}
}

func TestCheckWhenBetDueReturnsError(t *testing.T) {
	th := createMinimalHandWithBlind(t)

	if err := playBlind(th.h, th.p1); err != nil {
		t.Error(err)
	}
	if err := playCheck(th.h, th.p2); err == nil {
		t.Error("Expected an error for out of order check but none received")
	}
}

func TestBlindsPlayedFromDealer(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2} // p2 is not the first in order
	h, err := NewHand(players, p2, smallBlind)
	if err != nil {
		t.Error(err)
	}
	h.Begin()

	if err = playBlind(h, p2); err != nil {
		t.Error(err)
	}
	if p2.chips != initial-smallBlind {
		t.Errorf("Player has unexpected number of chips %d", p2.chips)
	}
}

func TestBlindsPlayedInWrongOrderReturnsError(t *testing.T) {
	th := createMinimalHandWithBlind(t)

	if err := playBlind(th.h, th.p2); err == nil {
		t.Error("Expected an error for out of order blind but none received")
	}
}

func TestSamePlayerCallsImmediatelyAfterBlindReturnsError(t *testing.T) {
	th := createMinimalHandWithBlind(t)

	if err := playBlind(th.h, th.p1); err != nil {
		t.Error(err)
	}

	if err := playCall(th.h, th.p1); err == nil {
		t.Error("Expected an error for out of order call but none received")
	}
}

func TestSecondCallInWrongOrderReturnsError(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	p3 := NewPlayer(initial)
	players := []*Player{p1, p2, p3}
	h, _ := NewHand(players, p1, smallBlind)
	h.Begin()

	if err := playBlind(h, p1); err != nil {
		t.Error(err)
	}

	if err := playCall(h, p2); err != nil {
		t.Error(err)
	}

	if err := playCall(h, p1); err == nil {
		t.Error("Expected an error for out of order call but none received")
	}
}

func TestProgressingThroughStagesIncrementsNumOfCardsInHand(t *testing.T) {
	th := createMinimalHandWithBlind(t)

	// preflop
	var numCards int
	numCards = len(th.h.cards)
	if numCards != 0 {
		t.Errorf("unexpected number of cards, got %d", numCards)
	}
	if err := playBlind(th.h, th.p1); err != nil {
		t.Error(err)
	}

	// flop
	numCards = len(th.h.cards)
	if numCards != 3 {
		t.Errorf("unexpected number of cards, got %d", numCards)
	}
	if err := playCall(th.h, th.p2); err != nil {
		t.Error(err)
	}
	// outstanding action in flop so should not advance stage
	numCards = len(th.h.cards)
	if numCards != 3 {
		t.Errorf("unexpected number of cards, got %d", numCards)
	}
	if err := playCheck(th.h, th.p1); err != nil {
		t.Error(err)
	}

	// turn
	numCards = len(th.h.cards)
	if numCards != 4 {
		t.Errorf("unexpected number of cards, got %d", numCards)
	}
	if err := playCheck(th.h, th.p1); err != nil {
		t.Error(err)
	}
	if err := playCheck(th.h, th.p2); err != nil {
		t.Error(err)
	}

	// river
	numCards = len(th.h.cards)
	if numCards != 5 {
		t.Errorf("unexpected number of cards, got %d", numCards)
	}
	if err := playCheck(th.h, th.p1); err != nil {
		t.Error(err)
	}
	if err := playCheck(th.h, th.p2); err != nil {
		t.Error(err)
	}
	// players hands evaluated

	want := FinishedHand{winner: th.p1, chips: smallBlind * len(th.h.players)}
	v := <-th.fin
	if v != want {
		t.Errorf("expected %v but got %v", want, v)
	}
}

func TestPlayerFoldsAfterRaise(t *testing.T) {
	th := createMinimalHandWithBlind(t)

	// preflop
	if err := playBlind(th.h, th.p1); err != nil {
		t.Error(err)
	}

	// flop
	if err := playRaise(th.h, th.p2, smallBlind+1); err != nil {
		t.Error(err)
	}
	if err := playFold(th.h, th.p1); err != nil {
		t.Error(err)
	}

	want := FinishedHand{winner: th.p2, chips: 3}
	v := <-th.fin
	if v != want {
		t.Errorf("expected %v but got %v", want, v)
	}
}

func TestPlayerFoldsAfterReraise(t *testing.T) {
	th := createMinimalHandWithBlind(t)

	// preflop
	if err := playBlind(th.h, th.p1); err != nil {
		t.Error(err)
	}

	// flop
	if err := playRaise(th.h, th.p2, smallBlind+1); err != nil {
		t.Error(err)
	}
	if err := playRaise(th.h, th.p1, 2); err != nil {
		t.Error(err)
	}
	if err := playFold(th.h, th.p2); err != nil {
		t.Error(err)
	}
	want := FinishedHand{winner: th.p1, chips: (2 * smallBlind) + 3}
	v := <-th.fin
	if v != want {
		t.Errorf("expected %v but got %v", want, v)
	}
}

func TestRaiseByLessThanRequiredBetDueReturnsError(t *testing.T) {
	th := createMinimalHand(t)

	if err := playRaise(th.h, th.p1, 2); err != nil {
		t.Error(err)
	}
	if err := playRaise(th.h, th.p2, 1); err == nil {
		t.Error("Expected error for raise less than required bet but none received")
	}
}

func TestRaiseAtSameValueAsRequiredBetDueReturnsError(t *testing.T) {
	th := createMinimalHand(t)

	if err := playRaise(th.h, th.p1, 2); err != nil {
		t.Error(err)
	}
	if err := playRaise(th.h, th.p2, 2); err == nil {
		t.Error()
	}
}

func TestValidMovesInPreflopReturnsBlind(t *testing.T) {
	th := createMinimalHandWithBlind(t)

	got := th.h.ValidMoves()

	want := make(map[string][]Move)
	moves := []Move{NewMove(Blind, NewExactBet(smallBlind))}
	want[th.p1.id] = moves
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v but got %v", want, got)
	}
}

func TestValidMovesInFlopReturnsBettableMoves(t *testing.T) {
	th := createMinimalHand(t)

	got := th.h.ValidMoves()

	want := make(map[string][]Move)
	moves := []Move{
		NewMove(Fold, RequiredBet{}),
		NewMove(Check, RequiredBet{}),
		NewMove(Raise, NewMinumumBet(0)),
	}
	want[th.p1.id] = moves
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v but got %v", want, got)
	}
}

func TestValidMovesInFlopWithOutstandingBetReturnsBettableMoves(t *testing.T) {
	th := createMinimalHand(t)
	if err := playRaise(th.h, th.p1, 1); err != nil {
		t.Error(err)
	}

	got := th.h.ValidMoves()

	want := make(map[string][]Move)
	moves := []Move{
		NewMove(Fold, RequiredBet{}),
		NewMove(Call, NewExactBet(1)),
		NewMove(Raise, NewMinumumBet(1)),
	}
	want[th.p2.id] = moves
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v but got %v", want, got)
	}
}

func TestNoValidMovesWhenGameWon(t *testing.T) {
	th := createMinimalHand(t)
	if err := playFold(th.h, th.p1); err != nil {
		t.Error(err)
	}

	got := th.h.ValidMoves()

	want := make(map[string][]Move)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v but got %v", want, got)
	}
}

func TestNoValidMovesBeforeGameBegins(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}
	h, _ := NewHand(players, p1)

	got := h.ValidMoves()

	want := make(map[string][]Move)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v but got %v", want, got)
	}
}

type testHand struct {
    h    *Hand
    p1   *Player
    p2  *Player
	fin chan FinishedHand
}

func createMinimalHand(t *testing.T) testHand {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}
	h, err := NewHand(players, p1)
	if err != nil {
		t.Error(err)
	}
	fin := h.Begin()

	return testHand{h, p1, p2, fin}
}

func createMinimalHandWithBlind(t *testing.T) testHand {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}
	h, err := NewHand(players, p1, smallBlind)
	if err != nil {
		t.Error(err)
	}
	fin := h.Begin()

	return testHand{h, p1, p2, fin}
}

func playBlind(h *Hand, p *Player) error {
	req := h.stage.requiredBet(h, p)
	return h.HandleInput(p, Input{Action: Blind, Chips: req})
}

func playFold(h *Hand, p *Player) error {
	return h.HandleInput(p, Input{Action: Fold, Chips: 0})
}

func playCheck(h *Hand, p *Player) error {
	return h.HandleInput(p, Input{Action: Check, Chips: 0})
}

func playCall(h *Hand, p *Player) error {
	req := h.stage.requiredBet(h, p)
	return h.HandleInput(p, Input{Action: Call, Chips: req})
}

func playRaise(h *Hand, p *Player, amount int) error {
	return h.HandleInput(p, Input{Action: Raise, Chips: amount})
}

func checkPlayers(t *testing.T, ps []*Player, rem ...*Player) {
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
