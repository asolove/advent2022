package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math/bits"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const MAX_TIME = 30

func main() {
	valves, aa_id := read(os.Stdin)
	valvesWorthTurning := len(valves)
	for _, v := range valves {
		if v.rate == 0 {
			valvesWorthTurning--
		}
	}
	fmt.Printf("%v, %d, %v", valves, aa_id, valves[aa_id])

	max := findMaxPressure(valves, aa_id, valvesWorthTurning)
	fmt.Printf("Max: %d\n", max)
}

func findMaxPressure(vs Valves, aa_id int, valvesWorthTurning int) int {
	return search(
		State{loc: aa_id, minute: 0, opened: 0, rate: 0, total: 0},
		func(s *State) bool { return s.minute >= MAX_TIME },
		func(s *State) int { return s.total },
		func(s *State) []*State { return nextStates(s, vs, valvesWorthTurning) },
	)
}

func nextStates(s *State, vs Valves, vwt int) []*State {
	ss := make([]*State, 0)

	// Done as much as we can, give up
	if bits.OnesCount64(s.opened) >= vwt {
		ss = append(ss, &State{
			loc:    s.loc,
			minute: MAX_TIME,
			opened: s.opened,
			rate:   s.rate,
			total:  s.total + (MAX_TIME-s.minute)*s.rate,
		})
		return ss
	}

	// open a valve
	if s.opened&(1<<s.loc) == 0 {
		ss = append(ss, &State{
			loc:    s.loc,
			minute: s.minute + 1,
			opened: s.opened | (1 << s.loc),
			rate:   s.rate + vs[s.loc].rate,
			total:  s.total + s.rate,
		})
	}

	// walk down each path
	for _, loc := range vs[s.loc].next {
		ss = append(ss, &State{
			loc:    loc,
			minute: s.minute + 1,
			opened: s.opened,
			rate:   s.rate,
			total:  s.total + s.rate,
		})
	}
	return ss
}

func search(start State, done func(*State) bool, score func(*State) int, next func(*State) []*State) int {
	max := 0
	i := 0

	var pq PriorityQueue
	pq.Push(&start)
	heap.Init(&pq)

	for !pq.Empty() {
		i++
		if i%100000000 == 0 {
			fmt.Printf("After %d steps, max: %d, candidates: %d\n", i, max, len(pq))
		}

		current := heap.Pop(&pq).(*State)
		if done(current) {
			if current.total > max {
				max = current.total
			}
			continue
		}

		for _, s := range next(current) {
			heap.Push(&pq, s)
		}
	}

	return max
}

type PriorityQueue []*State

func (s *State) Score() int {
	return s.total + ((MAX_TIME - s.minute) * s.rate)
}

func (pq PriorityQueue) Len() int    { return len(pq) }
func (pq PriorityQueue) Empty() bool { return pq.Len() == 0 }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Score() > pq[j].Score()
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

type State struct {
	loc    int
	minute int
	opened uint64
	rate   int
	total  int
	index  int
}

type Valves map[int]*Valve

type Valve struct {
	rate int
	next []int
}

func read(f *os.File) (Valves, int) {
	vs := make(Valves, 0)
	s := bufio.NewScanner(f)
	id := 0
	aa_id := 0
	ids := make(map[string]int)
	re := regexp.MustCompile(`^Valve (\S\S) has flow rate=(\d+); tunnels? leads? to valves? (.*)$`)
	for s.Scan() {
		match := re.FindAllStringSubmatch(s.Text(), -1)[0]
		name := match[1]
		id = intern(name, ids)
		if match[1] == "AA" {
			aa_id = id
		}
		next := make([]int, 0)
		for _, s := range strings.Split(match[3], ", ") {
			next_id := intern(s, ids)
			next = append(next, next_id)
		}
		vs[id] = &Valve{
			rate: atoi(match[2]),
			next: next,
		}
		id++
	}
	return vs, aa_id
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
