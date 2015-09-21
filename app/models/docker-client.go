package models

import (
	"fmt"
	"hash/fnv"
	"regexp"
	"strings"

	"github.com/pottava/docker-webui/app/config"
	"github.com/pottava/docker-webui/app/misc"
)

// DockerClient represents a Docker client
type DockerClient struct {
	ID       string `json:"id"`
	Endpoint string `json:"endpoint"`
	CertPath string `json:"certPath"`
	IsActive bool   `json:"isActive"`
}

// DockerClientSavePath is the path to save clients
var DockerClientSavePath string

func init() {
	r, _ := regexp.Compile("[^a-zA-Z0-9_\\.]")
	name := r.ReplaceAllString(strings.ToLower(config.NewConfig().Name), "-")
	DockerClientSavePath = "/tmp/" + name + "-docker-clients.json"
}

// LoadDockerClients returns registered clients
func LoadDockerClients() (clients []*DockerClient, err error) {
	err = misc.ReadFromFile(DockerClientSavePath, &clients)
	for _, client := range clients {
		client.ID = fmt.Sprint(Hash(client.Endpoint))
	}
	return
}

// RemoveDockerClient removes your specified configuration
func RemoveDockerClient(id string) bool {
	prev, _ := LoadDockerClients()
	next := []*DockerClient{}
	for _, client := range prev {
		if client.ID == id {
			continue
		}
		next = append(next, client)
	}
	misc.SaveAsFile(DockerClientSavePath, next)
	return true
}

// RemoveDockerClientByEndpoint removes your specified configuration
func RemoveDockerClientByEndpoint(endpoint string) {
	RemoveDockerClient(fmt.Sprint(Hash(endpoint)))
}

// Load select its configuration
func (c *DockerClient) Load() {
	clients, _ := LoadDockerClients()
	for _, client := range clients {
		if client.Endpoint == c.Endpoint {
			c.ID = client.ID
			c.CertPath = client.CertPath
			c.IsActive = client.IsActive
			break
		}
	}
}

// Save persists the client configuration
func (c *DockerClient) Save() {
	clients, _ := LoadDockerClients()
	found := false
	for _, client := range clients {
		if client.Endpoint == c.Endpoint {
			client.CertPath = c.CertPath
			client.IsActive = c.IsActive
			found = true
		}
	}
	if !found {
		clients = append(clients, c)
	}
	misc.SaveAsFile(DockerClientSavePath, clients)
}

// Hash returns its hashed value
func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
