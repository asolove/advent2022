package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	monkeys, rootId, humanId := read(os.Stdin)
	monkeys[humanId].human = true
	lhs := monkeys[monkeys[rootId].expr.lhs].eval(monkeys)
	fmt.Printf("Finding LHS: %v\n", lhs)
	rhs := monkeys[monkeys[rootId].expr.rhs].eval(monkeys)
	fmt.Printf("Finding RHS: %v\n", rhs)

	eq := Equality{lhs: lhs, rhs: rhs}
	eq.Rotate()
	fmt.Printf("humn = %d\n", eq.rhs.val)
}

type Equality struct {
	lhs *Expr
	rhs *Expr
}

func (eq *Equality) Rotate() {
	// Move unknown to LHS
	if eq.lhs.known {
		eq.lhs, eq.rhs = eq.rhs, eq.lhs
	}
	for eq.lhs.variable == "" {
		fmt.Printf("Before rotation:\n  %v\n", eq)
		eq.RotateOnce()
		fmt.Printf("After rotation:\n  %v\n\n", eq)
	}
}

func (eq *Equality) RotateOnce() {
	// Unknown in left sub-branch: (X + 2) = 10
	if eq.lhs.rhs.known {
		eq.rhs = NewExpr(eq.rhs, eq.lhs.op.Inverse(), eq.lhs.rhs)
		eq.lhs = eq.lhs.lhs
	} else {
		// (2 + x) == 10  -> x = 10 - 2
		// (2 - X) == 10  -> x = 2 - 10
		// (10 * X) == 20  -> x = 20 / 10
		// (10 / X) == 2  -> x = 10 / 2
		switch eq.lhs.op {
		case PLUS, TIMES:
			eq.rhs = NewExpr(eq.rhs, eq.lhs.op.Inverse(), eq.lhs.lhs)
		case MINUS, DIVIDE:
			eq.rhs = NewExpr(eq.lhs.lhs, eq.lhs.op, eq.rhs)
		}
		eq.lhs = eq.lhs.rhs
	}
}

func (eq *Equality) String() string {
	return fmt.Sprintf("%v == %v", eq.lhs, eq.rhs)
}

type Expr struct {
	known    bool
	val      int
	variable string
	op       Op
	lhs      *Expr
	rhs      *Expr
}

func NewExpr(lhs *Expr, op Op, rhs *Expr) *Expr {
	known := lhs.known && rhs.known
	val := 0
	if known {
		val = eval(lhs.val, op, rhs.val)
	}
	return &Expr{
		lhs:   lhs,
		op:    op,
		rhs:   rhs,
		known: known,
		val:   val,
	}
}

func (e *Expr) String() string {
	if e.variable != "" {
		return e.variable
	} else if e.known {
		return fmt.Sprintf("%d", e.val)
	} else {
		return fmt.Sprintf("(%v %v %v)", e.lhs.String(), e.op, e.rhs.String())
	}
}

type Monkey struct {
	id    int
	known bool
	human bool
	val   *Expr
	expr  *SimpleExpr
}

func (m *Monkey) eval(env map[int]*Monkey) *Expr {
	if m.human {
		return &Expr{known: false, variable: "humn"}
	} else if m.known {
		return m.val
	} else {
		lhs := env[m.expr.lhs].eval(env)
		rhs := env[m.expr.rhs].eval(env)
		if lhs.known && rhs.known {
			val := eval(lhs.val, m.expr.op, rhs.val)
			expr := &Expr{known: true, val: val}
			m.known = true
			m.val = expr
			return expr
		} else {
			expr := &Expr{known: false, op: m.expr.op, lhs: lhs, rhs: rhs}
			m.known = true
			m.val = expr
			return expr
		}
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

func (op Op) Inverse() Op {
	switch op {
	case PLUS:
		return MINUS
	case MINUS:
		return PLUS
	case TIMES:
		return DIVIDE
	case DIVIDE:
		return TIMES
	}

	panic(fmt.Sprintf("Bad case in op.Inverse: %d", op))
}

func (op Op) String() string {
	switch op {
	case PLUS:
		return "+"
	case MINUS:
		return "-"
	case TIMES:
		return "*"
	case DIVIDE:
		return "/"
	}
	panic(fmt.Sprintf("Invalid op int: %d", op))
}

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

type SimpleExpr struct {
	lhs int
	rhs int
	op  Op
}

func read(f *os.File) (map[int]*Monkey, int, int) {
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
			monkey := Monkey{id: id, known: true, val: &Expr{known: true, val: atoi(valMatch[0][2])}}
			monkeys[id] = &monkey
		} else {
			opMatch := opRe.FindAllStringSubmatch(s.Text(), -1)
			// fmt.Printf("Matching op expr:\n%v\n", opMatch)
			id := intern(opMatch[0][1], names)
			lhs := intern(opMatch[0][2], names)
			rhs := intern(opMatch[0][4], names)
			op := OpVal(opMatch[0][3])
			expr := SimpleExpr{lhs: lhs, rhs: rhs, op: op}
			monkey := Monkey{id: id, known: false, expr: &expr}
			monkeys[id] = &monkey
		}
	}
	return monkeys, names["root"], names["humn"]
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
