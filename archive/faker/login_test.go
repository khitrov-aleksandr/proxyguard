package faker

import (
	"testing"
)

func TestGetLoginResponse(t *testing.T) {
	login := GetLoginResponse()
	if login.Success == false {
		t.Errorf("Expected success to be true")
	}

	if login.DelaySec != 0 {
		t.Errorf("Expected delaySec to be 0")
	}
}
