package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strconv"
)

const (
	MIN_X = 0
	MAX_X = 4000000
	MIN_Y = 0
	MAX_Y = 4000000
)

func main() {
	sensors := readSensors(os.Stdin)
	//fmt.Printf("%v\n", sensors)

	findUncoveredSpot(sensors)
}

func findUncoveredSpot(sensors []Sensor) {
	allBits := bitString(MIN_X, MAX_X)

	ys := make(chan int, 20)
	stop := make(chan bool)

	for i := 0; i < 20; i++ {
		go func() {
			for {
				select {
				case y := <-ys:
					r := checkRow(sensors, y, allBits)
					if r {
						stop <- true
					}
				case <-stop:
					return
				}
			}
		}()
	}

	go func() {
		for y := MIN_Y; y <= MAX_Y; y++ {
			ys <- y
		}
	}()

	s := <-stop
	fmt.Printf("Stopped %d", s)
}

func checkRow(sensors []Sensor, y int, allBits *big.Int) bool {
	covered := new(big.Int).Set(allBits)
	if y%1000 == 0 {
		fmt.Printf("Considering row %d\n", y)
	}

	for _, s := range sensors {
		addCovered(s, y, covered)
		if covered.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	if covered.Cmp(big.NewInt(0)) > 0 {
		for x := MIN_X; x <= MAX_X; x++ {
			if covered.Bit(x) == 1 {
				fmt.Printf("Found at %d, %d = %d\n", x, y, x*MAX_X+y)
				return true
			}
		}
		fmt.Printf("  Found in row %d but no x found.d\n", y)
	}
	return false
}

// Makes a bit string with values [min, max] set to 1, and then trailing 0s
func bitString(minSet, maxSet int) *big.Int {
	s := big.NewInt(1)
	s.Lsh(s, uint(maxSet-minSet+1))
	s.Sub(s, big.NewInt(1))
	s.Lsh(s, uint(minSet))
	return s
}

// For a sensor and row, marks the spaces within its boundary in the dict
func addCovered(s Sensor, y int, covered *big.Int) {
	d := abs(s.bx-s.sx) + abs(s.by-s.sy)

	dy := abs(s.sy - y)
	dx := d - dy
	if dx < 0 {
		return
	}
	minX := max(MIN_X, s.sx-dx)
	maxX := min(MAX_X, s.sx+dx)

	mask := bitString(minX, maxX)

	covered.AndNot(covered, mask)
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

func min(i, j int) int {
	if i > j {
		return j
	}
	return i
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}
