package httptest

import "testing"

func TestGenUUID(t *testing.T) {
	uuid := genUUID()
	if len(uuid) != 36 {
		t.Errorf("genUUID: want 16-byte string, got %d-byte string", len(uuid))
	}
}

func TestIntInSlice(t *testing.T) {
	s := []int{0, -3, 100, 42, 27, 3}

	n := 42
	if !intInSlice(n, s) {
		t.Error("intInSlice: want true, got false")
	}

	n = 1000
	if intInSlice(n, s) {
		t.Error("intInSlice: want false, got true")
	}
}
