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

func createPlayerViewModel(self *hand.Player, tableId string, handId string) templates.PlayerViewModel {
	allMoves := h.ValidMoves()
	mvs := allMoves[self.Id]
	return templates.PlayerViewModel{
		Id:      self.Id,
		TableId: tableId,
		HandId:  handId,
		Entrant: templates.Entrant{
			Name:   me.Name,
			Chips:  me.Chips,
			Folded: me.Folded,
		},
		Cards: self.Cards,
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
