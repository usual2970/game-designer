package fake

import (
	"context"
	"testing"

	"github.com/example/game-designer-cli/internal/provider"
)

func TestFakeProvider_Deploy(t *testing.T) {
	p := New()
	result, err := p.Deploy(context.Background(), provider.DeployConfig{
		AppName: "test-app",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.AppName != "test-app" {
		t.Errorf("expected appName=test-app, got %s", result.AppName)
	}
	if result.Provider != "fake" {
		t.Errorf("expected provider=fake, got %s", result.Provider)
	}
	if !p.DeployCalled {
		t.Error("expected DeployCalled=true")
	}
}

func TestFakeProvider_DeployFailure(t *testing.T) {
	p := New()
	p.DeployShouldFail = true

	_, err := p.Deploy(context.Background(), provider.DeployConfig{})
	if err == nil {
		t.Error("expected error when DeployShouldFail=true")
	}
}

func TestFakeProvider_Status(t *testing.T) {
	p := New()
	result, err := p.Status(context.Background(), provider.DeployConfig{AppName: "test-app"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != "running" {
		t.Errorf("expected status=running, got %s", result.Status)
	}
	if !result.Healthy {
		t.Error("expected healthy=true")
	}
	if !p.StatusCalled {
		t.Error("expected StatusCalled=true")
	}
}

func TestFakeProvider_HealthCheck(t *testing.T) {
	p := New()
	result, err := p.HealthCheck(context.Background(), "https://test.fake.local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Healthy {
		t.Error("expected healthy=true")
	}
	if !p.HealthCheckCalled {
		t.Error("expected HealthCheckCalled=true")
	}
}
