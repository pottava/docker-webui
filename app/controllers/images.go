package controllers

import (
	"net/http"

	"github.com/pottava/docker-webui/app/docker"
	util "github.com/pottava/docker-webui/app/http"
	"github.com/pottava/docker-webui/app/logs"
	"github.com/pottava/docker-webui/app/models"
)

func init() {
	docker, err := engine.NewDockerClient()
	if err != nil {
		logs.Fatal.Printf("@docker.NewClient %v", err)
	}

	http.Handle("/images", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		util.RenderHTML(w, []string{"images/index.tmpl"}, nil, nil)
	}))
	http.Handle("/image/history/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/image/history/"):]
		util.RenderHTML(w, []string{"images/history.tmpl"}, struct{ ID string }{id}, nil)
	}))

	/**
	 * Images
	 * @param q string search words
	 * @return []model.DockerImage
	 */
	http.Handle("/api/images", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		images := docker.ListImages()
		var words []string
		if q, found := util.RequestGetParam(r, "q"); found {
			words = util.SplittedUpperStrings(q)
		}
		util.RenderJSON(w, models.SearchImages(images, words), err)
	}))

	// inspect
	http.Handle("/api/image/inspect/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/api/image/inspect/"):]
		meta := docker.InspectImage(id)
		util.RenderJSON(w, meta.Image, meta.Error)
	}))
	// history
	http.Handle("/api/image/history/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/api/image/history/"):]
		util.RenderJSON(w, docker.History(id), nil)
	}))

	// pull
	http.Handle("/api/image/pull/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		meta := docker.Pull(r.URL.Path[len("/api/image/pull/"):])
		if meta.Error != nil {
			util.RenderJSON(w, meta.Error.Error(), nil)
			return
		}
		util.RenderJSON(w, meta.Image, nil)
	}))
	// run
	http.Handle("/api/image/run/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		meta := docker.Run(models.ParseCreateContainerOption(r))
		util.RenderJSON(w, meta.Container, meta.Error)
	}))
	// rmi
	http.Handle("/api/image/rmi/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		err := docker.Rmi(r.URL.Path[len("/api/image/rmi/"):])
		message := "removed successfully."
		if err != nil {
			message = err.Error()
		}
		util.RenderJSON(w, message, nil)
	}))

}
