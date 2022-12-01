package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	var currentTotal, maxTotal int
	var currentIndex, maxIndex int

	s := bufio.NewScanner(os.Stdin);
	for s.Scan() {
		if t := s.Text(); t == "" {
			if currentTotal > maxTotal {
				maxTotal = currentTotal;
				maxIndex = currentIndex;
			}
			currentTotal = 0;
			currentIndex++;
		} else {
			if n, err := strconv.Atoi(t); err == nil {
				currentTotal += n;
			}
		}

	}

	fmt.Printf("%d %d\n", maxTotal, maxIndex);
}
