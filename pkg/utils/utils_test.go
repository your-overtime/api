package utils

import (
	"testing"
	"time"
)

func TestRandomString(t *testing.T) {
	r := RandString(10)
	if len(r) != 10 {
		t.Errorf("got %v want %v", len(r), 10)
	}
}

func TestDayStart(t *testing.T) {
	s, err := time.Parse(time.RFC3339, "2021-01-02T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}
	ds := DayStart(s)

	if ds.Day() != s.Day() {
		t.Errorf("got %v want %v", ds, s)
	}

	if ds.Hour() != 0 {
		t.Errorf("got %v want %v", ds, 0)
	}
}

func TestDayEnd(t *testing.T) {
	s, err := time.Parse(time.RFC3339, "2021-01-02T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}
	de := DayEnd(s)

	edX, err := time.Parse(time.RFC3339, "2021-01-03T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}

	if de.Day() == s.Day() {
		t.Errorf("got %v want %v", de, s)
	}

	if de.Day() != edX.Day() {
		t.Errorf("got %v want %v", de, edX)
	}

	if de.Hour() != 0 {
		t.Errorf("got %v want %v", de, 0)
	}
}
