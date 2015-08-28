package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pottava/docker-webui/app/config"
	"github.com/pottava/docker-webui/app/misc"
)

func TestIndex(t *testing.T) {
	ts := httptest.NewServer(http.Handler(index()))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if ok := misc.IsSuccessHTTPRequest(t, res, err); !ok {
		return
	}
	actual, _ := misc.ParseHTTPResponse(res)
	expected := config.NewConfig().Name
	if !strings.Contains(actual, expected) {
		t.Errorf("Invalid response. Expected %v, but got %v", expected, actual)
		return
	}
}

func TestAlive(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(alive))
	defer ts.Close()
	res, err := http.Get(ts.URL)
	misc.IsSuccessHTTPRequest(t, res, err)
}

func TestVersion(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(version))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if ok := misc.IsSuccessHTTPRequest(t, res, err); !ok {
		return
	}
	actual, _ := misc.ParseHTTPResponse(res)
	expected := fmt.Sprint(misc.Version)
	if !strings.Contains(actual, expected) {
		t.Errorf("Invalid response. Expected %v, but got %v", expected, actual)
		return
	}
}

func TestAssets(t *testing.T) {
	cfg := config.NewConfig()
	ts := httptest.NewServer(http.Handler(assets(cfg)))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/assets/css/main.css")
	if ok := misc.IsSuccessHTTPRequest(t, res, err); !ok {
		return
	}
	actual, _ := misc.ParseHTTPResponse(res)
	expected := "monospace"
	if !strings.Contains(actual, expected) {
		t.Errorf("Invalid response. Expected %v, but got %v", expected, actual)
		return
	}
}
