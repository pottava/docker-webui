package http

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestRequestGetParam(t *testing.T) {
	req := &http.Request{Method: "GET"}
	req.URL, _ = url.Parse("http://www.google.com/search?q=foo")

	actual, found := RequestGetParam(req, "q")
	if !found {
		t.Error("Could not find the parameter.")
		return
	}
	expected := "foo"
	if actual != expected {
		t.Errorf("Unexpected result. Expected %v, but got %v", expected, actual)
		return
	}
}

func TestRequestGetParamS(t *testing.T) {
	req := &http.Request{Method: "GET"}
	req.URL, _ = url.Parse("http://www.google.com/search?q=foo")

	actual := RequestGetParamS(req, "q", "bar")
	expected := "foo"
	if actual != expected {
		t.Errorf("Unexpected result. Expected %v, but got %v", expected, actual)
		return
	}
	actual = RequestGetParamS(req, "r", "bar")
	expected = "bar"
	if actual != expected {
		t.Errorf("Unexpected result. Expected %v, but got %v", expected, actual)
		return
	}
}

func TestRequestGetParamI(t *testing.T) {
	req := &http.Request{Method: "GET"}
	req.URL, _ = url.Parse("http://www.google.com/search?q=1")

	actual := RequestGetParamI(req, "q", 2)
	expected := 1
	if actual != expected {
		t.Errorf("Unexpected result. Expected %v, but got %v", expected, actual)
		return
	}
	actual = RequestGetParamI(req, "r", 2)
	expected = 2
	if actual != expected {
		t.Errorf("Unexpected result. Expected %v, but got %v", expected, actual)
		return
	}
}

func TestSplittedUpperStrings(t *testing.T) {
	actual := SplittedUpperStrings("")
	expected := []string{""}
	if reflect.DeepEqual(actual, expected) {
		t.Errorf("Unexpected result. Expected %v, but got %v", expected, actual)
		return
	}
	actual = SplittedUpperStrings("foo")
	expected = []string{"foo"}
	if reflect.DeepEqual(actual, expected) {
		t.Errorf("Unexpected result. Expected %v, but got %v", expected, actual)
		return
	}
	actual = SplittedUpperStrings("foo bar")
	expected = []string{"foo", "bar"}
	if reflect.DeepEqual(actual, expected) {
		t.Errorf("Unexpected result. Expected %v, but got %v", expected, actual)
		return
	}
	actual = SplittedUpperStrings("foo@bar*baz?qux")
	expected = []string{"foo", "bar", "baz", "qux"}
	if reflect.DeepEqual(actual, expected) {
		t.Errorf("Unexpected result. Expected %v, but got %v", expected, actual)
		return
	}
}
