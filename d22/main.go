package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	m, directions := read(os.Stdin)
	fmt.Printf("Grid:\n%v\n\nDirections:\n%v\n", m, directions)

	endState := MoveAll(m, directions)
	fmt.Printf("Finish state: %v\n", endState)
	fmt.Printf("Score: %v\n", Score(endState))
}

func StartState(m *Map) *State {
	y := 0
	x := 0
	for ; x <= m.maxX; x++ {
		if _, ok := m.grid[key(x, y)]; ok {
			break
		}
	}
	return &State{x: x, y: y, heading: East}
}

func MoveAll(m *Map, directions []int) *State {
	s := StartState(m)
	for i := 0; i < len(directions); i++ {
		s = Move(m, s, directions[i])
		i++
		s = Turn(s, directions[i])
	}
	return s
}

func Score(s *State) int {
	return 1000*(s.y+1) + 4*(s.x+1) + int(s.heading)
}

func Move(m *Map, state *State, n int) *State {
	x := state.x
	y := state.y
	for i := 0; i < n; i++ {
		nextX, nextY := Step(m, state.heading, x, y)
		blocked, ok := m.grid[key(nextX, nextY)]
		if !ok {
			fmt.Errorf("Step resulted in an invalid space: %v, %v", nextX, nextY)
		}
		if blocked {
			break
		} else {
			x, y = nextX, nextY
		}
	}
	return &State{x: x, y: y, heading: state.heading}
}

func Turn(state *State, direction int) *State {
	heading := state.heading
	if direction == Left {
		heading = heading - 1
		if heading < 0 {
			heading = 3
		}
	} else {
		heading = heading + 1
		if heading > 3 {
			heading = 0
		}
	}
	return &State{x: state.x, y: state.y, heading: heading}
}

func Step(m *Map, h Heading, x, y int) (int, int) {
	nextX, nextY := x, y
	switch h {
	case North:
		nextY -= 1
		if nextY < 0 {
			nextY = m.maxY
		}
	case South:
		nextY += 1
		if nextY > m.maxY {
			nextY = 0
		}
	case East:
		nextX += 1
		if nextX > m.maxX {
			nextX = 0
		}
	case West:
		nextX -= 1
		if nextX < 0 {
			nextX = m.maxX
		}
	}

	if _, ok := m.grid[key(nextX, nextY)]; ok {
		return nextX, nextY
	} else {
		return Step(m, h, nextX, nextY)
	}
}

type State struct {
	x       int
	y       int
	heading Heading
}

type Heading int

const (
	East  Heading = iota
	South         = iota
	West          = iota
	North         = iota
)

type Map struct {
	maxX int
	maxY int
	grid map[string]bool // true: wall, false: free space, not present: no space
}

func (m *Map) String() string {
	s := ""
	for y := 0; y <= m.maxY; y++ {
		for x := 0; x <= m.maxX; x++ {
			if wall, found := m.grid[key(x, y)]; !found {
				s += fmt.Sprintf(" ")
			} else if wall {
				s += fmt.Sprintf("#")
			} else {
				s += fmt.Sprintf(".")
			}
		}
		s += fmt.Sprintf("\n")
	}
	return s
}

const (
	Left  = -1
	Right = -2
)

func key(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

func read(f *os.File) (*Map, []int) {
	s := bufio.NewScanner(f)
	m := Map{
		grid: make(map[string]bool, 0),
	}
	y := 0
	for s.Scan() {
		text := s.Text()
		if text == "" {
			break
		}
		for x, c := range text {
			if c == '.' {
				m.grid[key(x, y)] = false
			} else if c == '#' {
				m.grid[key(x, y)] = true
			}
			if x > m.maxX {
				m.maxX = x
			}
		}
		if y > m.maxY {
			m.maxY = y
		}
		y++
	}
	s.Scan()

	ds := make([]int, 0)
	number := ""
	for _, r := range s.Text() {
		if r >= '0' && r <= '9' {
			number += string(r)
		} else {
			direction := Left
			if r == 'R' {
				direction = Right
			}
			if len(number) > 0 {
				ds = append(ds, atoi(number))
				number = ""
			}
			ds = append(ds, direction)
		}
	}

	return &m, ds
}

func atoi(i string) int {
	n, e := strconv.Atoi(i)
	if e != nil {
		panic(fmt.Sprintf("Can't parse to digit: %s", i))
	}
	return n
}
