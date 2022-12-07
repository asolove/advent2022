package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	fileTree := readFileTree(os.Stdin)
	fmt.Printf("Answer: %v\n", calculateAnswer(fileTree))
}

type FileTree struct {
	parent      *FileTree
	directories map[string]*FileTree
	files       map[string]int
}

func NewFileTree(parent *FileTree) *FileTree {
	var ft FileTree
	ft.parent = parent
	ft.directories = make(map[string]*FileTree)
	ft.files = make(map[string]int)
	return &ft
}

func (ft *FileTree) Size() int {
	size := 0
	for _, dir := range ft.directories {
		size += dir.Size()
	}
	for _, fileSize := range ft.files {
		size += fileSize
	}
	return size
}

const TOTAL_SPACE = 70000000
const SPACE_NEEDED = 30000000

func calculateAnswer(ft *FileTree) int {
	rootSize := ft.Size()
	toFree := rootSize - (TOTAL_SPACE - SPACE_NEEDED)

	return smallestSizeGreaterThan(ft, toFree)
}

func smallestSizeGreaterThan(ft *FileTree, toFree int) int {
	result := math.MaxInt
	pwdSize := ft.Size()
	if pwdSize > toFree {
		result = min(result, pwdSize)
	}
	for _, dir := range ft.directories {
		result = min(result, smallestSizeGreaterThan(dir, toFree))
	}

	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func readFileTree(f *os.File) *FileTree {
	s := bufio.NewScanner(os.Stdin)
	ft := NewFileTree(nil)
	pwd := ft
	for s.Scan() {
		t := s.Text()
		switch {
		case t[0:1] != "$":
			words := strings.Split(t, " ")
			if words[0] == "dir" {
				pwd.directories[words[1]] = NewFileTree(pwd)
			} else {
				pwd.files[words[1]] = atoi(words[0])
			}
		case t[1:4] == " cd":
			dirName := t[5:]
			if dirName == "/" {
				pwd = ft
			} else if dirName == ".." {
				pwd = pwd.parent
			} else {
				pwd = pwd.directories[dirName]
			}
			if pwd == nil {
				panic("Present working directory is nil after cd")
			}
		case t[1:4] == " ls":
			continue
		}
	}
	return ft
}

func atoi(i string) int {
	n, e := strconv.Atoi(i)
	if e != nil {
		panic("Can't parse to digit")
	}
	return n
}
