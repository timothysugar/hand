package hand

import (
	"testing"
)

func checkPlayers(t *testing.T, ps []player, rem ...player) {
	if (len(ps) != len(rem)) { t.Errorf("Player count should reduce by one; got: %v, want: %v", len(ps), len(rem))}
	for i, v := range(rem) {
		p := ps[i]
		if (p != v) { t.Errorf("Player %v should remain but players is %v", p, ps)}
	}
}

func TestPlayerOfManyFolds(t *testing.T) {
	p1 := newPlayer()
	p2 := newPlayer()
	p3 := newPlayer()
	players := []player{ p1, p2, p3 }
	h, _ := newHand(players)

	result, err := h.fold(p1)
	
	if (err != nil) { t.Error("Error should be nil")}

	checkPlayers(t, h.players, p2, p3)
	checkPlayers(t, result, p2, p3)
}

func TestPenultimatePlayerFolds(t *testing.T) {
	p1 := newPlayer()
	p2 := newPlayer()
	players := []player{ p1, p2 }
	h, _ := newHand(players)

	result, err := h.fold(p1)
	
	if (err != nil) { t.Error("Error should be nil" )}
	winner := h.winner()
	if (*winner != p2) { t.Errorf("Hand should have been won by %v but wasn't", p1)}
	checkPlayers(t, h.players, p2)
	checkPlayers(t, result, p2)
}

func TestFinalPlayerCannotFold(t *testing.T) {
	p1 := newPlayer()
	p2 := newPlayer()
	players := []player{ p1, p2 }
	h, _ := newHand(players)

	h.fold(p1)
	_, err := h.fold(p2)

	if (err == nil) { t.Errorf("Last player folding should return error but did not; err: %v", err)}
}

func TestHandRequiresTwoPlayers(t *testing.T) {
	players := []player{ newPlayer() }

	_, err := newHand(players)

	if (err == nil) { t.Error("New hand with less than one player should return error but did not")}
}
