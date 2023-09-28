package hand

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

const (
	initial    = 10
	bigBlind   = 2
	smallBlind = 1
)

var source = rand.NewSource(time.Now().UnixNano())

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
	p1 := createPlayer()
	players := []*Player{p1}

	_, err := NewHand(players, p1, source)

	if err == nil {
		t.Error("New hand with less than one player should return error but did not")
	}
}

func TestPlayerNotInHandCannotPlayMove(t *testing.T) {
	th := createMinimalHand(t)

	p := createPlayer() // not in hand
	err := th.h.PlayBlind(p.Id, 1)

	if err == nil {
		t.Error("Player not in hand playing move should return error but did not")
	}
}

func TestPlayerCanJoinAHand(t *testing.T) {
	p1 := createPlayer()
	p2 := createPlayer()
	players := []*Player{p1, p2}
	h, err := NewHand(players, p1, source)
	if err != nil {
		t.Error(err)
	}

	p3 := createPlayer()

	if err := h.Join(p3, initial); err != nil {
		t.Error(err)
	}
}

func TestPlayerJoiningAHandWithADuplicateNameReturnsError(t *testing.T) {
	p1 := createPlayer()
	p2 := createPlayer()
	players := []*Player{p1, p2}
	h, err := NewHand(players, p1, source)
	if err != nil {
		t.Error(err)
	}

	p3 := NewPlayer("Twin", initial)
	p4 := NewPlayer("Twin", initial)

	if err := h.Join(p3, initial); err != nil {
		t.Error(err)
	}
	if err := h.Join(p4, initial); err == nil {
		t.Error(err)
	}
}

func TestAllPlayersCallTheBlind(t *testing.T) {
	p1 := createPlayer()
	p2 := createPlayer()
	p3 := createPlayer()
	players := []*Player{p1, p2, p3}

	var err error
	h, err := NewHand(players, p1, source, smallBlind, bigBlind)
	if err != nil {
		t.Error(err)
	}
	h.Begin()

	// preflop
	if p1.Chips != initial {
		t.Errorf("Player 1 should have %d chips before playing blind but has %d", initial, p1.Chips)
	}
	if p2.Chips != initial {
		t.Errorf("Player 2 should have %d chips before playing blind but has %d", initial, p2.Chips)
	}
	if p3.Chips != initial {
		t.Errorf("Player 3 should have %d chips before playing blind but has %d", initial, p3.Chips)
	}
	if err = playBlind(h, p1); err != nil {
		t.Error(err)
	}
	if p1.Chips != (initial - smallBlind) {
		t.Errorf("Player 1 should have %d chips after playing blind but has %d", initial-smallBlind, p1.Chips)
	}
	err = playBlind(h, p2)
	if err != nil {
		t.Error(err)
	}
	if p2.Chips != (initial - bigBlind) {
		t.Errorf("Player 2 should have %d chips after playing blind but has %d", initial-bigBlind, p2.Chips)
	}

	// flop
	err = playCall(h, p3)
	if err != nil {
		t.Error(err)
	}
	if p3.Chips != (initial - bigBlind) {
		t.Errorf("Player 3 should have %d chips after calling but has %d", initial-bigBlind, p3.Chips)
	}
	err = playCall(h, p1)
	if err != nil {
		t.Error(err)
	}
	if p1.Chips != (initial - bigBlind) {
		t.Errorf("Player 1 should have %d chips after calling but has %d", initial-bigBlind, p1.Chips)
	}
	err = playCheck(h, p2)
	if err != nil {
		t.Error(err)
	}
	if p2.Chips != (initial - bigBlind) {
		t.Errorf("Player 2 should have %d chips after calling but has %d", initial-bigBlind, p2.Chips)
	}
}

func TestOneFoldsAndOneCallsBlind(t *testing.T) {
	p1 := createPlayer()
	p2 := createPlayer()
	p3 := createPlayer()
	players := []*Player{p1, p2, p3}
	var err error
	h, err := NewHand(players, p1, source, smallBlind)
	if err != nil {
		t.Error(err)
	}
	h.Begin()

	if p1.Chips != initial {
		t.Errorf("Player 1 should have %d chips before playing blind but has %d", initial, p1.Chips)
	}
	if p2.Chips != initial {
		t.Errorf("Player 2 should have %d chips before playing blind but has %d", initial, p2.Chips)
	}
	if p3.Chips != initial {
		t.Errorf("Player 3 should have %d chips before playing blind but has %d", initial, p3.Chips)
	}
	if err = playBlind(h, p1); err != nil {
		t.Error(err)
	}
	if p1.Chips != (initial - smallBlind) {
		t.Errorf("Player 1 should have %d chips after playing blind but has %d", initial-smallBlind, p1.Chips)
	}
	if err = playFold(h, p2); err != nil {
		t.Error(err)
	}
	if err = playCall(h, p3); err != nil {
		t.Error(err)
	}
	if p3.Chips != (initial - smallBlind) {
		t.Errorf("Player 3 should have %d chips after calling but has %d", initial-smallBlind, p3.Chips)
	}
	if err = playCall(h, p1); err != nil {
		t.Error(err)
	}
	if p1.Chips != (initial - smallBlind) {
		t.Errorf("Player 1 should have %d chips after calling but has %d", initial-smallBlind, p1.Chips)
	}
}

func TestOneFoldsAndOneChecksBlind(t *testing.T) {
	p1 := createPlayer()
	p2 := createPlayer()
	p3 := createPlayer()
	players := []*Player{p1, p2, p3}

	var err error
	h, err := NewHand(players, p1, source, smallBlind)
	if err != nil {
		t.Error(err)
	}
	h.Begin()

	if p1.Chips != initial {
		t.Errorf("Player 1 should have %d chips before playing blind but has %d", initial, p1.Chips)
	}
	if p2.Chips != initial {
		t.Errorf("Player 2 should have %d chips before playing blind but has %d", initial, p2.Chips)
	}
	if p3.Chips != initial {
		t.Errorf("Player 3 should have %d chips before playing blind but has %d", initial, p3.Chips)
	}
	if err = playBlind(h, p1); err != nil {
		t.Error(err)
	}
	if p1.Chips != (initial - smallBlind) {
		t.Errorf("Player 1 should have %d chips after playing blind but has %d", initial-smallBlind, p1.Chips)
	}
	if err = playFold(h, p2); err != nil {
		t.Error(err)
	}
	if err = playCall(h, p3); err != nil {
		t.Error(err)
	}
	if p3.Chips != (initial - smallBlind) {
		t.Errorf("Player 3 should have %d chips after calling but has %d", initial-smallBlind, p3.Chips)
	}
	if err = playCheck(h, p1); err != nil {
		t.Error(err)
	}
	if p1.Chips != (initial - smallBlind) {
		t.Errorf("Player 1 should have %d chips after checking but has %d", initial-smallBlind, p1.Chips)
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
	p1 := createPlayer()
	p2 := createPlayer()
	players := []*Player{p1, p2} // p2 is not the first in order
	h, err := NewHand(players, p2, source, smallBlind)
	if err != nil {
		t.Error(err)
	}
	h.Begin()

	if err = playBlind(h, p2); err != nil {
		t.Error(err)
	}
	if p2.Chips != initial-smallBlind {
		t.Errorf("Player has unexpected number of chips %d", p2.Chips)
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
	p1 := createPlayer()
	p2 := createPlayer()
	p3 := createPlayer()
	players := []*Player{p1, p2, p3}
	h, _ := NewHand(players, p1, source, smallBlind)
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
	numCards = len(th.h.Cards)
	if numCards != 0 {
		t.Errorf("unexpected number of cards, got %d", numCards)
	}
	if err := playBlind(th.h, th.p1); err != nil {
		t.Error(err)
	}

	// flop
	numCards = len(th.h.Cards)
	if numCards != 3 {
		t.Errorf("unexpected number of cards, got %d", numCards)
	}
	if err := playCall(th.h, th.p2); err != nil {
		t.Error(err)
	}
	// outstanding action in flop so should not advance stage
	numCards = len(th.h.Cards)
	if numCards != 3 {
		t.Errorf("unexpected number of cards, got %d", numCards)
	}
	if err := playCheck(th.h, th.p1); err != nil {
		t.Error(err)
	}

	// turn
	numCards = len(th.h.Cards)
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
	numCards = len(th.h.Cards)
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
	want[th.p1.Id] = moves
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
	want[th.p1.Id] = moves
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
	want[th.p2.Id] = moves
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
	p1 := createPlayer()
	p2 := createPlayer()
	players := []*Player{p1, p2}
	h, _ := NewHand(players, p1, source)

	got := h.ValidMoves()

	want := make(map[string][]Move)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v but got %v", want, got)
	}
}

func TestCallingBeginTwiceReturnsError(t *testing.T) {
	th := createMinimalHand(t)

	if _, err := th.h.Begin(); err == nil {
		t.Error("Expected error for calling begin twice but none received")
	}
}

func TestGetPlayersReturnsPlayerAndOpponents(t *testing.T) {
	th := createMinimalHand(t)

	self, opponents := th.h.Players(th.p1.Id)
	if self != th.p1 {
		t.Errorf("expected %v but got %v", th.p1, self)
	}
	if len(opponents) != 1 && opponents[0] != th.p2 {
		t.Errorf("expected %v but got %v", th.p2, opponents)
	}
}

func TestIsNextToPlayReturnsFalseBeforeGameBegins(t *testing.T) {
	p1 := createPlayer()
	p2 := createPlayer()
	players := []*Player{p1, p2}
	h, err := NewHand(players, p1, source)
	if err != nil {
		t.Error(err)
	}

	var p1Next = h.IsNextToPlay(p1.Id)
	var p2Next = h.IsNextToPlay(p2.Id)

	if !(!p1Next && !p2Next) {
		t.Error("expected no one is next to play before hand begins")
	}
}

func TestIsNextToPlayIteratesFromDealerWhenGameBegins(t *testing.T) {
	th := createMinimalHand(t)

	var p1Next = th.h.IsNextToPlay(th.p1.Id)
	var p2Next = th.h.IsNextToPlay(th.p2.Id)

	if !(p1Next && !p2Next) {
		t.Error("expected play to start from dealer but did not")
	}

	playCheck(th.h, th.p2) // out of order move should not change next to play
	p1Next = th.h.IsNextToPlay(th.p1.Id)
	p2Next = th.h.IsNextToPlay(th.p2.Id)

	if !(p1Next && !p2Next) {
		t.Error("expected play invalid move not to increment next to play")
	}

	playCheck(th.h, th.p1)
	p1Next = th.h.IsNextToPlay(th.p1.Id)
	p2Next = th.h.IsNextToPlay(th.p2.Id)

	if !(p2Next && !p1Next) {
		t.Error("expected play to increment from dealer for each move played")
	}
}

func TestIsActiveReturnsFalseBeforeGameBegins(t *testing.T) {
	p1 := createPlayer()
	p2 := createPlayer()
	players := []*Player{p1, p2}
	h, err := NewHand(players, p1, source)
	if err != nil {
		t.Error(err)
	}

	if h.IsActive() {
		t.Error("Expected hand to be inactive before game begins")
	}
}

func TestIsActiveReturnsTrueAfterGameBegins(t *testing.T) {
	th := createMinimalHand(t)

	if !th.h.IsActive() {
		t.Error("Expected hand to be active after game begins")
	}
}

func TestIsActiveReturnsFalseAfterGameEnds(t *testing.T) {
	th := createMinimalHand(t)

	if err := playFold(th.h, th.p1); err != nil {
		t.Error("Error should be nil")
	}

	if !th.h.IsActive() {
		t.Error("Expected hand to be active after game ends")
	}
}

type testHand struct {
	h   *Hand
	p1  *Player
	p2  *Player
	fin chan FinishedHand
}

func createMinimalHand(t *testing.T) testHand {
	p1 := createPlayer()
	p2 := createPlayer()
	players := []*Player{p1, p2}
	h, err := NewHand(players, p1, source)
	if err != nil {
		t.Error(err)
	}
	fin, err := h.Begin()
	if err != nil {
		t.Error(err)
	}

	return testHand{h, p1, p2, fin}
}

func createMinimalHandWithBlind(t *testing.T) testHand {
	p1 := createPlayer()
	p2 := createPlayer()
	players := []*Player{p1, p2}
	h, err := NewHand(players, p1, source, smallBlind)
	if err != nil {
		t.Error(err)
	}
	fin, err := h.Begin()
	if err != nil {
		t.Error(err)
	}

	return testHand{h, p1, p2, fin}
}

func createPlayer() *Player {
	name := randomString(10)
	return NewPlayer(name, initial)
}

func randomString(length int) string {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    b := make([]byte, length+2)
	r.Read(b)
    return fmt.Sprintf("%x", b)[2 : length+2]
}

func playBlind(h *Hand, p *Player) error {
	req := h.stage.requiredBet(h, p)
	return h.PlayBlind(p.Id, req)
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
