package stringutils

import "testing"

func TestRandomString(t *testing.T) {
	r := RandString(10)
	if len(r) != 10 {
		t.Errorf("got %v want %v", len(r), 10)
	}
}
