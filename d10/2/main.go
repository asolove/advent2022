package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	sum := 0
	x := 1
	addInTwo := 0
	addInOne := 0
	wait := 0

	image := make([][]bool, 6)
	for row := 0; row < 6; row++ {
		image[row] = make([]bool, 40)
	}

	s := bufio.NewScanner(os.Stdin)
	for i := 1; true; i++ {
		row := (i - 1) / 40
		col := (i - 1) % 40
		if x >= col-1 && x <= col+1 {
			image[row][col] = true
		}

		if i == 20 || (i-20)%40 == 0 {
			fmt.Printf("In cycle %d, signal is %d\n", i, x)
			sum += x * i
		}
		if wait > 0 {
			wait--

		} else {
			ok := s.Scan()
			if !ok {
				break
			}

			words := strings.Split(s.Text(), " ")
			if words[0] == "noop" {
			} else {
				addInTwo = atoi(words[1])
				wait = 1
			}
		}

		x += addInOne
		addInOne = addInTwo
		addInTwo = 0
	}

	for _, row := range image {
		for _, pixel := range row {
			if pixel {
				fmt.Printf("#")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("\n")
	}
}

func atoi(i string) int {
	n, e := strconv.Atoi(i)
	if e != nil {
		panic("Can't parse to digit")
	}
	return n
}
