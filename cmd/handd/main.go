package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-faker/faker/v4"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
	"github.com/timothysugar/hand/cmd/handd/templates"
	"github.com/timothysugar/hand/pkg/hand"
)

type PlayerViewModel struct {
	*hand.Player
	Moves  []hand.Move
	HandId string
}

const (
	initialChips      = 1000
	initialTableCount = 3
	tableId           = "table1"
	playerLimit       = 8
)

type GameState int

const (
	Lobby GameState = iota
	ActiveHand
)

type Table struct {
	Id      string
	Name    string
	players Players
	Status  GameState
}

func NewTable(name string) Table {
	return Table{
		Id:      xid.New().String(),
		Name:    name,
		players: NewPlayers(playerLimit),
		Status:  Lobby,
	}
}

var tables []Table = make([]Table, 0)
var h *hand.Hand
var me *hand.Player
var ts *templates.Template

func init() {
	log.Println("Initializing hand")
	ts = templates.New()
	for i := 0; i < initialTableCount; i++ {
		tables = append(tables, NewTable(fmt.Sprintf("table%d", i)))
	}

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
		{Suit: "Diamonds", Rank: "Queen"},
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

func serve() error {
	r := mux.NewRouter()
	assetsPath := "cmd/handd/static"
	listed, _ := os.ReadDir(assetsPath)
	log.Printf("listing dir %s: %v", assetsPath, listed)
	assets := http.Dir(assetsPath)
	log.Printf("found assets %v", assets)
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets", http.FileServer(assets)))
	r.HandleFunc("/table/{tableId}/hand/{handId}", getHandHandler).Name("get-hand").Methods("GET")
	r.HandleFunc("/hand/{handId}/player/{playerId}/moves/Blind/{amount}", blindHandler).Name("play-blind").Methods("POST")
	r.HandleFunc("/", getTablesHandler).Name("get-hands").Methods("GET")
	r.HandleFunc("/table", newTableHandler).Name(("new-table")).Methods("POST")
	r.HandleFunc("/table/{tableId}", getHandHandler).Name("get-game").Methods("GET")

	port := ":8070"
	fmt.Printf("listening at http://localhost%s\n", port)
	return http.ListenAndServe(port, r)
}

func main() {
	renderTemplateCmd := flag.NewFlagSet("tmpl", flag.ExitOnError)
	renderTemplateCmd.Usage = func() { fmt.Printf("render a named template to an HTML output file\n") }
	renderTemplateCmd.String("name", "", "template name to render")

	flag.Parse()
	flag.PrintDefaults()
	if len(os.Args) < 2 {
		log.Println("running game server")
		if err := serve(); err != nil {
			log.Fatalf("unexpected error running game serving so exiting: %s", err)
		}
		os.Exit(0)
	}
	switch os.Args[1] {
	case "tmpl":
		renderTemplateCmd.Parse(os.Args[2:])
		log.Println("rendering template")
		os.Exit(0)
	default:
		log.Fatalf("unexpected subcommand in first argument")
	}
}

var classSuffixLookup = map[hand.Card]string{
	{Rank: "Ace", Suit: "Clubs"}:      "ac",
	{Rank: "2", Suit: "Clubs"}:        "2c",
	{Rank: "3", Suit: "Clubs"}:        "3c",
	{Rank: "4", Suit: "Clubs"}:        "4c",
	{Rank: "5", Suit: "Clubs"}:        "5c",
	{Rank: "6", Suit: "Clubs"}:        "6c",
	{Rank: "7", Suit: "Clubs"}:        "7c",
	{Rank: "8", Suit: "Clubs"}:        "8c",
	{Rank: "9", Suit: "Clubs"}:        "9c",
	{Rank: "10", Suit: "Clubs"}:       "10c",
	{Rank: "Jack", Suit: "Clubs"}:     "jc",
	{Rank: "Queen", Suit: "Clubs"}:    "qc",
	{Rank: "King", Suit: "Clubs"}:     "kc",
	{Rank: "Ace", Suit: "Diamonds"}:   "ad",
	{Rank: "2", Suit: "Diamonds"}:     "2d",
	{Rank: "3", Suit: "Diamonds"}:     "3d",
	{Rank: "4", Suit: "Diamonds"}:     "4d",
	{Rank: "5", Suit: "Diamonds"}:     "5d",
	{Rank: "6", Suit: "Diamonds"}:     "6d",
	{Rank: "7", Suit: "Diamonds"}:     "7d",
	{Rank: "8", Suit: "Diamonds"}:     "8d",
	{Rank: "9", Suit: "Diamonds"}:     "9d",
	{Rank: "10", Suit: "Diamonds"}:    "10d",
	{Rank: "Jack", Suit: "Diamonds"}:  "jd",
	{Rank: "Queen", Suit: "Diamonds"}: "qd",
	{Rank: "King", Suit: "Diamonds"}:  "kd",
	{Rank: "Ace", Suit: "Hearts"}:     "ah",
	{Rank: "2", Suit: "Hearts"}:       "2h",
	{Rank: "3", Suit: "Hearts"}:       "3h",
	{Rank: "4", Suit: "Hearts"}:       "4h",
	{Rank: "5", Suit: "Hearts"}:       "5h",
	{Rank: "6", Suit: "Hearts"}:       "6h",
	{Rank: "7", Suit: "Hearts"}:       "7h",
	{Rank: "8", Suit: "Hearts"}:       "8h",
	{Rank: "9", Suit: "Hearts"}:       "9h",
	{Rank: "10", Suit: "Hearts"}:      "10h",
	{Rank: "Jack", Suit: "Hearts"}:    "jh",
	{Rank: "Queen", Suit: "Hearts"}:   "qh",
	{Rank: "King", Suit: "Hearts"}:    "kh",
	{Rank: "Ace", Suit: "Spades"}:     "as",
	{Rank: "2", Suit: "Spades"}:       "2s",
	{Rank: "3", Suit: "Spades"}:       "3s",
	{Rank: "4", Suit: "Spades"}:       "4s",
	{Rank: "5", Suit: "Spades"}:       "5s",
	{Rank: "6", Suit: "Spades"}:       "6s",
	{Rank: "7", Suit: "Spades"}:       "7s",
	{Rank: "8", Suit: "Spades"}:       "8s",
	{Rank: "9", Suit: "Spades"}:       "9s",
	{Rank: "10", Suit: "Spades"}:      "10s",
	{Rank: "Jack", Suit: "Spades"}:    "js",
	{Rank: "Queen", Suit: "Spades"}:   "qs",
	{Rank: "King", Suit: "Spades"}:    "ks",
}

func makeCardClass(card hand.Card) string {
	suffix := classSuffixLookup[card]
	log.Printf("finding the class suffix for card %v: %s", card, suffix)
	return fmt.Sprintf("pcard-%s", suffix)
}

func createPlayerViewModel(self *hand.Player, tableId string, handId string) templates.PlayerViewModel {
	allMoves := h.ValidMoves()
	mvs := allMoves[self.Id]
	cardsVM := make([]struct {
		Card  hand.Card
		Class string
	}, len(self.Cards))
	for i, c := range self.Cards {
		cardsVM[i] = struct {
			Card  hand.Card
			Class string
		}{
			c,
			makeCardClass(c),
		}
	}
	return templates.PlayerViewModel{
		Id:      self.Id,
		TableId: tableId,
		HandId:  handId,
		Entrant: templates.Entrant{
			Name:   me.Name,
			Chips:  me.Chips,
			Folded: me.Folded,
		},
		Cards: cardsVM,
		Moves: mvs,
	}
}

func createHandViewModel(playerId string, tableId string, handId string) (templates.HandViewModel, error) {
	self, opponents := h.Players(playerId)
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

func getTablesHandler(w http.ResponseWriter, req *http.Request) {
	var links = make([]string, len(tables))
	for i, v := range tables {
		links[i] = fmt.Sprintf("/table/%s/hand/%s", tableId, v.Id)
	}

	name := "hands.go.html"
	err := ts.Render(w, name, links)
	if err != nil {
		log.Printf("Error rendering template, err: %v, template name: %s, data: %v", err, name, links)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func newTableHandler(w http.ResponseWriter, req *http.Request) {
	tables = append(tables, NewTable(faker.Name()))
}

func getHandHandler(w http.ResponseWriter, req *http.Request) {
	pathVars := mux.Vars(req)
	tableId := pathVars["tableId"]
	handId := pathVars["handId"]

	vm, err := createHandViewModel(me.Id, tableId, handId)
	if err != nil {
		log.Printf("Error creating hand view model for request with URL:%s, %v\n", req.URL, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	name := "hand.go.html"
	err = ts.Render(w, name, vm)
	if err != nil {
		log.Printf("Error rendering template, err: %v, template name: %s, data: %v", err, name, vm)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func watchHandHandler(w http.ResponseWriter, req *http.Request) {
	pathVars := mux.Vars(req)
	tableId := pathVars["tableId"]
	handId := pathVars["handId"]

	vm, err := createHandViewModel(me.Id, tableId, handId)
	if err != nil {
		log.Printf("Error creating hand view model for request with URL:%s, %v\n", req.URL, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	name := "hand.go.html"
	err = ts.Render(w, name, vm)
	if err != nil {
		log.Printf("Error rendering template, err: %v, template name: %s, data: %v", err, name, vm)
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

	name := "player.go.html"
	err = ts.Render(w, name, vm)
	if err != nil {
		log.Printf("Error rendering template, err: %v, template name: %s, data: %v", err, name, vm)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
