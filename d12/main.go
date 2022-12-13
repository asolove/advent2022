package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	m := readMap(os.Stdin)
	g := makeGraph(m)
	// d := g.nodes[20][135].shortestPath()
	d := shortestPath(g)
	fmt.Printf("Shortest path: %d\n", d)
}

func shortestPath(g *Graph) int {
	best := math.MaxInt
	for _, start := range g.starts {
		var pq PriorityQueue
		start._shortestPath = 0
		pq.Push(&PathItem{fScore: 0, gScore: guessDist(start, g.end), node: start})

		for !pq.Empty() {
			current := pq.Pop()
			if current.node == g.end {
				if current.fScore < best {
					best = current.fScore
				}
				continue
			}

			for _, n2 := range current.node.links {
				if n2._shortestPath > current.fScore+1 {
					fmt.Printf("Considering %d %d\n", n2.row, n2.col)
					n2._shortestPath = current.fScore + 1
					pq.Push(&PathItem{
						fScore: n2._shortestPath,
						gScore: n2._shortestPath + guessDist(n2, g.end),
						node:   n2,
					})
				}
			}
		}
	}
	return best
}

func reset(g *Graph) {
	for _, row := range g.nodes {
		for _, node := range row {
			node._shortestPath = math.MaxInt
		}
	}
}

func guessDist(from, to *Node) int {
	return (to.height - from.height) + abs(to.col-from.col) + abs(to.row-from.row)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type PathItem struct {
	fScore int
	gScore int
	node   *Node
	index  int
}

type PriorityQueue []*PathItem

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

func (pq *PriorityQueue) Push(x *PathItem) {
	n := len(*pq)
	item := x
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() *PathItem {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

type Graph struct {
	nodes  [][]*Node
	starts []*Node
	end    *Node
}

type Node struct {
	row           int
	col           int
	height        int
	links         []*Node
	_shortestPath int
	_step         *Node
}

func makeGraph(m [][]rune) *Graph {
	g := Graph{nodes: make([][]*Node, len(m))}
	for r, row := range m {
		g.nodes[r] = make([]*Node, len(row))
		for c, height := range row {
			n := Node{
				row:           r,
				col:           c,
				_shortestPath: math.MaxInt,
			}
			if height == 'S' || height == 'a' {
				g.starts = append(g.starts, &n)
				height = 'a'
			}
			if height == 'E' {
				g.end = &n
				height = 'z'
			}
			n.height = int(height)
			n.links = make([]*Node, 0)
			g.nodes[r][c] = &n
		}
	}

	for r, row := range g.nodes {
		for c, _ := range g.nodes[r] {
			n := g.nodes[r][c]

			if r > 0 && g.nodes[r-1][c].height <= n.height+1 {
				n.links = append(n.links, g.nodes[r-1][c])
			}
			if r < len(g.nodes)-1 && g.nodes[r+1][c].height <= n.height+1 {
				n.links = append(n.links, g.nodes[r+1][c])
			}
			if c > 0 && g.nodes[r][c-1].height <= n.height+1 {
				n.links = append(n.links, g.nodes[r][c-1])
			}
			if c < len(row)-1 && g.nodes[r][c+1].height <= n.height+1 {
				n.links = append(n.links, g.nodes[r][c+1])
			}
		}
	}

	return &g
}

func readMap(f *os.File) [][]rune {
	m := make([][]rune, 0)
	s := bufio.NewScanner(f)
	for s.Scan() {
		m = append(m, []rune(s.Text()))
	}
	return m
}
