package main

import "testing"

func TestSplitRucksack(t *testing.T) {
  rucksack := "jqHRNqRjqzjGDLGLrsFMfFZSrLrFZsSL";
  left, right := splitRucksack(rucksack);
  expectedLeft := "jqHRNqRjqzjGDLGL";
  expectedRight := "rsFMfFZSrLrFZsSL";

  if left != expectedLeft {
    t.Errorf("got %q, wanted %q", left, expectedLeft);
  }

  if right != expectedRight {
    t.Errorf("got %q, wanted %q", right, expectedRight);
  }
}

func TestInBothSides(t *testing.T) {
  left := "vJrwpWtwJgWr";
  right := "hcsFMMfFFhFp";
  expected := 'p';
  actual := inBothSides(left, right);

  if expected != actual {
    t.Errorf("got %c, expected %c", actual, expected);
  }
}

func TestItemPriority(t *testing.T) {
  table := map[rune]int {
    'a': 1,
    'z': 26,
    'A': 27,
    'Z': 52,
  };

  for item, expected := range table {
    if actual := itemPriority(item); actual != expected {
      t.Errorf("ItemPriority(%c): got %d, expected %d", item, actual, expected);
    }
  }
}
