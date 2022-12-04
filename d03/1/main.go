package main

import (
	"bufio"
	"fmt"
	"os"
)

func splitRucksack(content string) (string, string) {
	splitAt := len(content) / 2
	return content[0:splitAt], content[splitAt:]
}

func byteSet(content string) map[rune]bool {
	s := make(map[rune]bool)
	for _, b := range content {
		s[b] = true
	}
	return s
}

func overlap(s1, s2 map[rune]bool) rune {
	for k := range s1 {
		if s2[k] {
			return k
		}
	}
	return rune(0)
}

func inBothSides(l, r string) rune {
	lSet := byteSet(l)
	rSet := byteSet(r)
	inBoth := overlap(lSet, rSet)
	if inBoth == 0 {
		panic("Character in both sides is invalid.")
	}
	return rune(inBoth)
}

func itemPriority(item rune) int {
	if item >= 'a' && item <= 'z' {
		return int(item) - 96
	} else if item >= 'A' && item <= 'Z' {
		return int(item) - 64 + 26
	}
	return 0
}

func main() {
	totalPriority := 0
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line := s.Text()
		left, right := splitRucksack(line)
		duplicatedItem := inBothSides(left, right)
		totalPriority += itemPriority(duplicatedItem)
	}
	fmt.Printf("Total priority of duplicated items: %d\n", totalPriority)
}
