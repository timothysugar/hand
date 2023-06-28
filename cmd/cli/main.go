package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/timothysugar/hand/pkg/hand"
)

func main() {
	// Setup signal handlers.
	ctx, cancel := context.WithCancel(context.Background())
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	initial := 10
	p1 := hand.NewPlayer("p1", initial)
	p2 := hand.NewPlayer("p2", initial)
	players := []*hand.Player{p1, p2}
	var fin chan hand.FinishedHand
	h, err := hand.NewHand(players, p1, 1)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	go func() {
		fin, err = h.Begin()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()
	ls := play(os.Stdin, h)

	// Wait for CTRL-C or hand to finish
out:
	for {
		time.Sleep(time.Millisecond * 100) //TODO: use a channel to send valid move updates
		mvs := h.ValidMoves()
		fmt.Println("Valid moves: ", mvs)
		fmt.Printf("Enter an action: [<player index><action><chips>]\ne.g. 0b1⏎ 1f0⏎\n")
		select {
		case l := <-ls:
			fmt.Println("Received line")
			pIdx, inp, err := parseLine(l)
			if err != nil {
				fmt.Printf("Could not parse line %s", l)
			}
			p := players[pIdx]
			err = h.HandleInput(p, inp)
			if err != nil {
				fmt.Printf("Could not handle input, %v, %v\n", inp, err)
			} else {
				fmt.Printf("Received input %v\n", l)
			}
		case result := <-fin:
			fmt.Printf("Hand finished with %v\n", result)
			break out
		case <-c:
			fmt.Println("Received Interrupt")
			cancel()
		case <-ctx.Done():
			fmt.Println("Done")
			break out
		}
	}

	fmt.Println("Exiting")
}

func parseLine(l string) (int, hand.Input, error) {
	rs := []rune(l)

	var err error
	var pIdx int64
	if pIdx, err = strconv.ParseInt(string(rs[0]), 10, 0); err != nil {
		return 0, hand.Input{}, err
	}

	var a hand.Action
	if a, err = parseAction(rs[1]); err != nil {
		return 0, hand.Input{}, err
	}

	var c int64
	if c, err = strconv.ParseInt(string(rs[2]), 10, 0); err != nil {
		return 0, hand.Input{}, err
	}

	return int(pIdx), hand.Input{Action: a, Chips: int(c)}, nil
}

func parseAction(r rune) (hand.Action, error) {
	switch {
	case r == 'b':
		return hand.Blind, nil
	case r == 'c':
		return hand.Check, nil
	case r == 'f':
		return hand.Fold, nil
	case r == 'l':
		return hand.Call, nil
	case r == 'r':
		return hand.Raise, nil
	default:
		return hand.Undefined, errors.New("unsupported action")
	}
}

func play(r io.Reader, h *hand.Hand) <-chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)
		scan := bufio.NewScanner(r)
		for scan.Scan() {
			s := scan.Text()
			lines <- s
		}
	}()
	return lines
}
