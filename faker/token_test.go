package faker

import (
	"testing"
)

func TestGetTokenResponse(t *testing.T) {
	token := GetTokenResponse()
	if token.Token == "" {
		t.Errorf("Expected token to be not empty")
	}
}
