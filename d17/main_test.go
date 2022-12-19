package main

import (
	"fmt"
	"os"
	"testing"
)

var testWind []int

func init() {
	file, e := os.Open("./data.test")
	if e != nil {
		panic("Didn't find test data")
	}

	testWind = readWind(file)
}

func TestPerf(t *testing.T) {
	run(testWind, GOAL_STEPS/1000)
	fmt.Printf("Done with perf test\n")
}

func TestKnown(t *testing.T) {
	height10 := run(testWind, 10)
	if height10 != 17 {
		t.Errorf("After 10 blocks, expected 17, got %d", height10)
	}
	height2022 := run(testWind, 2022)
	if height2022 != 3068 {
		t.Errorf("After 2022 blocks, expected 3068, got %d", height2022)
	}
}
