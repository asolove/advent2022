package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const ROW = 2000000

func main() {
	sensors := readSensors(os.Stdin)
	//fmt.Printf("%v\n", sensors)

	covered := make(map[int]bool)
	for _, s := range sensors {
		addCovered(s, ROW, covered)
	}
	for _, s := range sensors {
		if s.by == ROW {
			delete(covered, s.bx)
		}
	}
	fmt.Printf("Coverd on row %d: %d\n", ROW, len(covered))
	// fmt.Printf("%v\n", covered)
}

// For a sensor and row, marks the spaces within its boundary in the dict
func addCovered(s Sensor, y int, covered map[int]bool) {
	d := abs(s.bx-s.sx) + abs(s.by-s.sy)

	dy := abs(s.sy - y)
	dx := d - dy
	if dx < 0 {
		return
	}
	minX := s.sx - dx
	maxX := s.sx + dx

	for x := minX; x <= maxX; x++ {
		covered[x] = true
	}
}

func readSensors(f *os.File) []Sensor {
	r := make([]Sensor, 0)
	re := regexp.MustCompile(`Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)`)
	s := bufio.NewScanner(f)
	for s.Scan() {
		cs := re.FindAllStringSubmatch(s.Text(), -1)[0]
		r = append(r, Sensor{sx: atoi(cs[1]), sy: atoi(cs[2]), bx: atoi(cs[3]), by: atoi(cs[4])})
	}
	return r
}

type Sensor struct {
	sx int
	sy int
	bx int
	by int
}

func atoi(i string) int {
	n, e := strconv.Atoi(i)
	if e != nil {
		panic("Can't parse to digit")
	}
	return n
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
