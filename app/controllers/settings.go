package controllers

import (
	"net/http"

	"github.com/pottava/docker-webui/app/engine"
	util "github.com/pottava/docker-webui/app/http"
)

func init() {

	http.Handle("/set/docker/client", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if endpoint, found := util.RequestPostParam(r, "endpoint"); found {
			util.RenderJSON(w, "", engine.SetDockerClient(endpoint, ""))
			return
		}
		http.NotFound(w, r)
	}))

}
