package main

import "testing"

func TestFindMarker(t *testing.T) {
	table := map[string]int{
		"bvwbjplbgvbhsrlpgdmjqwftvncz":      5,
		"nppdvjthqldpwncqszvftbrmjlhg":      6,
		"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg": 10,
		"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw":  11,
	}

	for item, expected := range table {
		if actual := findMarker(item, 4); actual != expected {
			t.Errorf("findMarker(%s): got %d, expected %d", item, actual, expected)
		}
	}
}

func TestFindMessage(t *testing.T) {
	table := map[string]int{
		"mjqjpqmgbljsphdztnvjfqwrcgsmlb":    19,
		"bvwbjplbgvbhsrlpgdmjqwftvncz":      23,
		"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg": 29,
		"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw":  26,
	}

	for item, expected := range table {
		if actual := findMarker(item, 14); actual != expected {
			t.Errorf("findMarker(%s): got %d, expected %d", item, actual, expected)
		}
	}
}
