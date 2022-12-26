package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	runeLines := read(os.Stdin)
	numbers := []int{}
	for _, runes := range runeLines {
		numbers = append(numbers, eval(runes))
	}
	fmt.Printf("Read numbers: %v\n", numbers)
	sum := 0
	for _, n := range numbers {
		sum += n
	}
	fmt.Printf("Result is %d or %v\n", sum, string(lave(sum)))

	fmt.Printf("Expected: 2=-1=0 %d, got 1022=0 %d", eval([]rune("2=-1=0")), eval([]rune("1022=0")))
}

func lave(num int) []rune {
	r := []rune{}
	mod := 5
	for num != 0 {
		fmt.Printf("Trying to print %d\n", num)
		digit := num % mod
		c := ' '
		switch digit {
		case 0:
			c = '0'
		case 1:
			c = '1'
			num -= 1
		case 2:
			c = '2'
			num -= 2
		case 3:
			c = '='
			num += 2
		case 4:
			c = '-'
			num += 1
		default:
			fmt.Printf("Not sure what to do with mod %d\n", digit)
		}
		num = num / 5
		r = append([]rune{c}, r...)
	}
	return r
}

func eval(runes []rune) int {
	r := 0
	for _, c := range runes {
		r = r*5 + val(c)
	}
	return r
}

func val(c rune) int {
	switch c {
	case '=':
		return -2
	case '-':
		return -1
	case '0':
		return 0
	case '1':
		return 1
	case '2':
		return 2
	}
	panic(fmt.Sprintf("Invalid rune to eval: %c", c))
}

func read(f *os.File) [][]rune {
	r := make([][]rune, 0)
	s := bufio.NewScanner((f))
	for s.Scan() {
		r = append(r, []rune(s.Text()))
	}
	return r
}
