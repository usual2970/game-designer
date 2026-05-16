package provider

import "context"

type DeployConfig struct {
	ServerPath string `json:"serverPath"`
	Env        string `json:"env"`
	Region     string `json:"region"`
	AppName    string `json:"appName"`
}

type DeployResult struct {
	URL      string `json:"url"`
	Version  string `json:"version"`
	AppName  string `json:"appName"`
	Provider string `json:"provider"`
}

type StatusResult struct {
	Status   string `json:"status"`
	Healthy  bool   `json:"healthy"`
	URL      string `json:"url"`
	Version  string `json:"version"`
}

type HealthResult struct {
	Healthy bool   `json:"healthy"`
	Latency string `json:"latency"`
	URL     string `json:"url"`
}

type Provider interface {
	Name() string
	Deploy(ctx context.Context, config DeployConfig) (*DeployResult, error)
	Status(ctx context.Context, config DeployConfig) (*StatusResult, error)
	HealthCheck(ctx context.Context, url string) (*HealthResult, error)
}
