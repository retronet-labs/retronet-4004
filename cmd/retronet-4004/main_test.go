package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestReadDigit(t *testing.T) {
	// Cifre intervallate da caratteri non-cifra (spazi, a-capo, lettere).
	r := bufio.NewReader(strings.NewReader("7 5\nx9"))
	want := []uint8{7, 5, 9}
	for i, w := range want {
		got, ok := readDigit(r)
		if !ok {
			t.Fatalf("lettura %d: ok=false inatteso", i)
		}
		if got != w {
			t.Errorf("lettura %d = %d, atteso %d", i, got, w)
		}
	}
	if _, ok := readDigit(r); ok {
		t.Error("a fine input atteso ok=false")
	}
}
