package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	monkeys, rootId := read(os.Stdin)
	fmt.Printf("Root monkey: %v\n", monkeys[rootId])
	fmt.Printf("Root value: %v\n", monkeys[rootId].eval(monkeys))
}

type Monkey struct {
	id    int
	known bool
	val   int
	expr  *Expr
}

func (m *Monkey) eval(env map[int]*Monkey) int {
	if m.known {
		return m.val
	} else {
		lhs := env[m.expr.lhs].eval(env)
		rhs := env[m.expr.rhs].eval(env)
		val := eval(lhs, m.expr.op, rhs)
		m.known = true
		m.val = val
		return val
	}
}

func eval(lhs int, op Op, rhs int) int {
	switch op {
	case PLUS:
		return lhs + rhs
	case MINUS:
		return lhs - rhs
	case TIMES:
		return lhs * rhs
	case DIVIDE:
		return lhs / rhs
	}
	panic(fmt.Sprintf("Bad case in eval: %d", op))
}

type Op int

const (
	PLUS Op = iota
	MINUS
	TIMES
	DIVIDE
)

func OpVal(s string) Op {
	switch s {
	case "+":
		return PLUS
	case "-":
		return MINUS
	case "*":
		return TIMES
	case "/":
		return DIVIDE
	}
	panic("Invalid op type: " + s)
}

type Expr struct {
	lhs int
	rhs int
	op  Op
}

func read(f *os.File) (map[int]*Monkey, int) {
	names := make(map[string]int, 0)
	monkeys := make(map[int]*Monkey, 0)
	s := bufio.NewScanner(f)
	for s.Scan() {
		valRe := regexp.MustCompile(`(\w+): (\d+)`)
		opRe := regexp.MustCompile(`(\w+): (\w+) (.) (\w+)`)

		valMatch := valRe.FindAllStringSubmatch(s.Text(), -1)
		if len(valMatch) > 0 {
			// fmt.Printf("Matching direct val:\n%v\n", valMatch)
			id := intern(valMatch[0][1], names)
			monkey := Monkey{id: id, known: true, val: atoi(valMatch[0][2])}
			monkeys[id] = &monkey
		} else {
			opMatch := opRe.FindAllStringSubmatch(s.Text(), -1)
			// fmt.Printf("Matching op expr:\n%v\n", opMatch)
			id := intern(opMatch[0][1], names)
			lhs := intern(opMatch[0][2], names)
			rhs := intern(opMatch[0][4], names)
			op := OpVal(opMatch[0][3])
			expr := Expr{lhs: lhs, rhs: rhs, op: op}
			monkey := Monkey{id: id, known: false, val: 0, expr: &expr}
			monkeys[id] = &monkey
		}
	}
	return monkeys, names["root"]
}

func intern(name string, names map[string]int) int {
	if v, found := names[name]; found {
		return v
	}
	id := len(names)
	names[name] = id
	return id
}

func atoi(i string) int {
	n, e := strconv.Atoi(i)
	if e != nil {
		panic("Can't parse to digit")
	}
	return n
}
