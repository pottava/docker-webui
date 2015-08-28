// Package misc defines miscellaneous useful functions
package misc

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

// NVL is null value logic
func NVL(str string, def string) string {
	if len(str) == 0 {
		return def
	}
	return str
}

// ZeroOrNil checks if the argument is zero or null
func ZeroOrNil(obj interface{}) bool {
	value := reflect.ValueOf(obj)
	if !value.IsValid() {
		return true
	}
	if obj == nil {
		return true
	}
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		return value.Len() == 0
	}
	zero := reflect.Zero(reflect.TypeOf(obj))
	if obj == zero.Interface() {
		return true
	}
	return false
}

// Atoi returns casted int
func Atoi(candidate string) int {
	result := 0
	if candidate != "" {
		if i, err := strconv.Atoi(candidate); err == nil {
			result = i
		}
	}
	return result
}

// ParseUint16 returns casted uint16
func ParseUint16(candidate string) uint16 {
	var result uint16
	if candidate != "" {
		if u, err := strconv.ParseUint(candidate, 10, 16); err == nil {
			result = uint16(u)
		}
	}
	return result
}

// ParseDuration returns casted time.Duration
func ParseDuration(candidate string) time.Duration {
	var result time.Duration
	if candidate != "" {
		if d, err := time.ParseDuration(candidate); err == nil {
			result = d
		}
	}
	return result
}

// ParseBool returns casted bool
func ParseBool(candidate string) bool {
	result := false
	if candidate != "" {
		if b, err := strconv.ParseBool(candidate); err == nil {
			result = b
		}
	}
	return result
}

// ParseCsvLine returns comma splitted strings
func ParseCsvLine(data string) []string {
	splitted := strings.SplitN(data, ",", -1)

	parsed := make([]string, len(splitted))
	for i, val := range splitted {
		parsed[i] = strings.TrimSpace(val)
	}
	return parsed
}

// TimeToJST changes time.Time to Tokyo time zone
func TimeToJST(t time.Time) time.Time {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return t
	}
	return t.In(jst)
}

// TimeToString changes time.Time to string
func TimeToString(t time.Time) string {
	timeformat := "2006-01-02T15:04:05Z0700"
	return t.Format(timeformat)
}

// StringToTime changes string to time.Time
func StringToTime(t string) time.Time {
	timeformat := "2006-01-02T15:04:05Z0700"
	candidate, _ := time.Parse(timeformat, t)
	return candidate
}
