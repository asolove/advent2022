package main

import (
	"bufio"
	"fmt"
	"os"
)

const WIDTH = 7

// Best: 45s for 1/1000 of total
const GOAL_STEPS = 1000000000000
const MAX_SHAPES = GOAL_STEPS / 1000

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
			// fmt.Printf("After step %d, adding shape %d:\n%v\n", steps, i, board)
		}

	}
	// fmt.Printf("After shape %d:\n", maxShapes)
	// fmt.Printf("Board:\n%v\n", board)
	// fmt.Printf("At end, height: %d\n", board.height)
	return board.height
}

const MAX_GRID_BLOCKS = 1000000

type Board struct {
	height int
	shape  *Shape
	sx     int
	sy     int
	grid   [MAX_GRID_BLOCKS]uint64
	// how much we have deleted and need to offset y coords
}

func NewBoard() *Board {
	return &Board{height: 0, grid: [MAX_GRID_BLOCKS]uint64{}}
}

func (b *Board) Step(wind int) bool {
	done := false

	// try wind blowing
	b.sx += wind
	if !b.validSpot() {
		b.sx -= wind
		// fmt.Printf("Wind being ignored\n")
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

func (b *Board) blockIndex(y int) int {
	return (y / 8) % MAX_GRID_BLOCKS
}

func (b *Board) finishShape() {
	shapeMaxY := b.sy + b.shape.height - 1

	if shapeMaxY >= b.height {
		b.height = shapeMaxY + 1
	}

	blockIndex := b.blockIndex(b.sy)
	// fmt.Printf("Finishing shape at %d,%d\n", b.sx, b.sy)
	// fmt.Printf("First block\nBefore | Shape | Shape positioned | Result\n%v\n", showBlocks(b.grid[blockIndex], b.shape.mask, b.shape.mask<<(b.sx+8*(b.sy%8)), b.grid[blockIndex]|(b.shape.mask<<(b.sx+8*(b.sy%8)))))
	b.grid[blockIndex] |= (b.shape.mask << (b.sx + 8*(b.sy%8)))
	// fmt.Printf("Block above\nBefore | Shape | Shape positioned | Result\n%v\n", showBlocks(b.grid[blockIndex+1], b.shape.mask, (b.shape.mask>>(-b.sx+-8*(b.sy%8-8)))))
	b.grid[(blockIndex+1)%MAX_GRID_BLOCKS] |= (b.shape.mask >> (-b.sx + -8*(b.sy%8-8)))

	b.shape = nil
}

func getBit(x int, y int, b uint64) bool {
	return b&(1<<(x+y*8)) > 0
}

func (b *Board) get(x, y int) bool {
	if y >= b.height {
		return false
	}
	blockInGrid := (y / 8) % MAX_GRID_BLOCKS
	yInBlock := y % 8
	// fmt.Printf("get(%d, %d): yInGrid: %d, yInBlock: %d\n", x, y, yInGrid, yInBlock)
	return getBit(x, yInBlock, b.grid[blockInGrid])
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

	blockIndex := b.blockIndex(b.sy)
	if b.grid[blockIndex]&(b.shape.mask<<(b.sx+8*(b.sy%8))) > 0 {
		return false
	}

	// fmt.Printf("Block above\nBefore | Shape | Shape positioned | Result\n%v\n", showBlocks(b.grid[blockIndex+1], b.shape.mask, (b.shape.mask>>(-b.sx+-8*(b.sy%8-8)))))
	if b.grid[(blockIndex+1)%MAX_GRID_BLOCKS]&(b.shape.mask>>(-b.sx+-8*(b.sy%8-8))) > 0 {
		return false
	}

	return true
}

func (b *Board) AddShape(s *Shape) {
	b.shape = s
	b.sx = 2
	b.sy = b.Height() + 3
	b.clearNextBlock(b.sy)
}

func (b *Board) clearNextBlock(y int) {
	blockIndex := b.blockIndex(y)
	b.grid[(blockIndex+1)%MAX_GRID_BLOCKS] = 0
}

func (b *Board) Height() int {
	return b.height
}

func (b *Board) String() string {
	var s = "+-------+"
	pixels := map[bool]string{true: "#", false: "."}

	maxY := b.height
	minY := max(0, b.height-(MAX_GRID_BLOCKS*8))
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
	for i := range shapes {
		s := &shapes[i]
		s.mask = uint64(0)
		for _, coord := range shapes[i].coords {
			s.mask |= coordToMask(coord[0], coord[1])
		}
	}
	return shapes
}

func showBlocks(bs ...uint64) string {
	r := ""
	for i := 7; i >= 0; i-- {
		for _, b := range bs {
			r += fmt.Sprintf("%8b", byte(0xff&(b>>(8*i))))
		}
		r += "\n"
	}
	r += "\n"
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

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
