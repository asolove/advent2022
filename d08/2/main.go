package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	grid := parseGrid(os.Stdin)

	highScore := 0
	for r, row := range grid {
		for c, _ := range row {
			if score := scoreLocation(r, c, grid); score > highScore {
				highScore = score
			}
		}
	}
	fmt.Printf("High score: %d\n", highScore)
}

func scoreLocation(siteR, siteC int, grid [][]int) int {
	siteHeight := grid[siteR][siteC]
	upScore := siteR
	for r := siteR - 1; r >= 0; r-- {
		if grid[r][siteC] >= siteHeight {
			upScore = siteR - r
			break
		}
	}
	downScore := len(grid) - siteR - 1
	for r := siteR + 1; r < len(grid); r++ {
		if grid[r][siteC] >= siteHeight {
			downScore = r - siteR
			break
		}
	}

	leftScore := siteC
	for c := siteC - 1; c >= 0; c-- {
		if grid[siteR][c] >= siteHeight {
			leftScore = siteC - c
			break
		}
	}

	rightScore := len(grid[siteR]) - siteC - 1
	for c := siteC + 1; c < len(grid[siteR]); c++ {
		if grid[siteR][c] >= siteHeight {
			rightScore = c - siteC
			break
		}
	}

	// fmt.Printf("Scoring %d, %d: up %d, down %d, left %d, right %d\n", siteR, siteC, upScore, downScore, leftScore, rightScore)
	return upScore * downScore * leftScore * rightScore
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

func rtoi(i rune) int {
	n, e := strconv.Atoi(string(i))
	if e != nil {
		panic("Can't parse to digit")
	}
	return n
}
