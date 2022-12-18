package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	ds := readDroplets(os.Stdin)
	c := countUncoveredFaces(ds)
	fmt.Printf("Uncovered faces: %d\n", c)
}

func countUncoveredFaces(ds map[string]bool) int {
	c := 0
	for k, v := range ds {
		if !v {
			continue
		}
		coords := strings.Split(k, ",")
		for i, _ := range coords {
			if !ds[neighborKey(coords, i, +1)] {
				c += 1
			}
			if !ds[neighborKey(coords, i, -1)] {
				c += 1
			}
		}
	}
	return c
}

func neighborKey(coords []string, i int, di int) string {
	neighbor := make([]string, 3)
	for j, _ := range coords {
		neighbor[j] = coords[j]
	}
	neighbor[i] = fmt.Sprintf("%d", atoi(neighbor[i])+di)
	return strings.Join(neighbor, ",")
}

func readDroplets(f *os.File) map[string]bool {
	ds := make(map[string]bool, 0)
	s := bufio.NewScanner(f)
	for s.Scan() {
		ds[s.Text()] = true
	}
	return ds
}

func atoi(i string) int {
	n, e := strconv.Atoi(i)
	if e != nil {
		panic("Can't parse to digit")
	}
	return n
}
