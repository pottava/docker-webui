package misc

import (
	"io/ioutil"
	"net/http"
	"testing"
)

// IsSuccessHTTPRequest checks whether the HTTP request was success.
func IsSuccessHTTPRequest(t *testing.T, actual *http.Response, err error) bool {
	if err != nil {
		t.Error("Unexpected error occered")
		return false
	}
	expected := http.StatusOK
	if actual.StatusCode != expected {
		t.Errorf("Status code error. Expected %v, but got %v", expected, actual.StatusCode)
		return false
	}
	return true
}

// ParseHTTPResponse parses a HTTP response to string.
func ParseHTTPResponse(res *http.Response) (string, int) {
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return string(contents), res.StatusCode
}
