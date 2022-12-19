package main

import (
	"bufio"
	"fmt"
	"os"
)

const WIDTH = 7

// Best: 65s for 1/1000 of total
const GOAL_STEPS = 1000000000000
const MAX_SHAPES = 2022

func main() {
	wind := readWind(os.Stdin)
	run(wind, MAX_SHAPES)
}

func run(wind []int, maxShapes int) int {
	shapes := allShapes()
	steps := 0
	board := NewBoard()

	for i := 0; i < maxShapes; i++ {
		// if i%(maxShapes/10) == 0 {
		// 	fmt.Printf("At shape %d\n", i)
		// }
		board.AddShape(&shapes[i%len(shapes)])
		// fmt.Printf("At step %d, adding shape %d:\n%v\n", steps, i, board)
		falling := true

		for falling {
			falling = board.Step(wind[steps%len(wind)])
			steps++
		}
	}
	// fmt.Printf("After shape %d:\n", maxShapes)
	// fmt.Printf("Board:\n%v\n", board)
	// fmt.Printf("At end, height: %d\n", board.height)
	return board.height
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
	maxY := len(b.grid) - 1
	if y > maxY {
		b.addRows(y - maxY)
	}
	if y >= b.height {
		b.height = y + 1
	}
	b.grid[y-b.gridOffset] = setBit(x, b.grid[y-b.gridOffset])
}

func (b *Board) validSpot() bool {
	if b.sy < 0 {
		return false
	}
	if b.sy >= b.height {
		return b.sx >= 0 && b.sx+b.shape.width <= WIDTH
	}
	if b.sx < 0 || b.sx+b.shape.width > WIDTH {
		return false
	}
	for _, coord := range b.shape.coords {
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
	coords [][2]int
	mask   uint64
	width  int
	height int
}

func allShapes() []Shape {
	shapes := []Shape{
		// horizontal line
		{
			coords: [][2]int{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
			width:  4,
			height: 1,
		},
		// plus
		{
			coords: [][2]int{{1, 0}, {0, 1}, {1, 1}, {2, 1}, {1, 2}},
			width:  3,
			height: 3,
		},

		// flipped L
		{
			coords: [][2]int{{0, 0}, {1, 0}, {2, 0}, {2, 1}, {2, 2}},
			width:  3,
			height: 3,
		},
		// vertical line
		{
			coords: [][2]int{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
			width:  1,
			height: 4,
		},
		// square
		{
			coords: [][2]int{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
			width:  2,
			height: 2,
		},
	}
	for _, s := range shapes {
		s.mask = uint64(0)
		for _, coord := range s.coords {
			s.mask |= coordToMask(coord[0], coord[1])
		}
	}
	return shapes
}

func showBlock(bs ...uint64) string {
	r := ""
	for i := 7; i >= 0; i-- {
		for _, b := range bs {
			r += fmt.Sprintf("%8b\n", byte(0xff&(b>>(8*i))))
		}
		r += "\n"
	}
	return r
}

func coordToMask(x, y int) uint64 {
	return 1 << (x + y*8)
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
