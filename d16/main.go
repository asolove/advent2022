package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const MAX_TIME = 30

func main() {
	valves := read(os.Stdin)
	valves_to_open := make([]int, 0)
	for i, v := range valves {
		if v.rate == 0 {
			continue
		}
		valves_to_open = append(valves_to_open, i)
	}

	ds := distances(valves)

	ps := make(chan []int)
	r := make(chan int)
	// go func() { ps <- []int{1, 3, 9, 8, 5, 4}; close(ps) }()
	go permuteValves(valves_to_open, ps)
	go scorePermutations(ps, r, ds, valves)
	go scorePermutations(ps, r, ds, valves)
	go scorePermutations(ps, r, ds, valves)
	go scorePermutations(ps, r, ds, valves)
	go scorePermutations(ps, r, ds, valves)
	go scorePermutations(ps, r, ds, valves)
	go scorePermutations(ps, r, ds, valves)
	go scorePermutations(ps, r, ds, valves)
	go scorePermutations(ps, r, ds, valves)

	for score := range r {
		fmt.Printf("One high score: %d\n", score)
	}
}

func scorePermutations(ps chan []int, r chan int, ds [][]int, valves Valves) {
	max := 0
	var best []int
	for p := range ps {
		score := scorePermutation(p, ds, valves)
		if score > max {
			max = score
			best = p
			fmt.Printf("Best score: %d via %s\n", max, permutationString(best, valves))
		}
	}
	fmt.Printf("Best path: %v\n", permutationString(best, valves))
	r <- max
}

func permutationString(p []int, valves Valves) string {
	s := fmt.Sprintf("%v: ", p)
	for _, id := range p {
		s += valves[id].name + " - "
	}
	return s
}

// optimum for test: DD - BB - JJ - HH - EE - CC
func scorePermutation(perm []int, ds [][]int, valves Valves) int {
	time := 1
	total := 0
	rate := 0
	loc := 0

	for _, goal := range perm {
		// fmt.Printf("Time %d: at %s, releasing %d pressure\n", time, valves[loc].name, rate)
		dtime := ds[loc][goal]
		if time+dtime > MAX_TIME {
			break
		}
		// fmt.Printf("  Walking %dm to %s\n", dtime, valves[goal].name)
		total += rate * dtime
		time += dtime

		// fmt.Printf("  Opening valve %s\n", valves[goal].name)
		rate += valves[goal].rate
		total += rate
		time += 1
		loc = goal
	}

	if time < MAX_TIME {
		total += rate * (MAX_TIME - time)
		time = MAX_TIME
		// fmt.Printf("Time %d: releasing %d pressure\n", time, rate)
	}
	return total
}

func distances(valves Valves) [][]int {
	ds := make([][]int, len(valves))
	for i, v := range valves {
		ds[i] = make([]int, len(valves))
		for j, _ := range valves {
			ds[i][j] = math.MaxInt
		}
		for _, v2 := range v.next {
			ds[i][v2] = 1
		}
		ds[i][i] = 0
	}

	for _, _ = range valves {
		for from, _ := range valves {
			for to, _ := range valves {
				if ds[from][to] <= 2 {
					continue
				}

				for via, _ := range valves {
					if via == from || via == to {
						continue
					}
					d1 := ds[from][via]
					d2 := ds[via][to]
					if d1 == math.MaxInt || d2 == math.MaxInt {
						continue
					}
					vd := d1 + d2
					if vd < ds[from][to] {
						ds[from][to] = vd
					}
				}
			}
		}
	}

	// for fromId, fromV := range valves {
	// 	for toId, toV := range valves {
	// 		fmt.Printf("%s->%s: %d\n", fromV.name, toV.name, ds[fromId][toId])
	// 	}
	// }

	return ds
}

func permuteValves(values []int, r chan []int) {
	i := 0
	for p := make([]int, len(values)); p[0] < len(p); nextPerm(p) {
		if i%10000 == 0 {
			// fmt.Printf("Generating permutation %d\n", i)
		}
		i++
		r <- getPerm(values, p)
	}
	close(r)
}

func nextPerm(p []int) {
	for i := len(p) - 1; i >= 0; i-- {
		if i == 0 || p[i] < len(p)-i-1 {
			p[i]++
			return
		}
		p[i] = 0
	}
}

func getPerm(orig, p []int) []int {
	result := append([]int{}, orig...)
	for i, v := range p {
		result[i], result[i+v] = result[i+v], result[i]
	}
	return result
}

type Valves map[int]*Valve

type Valve struct {
	name string
	rate int
	next []int
}

func read(f *os.File) Valves {
	vs := make(Valves, 0)
	s := bufio.NewScanner(f)
	re := regexp.MustCompile(`^Valve (\S\S) has flow rate=(\d+); tunnels? leads? to valves? (.*)$`)

	// Track mapping of names to ids
	id := 0
	ids := make(map[string]int)

	// Ensure starting point has id 0
	intern("AA", ids)

	for s.Scan() {
		match := re.FindAllStringSubmatch(s.Text(), -1)[0]
		name := match[1]
		id = intern(name, ids)
		next := make([]int, 0)
		for _, s := range strings.Split(match[3], ", ") {
			next_id := intern(s, ids)
			next = append(next, next_id)
		}
		vs[id] = &Valve{
			name: name,
			rate: atoi(match[2]),
			next: next,
		}
		id++
	}
	return vs
}

func intern(name string, lookup map[string]int) int {
	id, found := lookup[name]
	if !found {
		id = len(lookup)
		lookup[name] = id
	}
	return id
}

func atoi(i string) int {
	n, e := strconv.Atoi(i)
	if e != nil {
		panic("Can't parse to digit")
	}
	return n
}
