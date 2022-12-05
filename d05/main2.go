package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	s := bufio.NewScanner(os.Stdin)
	data := parseData(s)

	stacks := data.stacks
	fmt.Printf("after parsing:\n%v", printStacks(stacks))
	for _, step := range data.steps {
		stacks = runStep(step, stacks)
		fmt.Printf("after step:\n%v", printStacks(stacks))
	}
	fmt.Printf("top of stacks: %v\n", topOfStacks(stacks))
}

type stacks [][]rune

type step struct {
	from, to, count int
}

type data struct {
	stacks stacks
	steps  []step
}

func runStep(step step, stacks stacks) stacks {
	from := stacks[step.from]
	stacks[step.to] = append(stacks[step.to], from[len(from)-step.count:]...)
	stacks[step.from] = from[:len(from)-step.count]
	return stacks
}

func topOfStacks(stacks stacks) string {
	r := []rune{}
	for _, val := range stacks[1:] {
		r = append(r, val[len(val)-1])
	}
	return string(r)
}

func printStacks(stacks stacks) string {
	r := "\n"
	for i, val := range stacks[1:] {
		r += fmt.Sprintf("%d\t%v\n", i, string(val))
	}
	return r
}

func parseData(s *bufio.Scanner) data {
	stacks := parseStacks(s)
	s.Scan() // empty line
	steps := parseSteps(s)

	return data{stacks: stacks, steps: steps}
}

func parseStacks(scanner *bufio.Scanner) stacks {
	var s stacks
top:
	for scanner.Scan() {
		t := scanner.Text()
		for i := 1; i < 10; i++ {
			idx := 4*(i-1) + 1
			if idx >= len(t) {
				break
			}
			r := rune(t[idx])

			// Label line below stacks
			if r == '1' {
				break top
			}
			// Empty space
			if r == ' ' {
				continue
			}

			for len(s) <= i {
				s = append(s, make([]rune, 0))
			}

			s[i] = append([]rune{r}, s[i]...)
		}
	}
	return s
}

func parseSteps(s *bufio.Scanner) []step {
	var r []step
	for s.Scan() {
		words := strings.Split(s.Text(), " ")
		r = append(r, step{
			from:  atoi(words[3]),
			to:    atoi(words[5]),
			count: atoi(words[1]),
		})
	}
	return r
}

// Gotta find a fread or convenient parser library
func atoi(i string) int {
	n, e := strconv.Atoi(i)
	if e != nil {
		panic("Can't parse to digit")
	}
	return n
}
