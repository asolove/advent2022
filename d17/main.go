package main

import (
	"bufio"
	"fmt"
	"os"
)

const WIDTH = 7

// Best: 65s for 1/1000 of total
const GOAL_STEPS = 1000000000000
const MAX_SHAPES = 100000000

func main() {
	run(os.Stdin)
}

func run(f *os.File) {
	shapes := allShapes()
	wind := readWind(f)
	steps := 0
	board := NewBoard()

	for i := 0; i < MAX_SHAPES; i++ {
		if i%(MAX_SHAPES/10) == 0 {
			fmt.Printf("At shape %d\n", i)
		}
		board.AddShape(&shapes[i%len(shapes)])
		// fmt.Printf("At step %d, adding shape %d:\n%v\n", steps, i, board)
		falling := true

		for falling {
			falling = board.Step(wind[steps%len(wind)])
			steps++
		}
	}
	fmt.Printf("After shape %d:\n", MAX_SHAPES)
	// fmt.Printf("Board:\n%v\n", board)
	fmt.Printf("At end, height: %d\n", board.height)
}

type Board struct {
	height int
	shape  *Shape
	sx     int
	sy     int
	grid   []uint8
	// how much we have deleted and need to offset y coords
	gridOffset int
}

func NewBoard() *Board {
	return &Board{height: 0, grid: make([]uint8, 0)}
}

func (b *Board) Step(wind int) bool {
	done := false

	// try wind blowing
	b.sx += wind
	if !b.validSpot() {
		b.sx -= wind
	}

	// try falling
	b.sy -= 1
	if !b.validSpot() {
		b.sy += 1
		b.finishShape()
		done = true
	}
	return !done
}

func (b *Board) finishShape() {
	for _, coord := range b.shape.coords {
		b.set(coord[0]+b.sx, coord[1]+b.sy)
	}
	b.shape = nil
}

func getBit(n int, b byte) bool {
	return b&(1<<(n)) > 0
}
func setBit(n int, b byte) byte {
	return b | (1 << (n))
}

func (b *Board) get(x, y int) bool {
	if y >= b.height {
		return false
	}
	return getBit(x, b.grid[y-b.gridOffset])
}

func (b *Board) set(x, y int) {
	if y >= b.height {
		b.addRows(y - b.height + 1)
	}
	b.grid[y-b.gridOffset] = setBit(x, b.grid[y-b.gridOffset])
}

func (b *Board) validSpot() bool {
	if b.sy < 0 {
		return false
	}
	if b.sy >= b.height {
		return b.sx >= 0 && b.sx+b.shape.width < WIDTH
	}
	if b.sx < 0 || b.sx+b.shape.width > WIDTH {
		return false
	}
	for _, coord := range b.shape.collCoords {
		x := coord[0] + b.sx
		y := coord[1] + b.sy
		if b.get(x, y) {
			return false
		}
	}
	return true
}

func (b *Board) AddShape(s *Shape) {
	// if b.shape != nil {
	// 	panic("Adding shape before clearing previous\n")
	// }
	b.shape = s
	b.sx = 2
	b.sy = b.Height() + 3
}

func (b *Board) addRows(n int) {
	for ; n > 0; n-- {
		b.grid = append(b.grid, byte(0))
		b.height++
	}
	if len(b.grid) > 200 {
		b.grid = b.grid[100:]
		b.gridOffset += 100
	}
}

func (b *Board) Height() int {
	return b.height
}

func (b *Board) String() string {
	var s = "+-------+"
	pixels := map[bool]string{true: "#", false: "."}

	minY := b.gridOffset
	maxY := b.height
	if b.shape != nil {
		maxY += b.shape.height + 3
	}

	for y := minY; y <= maxY; y++ {
		r := "|"
		for x := 0; x < WIDTH; x++ {
			if b.ShapeAt(x, y) {
				r += "@"
			} else {
				r += pixels[b.get(x, y)]
			}
		}
		r += fmt.Sprintf("| %d\n", y)
		s = r + s
	}
	return s
}

func (b *Board) ShapeAt(x, y int) bool {
	if b.shape == nil {
		return false
	}
	for _, coord := range b.shape.coords {
		if x == coord[0]+b.sx && y == coord[1]+b.sy {
			return true
		}
	}
	return false
}

// A shape is defined by its coordinates relative to the bottom-left-most
// point on the shape's grid (which in some cases is not in the shape itself)
type Shape struct {
	coords     [][2]int
	collCoords [][2]int
	width      int
	height     int
}

func allShapes() []Shape {
	return []Shape{
		// horizontal line
		{
			coords:     [][2]int{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
			collCoords: [][2]int{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
			width:      4,
			height:     1,
		},
		// plus
		{
			coords:     [][2]int{{1, 0}, {0, 1}, {1, 1}, {2, 1}, {1, 2}},
			collCoords: [][2]int{{1, 0}, {0, 1}, {2, 1}},
			width:      3,
			height:     3,
		},

		// flipped L
		{
			coords:     [][2]int{{0, 0}, {1, 0}, {2, 0}, {2, 1}, {2, 2}},
			collCoords: [][2]int{{0, 0}, {1, 0}, {2, 0}},
			width:      3,
			height:     3,
		},
		// vertical line
		{
			coords:     [][2]int{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
			collCoords: [][2]int{{0, 0}},
			width:      1,
			height:     4,
		},
		// square
		{
			coords:     [][2]int{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
			collCoords: [][2]int{{0, 0}, {1, 0}},
			width:      2,
			height:     2,
		},
	}
}

func readWind(f *os.File) []int {
	r := make([]int, 0)

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanRunes)
	for s.Scan() {
		dir := +1
		if s.Text() == "<" {
			dir = -1
		}
		r = append(r, dir)
	}
	return r
}

func min(i, j int) int {
	if i > j {
		return j
	}
	return i
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}
