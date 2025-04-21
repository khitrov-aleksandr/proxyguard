package main

import "testing"

func TestMain(t *testing.T) {
	if 5 != 5 {
		t.Error("main test error")
	}
}
