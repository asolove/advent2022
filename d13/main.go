package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	pairs := readPairs(os.Stdin)
	// fmt.Printf("Read pairs:\n%v\n", pairs)

	rightSum := 0

	for i, pair := range pairs {
		// fmt.Printf("Comparing %v to %v\n", pair[0], pair[1])
		if compare(pair[0], pair[1]) > 0 {
			// fmt.Printf("  In right order\n")
			rightSum += i + 1
		}
	}

	fmt.Printf("Sum of correct indices: %d\n", rightSum)
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
