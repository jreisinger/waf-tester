package httptest

import "testing"

func TestGenUUID(t *testing.T) {
	uuid := genUUID()
	if len(uuid) != 36 {
		t.Errorf("genUUID: want 16-byte string, got %d-byte string", len(uuid))
	}
}
