package models

import (
	"fmt"
	"sort"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/pottava/docker-webui/app/config"
	"github.com/pottava/docker-webui/app/misc"
)

// DockerContainer represents a container
type DockerContainer struct {
	ID         string            `json:"id"`
	Image      string            `json:"image,omitempty"`
	Command    string            `json:"command,omitempty"`
	Created    int64             `json:"created,omitempty"`
	Status     string            `json:"status,omitempty"`
	Ports      []docker.APIPort  `json:"ports,omitempty"`
	SizeRw     int64             `json:"sizeRw,omitempty"`
	SizeRootFs int64             `json:"sizeRootFs,omitempty"`
	Names      []string          `json:"names,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
}

// DockerContainers represents list of DockerContainer
type DockerContainers []DockerContainer

// DockerStats represents a container's stats
type DockerStats struct {
	Name  string
	Stats []*docker.Stats
}

var cfg *config.Config

func init() {
	cfg = config.NewConfig()
}

// ListContainerOption returns docker.ListContainersOptions according to the flag
// @param flag int (0: all, 1: created, 2: restarting, 3: running, 4: paused, 5&6: exited)
func ListContainerOption(flag int) docker.ListContainersOptions {
	options := docker.ListContainersOptions{Limit: 100, Filters: map[string][]string{}}
	switch flag {
	case 0:
		options.All = true
		break
	case 1:
		options.All = false
		options.Filters["status"] = []string{"created"}
		break
	case 2:
		options.All = false
		options.Filters["status"] = []string{"restarting"}
		break
	case 3:
		options.All = false
		options.Filters["status"] = []string{"running"}
		break
	case 4:
		options.All = false
		options.Filters["status"] = []string{"paused"}
		break
	case 5:
		options.All = false
		options.Filters["status"] = []string{"exited"}
		break
	case 6:
		options.All = false
		options.Filters["exited"] = []string{"0"}
		break
	}
	return options
}

// SearchContainers checks whether it contains key word or not
func SearchContainers(containers []docker.APIContainers, words []string) DockerContainers {
	results := DockerContainers{}
	for _, c := range containers {
		container := convertContainer(c)

		// hide specified named containers
		if !misc.ZeroOrNil(cfg.HiddenContainers) {
			names := container.Names
			hide := false
			if !misc.ZeroOrNil(cfg.LabelOverrideNames) {
				if value, found := container.Labels[cfg.LabelOverrideNames]; found {
					names = []string{"*" + value}
				}
			}
			for _, name := range names {
				hide = hide || inStringArray(cfg.HiddenContainers, name[1:])
			}
			if hide {
				continue
			}
		}

		// search words
		if container.contains(words) {
			results = append(results, container)
		}
	}
	sort.Sort(results)
	return results
}

func convertContainer(c docker.APIContainers) DockerContainer {
	container := DockerContainer{
		ID:         c.ID,
		Image:      c.Image,
		Command:    c.Command,
		Created:    c.Created,
		Status:     c.Status,
		Ports:      c.Ports,
		SizeRw:     c.SizeRw,
		SizeRootFs: c.SizeRootFs,
		Names:      make([]string, len(c.Names)),
		Labels:     c.Labels,
	}
	for idx, name := range c.Names {
		container.Names[idx] = name
	}
	return container
}

func (c DockerContainer) contains(words []string) bool {
	container := c.toUpperFields()
	match := true
	for _, word := range words {
		match = match && (strings.Contains(container.ID, word) ||
			strings.Contains(container.Image, word) ||
			strings.Contains(container.Command, word) ||
			strings.Contains(container.Status, word) ||
			inAPIPorts(container.Ports, word) ||
			inStringArray(container.Names, word) ||
			inMapString(container.Labels, word))
	}
	return match
}

func inAPIPorts(array []docker.APIPort, word string) bool {
	match := false
	for _, port := range array {
		match = match || strings.Contains(fmt.Sprint(port.PrivatePort), word)
		match = match || strings.Contains(fmt.Sprint(port.PublicPort), word)
		match = match || strings.Contains(port.Type, word)
		match = match || strings.Contains(port.IP, word)
	}
	return match
}

func inStringArray(array []string, word string) bool {
	match := false
	for _, value := range array {
		match = match || strings.Contains(value, word)
	}
	return match
}

func inMapString(m map[string]string, word string) bool {
	match := false
	for _, value := range m {
		match = match || strings.Contains(value, word)
	}
	return match
}

func (c DockerContainer) toUpperFields() DockerContainer {
	container := DockerContainer{}
	container.ID = strings.ToUpper(c.ID)
	container.Image = strings.ToUpper(c.Image)
	container.Command = strings.ToUpper(c.Command)
	container.Created = c.Created
	container.Status = strings.ToUpper(c.Status)
	container.Ports = make([]docker.APIPort, len(c.Ports))
	for idx, port := range c.Ports {
		container.Ports[idx] = docker.APIPort{
			PrivatePort: port.PrivatePort,
			PublicPort:  port.PublicPort,
			Type:        strings.ToUpper(port.Type),
			IP:          strings.ToUpper(port.IP),
		}
	}
	container.SizeRw = c.SizeRw
	container.SizeRootFs = c.SizeRootFs
	container.Names = make([]string, len(c.Names))
	for idx, name := range c.Names {
		container.Names[idx] = strings.ToUpper(name)
	}
	container.Labels = map[string]string{}
	for key, value := range c.Labels {
		container.Labels[key] = strings.ToUpper(value)
	}
	return container
}

func (cons DockerContainers) Len() int {
	return len(cons)
}

func (cons DockerContainers) Swap(i, j int) {
	cons[i], cons[j] = cons[j], cons[i]
}

func (cons DockerContainers) Less(i, j int) bool {
	a, b := cons[i], cons[j]
	if len(a.Names) > 0 && len(b.Names) > 0 {
		return a.Names[0] < b.Names[0]
	}
	return a.Created < b.Created
}
