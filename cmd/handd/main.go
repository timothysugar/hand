package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	ts "github.com/timothysugar/hand/cmd/handd/templates"
	"github.com/timothysugar/hand/pkg/hand"
)

type PlayerViewModel struct {
	*hand.Player
	Moves  []hand.Move
	HandId string
}

const (
	initialChips = 1000
)

var h *hand.Hand
var me *hand.Player
var baseTmpl *template.Template
var templates *ts.Template

func init() {
	log.Println("Initializing hand")
	templates = ts.New()
	baseTmpl = template.Must(template.ParseFiles("web/base.go.html"))

	bill := hand.NewPlayer("Bill", initialChips)
	bill.Cards = []hand.Card{
		{Suit: "Spades", Rank: "Ace"},
		{Suit: "Clubs", Rank: "Ace"},
	}
	ben := hand.NewPlayer("Ben", initialChips)
	ben.Cards = []hand.Card{
		{Suit: "Hearts", Rank: "Ace"},
		{Suit: "Diamonds", Rank: "Ace"},
	}
	me = hand.NewPlayer("Zephyr", initialChips)
	me.Cards = []hand.Card{
		{Suit: "Hearts", Rank: "10"},
		{Suit: "Diamonds", Rank: "9"},
	}
	players := []*hand.Player{
		bill,
		ben,
		me,
	}
	var err error
	h, err = hand.NewHand(players, players[2], 10)
	h.Begin()
	if err != nil {
		log.Fatalf("Error initializing hand: %s", err)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/table/{tableId}/hand/{handId}", getHandHandler).Name("get-hand").Methods("GET")
	r.HandleFunc("/hand/{handId}/player/{playerId}/moves/Blind/{amount}", blindHandler).Name("play-blind").Methods("POST")
	r.HandleFunc("/", getHandsHandler).Name("get-hands").Methods("GET")
	r.HandleFunc("/table/{tableId}", getHandHandler).Name("get-game").Methods("GET")

	port := ":8090"
	fmt.Printf("listening on %s", port)
	http.ListenAndServe(port, r)
}

func createPlayerViewModel(self *hand.Player, tableId string, handId string) templates.PlayerViewModel {
	allMoves := h.ValidMoves()
	mvs := allMoves[self.Id]
	return templates.PlayerViewModel{
		Id:     self.Id,
		TableId: tableId,
		HandId: handId,
		Entrant: templates.Entrant{
			Name:   me.Name,
			Chips: me.Chips,
			Folded: me.Folded,
		},
		Cards:  self.Cards,
		Moves:  mvs,
	}
}

func createHandViewModel(playerId string, tableId string, handId string) (templates.HandViewModel, error) {
	self, opponents, err := h.Players(playerId)
	if err != nil {
		return templates.HandViewModel{}, err
	}
	opponentsVM := make([]templates.OpponentViewModel, len(opponents))
	for i, o := range opponents {
		opponentsVM[i] = templates.OpponentViewModel{
			Entrant: templates.Entrant{
				Name:   o.Name,
				Chips:  o.Chips,
				Folded: o.Folded,
				Active: h.IsNextToPlay(o.Id),
			},
			FaceDownCards: make([]struct{}, len(o.Cards)),
		}
	}
	playerVM := createPlayerViewModel(self, tableId, handId)

	return templates.HandViewModel{
		HandId:    h.Id,
		TableId:   tableId,
		Opponents: opponentsVM,
		Player:    playerVM,
	}, nil
}

func getHandsHandler(w http.ResponseWriter, req *http.Request) {
}

func getHandHandler(w http.ResponseWriter, req *http.Request) {
	pathVars := mux.Vars(req)
	tableId := pathVars["tableId"]
	handId := pathVars["handId"]

	vm, err := createHandViewModel(me.Id, tableId, handId)
	if err != nil {
		log.Printf("Error creating hand view model for request with URL:%s, %v\n", req.URL, err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	var tmplFiles = []string{
		"web/base.go.html",
		"web/hand.go.html",
		"web/player.go.html",
		"web/opponentPositions.go.html",
	}

	contentTmpl, err := template.New("main").ParseFiles(tmplFiles...)
	if err != nil {
		log.Printf("Error parsing template files %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = contentTmpl.Execute(w, vm)
	if err != nil {
		log.Printf("Error executing template with data, err: %v, template: %v, data: %v", err, contentTmpl, vm)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func blindHandler(w http.ResponseWriter, req *http.Request) {
	pathVars := mux.Vars(req)
	tableId := pathVars["tableId"]
	handId := pathVars["handId"]
	playerId := pathVars["playerId"]
	amount := pathVars["amount"]

	v, err := strconv.Atoi(amount)
	if err != nil {
		log.Printf("Error converting amount to int: %s\n", amount)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.PlayBlind(playerId, v); err != nil {
		log.Printf("Error playing blind: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vm := createPlayerViewModel(me, tableId, handId)
	if err != nil {
		log.Printf("Error creating hand view model for request with URL:%s, %v\n", req.URL, err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	var tmplFiles = []string{
		"web/player.go.html",
	}
	tmpl, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		log.Printf("Error parsing template files %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, vm)
	if err != nil {
		log.Printf("Error executing template with data, err: %v, template: %v, data: %v", err, tmpl, vm)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
