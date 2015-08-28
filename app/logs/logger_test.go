package logs

import "testing"

func TestLogLevel(t *testing.T) {
	expected := info
	if level != expected {
		t.Errorf("Unexpected result. Expected %v, but got %v", expected, level)
		return
	}
}
