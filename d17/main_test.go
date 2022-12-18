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

	run(file)

	fmt.Printf("Done with perf test\n")
}
