package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
)

func main() {
	monkeys := parseMonkeys(os.Stdin)
	for i := 0; i < 10000; i++ {
		// fmt.Printf("** Round %d **", i)
		runRound(monkeys)
	}

	ins := []int{}
	for _, m := range monkeys {
		// fmt.Printf("Monkey %d inspected %d items\n", _i, m.inspected)
		ins = append(ins, m.inspected)
	}
	sort.Ints(ins)
	fmt.Printf("Monkey business: %d\n", ins[len(ins)-1]*ins[len(ins)-2])
}

func runRound(ms []monkey) {
	for i, m := range ms {
		// fmt.Printf("Monkey %d:\n", i)
		ms[i] = runTurn(ms, m)
	}
}

func runTurn(ms []monkey, m monkey) monkey {
	for _, item := range m.items {
		// fmt.Printf("  Monkey inspects an item with worry level of %d\n", item)
		m.inspected = m.inspected + 1
		item = m.op(item)

		// fmt.Printf("    Worry level is changed to %d\n", item)
		// fmt.Printf("    Monkey gets bored with item. Worry level is divided by 3 to %d\n", item)

		test := m.test(item)
		recipient := m.ifTrue
		if !test {
			recipient = m.ifFalse
		}
		// fmt.Printf("    Current worry level is divisible? %v\n", test)
		// fmt.Printf("    Item %d thrown to monkey %d\n", item, recipient)
		ms[recipient].items = append(ms[recipient].items, item)
	}
	m.items = []int{}
	return m
}

type state struct {
	monkeys []monkey
}

type monkey struct {
	inspected int
	items     []int
	op        func(int) int
	test      func(int) bool
	ifTrue    int
	ifFalse   int
}

func parseMonkeys(f *os.File) []monkey {
	s := bufio.NewScanner(f)
	ms := make([]monkey, 0)

	itemsRe := regexp.MustCompile(`\d+`)
	opRe := regexp.MustCompile(`new = old ([*+]) (old|(\d*))`)
	testRe := regexp.MustCompile(`Test: divisible by (\d*)`)
	throwRe := regexp.MustCompile(`throw to monkey (\d*)`)

	sharedMod := 1

	for s.Scan() {
		m := monkey{}

		// "Monkey N:"

		s.Scan()
		// "  Starting items: 79, 98"
		items := itemsRe.FindAllStringSubmatch(s.Text(), -1)
		m.items = []int{}
		for _, item := range items {
			m.items = append(m.items, atoi(item[0]))
		}

		s.Scan()
		// "  Operation: new = old * 19"
		opMatch := opRe.FindAllStringSubmatch(s.Text(), -1)
		if opMatch[0][1] == "+" {
			if opMatch[0][2] == "old" {
				m.op = func(old int) int {
					return (old + old) % sharedMod
				}
			} else {
				n := atoi(opMatch[0][2])
				m.op = func(old int) int {
					return (old + n) % sharedMod
				}
			}
		} else if opMatch[0][1] == "*" {
			if opMatch[0][2] == "old" {
				m.op = func(old int) int {
					return (old * old) % sharedMod
				}
			} else {
				n := atoi(opMatch[0][2])
				m.op = func(old int) int {
					return (old * n) % sharedMod
				}
			}
		} else {
			panic(fmt.Errorf("Unexpected operation: %v", opMatch[0][1]))
		}

		s.Scan()
		// "  Test: divisible by 23"
		r := testRe.FindStringSubmatch(s.Text())
		div := atoi(r[1])
		sharedMod *= div
		m.test = func(old int) bool {
			return old%div == 0
		}

		s.Scan()
		// "    If true: throw to monkey 2"
		m.ifTrue = atoi(throwRe.FindStringSubmatch(s.Text())[1])

		s.Scan()
		// "    If false: throw to monkey 3"
		m.ifFalse = atoi(throwRe.FindStringSubmatch(s.Text())[1])

		s.Scan()
		// Empty line between monkeys

		ms = append(ms, m)
	}
	return ms
}

func atoi(i string) int {
	n, e := strconv.Atoi(i)
	if e != nil {
		panic("Can't parse to digit")
	}
	return n
}
