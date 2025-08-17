package main

import "testing"

func TestSample(t *testing.T) {
	if 2+2 != 4 {
		t.Fatal("math is broken")
	}
}