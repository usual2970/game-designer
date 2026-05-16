package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/example/game-designer-cli/internal/preflight"
	"github.com/example/game-designer-cli/internal/provider"
	"github.com/example/game-designer-cli/internal/reporting"
	"github.com/spf13/cobra"
)

var version = "0.1.0"

type DeployOptions struct {
	ServerPath string
	AppName    string
	Env        string
	Region     string
}

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "game-designer",
		Short: "Game Designer Server deploy CLI",
	}

	root.AddCommand(newVersionCmd())
	root.AddCommand(newPreflightCmd())
	root.AddCommand(newDeployCmd())
	return root
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("game-designer %s\n", version)
		},
	}
}

func newPreflightCmd() *cobra.Command {
	var serverPath string
	cmd := &cobra.Command{
		Use:   "preflight",
		Short: "Run pre-deploy checks",
		RunE: func(cmd *cobra.Command, args []string) error {
			results := preflight.RunChecks(serverPath)
			for _, r := range results {
				status := "PASS"
				if !r.Passed {
					status = "FAIL"
				}
				fmt.Printf("[%s] %s: %s\n", status, r.Name, r.Message)
			}

			if preflight.AllPassed(results) {
				fmt.Println(reporting.SuccessResult("All preflight checks passed", results).ToJSON())
				return nil
			}
			fmt.Println(reporting.FailResult(reporting.CodePreflightFailed, "Preflight checks failed", results).ToJSON())
			os.Exit(1)
			return nil
		},
	}
	cmd.Flags().StringVar(&serverPath, "server-path", ".", "Path to the server template")
	return cmd
}

func newDeployCmd() *cobra.Command {
	var opts DeployOptions
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy the game server to PaaS",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, _ := cmd.Flags().GetString("provider")
			prov := resolveProvider(p)

			ctx := context.Background()

			// Preflight
			results := preflight.RunChecks(opts.ServerPath)
			if !preflight.AllPassed(results) {
				r := reporting.FailResult(reporting.CodePreflightFailed, "Preflight checks failed", results)
				fmt.Println(r.ToJSON())
				os.Exit(1)
			}

			// Deploy
			deployConfig := provider.DeployConfig{
				ServerPath: opts.ServerPath,
				AppName:    opts.AppName,
				Env:        opts.Env,
				Region:     opts.Region,
			}
			deployResult, err := prov.Deploy(ctx, deployConfig)
			if err != nil {
				r := reporting.FailResult(reporting.CodeDeployFailed, fmt.Sprintf("Deploy failed: %v", err), nil)
				fmt.Println(r.ToJSON())
				os.Exit(1)
			}

			// Health check
			healthResult, err := prov.HealthCheck(ctx, deployResult.URL)
			if err != nil || !healthResult.Healthy {
				r := reporting.FailResult(reporting.CodeHealthCheckFail, fmt.Sprintf("Health check failed: %v", err), deployResult)
				fmt.Println(r.ToJSON())
				os.Exit(1)
			}

			// Status
			statusResult, _ := prov.Status(ctx, deployConfig)

			details := map[string]interface{}{
				"deploy": deployResult,
				"health": healthResult,
				"status": statusResult,
			}
			r := reporting.SuccessResult(fmt.Sprintf("Deployed to %s", deployResult.URL), details)
			fmt.Println(r.ToJSON())
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.ServerPath, "server-path", ".", "Path to the server template")
	cmd.Flags().StringVar(&opts.AppName, "app-name", "game-server", "Application name")
	cmd.Flags().StringVar(&opts.Env, "env", "production", "Deployment environment")
	cmd.Flags().StringVar(&opts.Region, "region", "default", "Deployment region")
	cmd.Flags().String("provider", "fake", "Deployment provider")
	return cmd
}

func resolveProvider(name string) provider.Provider {
	// MVP only has fake provider
	return &fakeProviderAdapter{}
}

type fakeProviderAdapter struct{}

func (f *fakeProviderAdapter) Name() string { return "fake" }
func (f *fakeProviderAdapter) Deploy(ctx context.Context, config provider.DeployConfig) (*provider.DeployResult, error) {
	return &provider.DeployResult{
		URL:      fmt.Sprintf("https://%s.fake.local", config.AppName),
		Version:  "v0.1.0",
		AppName:  config.AppName,
		Provider: "fake",
	}, nil
}
func (f *fakeProviderAdapter) Status(ctx context.Context, config provider.DeployConfig) (*provider.StatusResult, error) {
	return &provider.StatusResult{
		Status:  "running",
		Healthy: true,
		URL:     fmt.Sprintf("https://%s.fake.local", config.AppName),
		Version: "v0.1.0",
	}, nil
}
func (f *fakeProviderAdapter) HealthCheck(ctx context.Context, url string) (*provider.HealthResult, error) {
	return &provider.HealthResult{
		Healthy: true,
		Latency: "1ms",
		URL:     url,
	}, nil
}
