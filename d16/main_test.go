package main

import (
	"fmt"
	"os"
	"testing"
)

func TestPerf(t *testing.T) {

	file, e := os.Open("./data.test")
	if e != nil {
		t.Errorf("Didn't find data")
	}

	valves, aa_id := read(file)
	valvesWorthTurning := len(valves)
	for _, v := range valves {
		if v.rate == 0 {
			valvesWorthTurning--
		}
	}
	fmt.Printf("%v, %d, %v", valves, aa_id, valves[aa_id])

	max := findMaxPressure(valves, aa_id, valvesWorthTurning)
	fmt.Printf("Max: %d\n", max)
}
