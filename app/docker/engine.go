// Package engine wrapped docker client APIs
package engine

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/pkg/parsers"
	"github.com/pottava/docker-webui/app/config"
	"github.com/pottava/docker-webui/app/logs"
	"github.com/pottava/docker-webui/app/misc"

	docker "github.com/fsouza/go-dockerclient"
)

// DockerClient represents wrapped docker.Client
type DockerClient struct {
	*docker.Client
}

// DockerContainerMetadata represents docker response
type DockerContainerMetadata struct {
	Container *docker.Container
	Error     error
}

// DockerImageMetadata represents docker response
type DockerImageMetadata struct {
	Image *docker.Image
	Error error
}

var cfg *config.Config
var containerID string
var pullLock sync.Mutex

func init() {
	cfg = config.NewConfig()

	if cfg.PreventSelfStop {
		if candidate, err := misc.ShellExec([]string{"bash", "-c",
			"cat /proc/self/cgroup | grep -o -e 'docker-.*.scope' | head -n 1"}); err == nil {
			candidate = strings.Replace(candidate, "docker-", "", -1)
			candidate = strings.Replace(candidate, ".scope", "", -1)
			containerID = candidate[:64]
			logs.Debug.Printf("Docker container ID: %s", containerID)
		}
	}
}

// NewDockerClient returns DockerClient if it was generated successfully
func NewDockerClient() (*DockerClient, error) {
	client, err := docker.NewVersionedClient(cfg.DockerEndpoint, cfg.DockerAPIVersion)
	if err != nil {
		return nil, err
	}
	err = client.Ping()
	if err != nil {
		return nil, err
	}
	return &DockerClient{client}, nil
}

// InspectContainer inspects the docker container
func (c *DockerClient) InspectContainer(id string) DockerContainerMetadata {
	container, err := c.Client.InspectContainer(id)
	if err != nil {
		return DockerContainerMetadata{
			Container: &docker.Container{ID: ""},
			Error:     CannotXContainerError{"Inspect", err.Error()},
		}
	}
	return DockerContainerMetadata{Container: container}
}

// ListImages list docker images
func (c *DockerClient) ListImages() []docker.APIImages {
	images, err := c.Client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		images = []docker.APIImages{}
	}
	return images
}

// InspectImage inspects the docker image
func (c *DockerClient) InspectImage(id string) DockerImageMetadata {
	image, err := c.Client.InspectImage(id)
	if err != nil {
		return DockerImageMetadata{
			Image: &docker.Image{ID: ""},
			Error: CannotXContainerError{"Inspect", err.Error()},
		}
	}
	return DockerImageMetadata{Image: image}
}

// Top returns processes
func (c *DockerClient) Top(id, args string) docker.TopResult {
	processes, err := c.TopContainer(id, args)
	if err != nil {
		processes = docker.TopResult{}
	}
	return processes
}

// Stats returns container statistics
func (c *DockerClient) Stats(id string, count int) (result []*docker.Stats, err error) {
	e := make(chan error, 1)
	s := make(chan *docker.Stats)
	done := make(chan bool)
	go func() {
		e <- c.Client.Stats(docker.StatsOptions{
			ID:      id,
			Stats:   s,
			Stream:  true,
			Done:    done,
			Timeout: cfg.DockerStatTimeout,
		})
		close(e)
	}()

	for {
		stats, ok := <-s
		if !ok {
			break
		}
		// logs.Trace.Printf("%v", stats)
		result = append(result, stats)
		if len(result) >= count {
			done <- true
			return
		}
	}
	err = <-e
	return
}

// Logs returns containers logs
func (c *DockerClient) Logs(id string, since, line int64) (stdout, stderr string, err error) {
	sto := bytes.Buffer{}
	ste := bytes.Buffer{}
	tail := "all"
	if line > 0 {
		tail = fmt.Sprint(line)
	}
	err = c.Client.Logs(docker.LogsOptions{
		Container:    id,
		OutputStream: &sto,
		ErrorStream:  &ste,
		Follow:       false,
		Stdout:       true,
		Stderr:       true,
		Since:        since,
		Timestamps:   true,
		Tail:         tail,
		RawTerminal:  false,
	})
	return sto.String(), ste.String(), err
}

// Changes returns containers changed files
func (c *DockerClient) Changes(id string) []docker.Change {
	changes, err := c.ContainerChanges(id)
	if err != nil {
		changes = []docker.Change{}
	}
	return changes
}

// History returns its history
func (c *DockerClient) History(id string) []docker.ImageHistory {
	history, err := c.ImageHistory(id)
	if err != nil {
		history = []docker.ImageHistory{}
	}
	return history
}

// Pull pulls docker images more safely
func (c *DockerClient) Pull(image string) DockerImageMetadata {
	timeout := time.After(cfg.DockerPullTimeout)

	pullLock.Lock()
	defer pullLock.Unlock()

	response := make(chan DockerImageMetadata, 1)
	go func() { response <- c.pull(image) }()

	select {
	case resp := <-response:
		return resp
	case <-timeout:
		return DockerImageMetadata{
			Image: &docker.Image{ID: ""},
			Error: &DockerTimeoutError{cfg.DockerPullTimeout, "pulled"},
		}
	}
}

func (c *DockerClient) pull(image string) DockerImageMetadata {
	reader, writer := io.Pipe()
	defer writer.Close()

	repository, tag := parsers.ParseRepositoryTag(image)
	tag = misc.NVL(tag, "latest")
	opts := docker.PullImageOptions{
		Repository:   repository + ":" + tag,
		OutputStream: writer,
	}

	// check output goroutine
	began := make(chan bool, 1)
	once := sync.Once{}

	go func() {
		reader := bufio.NewReader(reader)
		var line string
		var err error
		for err == nil {
			line, err = reader.ReadString('\n')
			if err != nil {
				break
			}
			once.Do(func() {
				began <- true
			})
			if strings.Contains(line, "already being pulled by another client. Waiting.") {
				logs.Error.Printf("Image 'pull' status marked as already being pulled. image: %v, status: %v", opts.Repository, line)
			}
		}
		if err != nil && err != io.EOF {
			logs.Warn.Printf("Error reading pull image status. image: %v, err: %v", opts.Repository, err)
		}
	}()

	// pull the image
	timeout := time.After(cfg.DockerPullBeginTimeout)
	finished := make(chan error, 1)
	go func() {
		finished <- c.PullImage(opts, docker.AuthConfiguration{})
		logs.Debug.Printf("Pull completed for image: %v", opts.Repository)
	}()

	// wait for the pulling to begin
	select {
	case <-began:
		break
	case err := <-finished:
		if err != nil {
			return DockerImageMetadata{
				Image: &docker.Image{ID: ""},
				Error: CannotXContainerError{"Pull", err.Error()},
			}
		}
		return c.InspectImage(opts.Repository)
	case <-timeout:
		return DockerImageMetadata{
			Image: &docker.Image{ID: ""},
			Error: &DockerTimeoutError{cfg.DockerPullBeginTimeout, "pullBegin"},
		}
	}

	// wait for the completion
	err := <-finished
	if err != nil {
		return DockerImageMetadata{
			Image: &docker.Image{ID: ""},
			Error: CannotXContainerError{"Pull", err.Error()},
		}
	}
	return c.InspectImage(opts.Repository)
}

// Run runs docker containers more safely
func (c *DockerClient) Run(opt docker.CreateContainerOptions) DockerContainerMetadata {
	timeout := time.After(cfg.DockerStartTimeout)

	ctx, cancel := context.WithCancel(context.TODO())
	response := make(chan DockerContainerMetadata, 1)
	go func() { response <- c.run(ctx, opt) }()

	select {
	case resp := <-response:
		return resp
	case <-timeout:
		cancel()
		return DockerContainerMetadata{
			Container: &docker.Container{ID: ""},
			Error:     &DockerTimeoutError{cfg.DockerStartTimeout, "run"},
		}
	}
}

func (c *DockerClient) run(ctx context.Context, opt docker.CreateContainerOptions) DockerContainerMetadata {
	ch := make(chan DockerContainerMetadata, 1)
	go func() {
		container, err := c.CreateContainer(opt)
		ch <- DockerContainerMetadata{container, err}
	}()
	select {
	case meta := <-ch:
		return meta
	case <-ctx.Done():
		return DockerContainerMetadata{
			Container: &docker.Container{ID: ""},
		}
	}
}

// Start starts docker containers more safely
func (c *DockerClient) Start(id string) DockerContainerMetadata {
	timeout := time.After(cfg.DockerStartTimeout)

	ctx, cancel := context.WithCancel(context.TODO())
	response := make(chan DockerContainerMetadata, 1)
	go func() { response <- c.start(ctx, id) }()

	select {
	case resp := <-response:
		return resp
	case <-timeout:
		cancel()
		return DockerContainerMetadata{
			Container: &docker.Container{ID: ""},
			Error:     &DockerTimeoutError{cfg.DockerStartTimeout, "started"},
		}
	}
}

func (c *DockerClient) start(ctx context.Context, id string) DockerContainerMetadata {
	ch := make(chan error, 1)
	go func() { ch <- c.StartContainer(id, nil) }()

	select {
	case err := <-ch:
		meta := c.InspectContainer(id)
		if err != nil {
			meta.Error = CannotXContainerError{"Start", err.Error()}
		}
		return meta
	case <-ctx.Done():
		return DockerContainerMetadata{
			Container: &docker.Container{ID: id},
		}
	}
}

// Stop stops docker containers more safely
func (c *DockerClient) Stop(id string) DockerContainerMetadata {
	if cfg.PreventSelfStop {
		info := c.InspectContainer(id)
		if containerID == info.Container.ID[:64] {
			return DockerContainerMetadata{
				Container: &docker.Container{ID: info.Container.ID},
				Error:     &CannotXContainerError{" stop ", "Prevented this application itself for stopping"},
			}
		}
	}
	timeout := time.After(cfg.DockerStopTimeout)

	ctx, cancel := context.WithCancel(context.TODO())
	response := make(chan DockerContainerMetadata, 1)
	go func() { response <- c.stop(ctx, id) }()

	select {
	case resp := <-response:
		return resp
	case <-timeout:
		cancel()
		return DockerContainerMetadata{
			Container: &docker.Container{ID: ""},
			Error:     &DockerTimeoutError{cfg.DockerStopTimeout, "stopped"},
		}
	}
}

func (c *DockerClient) stop(ctx context.Context, id string) DockerContainerMetadata {
	ch := make(chan error, 1)
	go func() { ch <- c.StopContainer(id, 30) }()

	select {
	case err := <-ch:
		meta := c.InspectContainer(id)
		if err != nil {
			meta.Error = CannotXContainerError{"Stop", err.Error()}
		}
		return meta
	case <-ctx.Done():
		return DockerContainerMetadata{
			Container: &docker.Container{ID: id},
		}
	}
}

// Restart restarts docker containers more safely
func (c *DockerClient) Restart(id string, wait uint) DockerContainerMetadata {
	timeout := time.After(cfg.DockerRestartTimeout)

	ctx, cancel := context.WithCancel(context.TODO())
	response := make(chan DockerContainerMetadata, 1)
	go func() { response <- c.restart(ctx, id, wait) }()

	select {
	case resp := <-response:
		return resp
	case <-timeout:
		cancel()
		return DockerContainerMetadata{
			Container: &docker.Container{ID: ""},
			Error:     &DockerTimeoutError{cfg.DockerRestartTimeout, "restarted"},
		}
	}
}

func (c *DockerClient) restart(ctx context.Context, id string, wait uint) DockerContainerMetadata {
	if cfg.PreventSelfStop {
		info := c.InspectContainer(id)
		if containerID == info.Container.ID[:64] {
			return DockerContainerMetadata{
				Container: &docker.Container{ID: info.Container.ID},
				Error:     &CannotXContainerError{" restart ", "Prevented this application itself for restarting"},
			}
		}
	}
	ch := make(chan error, 1)
	go func() { ch <- c.RestartContainer(id, wait) }()

	select {
	case err := <-ch:
		meta := c.InspectContainer(id)
		if err != nil {
			meta.Error = CannotXContainerError{"Start", err.Error()}
		}
		return meta
	case <-ctx.Done():
		return DockerContainerMetadata{
			Container: &docker.Container{ID: id},
		}
	}
}

// Rm removes docker containers more safely
func (c *DockerClient) Rm(id string) error {
	timeout := time.After(cfg.DockerRmTimeout)
	response := make(chan error, 1)
	go func() {
		response <- c.RemoveContainer(docker.RemoveContainerOptions{
			ID:            id,
			RemoveVolumes: true,
			Force:         false,
		})
	}()
	select {
	case resp := <-response:
		return resp
	case <-timeout:
		return &DockerTimeoutError{cfg.DockerRmTimeout, "removing"}
	}
}

// Rmi removes docker images more safely
func (c *DockerClient) Rmi(id string) error {
	timeout := time.After(cfg.DockerRmTimeout)
	response := make(chan error, 1)
	go func() {
		response <- c.RemoveImageExtended(id, docker.RemoveImageOptions{
			Force:   false,
			NoPrune: false,
		})
	}()
	select {
	case resp := <-response:
		logs.Debug.Printf("Remove completed for image: %v", id)
		return resp
	case <-timeout:
		return &DockerTimeoutError{cfg.DockerRmTimeout, "removing"}
	}
}
