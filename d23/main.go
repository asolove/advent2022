package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const MAX_ROUNDS = 1000

func main() {
	elves := readElves(os.Stdin)
	// fmt.Printf("Elves at start:\n%v\n", elves)

	directions := [4][3][2]int{
		// N
		{{0, -1}, {-1, -1}, {1, -1}},
		// S
		{{0, 1}, {-1, 1}, {1, 1}},
		// W
		{{-1, 0}, {-1, -1}, {-1, 1}},
		// E
		{{1, 0}, {1, -1}, {1, 1}},
	}

	for round := 0; round < MAX_ROUNDS; round++ {
		firstDir := round % len(directions)

		// newLoc -> oldLoc
		proposedMoves := make(map[string]string)
		anyMoved := false

		for pos, _ := range elves.grid {
			// fmt.Printf("\nConsidering how to move elf at %s\n", pos)
			x, y := coords(pos)
			proposedMoves[pos] = pos

			shouldMove := false
		shouldMove:
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					if dx == 0 && dy == 0 {
						continue
					}
					if elves.grid[key(x+dx, y+dy)] {
						shouldMove = true
						break shouldMove
					}
				}
			}

			if !shouldMove {
				continue
			} else {
				anyMoved = true
			}

			// for each direction from here
			//   if not empty
			//     then it shouldMove

			for d := 0; d < 4; d++ {
				dir := directions[(firstDir+d)%len(directions)]
				canMove := true
				for _, offset := range dir {
					x2, y2 := x+offset[0], y+offset[1]
					if elves.grid[key(x2, y2)] {
						canMove = false
						break
					}
				}
				if !canMove {
					continue
				} else {
					// move to destination for current direction
					destination := key(x+dir[0][0], y+dir[0][1])
					if origKey, found := proposedMoves[destination]; found {
						// fmt.Printf("    Found collision at %v. Reverting both\n", destination)
						// Another elf already wanted to move to the same space, so move it back
						proposedMoves[origKey] = origKey
						delete(proposedMoves, destination)
					} else {
						// fmt.Printf("    Accepting direction %v\n", dir[0])
						proposedMoves[destination] = pos
						delete(proposedMoves, pos)
					}

					// don't consider any more directions
					break
				}
			}
		}

		elves.grid = make(map[string]bool)
		elves.minX = math.MaxInt
		elves.minY = math.MaxInt
		elves.maxX = math.MinInt
		elves.maxY = math.MinInt
		for dest, _ := range proposedMoves {
			x, y := coords(dest)

			if x > elves.maxX {
				elves.maxX = x
			}
			if x < elves.minX {
				elves.minX = x
			}
			if y > elves.maxY {
				elves.maxY = y
			}
			if y < elves.minY {
				elves.minY = y
			}
			elves.grid[dest] = true
		}

		// fmt.Printf("Elves after round %d:\n%v\n", round, elves)

		if !anyMoved {
			fmt.Printf("No moves needed after %d\n", round)
			break
		}

	}

	// fmt.Printf("Empty spaces after round %d: %d\n", MAX_ROUNDS, elves.EmptySpaces())

}

type Elves struct {
	maxX int
	maxY int
	minX int
	minY int
	grid map[string]bool
}

func NewElves() *Elves {
	return &Elves{
		minX: math.MaxInt,
		maxX: math.MinInt,
		minY: math.MaxInt,
		maxY: math.MinInt,
		grid: make(map[string]bool),
	}
}

func (e *Elves) String() string {
	s := fmt.Sprintf("Elves:\n      %d\n", e.minX)
	for y := e.minY; y <= e.maxY; y++ {
		s += fmt.Sprintf("%4d  ", y)
		for x := e.minX; x <= e.maxX; x++ {
			if e.grid[key(x, y)] {
				s += "#"
			} else {
				s += "."
			}
		}
		s += "\n"
	}
	return s
}

func (e *Elves) EmptySpaces() int {
	empty := 0
	for y := e.minY; y <= e.maxY; y++ {
		for x := e.minX; x <= e.maxX; x++ {
			if !e.grid[key(x, y)] {
				empty++
			}
		}
	}
	return empty
}

func readElves(f *os.File) *Elves {
	elves := NewElves()
	s := bufio.NewScanner(f)
	for y := 0; s.Scan(); y++ {
		for x, c := range s.Text() {
			if c == '#' {
				elves.grid[key(x, y)] = true
				if x > elves.maxX {
					elves.maxX = x
				}
				if x < elves.minX {
					elves.minX = x
				}
				if y > elves.maxY {
					elves.maxY = y
				}
				if y < elves.minY {
					elves.minY = y
				}
			}
		}
	}

	return elves
}

func coords(k string) (int, int) {
	coords := strings.Split(k, ",")
	return atoi(coords[0]), atoi(coords[1])
}

func key(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

func atoi(i string) int {
	n, e := strconv.Atoi(i)
	if e != nil {
		panic("Can't parse to digit")
	}
	return n
}
