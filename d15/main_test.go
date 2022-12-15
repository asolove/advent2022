package main

import (
	"fmt"
	"os"
	"testing"
)

func TestPerf(t *testing.T) {

	file, e := os.Open("./data.real")
	if e != nil {
		t.Errorf("Didn't find data")
	}
	sensors := readSensors(file)
	//fmt.Printf("%v\n", sensors)

	spot := findUncoveredSpot(sensors)
	fmt.Printf("Spot: %v\n", spot)
}
