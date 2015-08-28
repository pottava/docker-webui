package config

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	var expected interface{}
	actual := defaultConfig()

	expected = "docker web-ui"
	if actual.Name != expected {
		t.Errorf("Expected %v, but got %v", expected, actual.Name)
		return
	}
	expected = uint16(9000)
	if actual.Port != expected {
		t.Errorf("Expected %v, but got %v", expected, actual.Port)
		return
	}
}

func TestMerge(t *testing.T) {
	cfg := Config{
		Name:     "Test",
		Port:     8080,
		LogLevel: 6,
	}
	actual := cfg.merge(defaultConfig())
	gopath := os.Getenv("GOPATH")
	expected := &Config{
		Name:                   "Test",
		Port:                   8080,
		LogLevel:               6,
		DockerEndpoint:         "unix:///var/run/docker.sock",
		DockerAPIVersion:       "1.17",
		DockerPullBeginTimeout: 3 * time.Minute,
		DockerPullTimeout:      2 * time.Hour,
		DockerStartTimeout:     1 * time.Minute,
		DockerStopTimeout:      1 * time.Minute,
		DockerRestartTimeout:   1 * time.Minute,
		DockerRmTimeout:        5 * time.Minute,
		StaticFileHost:         "",
		StaticFilePath:         gopath + "/src/github.com/pottava/docker-webui/app",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
}

func TestComplete(t *testing.T) {
	actual := defaultConfig()
	if actual.complete() {
		t.Errorf("Unexpected result. %v", actual)
		return
	}
	actual = *actual.merge(Config{
		StaticFileHost: "cdn-host",
	})
	if !actual.complete() {
		t.Errorf("Unexpected result. %v", actual)
		return
	}
}

func TestTrimWhitespace(t *testing.T) {
	actual := defaultConfig()
	actual.Name = " 　a b 　　c "
	actual.trimWhitespace()

	expected := "a b 　　c"
	if actual.Name != expected {
		t.Errorf("Expected %v, but got %v", expected, actual.Name)
		return
	}
}
