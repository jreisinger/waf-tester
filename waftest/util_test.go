package waftest

import "testing"

func TestGenUUID(t *testing.T) {
	uuid := genUUID()
	if len(uuid) != 36 {
		t.Errorf("genUUID: want 16-byte string, got %d-byte string", len(uuid))
	}
}

func TestIntInSlice(t *testing.T) {
	is := []int{0, -3, 100, 42, 27, 3}

	i := 42
	if !intInSlice(i, is) {
		t.Error("intInSlice: want true, got false")
	}

	i = 1000
	if intInSlice(i, is) {
		t.Error("intInSlice: want false, got true")
	}
}

func TestStringInSlice(t *testing.T) {
	ss := []string{"repetitio", "est", "mater", "studorium"}

	s := "est"
	if !stringInSlice(s, ss) {
		t.Error("stringInSlice: want true, got false")
	}

	s = "veritas"
	if stringInSlice(s, ss) {
		t.Error("stringInSlice: want false, got true")
	}
}
