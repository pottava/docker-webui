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
		DockerEndpoints:        []string{"unix:///var/run/docker.sock"},
		DockerCertPath:         []string{""},
		DockerPullBeginTimeout: 3 * time.Minute,
		DockerPullTimeout:      2 * time.Hour,
		DockerStatTimeout:      5 * time.Second,
		DockerStartTimeout:     10 * time.Second,
		DockerStopTimeout:      10 * time.Second,
		DockerRestartTimeout:   10 * time.Second,
		DockerKillTimeout:      10 * time.Second,
		DockerRmTimeout:        5 * time.Minute,
		DockerCommitTimeout:    30 * time.Second,
		StaticFileHost:         "",
		StaticFilePath:         gopath + "/src/github.com/pottava/docker-webui/app",
		PreventSelfStop:        true,
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
		LabelOverrideNames: "com.github.pottava.name",
		DockerCertPath:     []string{"cert-path"},
		StaticFileHost:     "cdn-host",
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
