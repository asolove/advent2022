package main

import (
  "bufio"
  "fmt"
  "os"
)

func byteSet(content string) map[rune]bool {
  s := make(map[rune]bool);
  for _, b := range content {
    s[b] = true;
  }
  return s;
}

func intersect(s1, s2 map[rune]bool) map[rune]bool {
  r := make(map[rune]bool);
  for k, _ := range s1 {
    if(s2[k]) {
      r[k] = true;
    }
  }
  return r;
}

func badgeItem(r1, r2, r3 string) rune {
  s1 := byteSet(r1);
  s2 := byteSet(r2);
  s3 := byteSet(r3);

  for k, _ := range intersect(s1, intersect(s2, s3)) {
    return k;
  }
  return rune(0);
}

func itemPriority(item rune) int {
  if item >= 'a' && item <= 'z' {
    return int(item) - 96;
  } else if item >= 'A' && item <= 'Z' {
    return int(item) - 64 + 26;
  }
  return 0;
}

func main() {
  totalPriority := 0;
	s := bufio.NewScanner(os.Stdin);
	for s.Scan() {
    r1 := s.Text();
    s.Scan();
    r2 := s.Text();
    s.Scan();
    r3 := s.Text();

    item := badgeItem(r1, r2, r3);
    totalPriority += itemPriority(item);
  }
  fmt.Printf("Total priority of duplicated items: %d", totalPriority);
}

