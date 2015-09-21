package controllers

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"

	api "github.com/fsouza/go-dockerclient"
	"github.com/pottava/docker-webui/app/config"
	"github.com/pottava/docker-webui/app/engine"
	util "github.com/pottava/docker-webui/app/http"
	"github.com/pottava/docker-webui/app/misc"
	"github.com/pottava/docker-webui/app/models"
)

type information struct {
	Client  *models.DockerClient `json:"client,omitempty"`
	Info    *api.Env             `json:"info"`
	Version *api.Env             `json:"version"`
}

func init() {
	cfg := config.NewConfig()

	http.Handle("/clients", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		params := struct{ ViewOnly bool }{cfg.ViewOnly}
		util.RenderHTML(w, []string{"clients/index.tmpl"}, params, nil)
	}))
	http.Handle("/clients/export", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", "attachment; filename=docker-clients.json")
		w.Header().Set("Content-Type", "application/force-download")
		http.ServeFile(w, r, models.DockerClientSavePath)
	}))
	http.Handle("/clients/import", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		var err error
		if err = r.ParseMultipartForm(32 << 20); nil != err {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, headers := range r.MultipartForm.File {
			var in multipart.File
			if in, err = headers[0].Open(); nil != err {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer in.Close()

			var out *os.File
			if out, err = os.Create(models.DockerClientSavePath); nil != err {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer out.Close()

			if _, err = io.Copy(out, in); nil != err {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}))

	/**
	 * Docker client's API
	 */
	http.Handle("/api/clients", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		result := []information{}
		clients, err := models.LoadDockerClients()
		if err != nil {
			renderErrorJSON(w, err)
			return
		}
		for _, client := range clients {
			engine.Configure(client.Endpoint, client.CertPath)
			docker, err := engine.Docker()
			if err != nil {
				client.IsActive = false
				client.Save()
				result = append(result, information{client, nil, nil})
				continue
			}
			info, _ := docker.Info()
			version, _ := docker.Version()
			result = append(result, information{client, info, version})
		}
		util.RenderJSON(w, result, nil)
	}))

	http.Handle("/api/client/", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if endpoint, found := util.RequestPostParam(r, "endpoint"); found {
			cert, _ := util.RequestPostParam(r, "cert")
			engine.Configure(endpoint, cert)
			engine.Save()
			_, err := engine.Docker()
			if err != nil {
				models.RemoveDockerClientByEndpoint(endpoint)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/api/client/", http.StatusFound)
			return
		}
		if r.Method == "DELETE" {
			models.RemoveDockerClient(r.URL.Path[len("/api/client/"):])
			w.WriteHeader(http.StatusOK)
			return
		}
		docker, err := engine.Docker()
		if err != nil {
			renderErrorJSON(w, err)
			return
		}
		info, _ := docker.Info()
		version, _ := docker.Version()
		util.RenderJSON(w, information{nil, info, version}, nil)
	}))
}

func client(w http.ResponseWriter, id string) (client *engine.Client, ok bool) {
	if misc.ZeroOrNil(id) {
		client, err := engine.Docker()
		if err == nil {
			return client, true
		}
		return nil, false
	}
	masters, err := models.LoadDockerClients()
	if err != nil {
		renderErrorJSON(w, err)
		return nil, false
	}
	for _, master := range masters {
		if master.ID == id {
			engine.Configure(master.Endpoint, master.CertPath)
			client, err := engine.Docker()
			if err == nil {
				return client, true
			}
			break
		}
	}
	return nil, false
}

func clients(w http.ResponseWriter) (clients []*engine.Client, ok bool) {
	masters, err := models.LoadDockerClients()
	if err != nil {
		renderErrorJSON(w, err)
		return nil, false
	}
	for _, master := range masters {
		engine.Configure(master.Endpoint, master.CertPath)
		client, err := engine.Docker()
		if err != nil {
			continue
		}
		clients = append(clients, client)
	}
	return clients, true
}
