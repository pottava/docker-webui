// Package controllers implements functions to route user requests
package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	api "github.com/fsouza/go-dockerclient"

	"github.com/pottava/docker-webui/app/config"
	"github.com/pottava/docker-webui/app/engine"
	util "github.com/pottava/docker-webui/app/http"
	"github.com/pottava/docker-webui/app/logs"
	"github.com/pottava/docker-webui/app/misc"
	"github.com/pottava/docker-webui/app/models"
)

func init() {
	cfg := config.NewConfig()

	http.Handle("/container/top/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/container/top/"):]
		client, _ := util.RequestGetParam(r, "client")
		params := struct{ ID, Name, Client string }{id, _label(id, client, cfg.LabelOverrideNames), client}
		util.RenderHTML(w, []string{"containers/top.tmpl"}, params, nil)
	}))
	http.Handle("/container/statlog/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/container/statlog/"):]
		client, _ := util.RequestGetParam(r, "client")
		params := struct {
			ID, Name, Client string
			ViewOnly         bool
		}{id, _label(id, client, cfg.LabelOverrideNames), client, cfg.ViewOnly}
		util.RenderHTML(w, []string{"containers/statlog.tmpl"}, params, nil)
	}))
	http.Handle("/container/changes/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/container/changes/"):]
		client, _ := util.RequestGetParam(r, "client")
		params := struct{ ID, Name, Client string }{id, _label(id, client, cfg.LabelOverrideNames), client}
		util.RenderHTML(w, []string{"containers/changes.tmpl"}, params, nil)
	}))
	http.Handle("/logs", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		params := struct{ LabelFilters string }{strings.Join(cfg.LabelFilters, ",")}
		util.RenderHTML(w, []string{"containers/logs.tmpl"}, params, nil)
	}))
	http.Handle("/statistics", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		clients, err := models.LoadDockerClients()
		if err != nil {
			renderErrorJSON(w, err)
			return
		}
		params := struct {
			LabelFilters string
			Clients      int
		}{strings.Join(cfg.LabelFilters, ","), len(clients)}
		util.RenderHTML(w, []string{"containers/statistics.tmpl"}, params, nil)
	}))

	/**
	 * Containers' API
	 * @param limit int
	 * @param status int (0: all, 1: created, 2: restarting, 3: running, 4: paused, 5&6: exited)
	 * @param q string search words
	 * @return []model.DockerContainer
	 */
	http.Handle("/api/containers", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		type container struct {
			Client     *models.DockerClient     `json:"client"`
			Containers []models.DockerContainer `json:"containers"`
		}
		dockers, ok := clients(w)
		if !ok {
			return
		}
		options := models.ListContainerOption(util.RequestGetParamI(r, "status", 0))
		options.Limit = util.RequestGetParamI(r, "limit", 100)
		d := make(chan *container, len(dockers))
		result := []*container{}

		for _, docker := range dockers {
			go func(docker *engine.Client) {
				containers, err := docker.ListContainers(options)
				if err != nil {
					renderErrorJSON(w, err)
					return
				}
				var words []string
				if q, found := util.RequestGetParam(r, "q"); found {
					words = util.SplittedUpperStrings(q)
				}
				d <- &container{
					Client:     docker.Conf,
					Containers: models.SearchContainers(containers, words),
				}
			}(docker)
		}
		for i := 0; i < len(dockers); i++ {
			result = append(result, <-d)
		}
		close(d)
		util.RenderJSON(w, result, nil)
	}))

	http.Handle("/api/statistics", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		var dockers []*engine.Client
		if c, found := util.RequestGetParam(r, "client"); found {
			if docker, ok := client(w, c); ok {
				dockers = []*engine.Client{docker}
			}
		} else {
			var ok bool
			dockers, ok = clients(w)
			if !ok {
				return
			}
		}
		type statistics struct {
			Client *models.DockerClient               `json:"client"`
			Stats  map[string]map[string][]*api.Stats `json:"stats"`
		}
		stats := []statistics{}

		d := make(chan statistics, len(dockers))
		for _, docker := range dockers {
			go func(docker *engine.Client) {
				candidate, err := docker.ListContainers(models.ListContainerOption(3))
				if err != nil {
					renderErrorJSON(w, err)
					return
				}
				containers := models.SearchContainers(candidate, []string{})
				c := make(chan models.DockerStats, len(containers))
				stats := map[string]map[string][]*api.Stats{}
				count := util.RequestGetParamI(r, "count", 1)

				for _, container := range containers {
					go func(container models.DockerContainer) {
						stat, _ := docker.Stats(container.ID, count)

						name := strings.Join(container.Names, ",")
						if !misc.ZeroOrNil(cfg.LabelOverrideNames) {
							if label, found := container.Labels[cfg.LabelOverrideNames]; found {
								name = "*" + label
							}
						}
						c <- models.DockerStats{
							ID:    container.ID,
							Name:  name,
							Stats: stat,
						}
					}(container)
				}
				for i := 0; i < len(containers); i++ {
					ds := <-c
					inner := map[string][]*api.Stats{}
					inner[ds.Name] = ds.Stats
					stats[ds.ID] = inner
				}
				close(c)
				d <- statistics{
					Client: docker.Conf,
					Stats:  stats,
				}
			}(docker)
		}
		for i := 0; i < len(dockers); i++ {
			stats = append(stats, <-d)
		}
		close(d)
		util.RenderJSON(w, stats, nil)
	}))

	http.Handle("/api/logs", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		count := util.RequestGetParamI(r, "count", 100)
		var dockers []*engine.Client
		if c, found := util.RequestGetParam(r, "client"); found {
			if docker, ok := client(w, c); ok {
				dockers = []*engine.Client{docker}
			}
		} else {
			var ok bool
			dockers, ok = clients(w)
			if !ok {
				return
			}
		}
		type stdlogs struct {
			ID     string   `json:"id"`
			Stdout []string `json:"stdout"`
			Stderr []string `json:"stderr"`
		}
		type clientlogs struct {
			Client *models.DockerClient `json:"client"`
			Logs   []stdlogs            `json:"logs"`
		}
		logs := []clientlogs{}

		d := make(chan clientlogs, len(dockers))
		for _, docker := range dockers {
			go func(docker *engine.Client) {
				candidate, err := docker.ListContainers(models.ListContainerOption(3))
				if err != nil {
					renderErrorJSON(w, err)
					return
				}
				containers := models.SearchContainers(candidate, []string{})

				c := make(chan stdlogs, len(containers))
				inner := []stdlogs{}
				for _, container := range containers {
					go func(container models.DockerContainer) {

						stdout, stderr, err := docker.Logs(container.ID, count, 1*time.Second)
						if err != nil {
							renderErrorJSON(w, err)
							return
						}
						c <- stdlogs{container.ID, stdout, stderr}
					}(container)
				}
				for i := 0; i < len(containers); i++ {
					inner = append(inner, <-c)
				}
				close(c)
				d <- clientlogs{
					Client: docker.Conf,
					Logs:   inner,
				}
			}(docker)
		}
		for i := 0; i < len(dockers); i++ {
			logs = append(logs, <-d)
		}
		close(d)
		util.RenderJSON(w, logs, nil)
	}))

	/**
	 * A container's API
	 */
	// inspect
	http.Handle("/api/container/inspect/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if docker, ok := client(w, util.RequestGetParamS(r, "client", "")); ok {
			id := r.URL.Path[len("/api/container/inspect/"):]
			meta := docker.InspectContainer(id)
			util.RenderJSON(w, meta.Container, meta.Error)
		}
	}))
	// top
	http.Handle("/api/container/top/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if docker, ok := client(w, util.RequestGetParamS(r, "client", "")); ok {
			id := r.URL.Path[len("/api/container/top/"):]
			args := util.RequestGetParamS(r, "args", "aux")
			util.RenderJSON(w, docker.Top(id, args), nil)
		}
	}))
	// stats
	http.Handle("/api/container/stats/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if docker, ok := client(w, util.RequestGetParamS(r, "client", "")); ok {
			id := r.URL.Path[len("/api/container/stats/"):]
			result, err := docker.Stats(id, util.RequestGetParamI(r, "count", 1))
			if err != nil {
				renderErrorJSON(w, err)
				return
			}
			util.RenderJSON(w, result, nil)
		}
	}))
	// logs
	http.Handle("/api/container/logs/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if docker, ok := client(w, util.RequestGetParamS(r, "client", "")); ok {
			id := r.URL.Path[len("/api/container/logs/"):]
			count := util.RequestGetParamI(r, "count", 100)
			stdout, stderr, err := docker.Logs(id, count, 1*time.Second)
			if err != nil {
				renderErrorJSON(w, err)
				return
			}
			util.RenderJSON(w, struct {
				Stdout []string `json:"stdout"`
				Stderr []string `json:"stderr"`
			}{
				stdout,
				stderr,
			}, nil)
		}
	}))
	// diff
	http.Handle("/api/container/changes/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if docker, ok := client(w, util.RequestGetParamS(r, "client", "")); ok {
			id := r.URL.Path[len("/api/container/changes/"):]
			util.RenderJSON(w, docker.Changes(id), nil)
		}
	}))

	// restart
	http.Handle("/api/container/restart/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		if docker, ok := client(w, util.RequestPostParamS(r, "client", "")); ok {
			meta := docker.Restart(r.URL.Path[len("/api/container/restart/"):], 5)
			if meta.Error != nil {
				renderErrorJSON(w, meta.Error)
				return
			}
			util.RenderJSON(w, meta.Container, nil)
		}
	}))
	// start
	http.Handle("/api/container/start/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		if docker, ok := client(w, util.RequestPostParamS(r, "client", "")); ok {
			meta := docker.Start(r.URL.Path[len("/api/container/start/"):])
			if meta.Error != nil {
				renderErrorJSON(w, meta.Error)
				return
			}
			util.RenderJSON(w, meta.Container, nil)
		}
	}))
	// stop
	http.Handle("/api/container/stop/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		if docker, ok := client(w, util.RequestPostParamS(r, "client", "")); ok {
			meta := docker.Stop(r.URL.Path[len("/api/container/stop/"):])
			if meta.Error != nil {
				renderErrorJSON(w, meta.Error)
				return
			}
			util.RenderJSON(w, meta.Container, nil)
		}
	}))
	// kill
	http.Handle("/api/container/kill/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		if docker, ok := client(w, util.RequestPostParamS(r, "client", "")); ok {
			meta := docker.Kill(r.URL.Path[len("/api/container/kill/"):], 5)
			if meta.Error != nil {
				renderErrorJSON(w, meta.Error)
				return
			}
			util.RenderJSON(w, meta.Container, nil)
		}
	}))
	// rm
	http.Handle("/api/container/rm/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		if docker, ok := client(w, util.RequestPostParamS(r, "client", "")); ok {
			err := docker.Rm(r.URL.Path[len("/api/container/rm/"):])
			if err != nil {
				renderErrorJSON(w, err)
				return
			}
			util.RenderJSON(w, "removed successfully.", nil)
		}
	}))
	// rename
	http.Handle("/api/container/rename/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		if docker, ok := client(w, util.RequestPostParamS(r, "client", "")); ok {
			if name, found := util.RequestPostParam(r, "name"); found {
				err := docker.Rename(r.URL.Path[len("/api/container/rename/"):], name)
				message := "renamed successfully."
				if err != nil {
					message = err.Error()
				}
				util.RenderJSON(w, message, nil)
				return
			}
		}
	}))
	// commit
	http.Handle("/api/container/commit/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		if docker, ok := client(w, util.RequestPostParamS(r, "client", "")); ok {
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
		}
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
		repository, tag := api.ParseRepositoryTag(webhook.Repository.RepoName)
		tag = misc.NVL(tag, "latest")

		// TODO trying all docker clients!!
		docker, ok := client(w, util.RequestPostParamS(r, "client", ""))
		if !ok {
			return
		}

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

func _label(id, client, key string) string {
	var docker *engine.Client
	if misc.ZeroOrNil(client) {
		if c, err := engine.Docker(); err == nil {
			docker = c
		}
	}
	if docker == nil {
		masters, err := models.LoadDockerClients()
		if err != nil {
			return id
		}
		for _, master := range masters {
			if master.ID == client {
				engine.Configure(master.Endpoint, master.CertPath)
				if c, err := engine.Docker(); err == nil {
					docker = c
				}
				break
			}
		}
	}
	if docker != nil {
		meta := docker.InspectContainer(id)
		if meta.Error != nil {
			return id
		}
		if name, found := meta.Container.Config.Labels[key]; found {
			return name
		}
	}
	return id
}

func renderErrorJSON(w http.ResponseWriter, err error) {
	util.RenderJSON(w, struct {
		Error string `json:"error"`
	}{err.Error()}, nil)
}
