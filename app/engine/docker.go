package engine

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/pottava/docker-webui/app/config"
	"github.com/pottava/docker-webui/app/logs"
	"github.com/pottava/docker-webui/app/misc"
	"github.com/pottava/docker-webui/app/models"

	api "github.com/fsouza/go-dockerclient"
)

// Client represents wrapped api.Client
type Client struct {
	*api.Client
	Conf *models.DockerClient
}

// ContainerMetadata represents docker response
type ContainerMetadata struct {
	Container *api.Container
	Error     error
}

// ImageMetadata represents docker response
type ImageMetadata struct {
	Image *api.Image
	Error error
}

var current *models.DockerClient

var cfg *config.Config
var containerID string
var pullLock sync.Mutex

func init() {
	cfg = config.NewConfig()

	if cfg.PreventSelfStop {
		if candidate, err := misc.ShellExec([]string{"bash", "-c",
			"cat /proc/self/cgroup | grep -o -e 'docker.*' | head -n 1"}); err == nil {
			pattern := regexp.MustCompile(`^docker[^0-9a-zA-Z]*`)
			candidate = pattern.ReplaceAllString(candidate, "")
			if len(candidate) >= 64 {
				containerID = candidate[:64]
			}
		}
		logs.Debug.Printf("Docker container ID: %s", containerID)
	}
	for index, endpoint := range cfg.DockerEndpoints {
		if len(cfg.DockerCertPath) > index {
			Configure(endpoint, cfg.DockerCertPath[index], true)
		} else {
			Configure(endpoint, "", true)
		}
		Save()
	}
}

// Configure set client's configuration to current
func Configure(endpoint, certPath string, def bool) {
	current = &models.DockerClient{
		Endpoint:  endpoint,
		CertPath:  certPath,
		IsActive:  true,
		IsDefault: def,
	}
}

// Save persists current configuration
func Save() {
	current.Save()
}

// Docker generates a docker client
func Docker() (client *Client, err error) {
	var c *api.Client
	if misc.ZeroOrNil(current.CertPath) {
		c, err = api.NewClient(current.Endpoint)
	} else {
		cert := fmt.Sprintf("%s/cert.pem", current.CertPath) // X.509 Certificate
		key := fmt.Sprintf("%s/key.pem", current.CertPath)   // Private Key
		ca := fmt.Sprintf("%s/ca.pem", current.CertPath)     // Certificate authority
		c, err = api.NewTLSClient(current.Endpoint, cert, key, ca)
	}
	if !misc.ZeroOrNil(c) {
		err = c.Ping()
	}
	if misc.ZeroOrNil(err) {
		c.SkipServerVersionCheck = true
		return &Client{c, current}, nil
	}
	return nil, err
}

// InspectContainer inspects the docker container
func (c *Client) InspectContainer(id string) ContainerMetadata {
	container, err := c.Client.InspectContainer(id)
	if err != nil {
		return ContainerMetadata{
			Container: &api.Container{ID: ""},
			Error:     CannotXContainerError{"Inspect", err.Error()},
		}
	}
	return ContainerMetadata{Container: container}
}

// ListImages list docker images
func (c *Client) ListImages() []api.APIImages {
	images, err := c.Client.ListImages(api.ListImagesOptions{All: false})
	if err != nil {
		images = []api.APIImages{}
	}
	return images
}

// InspectImage inspects the docker image
func (c *Client) InspectImage(id string) ImageMetadata {
	image, err := c.Client.InspectImage(id)
	if err != nil {
		return ImageMetadata{
			Image: &api.Image{ID: ""},
			Error: CannotXContainerError{"Inspect", err.Error()},
		}
	}
	return ImageMetadata{Image: image}
}

// Top returns processes
func (c *Client) Top(id, args string) api.TopResult {
	processes, err := c.TopContainer(id, args)
	if err != nil {
		processes = api.TopResult{}
	}
	return processes
}

// Stats returns container statistics
func (c *Client) Stats(id string, count int) (result []*api.Stats, err error) {
	e := make(chan error, 1)
	s := make(chan *api.Stats)
	done := make(chan bool)
	go func() {
		e <- c.Client.Stats(api.StatsOptions{
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
func (c *Client) Logs(id string, line int, timeout time.Duration) (stdout []string, stderr []string, err error) {
	stdout = []string{}
	stderr = []string{}

	cStdOut := make(chan string)
	cStdErr := make(chan string)

	e := make(chan error, 1)
	done := make(chan bool)
	go func() {
		e <- c.LogStream(LogsOptions{
			ID:      id,
			Tail:    int64(line),
			Stdout:  cStdOut,
			Stderr:  cStdErr,
			Done:    done,
			Timeout: timeout,
		})
		close(e)
	}()

	for {
		select {
		case err = <-e:
			return
		case log, ok := <-cStdOut:
			if !ok {
				err = <-e
				return
			}
			stdout = append(stdout, log)
			if (line > 0) && ((len(stdout) + len(stderr)) >= line) {
				done <- true
				err = <-e
				return
			}
		case log, ok := <-cStdErr:
			if !ok {
				err = <-e
				return
			}
			stderr = append(stderr, log)
			if (line > 0) && ((len(stdout) + len(stderr)) >= line) {
				done <- true
				err = <-e
				return
			}
		}
	}
}

// LogsOptions can be
type LogsOptions struct {
	ID      string
	Tail    int64
	Stdout  chan<- string
	Stderr  chan<- string
	Done    <-chan bool
	Timeout time.Duration
}

// LogStream returns containers logs as streams
func (c *Client) LogStream(opts LogsOptions) (err error) {
	timeout := time.After(opts.Timeout)
	outReader, outWriter := io.Pipe()
	errReader, errWriter := io.Pipe()

	defer func() {
		close(opts.Stdout)
		close(opts.Stderr)
	}()

	done := make(chan bool, 1)
	go func() {
		tail := "all"
		if opts.Tail > 0 {
			tail = fmt.Sprint(opts.Tail)
		}
		err = c.Client.Logs(api.LogsOptions{
			Container:    opts.ID,
			OutputStream: outWriter,
			ErrorStream:  errWriter,
			// TODO
			// There's no way to stop this goroutine if this flag is true for now.
			// Have to turn to true if the issue (https://github.com/fsouza/go-Client/issues/298) is closed.
			Follow:     false,
			Stdout:     true,
			Stderr:     true,
			Timestamps: true,
			Tail:       tail,
		})
		outWriter.Close()
		errWriter.Close()
		done <- true
		close(done)
	}()

	outDone := make(chan bool, 1)
	go func() {
		reader := bufio.NewReader(outReader)
		for {
			line, _, e := reader.ReadLine()
			if e != nil {
				outDone <- true
				close(outDone)
				return
			}
			opts.Stdout <- string(line)
		}
	}()
	errDone := make(chan bool, 1)
	go func() {
		reader := bufio.NewReader(errReader)
		for {
			line, _, e := reader.ReadLine()
			if e != nil {
				errDone <- true
				close(errDone)
				return
			}
			opts.Stderr <- string(line)
		}
	}()

	// block here waiting for the signal to stop function
	select {
	case <-opts.Done:
	case <-timeout:
	}
	outReader.Close()
	errReader.Close()
	<-done
	<-outDone
	<-errDone
	return
}

// Changes returns containers changed files
func (c *Client) Changes(id string) []api.Change {
	changes, err := c.ContainerChanges(id)
	if err != nil {
		changes = []api.Change{}
	}
	return changes
}

// History returns its history
func (c *Client) History(id string) []api.ImageHistory {
	history, err := c.ImageHistory(id)
	if err != nil {
		history = []api.ImageHistory{}
	}
	return history
}

// Rename renames the container
func (c *Client) Rename(id, name string) error {
	return c.RenameContainer(api.RenameContainerOptions{
		ID:   id,
		Name: name,
	})
}

// Pull pulls docker images more safely
func (c *Client) Pull(image string) ImageMetadata {
	timeout := time.After(cfg.DockerPullTimeout)

	pullLock.Lock()
	defer pullLock.Unlock()

	response := make(chan ImageMetadata, 1)
	go func() { response <- c.pull(image) }()

	select {
	case resp := <-response:
		return resp
	case <-timeout:
		return ImageMetadata{
			Image: &api.Image{ID: ""},
			Error: &DockerTimeoutError{cfg.DockerPullTimeout, "pulled"},
		}
	}
}

func (c *Client) pull(image string) ImageMetadata {
	reader, writer := io.Pipe()
	defer writer.Close()

	repository, tag := api.ParseRepositoryTag(image)
	tag = misc.NVL(tag, "latest")
	opts := api.PullImageOptions{
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
		finished <- c.PullImage(opts, api.AuthConfiguration{})
		logs.Debug.Printf("Pull completed for image: %v", opts.Repository)
	}()

	// wait for the pulling to begin
	select {
	case <-began:
		break
	case err := <-finished:
		if err != nil {
			return ImageMetadata{
				Image: &api.Image{ID: ""},
				Error: CannotXContainerError{"Pull", err.Error()},
			}
		}
		return c.InspectImage(opts.Repository)
	case <-timeout:
		return ImageMetadata{
			Image: &api.Image{ID: ""},
			Error: &DockerTimeoutError{cfg.DockerPullBeginTimeout, "pullBegin"},
		}
	}

	// wait for the completion
	err := <-finished
	if err != nil {
		return ImageMetadata{
			Image: &api.Image{ID: ""},
			Error: CannotXContainerError{"Pull", err.Error()},
		}
	}
	return c.InspectImage(opts.Repository)
}

// Create creates docker containers more safely
func (c *Client) Create(name string, config *api.Config, host *api.HostConfig) ContainerMetadata {
	timeout := time.After(cfg.DockerStartTimeout)

	ctx, cancel := context.WithCancel(context.TODO())
	response := make(chan ContainerMetadata, 1)
	go func() {
		response <- c.create(ctx, api.CreateContainerOptions{
			Name:       name,
			Config:     config,
			HostConfig: host,
		})
	}()

	select {
	case resp := <-response:
		return resp
	case <-timeout:
		cancel()
		return ContainerMetadata{
			Container: &api.Container{ID: ""},
			Error:     &DockerTimeoutError{cfg.DockerStartTimeout, "run"},
		}
	}
}

func (c *Client) create(ctx context.Context, opt api.CreateContainerOptions) ContainerMetadata {
	ch := make(chan ContainerMetadata, 1)
	go func() {
		container, err := c.CreateContainer(opt)
		ch <- ContainerMetadata{container, err}
	}()
	select {
	case meta := <-ch:
		return meta
	case <-ctx.Done():
		return ContainerMetadata{
			Container: &api.Container{ID: ""},
		}
	}
}

// Restart restarts docker containers more safely
func (c *Client) Restart(id string, wait uint) ContainerMetadata {
	timeout := time.After(cfg.DockerRestartTimeout)

	ctx, cancel := context.WithCancel(context.TODO())
	response := make(chan ContainerMetadata, 1)
	go func() { response <- c.restart(ctx, id, wait) }()

	select {
	case resp := <-response:
		return resp
	case <-timeout:
		cancel()
		return ContainerMetadata{
			Container: &api.Container{ID: ""},
			Error:     &DockerTimeoutError{cfg.DockerRestartTimeout, "restarting"},
		}
	}
}

func (c *Client) restart(ctx context.Context, id string, wait uint) ContainerMetadata {
	if cfg.PreventSelfStop {
		info := c.InspectContainer(id)
		if containerID == info.Container.ID[:64] {
			return ContainerMetadata{
				Container: &api.Container{ID: info.Container.ID},
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
		return ContainerMetadata{
			Container: &api.Container{ID: id},
		}
	}
}

// Start starts docker containers more safely
func (c *Client) Start(id string) ContainerMetadata {
	timeout := time.After(cfg.DockerStartTimeout)

	ctx, cancel := context.WithCancel(context.TODO())
	response := make(chan ContainerMetadata, 1)
	go func() { response <- c.start(ctx, id) }()

	select {
	case resp := <-response:
		return resp
	case <-timeout:
		cancel()
		return ContainerMetadata{
			Container: &api.Container{ID: ""},
			Error:     &DockerTimeoutError{cfg.DockerStartTimeout, "starting"},
		}
	}
}

func (c *Client) start(ctx context.Context, id string) ContainerMetadata {
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
		return ContainerMetadata{
			Container: &api.Container{ID: id},
		}
	}
}

// Stop stops docker containers more safely
func (c *Client) Stop(id string) ContainerMetadata {
	if cfg.PreventSelfStop {
		info := c.InspectContainer(id)
		if containerID == info.Container.ID[:64] {
			return ContainerMetadata{
				Container: &api.Container{ID: info.Container.ID},
				Error:     &CannotXContainerError{" stop ", "Prevented this application itself for stopping"},
			}
		}
	}
	timeout := time.After(cfg.DockerStopTimeout)

	ctx, cancel := context.WithCancel(context.TODO())
	response := make(chan ContainerMetadata, 1)
	go func() { response <- c.stop(ctx, id) }()

	select {
	case resp := <-response:
		return resp
	case <-timeout:
		cancel()
		return ContainerMetadata{
			Container: &api.Container{ID: ""},
			Error:     &DockerTimeoutError{cfg.DockerStopTimeout, "stopping"},
		}
	}
}

func (c *Client) stop(ctx context.Context, id string) ContainerMetadata {
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
		return ContainerMetadata{
			Container: &api.Container{ID: id},
		}
	}
}

// Kill kills docker containers more safely
func (c *Client) Kill(id string, wait uint) ContainerMetadata {
	timeout := time.After(cfg.DockerKillTimeout)

	ctx, cancel := context.WithCancel(context.TODO())
	response := make(chan ContainerMetadata, 1)
	go func() { response <- c.kill(ctx, id, wait) }()

	select {
	case resp := <-response:
		return resp
	case <-timeout:
		cancel()
		return ContainerMetadata{
			Container: &api.Container{ID: ""},
			Error:     &DockerTimeoutError{cfg.DockerKillTimeout, "killing"},
		}
	}
}

func (c *Client) kill(ctx context.Context, id string, wait uint) ContainerMetadata {
	if cfg.PreventSelfStop {
		info := c.InspectContainer(id)
		if containerID == info.Container.ID[:64] {
			return ContainerMetadata{
				Container: &api.Container{ID: info.Container.ID},
				Error:     &CannotXContainerError{" kill ", "Prevented this application itself for killing"},
			}
		}
	}
	ch := make(chan error, 1)
	go func() {
		ch <- c.KillContainer(api.KillContainerOptions{
			ID:     id,
			Signal: api.SIGKILL,
		})
	}()

	select {
	case err := <-ch:
		meta := c.InspectContainer(id)
		if err != nil {
			meta.Error = CannotXContainerError{"Start", err.Error()}
		}
		return meta
	case <-ctx.Done():
		return ContainerMetadata{
			Container: &api.Container{ID: id},
		}
	}
}

// Commit commit docker containers more safely
func (c *Client) Commit(id, repository, tag, message, author string) ImageMetadata {
	timeout := time.After(cfg.DockerCommitTimeout)
	response := make(chan ImageMetadata, 1)
	go func() {
		image, err := c.CommitContainer(api.CommitContainerOptions{
			Container:  id,
			Repository: repository,
			Tag:        tag,
			Message:    message,
			Author:     author,
		})
		response <- ImageMetadata{Image: image, Error: err}
	}()
	select {
	case resp := <-response:
		return resp
	case <-timeout:
		return ImageMetadata{
			Image: &api.Image{ID: ""},
			Error: &DockerTimeoutError{cfg.DockerCommitTimeout, "committing"}}
	}
}

// Tag tags docker images more safely
func (c *Client) Tag(id, repository, tag string) error {
	return c.TagImage(id, api.TagImageOptions{
		Repo: repository,
		Tag:  tag,
	})
}

// Rm removes docker containers more safely
func (c *Client) Rm(id string) error {
	timeout := time.After(cfg.DockerRmTimeout)
	response := make(chan error, 1)
	go func() {
		response <- c.RemoveContainer(api.RemoveContainerOptions{
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
func (c *Client) Rmi(id string) error {
	timeout := time.After(cfg.DockerRmTimeout)
	response := make(chan error, 1)
	go func() {
		response <- c.RemoveImageExtended(id, api.RemoveImageOptions{
			Force:   false,
			NoPrune: false,
		})
	}()
	select {
	case resp := <-response:
		if resp == nil {
			logs.Debug.Printf("Remove completed for image: %v", id)
		}
		return resp
	case <-timeout:
		return &DockerTimeoutError{cfg.DockerRmTimeout, "removing"}
	}
}
