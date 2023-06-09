package hand

import (
	"reflect"
	"sync"
	"testing"
)

const (
	initial    = 10
	bigBlind   = 2
	smallBlind = 1
)

func TestPenultimatePlayerFolds(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}
	h, _ := NewHand(players, p1)

	done := h.Begin()
	var v FinishedHand
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		v = <-done
		wg.Done()
	}()

	if err := playFold(h, p1); err != nil {
		t.Error("Error should be nil")
	}
	wg.Wait()

	checkPlayers(t, h.players, p2)
	want := FinishedHand{winner: p2, chips: 0}
	if v != want {
		t.Errorf("expected %v but got %v", want, v)
	}
	_, ok := <-done
	if ok {
		t.Error("expected done channel to be closed")
	}
}

func TestPenultimatePlayerFoldsFromBlind(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}
	h, _ := NewHand(players, p1, smallBlind)

	done := h.Begin()
	var v FinishedHand
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		v = <-done
		wg.Done()
	}()

	var err error
	err = playBlind(h, p1)
	if err != nil {
		t.Error(err)
	}

	err = playFold(h, p2)
	if err != nil {
		t.Error(err)
	}
	wg.Wait()

	want := FinishedHand{winner: p1, chips: 1}
	if v != want {
		t.Errorf("expected %v but got %v", want, v)
	}
}

func TestFinalPlayerCannotFold(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}
	h, _ := NewHand(players, p1)

	done := h.Begin()
	var v FinishedHand
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		v = <-done
		wg.Done()
	}()

	var err error
	err = playFold(h, p1)
	if err != nil {
		t.Error()
	}

	err = playFold(h, p2)
	if err == nil {
		t.Error("last player folding should return error")
	}
	wg.Wait()

	want := FinishedHand{winner: p2}
	if v != want {
		t.Errorf("expected %v but got %v", want, v)
	}
}

func TestBlindPlayerCannotFold(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}
	h, _ := NewHand(players, p1, 1)

	h.Begin()
	err := playFold(h, p1)

	if err == nil {
		t.Error("Player playing blind cannot fold")
	}
}

func TestBlindPlayerCannotCall(t *testing.T) {
	p1 := NewPlayer(1)
	p2 := NewPlayer(1)
	players := []*Player{p1, p2}

	h, _ := NewHand(players, p1, 1)

	err := playCall(h, p1)

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

	h, _ := NewHand(players, p1, smallBlind, bigBlind)
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
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	p3 := NewPlayer(initial)
	players := []*Player{p1, p2, p3}
	h, _ := NewHand(players, p1, smallBlind)
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
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	p3 := NewPlayer(initial)
	players := []*Player{p1, p2, p3}

	h, _ := NewHand(players, p1, smallBlind)
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
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}

	h, _ := NewHand(players, p1, smallBlind)

	err := playCheck(h, p1)
	if err == nil {
		t.Error()
	}
}

func TestCheckWhenBetDueReturnsError(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}

	var err error
	h, err := NewHand(players, p1, smallBlind)
	if err != nil {
		t.Error(err)
	}
	h.Begin()

	err = playBlind(h, p1)
	if err != nil {
		t.Error(err)
	}
	err = playCheck(h, p2)
	if err == nil {
		t.Error(err)
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

	err = playBlind(h, p2)
	if err != nil {
		t.Error(err)
	}
	if p2.chips != initial-smallBlind {
		t.Errorf("Player has unexpected number of chips %d", p2.chips)
	}
}

func TestBlindsPlayedInWrongOrderReturnsError(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}
	h, _ := NewHand(players, p1, smallBlind)

	err := playBlind(h, p2)

	if err == nil {
		t.Error("Expected an error for out of order blind but none received")
	}
}

func TestSamePlayerCallsImmediatelyAfterBlindReturnsError(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}
	h, _ := NewHand(players, p1, smallBlind)
	h.Begin()

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
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	p3 := NewPlayer(initial)
	players := []*Player{p1, p2, p3}
	h, _ := NewHand(players, p1, smallBlind)
	h.Begin()

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
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}

	h, _ := NewHand(players, p1, smallBlind)
	done := h.Begin()
	var wg sync.WaitGroup
	var v FinishedHand
	wg.Add(1)
	go func() {
		v = <-done
		wg.Done()
	}()

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
	wg.Wait()

	want := FinishedHand{winner: p1, chips: smallBlind * len(players)}
	if v != want {
		t.Errorf("expected %v but got %v", want, v)
	}
}

func TestPlayerFoldsAfterReraise(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}

	h, _ := NewHand(players, p1, smallBlind)
	done := h.Begin()
	var wg sync.WaitGroup
	var v FinishedHand
	wg.Add(1)
	go func() {
		v = <-done
		wg.Done()
	}()

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
	if err := playRaise(h, p1, 2); err != nil {
		t.Error(err)
	}
	// outstanding action in flop so should not advance stage
	if len(h.cards) != 3 {
		t.Errorf("unexpected number of cards, got %d", len(h.cards))
	}
	if err := playFold(h, p2); err != nil {
		t.Error(err)
	}
	wg.Wait()
	want := FinishedHand{winner: p1, chips: (2 * smallBlind) + 3}
	if v != want {
		t.Errorf("expected %v but got %v", want, v)
	}
}

func TestBeginPlayingReturnsChannelForResult(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}
	h, _ := NewHand(players, p1, smallBlind)

	var wg sync.WaitGroup
	done := h.Begin()
	wg.Add(1)
	var v FinishedHand
	go func() {
		v = <-done
		wg.Done()
	}()
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
	wg.Wait()
	want := FinishedHand{winner: p2, chips: 3}
	if v != want {
		t.Errorf("expected %v but got %v", want, v)
	}
}

func TestPlayerFoldsAfterRaise(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}
	h, _ := NewHand(players, p1, smallBlind)

	done := h.Begin()
	var wg sync.WaitGroup
	var v FinishedHand
	wg.Add(1)
	go func() {
		v = <-done
		wg.Done()
	}()

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
	wg.Wait()

	want := FinishedHand{winner: p2, chips: 3}
	if v != want {
		t.Errorf("expected %v but got %v", want, v)
	}
}

func TestRaiseByLessThanRequiredBetDueReturnsError(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}
	h, _ := NewHand(players, p1)
	h.Begin()

	if err := playRaise(h, p1, 2); err != nil {
		t.Error(err)
	}
	if err := playRaise(h, p2, 1); err == nil {
		t.Error()
	}
}

func TestRaiseAtSameValueAsRequiredBetDueReturnsError(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}
	h, _ := NewHand(players, p1)
	h.Begin()

	if err := playRaise(h, p1, 2); err != nil {
		t.Error(err)
	}
	if err := playRaise(h, p2, 2); err == nil {
		t.Error()
	}
}

func TestValidMovesInPreflopReturnsBlind(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}
	h, _ := NewHand(players, p1, smallBlind)
	h.Begin()

	got := h.ValidMoves()

	want := make(map[string][]Move)
	moves := []Move{NewMove(Blind, NewExactBet(smallBlind))}
	want[p1.id] = moves
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v but got %v", want, got)
	}
}

func TestValidMovesInFlopReturnsBettableMoves(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}
	h, _ := NewHand(players, p1)
	h.Begin()

	got := h.ValidMoves()

	want := make(map[string][]Move)
	moves := []Move{
		NewMove(Fold, RequiredBet{}),
		NewMove(Check, RequiredBet{}),
		NewMove(Raise, NewMinumumBet(0)),
	}
	want[p1.id] = moves
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v but got %v", want, got)
	}
}

func TestValidMovesInFlopWithOutstandingBetReturnsBettableMoves(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}
	h, _ := NewHand(players, p1)
	h.Begin()
	if err := playRaise(h, p1, 1); err != nil {
		t.Error(err)
	}

	got := h.ValidMoves()

	want := make(map[string][]Move)
	moves := []Move{
		NewMove(Fold, RequiredBet{}),
		NewMove(Call, NewExactBet(1)),
		NewMove(Raise, NewMinumumBet(1)),
	}
	want[p2.id] = moves
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v but got %v", want, got)
	}
}

func TestNoValidMovesWhenGameWon(t *testing.T) {
	p1 := NewPlayer(initial)
	p2 := NewPlayer(initial)
	players := []*Player{p1, p2}
	h, _ := NewHand(players, p1)
	done := h.Begin()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		<-done
		wg.Done()
	}()
	if err := playFold(h, p1); err != nil {
		t.Error(err)
	}

	got := h.ValidMoves()

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
