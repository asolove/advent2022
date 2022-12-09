package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	moves := parseMoves()

	var state State
	tailSpots := make(map[string]bool)
	for _, move := range moves {
		for i := 0; i < move.steps; i++ {
			state = stepHead(state, move.dir)
			state = stepTail(state)
			tailSpots[fmt.Sprintf("%d,%d", state.tailX, state.tailY)] = true
		}

	}
	fmt.Printf("Visited by tail: %d\n", len(tailSpots))
}

func stepHead(state State, dir Dir) State {
	switch dir {
	case Up:
		state.headY -= 1
	case Down:
		state.headY += 1
	case Left:
		state.headX -= 1
	case Right:
		state.headX += 1
	}
	return state
}

func stepTail(state State) State {
	xd := state.headX - state.tailX
	xm := xd < -1 || xd > 1
	yd := state.headY - state.tailY
	ym := yd < -1 || yd > 1

	if !xm && !ym {
		return state
	}

	if xd >= 1 {
		state.tailX += 1
	}
	if xd <= -1 {
		state.tailX -= 1
	}
	if yd >= 1 {
		state.tailY += 1
	}
	if yd <= -1 {
		state.tailY -= 1
	}
	return state
}

type State struct {
	headX int
	headY int
	tailX int
	tailY int
}

type Dir string

const (
	Up    Dir = "U"
	Down      = "D"
	Left      = "L"
	Right     = "R"
)

type Move struct {
	dir   Dir
	steps int
}

func parseMoves() []Move {
	moves := make([]Move, 0)
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		words := strings.Split(s.Text(), " ")
		moves = append(moves, Move{dir: atod(words[0]), steps: atoi(words[1])})
	}
	return moves
}

func atod(s string) Dir {
	dirs := map[string]Dir{"U": Up, "D": Down, "L": Left, "R": Right}
	return dirs[s]
}

func atoi(i string) int {
	n, e := strconv.Atoi(i)
	if e != nil {
		panic("Can't parse to digit")
	}
	return n
}
