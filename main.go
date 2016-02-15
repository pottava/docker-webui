package main

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/pottava/docker-webui/app/config"
	_ "github.com/pottava/docker-webui/app/controllers"
	misc "github.com/pottava/docker-webui/app/http"
	"github.com/pottava/docker-webui/app/logs"
	_ "github.com/pottava/docker-webui/app/models"
)

var (
	buildVersion string
	buildDate    string
)

func main() {
	cfg := config.NewConfig()
	logs.Debug.Print("[config] " + cfg.String())

	http.Handle("/", index(cfg))
	http.HandleFunc("/alive", alive)
	http.HandleFunc("/version", version)
	http.Handle("/assets/", assets(cfg))

	logs.Info.Printf("[service] listening on port %v", cfg.Port)
	logs.Fatal.Print(http.ListenAndServe(":"+fmt.Sprint(cfg.Port), nil))
}

func index(cfg *config.Config) http.Handler {
	return misc.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		params := struct {
			LabelOverride string
			LabelFilters  string
			ViewOnly      bool
		}{cfg.LabelOverrideNames, strings.Join(cfg.LabelFilters, ","), cfg.ViewOnly}
		misc.RenderHTML(w, []string{"containers/index.tmpl"}, params, nil)
	})
}
func alive(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func version(w http.ResponseWriter, r *http.Request) {
	if len(buildVersion) > 0 && len(buildDate) > 0 {
		fmt.Fprintf(w, "version: %s (built at %s)", buildVersion, buildDate)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
func assets(cfg *config.Config) http.Handler {
	fs := http.FileServer(http.Dir(path.Join(cfg.StaticFilePath, "assets")))
	return misc.AssetsChain(http.StripPrefix("/assets/", fs).ServeHTTP)
}
