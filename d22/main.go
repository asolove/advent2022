package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	grid, directions := read(os.Stdin)
	fmt.Printf("Grid:\n%v\n\nDirections:\n%v\n", grid, directions)
}

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
