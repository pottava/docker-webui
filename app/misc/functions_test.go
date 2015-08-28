package misc

import (
	"strings"
	"testing"
	"time"
)

func TestNVL(t *testing.T) {
	actual := NVL("foo", "bar")
	expected := "foo"
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
	actual = NVL("", "foo")
	expected = "foo"
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
}

func TestZeroOrNil(t *testing.T) {
	actual := ZeroOrNil(0)
	expected := true
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
	actual = ZeroOrNil(nil)
	expected = true
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
	actual = ZeroOrNil([]string{})
	expected = true
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
	actual = ZeroOrNil(make([]string, 0))
	expected = true
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
	actual = ZeroOrNil("")
	expected = true
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
	actual = ZeroOrNil("nil")
	expected = false
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
}

func TestAtoi(t *testing.T) {
	actual := Atoi("foo")
	expected := 0
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
	actual = Atoi("1")
	expected = 1
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
}

func TestParseUint16(t *testing.T) {
	actual := ParseUint16("-1")
	expected := uint16(0)
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
	actual = ParseUint16("0")
	expected = uint16(0)
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
	actual = ParseUint16("65535")
	expected = uint16(65535)
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
	actual = ParseUint16("65536")
	expected = uint16(0)
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
}

func TestParseDuration(t *testing.T) {
	actual := ParseDuration("1h1m1s")
	expected := 1*time.Hour + 1*time.Minute + 1*time.Second
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
	actual = ParseDuration("-1s")
	expected = -1 * time.Second
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
	actual = ParseDuration("1S")
	expected = 0
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
	actual = ParseDuration("1m 1s")
	expected = 0
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
	actual = ParseDuration("1")
	expected = 0
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
}

func TestParseBool(t *testing.T) {
	actual := ParseBool("true")
	expected := true
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
	actual = ParseBool("True")
	expected = true
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
	actual = ParseBool("TRUE")
	expected = true
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
	actual = ParseBool("foo")
	expected = false
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
	actual = ParseBool("false")
	expected = false
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
}

func TestParseCsvLine(t *testing.T) {
	var actual, expected interface{}

	parsed := ParseCsvLine("1,2")
	actual = len(parsed)
	expected = 2
	if actual != expected {
		t.Errorf("Unexpected count. Expected %v, but got %v", expected, actual)
		return
	}
	actual = parsed[1]
	expected = "2"
	if actual != expected {
		t.Errorf("Unexpected count. Expected %v, but got %v", expected, actual)
		return
	}
	actual = len(ParseCsvLine("1 2"))
	expected = 1
	if actual != expected {
		t.Errorf("Unexpected count. Expected %v, but got %v", expected, actual)
		return
	}
	actual = len(ParseCsvLine(","))
	expected = 2
	if actual != expected {
		t.Errorf("Unexpected count. Expected %v, but got %v", expected, actual)
		return
	}
	actual = len(ParseCsvLine(",,"))
	expected = 3
	if actual != expected {
		t.Errorf("Unexpected count. Expected %v, but got %v", expected, actual)
		return
	}
}

func TestStringToTime(t *testing.T) {
	actual := StringToTime("2015-08-02T23:45:06Z")
	expected, _ := time.Parse("2006/01/02 15:04:05", "2015/08/02 23:45:06")
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
}

func TestTimeToString(t *testing.T) {
	actual := TimeToString(time.Now())
	if strings.Contains(actual, "0001-01-01") {
		t.Errorf("Could not cast to string. %v", actual)
		return
	}
	if !strings.Contains(actual, "T") {
		t.Errorf("Could not cast to string. %v", actual)
		return
	}
}

func TestTimeToJST(t *testing.T) {
	actual := TimeToJST(StringToTime("2015-08-02T23:45:06Z"))
	if !strings.Contains(TimeToString(actual), "+0900") {
		t.Errorf("Could not change the local. %v", actual)
		return
	}
}
