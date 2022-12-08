package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	grid := parseGrid(os.Stdin)
	visible := makeSet(len(grid), len(grid[0]))

	markVisibleLTR(grid, visible)
	for i := 0; i < 4; i++ {
		rotate(grid)
		rotate(visible)
		markVisibleLTR(grid, visible)
	}
	fmt.Printf("Visible: %d\n", countVisible(visible))
}

func countVisible(visible [][]bool) int {
	count := 0
	for r := 0; r < len(visible); r++ {
		for c := 0; c < len(visible[r]); c++ {
			if visible[r][c] {
				count += 1
			}
		}
	}
	return count
}

func markVisibleLTR(grid [][]int, visible [][]bool) {
	for r := 0; r < len(grid); r++ {
		row := grid[r]
		tallest := -1
		for c := 0; c < len(row); c++ {
			if row[c] > tallest {
				tallest = row[c]
				visible[r][c] = true
			}
		}
	}
}

func rotate[T bool | int](grid [][]T) {
	// Assumes a square grid for now
	max := len(grid) - 1
	for r := 0; r < max/2; r++ {
		for c := 0; c < max/2+1; c++ {
			grid[r][c], grid[c][max-r], grid[max-r][max-c], grid[max-c][r] = grid[max-c][r], grid[r][c], grid[c][max-r], grid[max-r][max-c]
		}
	}
}

func parseGrid(f *os.File) [][]int {
	grid := make([][]int, 0)
	s := bufio.NewScanner(f)

	for s.Scan() {
		row := make([]int, 0)
		for _, c := range s.Text() {
			row = append(row, rtoi(c))
		}
		grid = append(grid, row)
	}
	return grid
}

func makeSet(rows, cols int) [][]bool {
	set := make([][]bool, 0)

	for row := 0; row < rows; row++ {
		row := make([]bool, cols)
		set = append(set, row)
	}
	return set
}

func rtoi(i rune) int {
	n, e := strconv.Atoi(string(i))
	if e != nil {
		panic("Can't parse to digit")
	}
	return n
}
