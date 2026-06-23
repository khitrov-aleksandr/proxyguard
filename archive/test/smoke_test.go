package main

import "testing"

func TestSmoke(t *testing.T) {
	t.Skip("Skipping smoke test")

	if 5 != 5 {
		t.Error("main test error")
	}
}
