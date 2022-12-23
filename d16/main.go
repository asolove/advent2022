package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"math/bits"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const MAX_TIME = 30

func main() {
	valves := read(os.Stdin)
	valvesToOpen := make([]int, 0)
	for i, v := range valves {
		if v.rate == 0 {
			continue
		}
		valvesToOpen = append(valvesToOpen, i)
	}

	ds := distances(valves)
	best := findBest(valves, valvesToOpen, ds)
	fmt.Printf("Best: %v", best)
}

func findBest(valves Valves, valvesToOpen []int, ds [][]int) int {
	best := 0

	var pq PriorityQueue
	startState := State{}
	pq.Push(&startState)
	heap.Init(&pq)

	for !pq.Empty() {
		current := pq.Pop().(*State)
		if doneState(current, len(valvesToOpen)) {
			score := scoreState(current)
			if score > best {
				best = score
				fmt.Printf("New best: %d\n", best)
			}
		} else {
			next := nextStates(current, valves, valvesToOpen, ds)
			for _, ns := range next {
				heap.Push(&pq, ns)
			}
		}
		// if current is done
		// compare score

		// else
		// generate list of next states
		// push them into pq
	}
	return best
}

func nextStates(s *State, valves Valves, valvesToOpen []int, ds [][]int) []*State {
	nexts := make([]*State, 0)

	// FIXME: filter to just unopened ones
	nextValves := valvesToOpen

	for _, nv := range nextValves {
		openedWithNext := (1 << nv) | s.opened
		if openedWithNext == s.opened {
			continue
		}
		dt := ds[s.loc][nv]
		next := State{
			time:   s.time + dt + 1,
			opened: openedWithNext,
			rate:   s.rate + valves[nv].rate,
			total:  s.total + s.rate*(dt+1),
			loc:    nv,
		}
		next.gScore = scoreState(&next)
		nexts = append(nexts, &next)
	}

	return nexts
}

func scoreState(s *State) int {
	return s.total + (s.rate * (MAX_TIME - s.time))
}

func doneState(s *State, valveCount int) bool {
	return s.time >= MAX_TIME || bits.OnesCount64(s.opened) >= valveCount
}

type State struct {
	gScore int
	index  int

	time   int
	total  int
	rate   int
	loc    int
	opened uint64
}

type PriorityQueue []*State

func (pq PriorityQueue) Len() int    { return len(pq) }
func (pq PriorityQueue) Empty() bool { return pq.Len() == 0 }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].gScore < pq[j].gScore
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*State)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
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
