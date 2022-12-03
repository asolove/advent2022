package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func topThree(challenger, a, b, c int) (int, int, int) {
	if challenger > c {
		if challenger > b {
			if challenger > a {
				return challenger, a, b
			} else {
				return a, challenger, b
			}
		} else {
			return a, b, challenger
		}
	}
	return a, b, c
}

func main() {
	var a, b, c int
	var currentTotal int

	// Would be nice to parse at two layers: first to find the per-elf groups and then to consume the numbers inside them.
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		if t := s.Text(); t == "" {
			a, b, c = topThree(currentTotal, a, b, c)
			currentTotal = 0
		} else {
			if n, err := strconv.Atoi(t); err == nil {
				currentTotal += n
			}
		}
	}
	// Would be nice to not have to run this separately for final round without newline
	a, b, c = topThree(currentTotal, a, b, c)

	fmt.Printf("%d\n", a+b+c)
	return
}
