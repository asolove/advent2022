package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	ds := readDroplets(os.Stdin)
	steam := findSteam(ds)

	c := countSteamedFaces(ds, steam)
	fmt.Printf("Ds: %d. Steam: %d\n", len(ds), len(steam))
	fmt.Printf("Uncovered faces: %d\n", c)
}

func countSteamedFaces(ds map[string]bool, steam map[string]bool) int {
	c := 0
	for k, v := range ds {
		if !v {
			continue
		}
		nks := neighborKeys(k)
		for _, nk := range nks {
			if !ds[nk] {
				if steam[nk] {
					c += 1
				}
			}
		}
	}
	return c
}

func neighborKeys(k string) []string {
	r := make([]string, 0)
	cs := coords(k)
	for i := range cs {
		copy := [3]int{cs[0], cs[1], cs[2]}
		copy[i] -= 1
		r = append(r, key(copy))
		copy[i] += 2
		r = append(r, key(copy))
	}

	return r
}

func findSteam(ds map[string]bool) map[string]bool {
	steam := make(map[string]bool, 0)
	bounds := findBounds(ds)
	fmt.Printf("%v\n", bounds)

	var propagateSteam func(x, y, z int)
	propagateSteam = func(x, y, z int) {
		k := key([3]int{x, y, z})
		if outOfBounds(x, y, z, bounds) {
			return
		}
		if ds[k] {
			return
		}
		if steam[k] {
			return
		}
		steam[k] = true
		propagateSteam(x-1, y, z)
		propagateSteam(x+1, y, z)
		propagateSteam(x, y-1, z)
		propagateSteam(x, y+1, z)
		propagateSteam(x, y, z-1)
		propagateSteam(x, y, z+1)
	}
	// Steam source is at bottom-most corner in every direction
	// and due to bounds having a border space can definitely flow around
	// each dimension of the surface area.
	propagateSteam(bounds[0], bounds[2], bounds[4])

	return steam
}

func outOfBounds(x, y, z int, bounds [6]int) bool {
	if x < bounds[0] || x > bounds[1] {
		return true
	}
	if y < bounds[2] || y > bounds[3] {
		return true
	}
	if z < bounds[4] || z > bounds[5] {
		return true
	}
	return false
}

// Return a bounding rect around the lava,
// with one empty boundary space in each direction for steam to flow
func findBounds(ds map[string]bool) [6]int {
	bounds := [6]int{math.MaxInt, math.MinInt, math.MaxInt, math.MinInt, math.MaxInt, math.MinInt}
	for key, val := range ds {
		if !val {
			continue
		}
		cs := coords(key)
		if cs[0] < bounds[0] {
			bounds[0] = cs[0]
		} else if cs[0] > bounds[1] {
			bounds[1] = cs[0]
		}
		if cs[1] < bounds[2] {
			bounds[2] = cs[1]
		} else if cs[1] > bounds[3] {
			bounds[3] = cs[1]
		}
		if cs[2] < bounds[4] {
			bounds[4] = cs[2]
		} else if cs[2] > bounds[5] {
			bounds[5] = cs[2]
		}
	}
	bounds[0] -= 1
	bounds[1] += 1
	bounds[2] -= 1
	bounds[3] += 1
	bounds[4] -= 1
	bounds[5] += 1
	return bounds
}

func key(coords [3]int) string {
	return strings.Join([]string{strconv.Itoa(coords[0]), strconv.Itoa(coords[1]), strconv.Itoa(coords[2])}, ",")
}

func coords(key string) [3]int {
	ss := strings.Split(key, ",")
	return [3]int{atoi(ss[0]), atoi(ss[1]), atoi(ss[2])}
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
