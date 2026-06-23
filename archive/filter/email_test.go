package filter

import (
	"testing"
)

func TestBlockByEmail(t *testing.T) {
	email := "test@gmail.com"
	if !BlockByEmail(email) {
		t.Errorf("Expected %s to be blocked", email)
	}

	email = "test@yahoo.com"
	if BlockByEmail(email) {
		t.Errorf("Expected %s to be not blocked", email)
	}
}
