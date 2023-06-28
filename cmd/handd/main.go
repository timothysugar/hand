package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/timothysugar/hand/cmd/handd/templates"
	"github.com/timothysugar/hand/pkg/hand"
)

type PlayerViewModel struct {
	*hand.Player
	Moves  []hand.Move
	HandId string
}

type HandsViewModel struct {
	Name         string
	TableId      string
	HandListings []HandListingViewModel
}

type ctxKey string

const (
	initialChips            = 1000
	tableId                 = "table1"
	playerSessionKey ctxKey = "playerSession"
)

var h *hand.Hand
var me *hand.Player
var ts *templates.Template

var sessionStore *sessions.CookieStore
var pStore *playerStore = &playerStore{players: make(map[string]*hand.Player)}

type playerStore struct {
	players map[string]*hand.Player
	mux     sync.RWMutex
}

func (ps *playerStore) Has(name string) bool {
	ps.mux.Lock()
	defer ps.mux.Unlock()

	for _, p := range ps.players {
		if p.Name == name {
			return true
		}
	}
	return false
}

func (ps *playerStore) Save(p *hand.Player) {
	ps.mux.Lock()
	defer ps.mux.Unlock()

	ps.players[p.Id] = p
}

func (ps *playerStore) Get(id string) *hand.Player {
	ps.mux.Lock()
	defer ps.mux.Unlock()

	return ps.players[id]
}

func init() {
	log.Println("Initializing hand")
	ts = templates.New()
	sessionsKey := os.Getenv("SESSION_KEY")
	var err error
	sessionStore = sessions.NewCookieStore([]byte(sessionsKey))
	if err != nil {
		panic(err)
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
	h, err = hand.NewHand(players, players[2], 10)
	if err != nil {
		log.Fatalf("Error initializing hand: %s", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received: %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func getPlayer(ctx context.Context) (*hand.Player, error) {
	p, ok := ctx.Value(playerSessionKey).(*hand.Player)
	if !ok {
		return nil, errors.New("could not get player from context")
	}
	return p, nil
}

func setPlayer(req *http.Request, p *hand.Player) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), playerSessionKey, p))
}

func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var p *hand.Player
		session, err := sessionStore.Get(req, "hand-poker-session")
		if err != nil {
			log.Printf("Session could not be decoded from store, err: %v, req: %v\n", err, req)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if session.IsNew {
			log.Printf("No existing session found, session: %v\n", session)
			http.Redirect(w, req, "/signup", http.StatusFound)
			return
		}

		pId := session.Values["playerId"]
		if pId == nil {
			log.Printf("No player ID found in session, session: %v\n", session)
			http.Redirect(w, req, "/signup", http.StatusFound)
			return
		}

		p = pStore.Get(pId.(string))
		if p == nil {
			log.Printf("Player not found in store, pId: %s, session: %v, store: %v\n", pId, session, pStore)
			http.Redirect(w, req, "/signup", http.StatusFound)
			return
		}

		log.Printf("Request received from %s: player ID: %s", p.Name, p.Id)
		req = setPlayer(req, p)

		next.ServeHTTP(w, req)
	})
}

func main() {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/signup", signupHandler).Name("signup").Methods("GET")
	r.HandleFunc("/players", newPlayerHandler).Name("new-player").Methods("POST")

	// authenticated routes
	sm := r.PathPrefix("/").Subrouter()
	sm.Use(authenticationMiddleware)
	sm.HandleFunc("/", getHandsHandler).Name("get-hands").Methods("GET")
	sm.HandleFunc("/table/{tableId}/hand/{handId}/players/join", joinHandHandler).Name("join-hand").Methods("POST")
	sm.HandleFunc("/table/{tableId}/hand/{handId}", getHandHandler).Name("get-hand").Methods("GET")
	sm.HandleFunc("/table/{tableId}/hand/{handId}/begin", beginHandHandler).Name("begin-hand").Methods("POST")
	sm.HandleFunc("/hand/{handId}/player/{playerId}/moves/Blind/{amount}", blindHandler).Name("play-blind").Methods("POST")

	port := ":8090"
	fmt.Printf("listening on %s", port)
	http.ListenAndServe(port, r)
}

func signupHandler(w http.ResponseWriter, req *http.Request) {
	tmpl := "signup.go.html"
	err := ts.Render(w, tmpl, nil)
	if err != nil {
		log.Printf("Error rendering template, err: %v, template name: %s", err, tmpl)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func newPlayerHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Printf("Error parsing form, err: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session, err := sessionStore.Get(req, "hand-poker-session")
	if err != nil {
		log.Printf("Session could not be decoded from store, session: %v, err: %v\n", session, err)
	}

	name := req.Form.Get("name")
	if pStore.Has(name) {
		log.Printf("Player already exists with name: %s\n", name)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p := hand.NewPlayer(name, initialChips)
	pStore.Save(p)

	session.Values["name"] = name
	session.Values["playerId"] = p.Id
	err = session.Save(req, w)
	if err != nil {
		log.Printf("Error saving session for %s, session: %v, err: %v\n", name, session, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/", http.StatusFound)
}

func createPlayerViewModel(self *hand.Player, tableId string, handId string) templates.PlayerViewModel {
	allMoves := h.ValidMoves()
	mvs := allMoves[self.Id]
	return templates.PlayerViewModel{
		Id:      self.Id,
		TableId: tableId,
		HandId:  handId,
		Entrant: templates.Entrant{
			Name:   self.Name,
			Chips:  self.Chips,
			Folded: self.Folded,
		},
		Cards: self.Cards,
		Moves: mvs,
	}
}

func createHandViewModel(playerId string, tableId string, handId string) (templates.HandViewModel, error) {
	self, opponents := h.Players(playerId)
	if self == nil {
		return templates.HandViewModel{}, errors.New("player not found in hand")
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
		IsActive: h.IsActive(),
	}, nil
}

type HandListingViewModel struct {
	HandId  string
	Players []struct {
		Name string
		Id   string
	}
	Joined bool
	Active bool
}

func createHandsViewModel(player *hand.Player, hands []*hand.Hand) HandsViewModel {
	hs := make([]HandListingViewModel, len(hands))
	for i, v := range hands {
		self, ops := v.Players(player.Id)
		ps := make([]struct {
			Name string
			Id   string
		}, len(ops))
		log.Printf("ops: %v\n", ops)
		for i, p := range ops {
			ps[i] = struct {
				Name string
				Id   string
			}{
				Name: p.Name,
				Id:   p.Id,
			}
		}

		log.Printf("ps: %v\n", ps)
		hs[i] = HandListingViewModel{
			HandId:  v.Id,
			Players: ps,
			Joined:  self != nil,
			Active: v.IsActive(),
		}
	}
	return HandsViewModel{
		Name:         me.Name,
		TableId:      tableId,
		HandListings: hs,
	}
}

func getHandsHandler(w http.ResponseWriter, req *http.Request) {
	hands := []*hand.Hand{h}
	ctx := req.Context()
	p, err := getPlayer(ctx)
	if err != nil {
		log.Printf("Could not get player from request context, ctx: %v, err: %v", ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	vm := createHandsViewModel(p, hands)
	tmpl := "hands.go.html"
	err = ts.Render(w, tmpl, vm)
	if err != nil {
		log.Printf("Error rendering template, err: %v, template name: %s", err, tmpl)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func joinHandHandler(w http.ResponseWriter, req *http.Request) {
	// pathVars := mux.Vars(req)
	// tableId := pathVars["tableId"]
	// handId := pathVars["handId"]
	ctx := req.Context()
	p, err := getPlayer(ctx)
	if err != nil {
		log.Printf("Could not get player from request context, ctx: %v, err: %v", ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.Join(p, initialChips); err != nil {
		log.Printf("Error with %v joining hand: %v\n", p, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Player %v joined hand: %v\n", p, h.Id)
	w.Header().Set("HX-Redirect", fmt.Sprintf("/table/%s/hand/%s", tableId, h.Id))
}

func getHandHandler(w http.ResponseWriter, req *http.Request) {
	pathVars := mux.Vars(req)
	tableId := pathVars["tableId"]
	handId := pathVars["handId"]

	ctx := req.Context()
	p, err := getPlayer(ctx)
	if err != nil {
		log.Printf("Could not get player from request context, ctx: %v, err: %v", ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vm, err := createHandViewModel(p.Id, tableId, handId)
	if err != nil {
		log.Printf("Error creating hand view model for request with URL:%s, %v\n", req.URL, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl := "hand.go.html"
	err = ts.Render(w, tmpl, vm)
	if err != nil {
		log.Printf("Error rendering template, err: %v, template name: %s, data: %v", err, tmpl, vm)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func beginHandHandler(w http.ResponseWriter, req *http.Request) {
	pathVars := mux.Vars(req)
	tableId := pathVars["tableId"]
	handId := pathVars["handId"]

	ctx := req.Context()
	p, err := getPlayer(ctx)
	if err != nil {
		log.Printf("Could not get player from request context, ctx: %v, err: %v", ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := h.Begin(); err != nil {
		log.Printf("Error beginning hand: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Hand %v begun by %v\n", h.Id, p)

	vm, err := createHandViewModel(p.Id, tableId, handId)
	if err != nil {
		log.Printf("Error creating hand view model for request with URL:%s, %v\n", req.URL, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl := "hand"
	err = ts.RenderPartial(w, tmpl, vm)
	if err != nil {
		log.Printf("Error rendering template, err: %v, template name: %s, data: %v", err, tmpl, vm)
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

	tmpl := "player.go.html"
	err = ts.Render(w, tmpl, vm)
	if err != nil {
		log.Printf("Error rendering template, err: %v, template name: %s, data: %v", err, tmpl, vm)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
