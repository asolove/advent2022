package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	pairs := readPairs(os.Stdin)
	// fmt.Printf("Read pairs:\n%v\n", pairs)

	d1_t := &Tree{Node{n: 2, t: nil}}
	d1 := Tree{Node{t: d1_t}}
	d2_t := &Tree{Node{n: 6, t: nil}}
	d2 := Tree{Node{t: d2_t}}

	packets := []Tree{d1, d2}
	for _, pair := range pairs {
		packets = append(packets, pair[0], pair[1])
	}

	sort.SliceStable(packets, func(i, j int) bool {
		return compare(packets[i], packets[j]) > 0
	})

	key := 1
	for i, packet := range packets {
		if len(packet) == 1 && (packet[0].t == d1_t || packet[0].t == d2_t) {
			key *= i + 1
		}
	}

	fmt.Printf("Key:\n%v\n", key)
}

// Returns 1 if in right order, 0 if tied, -1 if wrong
func compare(t1, t2 Tree) int {
	for i, n1 := range t1 {
		if i >= len(t2) {
			return -1
		}
		if c := compareNodes(n1, t2[i]); c != 0 {
			return c
		}
	}
	if len(t2) > len(t1) {
		return 1
	}
	return 0
}

func compareNodes(n1, n2 Node) int {
	if n1.t == nil && n2.t == nil {
		return compareInts(n1.n, n2.n)
	}
	if n1.t == nil {
		return compare(Tree{n1}, *n2.t)
	}
	if n2.t == nil {
		return compare(*n1.t, Tree{n2})
	}
	return compare(*n1.t, *n2.t)
}

func compareInts(n1, n2 int) int {
	if n1 < n2 {
		return 1
	}
	if n1 == n2 {
		return 0
	}
	return -1
}

type Tree []Node
type Node struct {
	n int
	t *Tree
}

func (t Tree) String() string {
	var sb strings.Builder
	sb.WriteRune('[')
	for i, n := range t {
		sb.WriteString(fmt.Sprint(n))
		if i < len(t)-1 {
			sb.WriteRune(',')
		}
	}
	sb.WriteRune(']')
	return sb.String()
}

func (n Node) String() string {
	// 0 is empty value, not in real data
	if n.t != nil {
		return fmt.Sprintf("%v", n.t)
	} else {
		return fmt.Sprintf("%d", n.n)
	}
}

func writePairs(pairs [][2]Tree) {
	f, e := os.Create("data.out")
	if e != nil {
		panic(e)
	}
	for _, pair := range pairs {
		fmt.Fprintf(f, "%v\n%v\n\n", pair[0], pair[1])
	}
}

func readPairs(f *os.File) [][2]Tree {
	pairs := make([][2]Tree, 0)
	s := bufio.NewScanner(f)
	for s.Scan() {
		t1, _ := readTree(s.Text())
		s.Scan()
		t2, _ := readTree(s.Text())
		s.Scan()

		pairs = append(pairs, [2]Tree{t1, t2})
	}
	return pairs
}

func readTree(input string) (Tree, string) {
	// fmt.Printf("readTree for %s\n", input)

	if input[0] != '[' {
		panic("Bad input to readTree")
	}
	// consume open brace
	input = input[1:]

	t := make(Tree, 0)
	var n Node
	for input[0] != ']' {
		n, input = readNode(input)
		t = append(t, n)
		for input[0] == ' ' || input[0] == ',' {
			input = input[1:]
		}
	}
	return t, input[1:] // consume closing brace
}

func readNode(input string) (Node, string) {
	// fmt.Printf("readNode for %s\n", input)
	if unicode.IsDigit(rune(input[0])) {
		ds := ""
		for i := 0; i < len(input) && unicode.IsDigit(rune(input[i])); i++ {
			ds += string(input[i])
		}
		return Node{n: atoi(ds)}, input[len(ds):]
	} else {
		tree, rest := readTree(input)
		return Node{t: &tree}, rest
	}
}

func atoi(i string) int {
	n, e := strconv.Atoi(string(i))
	if e != nil {
		panic("Can't parse to digit")
	}
	return n
}
