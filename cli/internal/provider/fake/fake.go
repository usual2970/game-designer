package fake

import (
	"context"
	"fmt"

	"github.com/example/game-designer-cli/internal/provider"
)

type FakeProvider struct {
	DeployShouldFail   bool
	StatusShouldFail   bool
	HealthShouldFail   bool
	DeployCalled       bool
	StatusCalled       bool
	HealthCheckCalled  bool
}

func New() *FakeProvider {
	return &FakeProvider{}
}

func (f *FakeProvider) Name() string {
	return "fake"
}

func (f *FakeProvider) Deploy(ctx context.Context, config provider.DeployConfig) (*provider.DeployResult, error) {
	f.DeployCalled = true
	if f.DeployShouldFail {
		return nil, fmt.Errorf("fake deploy failure")
	}
	return &provider.DeployResult{
		URL:      fmt.Sprintf("https://%s.fake.local", config.AppName),
		Version:  "v0.1.0-fake",
		AppName:  config.AppName,
		Provider: "fake",
	}, nil
}

func (f *FakeProvider) Status(ctx context.Context, config provider.DeployConfig) (*provider.StatusResult, error) {
	f.StatusCalled = true
	if f.StatusShouldFail {
		return nil, fmt.Errorf("fake status failure")
	}
	return &provider.StatusResult{
		Status:  "running",
		Healthy: true,
		URL:     fmt.Sprintf("https://%s.fake.local", config.AppName),
		Version: "v0.1.0-fake",
	}, nil
}

func (f *FakeProvider) HealthCheck(ctx context.Context, url string) (*provider.HealthResult, error) {
	f.HealthCheckCalled = true
	if f.HealthShouldFail {
		return nil, fmt.Errorf("fake health check failure")
	}
	return &provider.HealthResult{
		Healthy: true,
		Latency: "1ms",
		URL:     url,
	}, nil
}
