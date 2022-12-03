package main

import "testing"

func TestBadgeItem(t *testing.T) {
	expected := 'Z'
	r1 := "wMqvLMZHhHMvwLHjbvcjnnSBnvTQFn"
	r2 := "ttgJtRGJQctTZtZT"
	r3 := "CrZsJsPPZsGzwwsLwLmpwMDw"
	actual := badgeItem(r1, r2, r3)

	if actual != expected {
		t.Errorf("Expected %c but got %c", expected, actual)
	}
}

func TestItemPriority(t *testing.T) {
	table := map[rune]int{
		'a': 1,
		'z': 26,
		'A': 27,
		'Z': 52,
	}

	for item, expected := range table {
		if actual := itemPriority(item); actual != expected {
			t.Errorf("ItemPriority(%c): got %d, expected %d", item, actual, expected)
		}
	}
}
