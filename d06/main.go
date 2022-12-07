package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	t := s.Text()

	fmt.Printf("Start of packet marker is at %d\n", findMarker(t, 14))
}

func findMarker(t string, length int) int {
	for i := range t {
		set := map[rune]bool{}
		for j := 0; j < length; j++ {
			set[rune(t[i+j])] = true
		}
		if len(set) == length {
			fmt.Printf("Found set: %v\n", set)
			return i + length
		}
	}

	return 0
}
