package models

import "testing"

func TestListAPI(t *testing.T) {
	var actual, expected interface{}
	apis := ListAPI()

	actual = len(apis)
	expected = 2
	if actual != expected {
		t.Errorf("Unexpected API count. Expected %v, but got %v", expected, actual)
		return
	}
	actual = apis[1].Name
	expected = "/reinvent-session"
	if actual != expected {
		t.Errorf("Unexpected name. Expected %v, but got %v", expected, actual)
		return
	}
	actual = len(apis[0].Parameters)
	expected = 2
	if actual != expected {
		t.Errorf("Unexpected parameters count. Expected %v, but got %v", expected, actual)
		return
	}
	actual = apis[1].Parameters[0].Necessary
	expected = true
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
}
