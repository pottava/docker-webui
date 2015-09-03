// Package controllers implements functions to route user requests
package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/docker/docker/pkg/parsers"

	"github.com/pottava/docker-webui/app/engine"
	util "github.com/pottava/docker-webui/app/http"
	"github.com/pottava/docker-webui/app/logs"
	"github.com/pottava/docker-webui/app/misc"
	"github.com/pottava/docker-webui/app/models"
)

func init() {
	docker := engine.Docker

	http.Handle("/container/top/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/container/top/"):]
		util.RenderHTML(w, []string{"containers/top.tmpl"}, struct{ ID string }{id}, nil)
	}))
	http.Handle("/container/statlog/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/container/statlog/"):]
		util.RenderHTML(w, []string{"containers/statlog.tmpl"}, struct{ ID string }{id}, nil)
	}))
	http.Handle("/container/changes/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/container/changes/"):]
		util.RenderHTML(w, []string{"containers/changes.tmpl"}, struct{ ID string }{id}, nil)
	}))

	/**
	 * Containers
	 * @param limit int
	 * @param status int (0: all, 1: created, 2: restarting, 3: running, 4: paused, 5&6: exited)
	 * @param q string search words
	 * @return []model.DockerContainer
	 */
	http.Handle("/api/containers", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		options := models.ListContainerOption(util.RequestGetParamI(r, "status", 0))
		options.Limit = util.RequestGetParamI(r, "limit", 100)

		containers, err := docker.ListContainers(options)
		var words []string
		if q, found := util.RequestGetParam(r, "q"); found {
			words = util.SplittedUpperStrings(q)
		}
		util.RenderJSON(w, models.SearchContainers(containers, words), err)
	}))

	// inspect
	http.Handle("/api/container/inspect/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/api/container/inspect/"):]
		meta := docker.InspectContainer(id)
		util.RenderJSON(w, meta.Container, meta.Error)
	}))
	// top
	http.Handle("/api/container/top/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/api/container/top/"):]
		args := util.RequestGetParamS(r, "args", "aux")
		util.RenderJSON(w, docker.Top(id, args), nil)
	}))
	// stats
	http.Handle("/api/container/stats/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/api/container/stats/"):]
		result, err := docker.Stats(id, util.RequestGetParamI(r, "count", 1))
		if err != nil {
			renderErrorJSON(w, err)
			return
		}
		util.RenderJSON(w, result, nil)
	}))
	// logs
	http.Handle("/api/container/logs/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/api/container/logs/"):]
		since := time.Now().Add(time.Duration(util.RequestGetParamI(r, "prev", 5)*-1) * time.Second).UnixNano()
		count := util.RequestGetParamI(r, "count", 100)

		stdout, stderr, err := docker.Logs(id, since, int64(count))
		util.RenderJSON(w, struct {
			Stdout string `json:"stdout"`
			Stderr string `json:"stderr"`
		}{
			stdout,
			stderr,
		}, err)
	}))
	// diff
	http.Handle("/api/container/changes/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/api/container/changes/"):]
		util.RenderJSON(w, docker.Changes(id), nil)
	}))

	// restart
	http.Handle("/api/container/restart/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		meta := docker.Restart(r.URL.Path[len("/api/container/restart/"):], 5)
		if meta.Error != nil {
			renderErrorJSON(w, meta.Error)
			return
		}
		util.RenderJSON(w, meta.Container, nil)
	}))
	// start
	http.Handle("/api/container/start/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		meta := docker.Start(r.URL.Path[len("/api/container/start/"):])
		if meta.Error != nil {
			renderErrorJSON(w, meta.Error)
			return
		}
		util.RenderJSON(w, meta.Container, nil)
	}))
	// stop
	http.Handle("/api/container/stop/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		meta := docker.Stop(r.URL.Path[len("/api/container/stop/"):])
		if meta.Error != nil {
			renderErrorJSON(w, meta.Error)
			return
		}
		util.RenderJSON(w, meta.Container, nil)
	}))
	// kill
	http.Handle("/api/container/kill/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		meta := docker.Kill(r.URL.Path[len("/api/container/kill/"):], 5)
		if meta.Error != nil {
			renderErrorJSON(w, meta.Error)
			return
		}
		util.RenderJSON(w, meta.Container, nil)
	}))
	// rm
	http.Handle("/api/container/rm/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		err := docker.Rm(r.URL.Path[len("/api/container/rm/"):])
		message := "removed successfully."
		if err != nil {
			message = err.Error()
		}
		util.RenderJSON(w, message, nil)
	}))
	// rename
	http.Handle("/api/container/rename/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if name, found := util.RequestPostParam(r, "name"); found {
			err := docker.Rename(r.URL.Path[len("/api/container/rename/"):], name)
			message := "renamed successfully."
			if err != nil {
				message = err.Error()
			}
			util.RenderJSON(w, message, nil)
			return
		}
		http.NotFound(w, r)
	}))
	// commit
	http.Handle("/api/container/commit/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		repository, _ := util.RequestPostParam(r, "repo")
		tag, _ := util.RequestPostParam(r, "tag")
		massage, _ := util.RequestPostParam(r, "msg")
		author, _ := util.RequestPostParam(r, "author")

		meta := docker.Commit(
			r.URL.Path[len("/api/container/commit/"):],
			repository, tag, massage, author)
		if meta.Error != nil {
			renderErrorJSON(w, meta.Error)
			return
		}
		util.RenderJSON(w, meta.Image, nil)
	}))

	// @see https://docs.docker.com/docker-hub/builds/#webhooks
	type Repository struct {
		Name     string `json:"name"`
		Owner    string `json:"owner"`
		RepoName string `json:"repo_name"`
	}
	type Webhook struct {
		CallbackURL string     `json:"callback_url"`
		Repository  Repository `json:"repository"`
	}

	/**
	 * Update by DockerHub
	 * which pull an image again & restart the container to update its service
	 */
	http.Handle("/api/container/update", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		// parse parameters
		webhook := Webhook{}
		err := json.NewDecoder(r.Body).Decode(&webhook)
		if err != nil {
			logs.Warn.Print(err)
			http.NotFound(w, r)
			return
		}
		repository, tag := parsers.ParseRepositoryTag(webhook.Repository.RepoName)
		tag = misc.NVL(tag, "latest")

		// pull the latest image
		var imageRepoTag string
		for _, image := range docker.ListImages() {
			for _, repotag := range image.RepoTags {
				if repotag == repository+":"+tag {
					imageRepoTag = repotag
				}
			}
		}
		if misc.ZeroOrNil(imageRepoTag) {
			http.NotFound(w, r)
			return
		}
		if meta := docker.Pull(imageRepoTag); meta.Error != nil {
			util.RenderJSON(w, meta.Error.Error(), nil)
			return
		}

		// list running containers
		containers, err := docker.ListContainers(models.ListContainerOption(3))
		if err != nil {
			util.RenderJSON(w, err.Error(), nil)
			return
		}

		// restart its container
		restarted := []string{}
		for _, container := range containers {
			if container.Image == imageRepoTag {
				meta := docker.InspectContainer(container.ID)
				if meta.Error != nil {
					continue
				}
				// remove the existing container
				if meta := docker.Stop(container.ID); meta.Error != nil {
					renderErrorJSON(w, meta.Error)
					return
				}
				if err := docker.Rm(container.ID); err != nil {
					renderErrorJSON(w, err)
					return
				}
				// create a new container using the dead container configurations.
				// because if we just restart the container, its image would not
				// reference the new one.
				c := meta.Container
				if meta := docker.Create(c.Name, c.Config, c.HostConfig); meta.Error != nil {
					renderErrorJSON(w, meta.Error)
					return
				}
				if meta := docker.Start(c.Name[1:]); meta.Error != nil {
					renderErrorJSON(w, meta.Error)
					return
				}
				restarted = append(restarted, c.Name[1:])
			}
		}
		util.RenderJSON(w, restarted, nil)
	}))
}

func renderErrorJSON(w http.ResponseWriter, err error) {
	util.RenderJSON(w, struct {
		Error string `json:"error"`
	}{err.Error()}, nil)
}
