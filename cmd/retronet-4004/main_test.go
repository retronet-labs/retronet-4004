package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestReadKey(t *testing.T) {
	// Cifre e operatori, intervallati da caratteri da saltare (spazi, a-capo).
	r := bufio.NewReader(strings.NewReader("7 +\n5 - * / = 9"))
	want := []uint8{7, 10, 5, 11, 12, 13, 14, 9} // 7 + 5 - * / = 9
	for i, w := range want {
		got, ok := readKey(r)
		if !ok {
			t.Fatalf("lettura %d: ok=false inatteso", i)
		}
		if got != w {
			t.Errorf("lettura %d = %d, atteso %d", i, got, w)
		}
	}
	if _, ok := readKey(r); ok {
		t.Error("a fine input atteso ok=false")
	}
}
