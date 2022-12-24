package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const STEPS = 18

func main() {
	valley := read(os.Stdin)
	fmt.Printf("Valley:\n%v\n", valley)

	for i := 1; i <= STEPS; i++ {
		valley.Step()
		fmt.Printf("Valley after %d steps:\n%v\n", i, valley)
	}
}

type Valley struct {
	width     int
	height    int
	blizzards []*Blizzard
}

func (v *Valley) Step() {
	for _, b := range v.blizzards {
		b.Step(v.width, v.height)
	}
}

func (v *Valley) Locations() map[string]rune {
	ls := make(map[string]rune)
	for _, b := range v.blizzards {
		k := key(b.x, b.y)
		if _, found := ls[k]; !found {
			ls[k] = b.dir
		} else {
			// FIXME: calculate actual number
			// fmt.Printf("At key=%s, found %c and adding %c\n", k, v, b.dir)
			ls[k] = '2'
		}
	}
	return ls
}

func (v *Valley) String() string {
	s := ""
	ls := v.Locations()
	for y := 0; y < v.height; y++ {
		for x := 0; x < v.width; x++ {
			if v, found := ls[key(x, y)]; found {
				s += string(v)
			} else {
				s += "."
			}
		}
		s += "\n"
	}
	return s
}

// A wall is a blizzard that doesn't move
type Blizzard struct {
	x   int
	y   int
	dir rune
}

func (b *Blizzard) Step(width, height int) {
	switch b.dir {
	case '.':
		return
	case '^':
		if b.y <= 1 {
			b.y = height - 2
		} else {
			b.y -= 1
		}
	case 'v':
		if b.y >= height-2 {
			b.y = 1
		} else {
			b.y += 1
		}
	case '<':
		if b.x <= 1 {
			b.x = width - 2
		} else {
			b.x -= 1
		}
	case '>':
		if b.x >= width-2 {
			b.x = 1
		} else {
			b.x += 1
		}
	}
}

func read(f *os.File) *Valley {
	v := &Valley{width: 0, height: 0, blizzards: make([]*Blizzard, 0)}
	s := bufio.NewScanner(f)
	for y := 0; s.Scan(); y++ {
		if y >= v.height {
			v.height = y + 1
		}
		for x, c := range s.Text() {
			if x >= v.width {
				v.width = x + 1
			}
			if c == '.' {
				continue
			}
			b := &Blizzard{x: x, y: y, dir: c}
			v.blizzards = append(v.blizzards, b)
		}
	}
	return v
}

func coords(k string) (int, int) {
	coords := strings.Split(k, ",")
	return atoi(coords[0]), atoi(coords[1])
}

func key(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

func atoi(i string) int {
	n, e := strconv.Atoi(i)
	if e != nil {
		panic("Can't parse to digit")
	}
	return n
}
