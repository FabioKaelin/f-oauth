package auth

import (
	"strings"
	"testing"
)

func TestNewPasswordByteBoundaries(t *testing.T) {
	for _, test := range []struct {
		password string
		valid    bool
	}{{"", false}, {"1234567", false}, {"12345678", true}, {strings.Repeat("a", 72), true}, {strings.Repeat("a", 73), false}, {strings.Repeat("ä", 37), false}} {
		if (ValidateNewPassword(test.password) == nil) != test.valid {
			t.Fatal("password validation boundary failed")
		}
	}
}
