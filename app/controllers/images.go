package controllers

import (
	"net/http"

	"github.com/pottava/docker-webui/app/config"
	"github.com/pottava/docker-webui/app/engine"
	util "github.com/pottava/docker-webui/app/http"
	"github.com/pottava/docker-webui/app/models"
)

func init() {
	cfg := config.NewConfig()

	http.Handle("/images", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		params := struct{ ViewOnly bool }{cfg.ViewOnly}
		util.RenderHTML(w, []string{"images/index.tmpl"}, params, nil)
	}))
	http.Handle("/image/history/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/image/history/"):]
		client, _ := util.RequestGetParam(r, "client")
		util.RenderHTML(w, []string{"images/history.tmpl"}, struct{ ID, Client string }{id, client}, nil)
	}))

	/**
	 * Images' API
	 * @param q string search words
	 * @return []model.DockerImage
	 */
	http.Handle("/api/images", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		type image struct {
			Client *models.DockerClient `json:"client"`
			Images []models.DockerImage `json:"images"`
		}
		if dockers, ok := clients(w); ok {
			result := []*image{}

			d := make(chan *image, len(dockers))
			for _, docker := range dockers {
				go func(docker *engine.Client) {
					images := docker.ListImages()
					var words []string
					if q, found := util.RequestGetParam(r, "q"); found {
						words = util.SplittedUpperStrings(q)
					}
					d <- &image{
						Client: docker.Conf,
						Images: models.SearchImages(images, words),
					}
				}(docker)
			}
			for i := 0; i < len(dockers); i++ {
				images := <-d
				result = append(result, images)
			}
			close(d)
			util.RenderJSON(w, result, nil)
		}
	}))

	/**
	 * An image's API
	 */
	// inspect
	http.Handle("/api/image/inspect/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if docker, ok := client(w, util.RequestGetParamS(r, "client", "")); ok {
			id := r.URL.Path[len("/api/image/inspect/"):]
			meta := docker.InspectImage(id)
			util.RenderJSON(w, meta.Image, meta.Error)
		}
	}))
	// history
	http.Handle("/api/image/history/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if docker, ok := client(w, util.RequestGetParamS(r, "client", "")); ok {
			id := r.URL.Path[len("/api/image/history/"):]
			util.RenderJSON(w, docker.History(id), nil)
		}
	}))

	// pull
	http.Handle("/api/image/pull/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if docker, ok := client(w, util.RequestGetParamS(r, "client", "")); ok {
			meta := docker.Pull(r.URL.Path[len("/api/image/pull/"):])
			if meta.Error != nil {
				util.RenderJSON(w, meta.Error.Error(), nil)
				return
			}
			util.RenderJSON(w, meta.Image, nil)
		}
	}))
	// rmi
	http.Handle("/api/image/rmi/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		if docker, ok := client(w, util.RequestPostParamS(r, "client", "")); ok {
			err := docker.Rmi(r.URL.Path[len("/api/image/rmi/"):])
			message := "removed successfully."
			if err != nil {
				message = err.Error()
			}
			util.RenderJSON(w, message, nil)
		}
	}))
	// tag
	http.Handle("/api/image/tag/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		if docker, ok := client(w, util.RequestPostParamS(r, "client", "")); ok {
			repository, _ := util.RequestPostParam(r, "repo")
			tag, _ := util.RequestPostParam(r, "tag")
			err := docker.Tag(r.URL.Path[len("/api/image/tag/"):], repository, tag)
			message := "tagged successfully."
			if err != nil {
				message = err.Error()
			}
			util.RenderJSON(w, message, nil)
		}
	}))

}
