package config

import "time"

// Config defines the application configurations
type Config struct {
	Name                   string `trim:"true"`
	Port                   uint16
	LogLevel               int
	DockerEndpoint         string `trim:"true"`
	DockerAPIVersion       string `trim:"true"`
	DockerCertPath         string `trim:"true"`
	DockerPullBeginTimeout time.Duration
	DockerPullTimeout      time.Duration
	DockerStatTimeout      time.Duration
	DockerStartTimeout     time.Duration
	DockerStopTimeout      time.Duration
	DockerRestartTimeout   time.Duration
	DockerKillTimeout      time.Duration
	DockerRmTimeout        time.Duration
	DockerCommitTimeout    time.Duration
	StaticFileHost         string `trim:"true"`
	StaticFilePath         string `trim:"true"`
	PreventSelfStop        bool
}
