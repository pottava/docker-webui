package config

import "time"

// Config defines the application configurations
type Config struct {
	Name                   string `trim:"true"`
	Port                   uint16
	LogLevel               int
	LabelOverrideNames     string `trim:"true"`
	DockerEndpoints        []string
	DockerCertPath         []string
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
