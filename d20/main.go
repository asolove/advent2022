package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	state := read(os.Stdin)

	state.Mix()

	var zero *Node
	for _, n := range state.original {
		if n.val == 0 {
			zero = n
			break
		}
	}

	offsets := []int{1000, 2000, 3000}
	sum := 0
	for _, offset := range offsets {
		node := zero
		for i := 0; i < offset; i++ {
			node = node.next
		}
		fmt.Printf("Adding %d\n", node.val)
		sum += node.val
	}
	fmt.Printf("Sum: %d\n", sum)
}

type Node struct {
	val  int
	next *Node
	prev *Node
}
type State struct {
	original []*Node
	head     *Node
	step     int
	size     int
}

func (s *State) String() string {
	r := ""
	node := s.head
	for i := 0; i < s.size; i++ {
		r += fmt.Sprintf("%d, ", node.val)
		node = node.next
	}
	return r + "\n"
}

func (s *State) Mix() {
	for i := 0; i < s.size; i++ {
		s.Step()
		// fmt.Printf("%v\n", s)
	}
}

func (s *State) Step() {
	nodeToMove := s.original[s.step]
	moves := nodeToMove.val % (s.size - 1)
	for moves < 0 {
		moves += s.size - 1
	}
	for i := 0; i < moves; i++ {
		s.Advance(nodeToMove)
	}
	// fmt.Printf("%d moves between %d and %d:\n", nodeToMove.val, nodeToMove.prev.val, nodeToMove.next.val)
	// fmt.Printf("%d moved by %d, resulting in:\n%v\n", nodeToMove.val, moves, s)
	s.step++
}

func (s *State) Advance(node *Node) {
	next := node.next
	prev := node.prev

	node.next = next.next
	next.next.prev = node

	node.prev = next
	next.next = node

	next.prev = prev
	prev.next = next
}

func read(f *os.File) *State {
	r := &State{step: 0, original: make([]*Node, 0)}
	s := bufio.NewScanner(f)
	i := 0
	var prev *Node
	for s.Scan() {
		node := Node{val: atoi(s.Text()), prev: prev, next: nil}

		if prev == nil {
			r.head = &node
		} else {
			prev.next = &node
		}
		r.original = append(r.original, &node)
		r.size++
		prev = &node
		i++
	}
	prev.next = r.head
	r.head.prev = prev

	return r
}

func atoi(i string) int {
	n, e := strconv.Atoi(i)
	if e != nil {
		panic("Can't parse to digit")
	}
	return n
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

// Start: moving -7:
// length: 4, startIndex: 2, endIndex: 3
// 1 2 -7 3
// 1 -7 2 3 <- after 1
// -7 1 2 3 <- after 2
// 1 2 3 -7  <- after 3
// 1 2 -7 3 <- after 4
// 1 -7 2 3 <- after 5
// -7 1 2 3 <- after 6
// 1 2 3 -7 <- after 7, index: 3

// -7 % 4 == 1
