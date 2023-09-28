package hand

import "math/rand"



type Suit int

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

type Rank int

const (
	Two Rank = iota
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)


type Card struct {
	Suit Suit
	Rank Rank
}
type Deck struct {
	cards []Card
}

func newDeck(source rand.Source) Deck {
	cards := make([]Card, 0)
	for i := 0; i < 4; i++ {
		for j := 0; j < 13; j++ {
			cards = append(cards, Card{Suit(i), Rank(j)})
		}
	}
	r := rand.New(source)
	r.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	return Deck{cards}
}


func (d *Deck) pop() Card {
	var c Card
	c, d.cards = d.cards[0], d.cards[1:]
	return c
}

