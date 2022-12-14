package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	lines := readLines(os.Stdin)
	fmt.Printf("lines: %v\n", lines)

	minX := 300
	maxX := 700
	minY := 0
	maxY := 0
	for _, l := range lines {
		for _, coord := range l {
			if coord[0] < minX {
				minX = coord[0]
			}
			if coord[0] > maxX {
				maxX = coord[0]
			}
			if coord[1] < minY {
				minY = coord[1]
			}
			if coord[1] > maxY {
				maxY = coord[1]
			}
		}
	}
	fmt.Printf("%v-%v x %v-%v\n", minX, maxX, minY, maxY)

	grid := make([][]rune, maxX-minX+1)
	for i := 0; i < maxX-minX+1; i++ {
		grid[i] = make([]rune, maxY-minY+3)
		for j := 0; j < maxY-minY+3; j++ {
			if j == maxY-minY+2 {
				grid[i][j] = '#'
			} else {
				grid[i][j] = ' '
			}
		}
	}
	fmt.Printf("%d x %d\n", len(grid), len(grid[0]))

	for _, line := range lines {
		for i, curr := range line {
			if i == 0 {
				continue
			}
			prev := line[i-1]

			dx, dy := 0, 0
			if curr[0] > prev[0] {
				dx = +1
			}
			if curr[0] < prev[0] {
				dx = -1
			}
			if curr[1] > prev[1] {
				dy = +1
			}
			if curr[1] < prev[1] {
				dy = -1
			}

			for x, y := prev[0], prev[1]; true; x, y = x+dx, y+dy {
				r := grid[x-minX]
				r[y-minY] = '#'

				if x == curr[0] && y == curr[1] {
					break
				}
			}
		}
	}

	for i := 0; true; i++ {
		dropSand(grid, 500-minX, 0-minY)
		// fmt.Printf("After turn %d:\n", i)
		// printGrid(grid)
		if grid[500-minX][-minY] != ' ' {
			fmt.Printf("Stacked up to root after %d units\n", i)
			break
		}
	}
}

func printGrid(grid [][]rune) {
	for y, _ := range grid[0] {
		for x, _ := range grid {
			fmt.Printf("%c", grid[x][y])
		}
		fmt.Printf("\n")
	}
}

// Returns true iff the sand runs off the grid
func dropSand(grid [][]rune, x, y int) bool {
	if x < 0 || x >= len(grid) {
		return true
	}
	if y < 0 || y >= len(grid[0])-1 {
		return true
	}

	below := grid[x][y+1]
	if below == ' ' {
		return dropSand(grid, x, y+1)
	}
	if below == 'O' || below == '#' {
		if x == 0 {
			return true
		}
		if grid[x-1][y+1] == ' ' {
			return dropSand(grid, x-1, y+1)
		}
		if x == len(grid)-1 {
			return true
		}
		if grid[x+1][y+1] == ' ' {
			return dropSand(grid, x+1, y+1)
		}
	}
	grid[x][y] = 'O'
	return false
}

/*
498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9
*/
func readLines(f *os.File) [][][2]int {
	r := make([][][2]int, 0)
	s := bufio.NewScanner(f)
	for s.Scan() {
		l := make([][2]int, 0)
		for _, coord := range strings.Split(s.Text(), " -> ") {
			cs := strings.Split(coord, ",")
			l = append(l, [2]int{atoi(cs[0]), atoi(cs[1])})
		}
		r = append(r, l)
	}
	return r
}

func atoi(i string) int {
	n, e := strconv.Atoi(i)
	if e != nil {
		panic("Can't parse to digit")
	}
	return n
}
