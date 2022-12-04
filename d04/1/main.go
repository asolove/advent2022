package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Gotta find a fread or convenient parser library
func atoi(i string) int {
	n, e := strconv.Atoi(i)
	if e != nil {
		panic("Can't parse to digit")
	}
	return n
}

func parse(line string) (int, int, int, int) {
	parts := strings.Split(line, ",")
	a := strings.Split(parts[0], "-")
	b := strings.Split(parts[1], "-")
	return atoi(a[0]), atoi(a[1]), atoi(b[0]), atoi(b[1])
}

func contains(aMin, aMax, bMin, bMax int) bool {
	if aMin <= bMin && aMax >= bMax {
		return true
	}
	if bMin <= aMin && bMax >= aMax {
		return true
	}
	return false
}

func overlaps(aMin, aMax, bMin, bMax int) bool {
	return (aMax >= bMin && aMin <= bMax) || (bMax >= aMin && bMin <= aMax)
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	count := 0
	for s.Scan() {
		aMin, aMax, bMin, bMax := parse(s.Text())
		if overlaps(aMin, aMax, bMin, bMax) {
			count += 1
		}
	}
	fmt.Printf("%d contains", count)
}
