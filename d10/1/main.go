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
	s := bufio.NewScanner(os.Stdin)
	for i := 1; true; i++ {
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

	fmt.Printf("Total is %d\n", sum)
}

func atoi(i string) int {
	n, e := strconv.Atoi(i)
	if e != nil {
		panic("Can't parse to digit")
	}
	return n
}
