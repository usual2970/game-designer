package commands

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/example/game-designer-cli/internal/preflight"
	"github.com/example/game-designer-cli/internal/provider"
	"github.com/example/game-designer-cli/internal/provider/fake"
	"github.com/example/game-designer-cli/internal/provider/threeos"
	"github.com/example/game-designer-cli/internal/reporting"
	"github.com/spf13/cobra"
)

var version = "0.1.0"

type DeployOptions struct {
	ServerPath string
	AppName    string
	Env        string
	Region     string

	BaseURL      string
	Identifier   string
	Password     string
	Mode         string
	GameURI      string
	GameName     string
	GameDesc     string
	GameLogo     string
	PackagePath  string
	SQLPath      string
	VersionStr   string
	ChangeLog    string
	ReviewURI    string
	Page         int
	PageSize     int
	ScreenType   int
	HalfSupport  int
	HalfRatio    string
	BackendDir   string
	BackendCmd   string
	FrontendDir  string
	FrontendCmd  string
	SocketDir    string
	SocketCmd    string
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
			pFlag, _ := cmd.Flags().GetString("provider")
			deployConfig := buildDeployConfig(opts)

			if pFlag == "3os" {
				if err := validateProductionConfig(deployConfig); err != nil {
					r := reporting.FailResult(reporting.CodeConfigError, err.Error(), nil)
					fmt.Println(r.ToJSON())
					os.Exit(1)
				}
			}

			prov, err := resolveProvider(pFlag, deployConfig)
			if err != nil {
				r := reporting.FailResult(reporting.CodeConfigError, err.Error(), nil)
				fmt.Println(r.ToJSON())
				os.Exit(1)
			}

			ctx := context.Background()

			results := preflight.RunChecks(opts.ServerPath)
			if !preflight.AllPassed(results) {
				r := reporting.FailResult(reporting.CodePreflightFailed, "Preflight checks failed", results)
				fmt.Println(r.ToJSON())
				os.Exit(1)
			}

			deployResult, err := prov.Deploy(ctx, deployConfig)
			if err != nil {
				code := classifyDeployError(err)
				r := reporting.FailResult(code, fmt.Sprintf("Deploy failed: %v", err), deployResult)
				fmt.Println(r.ToJSON())
				os.Exit(1)
			}

			healthResult, err := prov.HealthCheck(ctx, deployResult.URL)
			if err != nil || !healthResult.Healthy {
				r := reporting.FailResult(reporting.CodeHealthCheckFail, fmt.Sprintf("Health check failed: %v", err), deployResult)
				fmt.Println(r.ToJSON())
				os.Exit(1)
			}

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
	cmd.Flags().String("provider", "fake", "Deployment provider (fake or 3os)")

	cmd.Flags().StringVar(&opts.BaseURL, "base-url", envOrDefault("GD_BASE_URL", "https://api.3sdk.yu3.co"), "API base URL")
	cmd.Flags().StringVar(&opts.Identifier, "identifier", envOrDefault("GD_IDENTIFIER", ""), "Login identifier (email)")
	cmd.Flags().StringVar(&opts.Password, "password", envOrDefault("GD_PASSWORD", ""), "Login password")

	cmd.Flags().StringVar(&opts.Mode, "mode", "", "Publish mode: create, update-info, update-version, list, apply-review")

	cmd.Flags().StringVar(&opts.GameURI, "game-uri", envOrDefault("GD_GAME_URI", ""), "Game URI (for update/review)")
	cmd.Flags().StringVar(&opts.GameName, "game-name", "", "Game name (for create)")
	cmd.Flags().StringVar(&opts.GameDesc, "game-desc", "", "Game description")
	cmd.Flags().StringVar(&opts.GameLogo, "game-logo", "", "Game logo URL")

	cmd.Flags().StringVar(&opts.PackagePath, "package-path", "", "Local path to game package (.zip)")
	cmd.Flags().StringVar(&opts.SQLPath, "sql-path", "", "Local path to init SQL file")
	cmd.Flags().StringVar(&opts.VersionStr, "version", "", "Game version (e.g. 1.0.0)")
	cmd.Flags().StringVar(&opts.ChangeLog, "change-log", "", "Version change log")

	cmd.Flags().IntVar(&opts.ScreenType, "screen-type", 0, "Screen type: 1=vertical, 2=horizontal")
	cmd.Flags().IntVar(&opts.HalfSupport, "half-support", 0, "Half screen support: 1=none, 2=supported")
	cmd.Flags().StringVar(&opts.HalfRatio, "half-ratio", "", "Half screen ratio (e.g. 0.75)")

	cmd.Flags().StringVar(&opts.BackendDir, "backend-dir", "", "Backend working directory")
	cmd.Flags().StringVar(&opts.BackendCmd, "backend-cmd", "", "Backend start command")
	cmd.Flags().StringVar(&opts.FrontendDir, "frontend-dir", "", "Frontend working directory")
	cmd.Flags().StringVar(&opts.FrontendCmd, "frontend-cmd", "", "Frontend start command")
	cmd.Flags().StringVar(&opts.SocketDir, "socket-dir", "", "Socket working directory")
	cmd.Flags().StringVar(&opts.SocketCmd, "socket-cmd", "", "Socket start command")

	cmd.Flags().IntVar(&opts.Page, "page", 1, "Page number (list mode)")
	cmd.Flags().IntVar(&opts.PageSize, "page-size", 10, "Page size (list mode)")

	cmd.Flags().StringVar(&opts.ReviewURI, "review-uri", "", "Review URI (apply-review mode)")

	return cmd
}

func resolveProvider(name string, cfg provider.DeployConfig) (provider.Provider, error) {
	switch name {
	case "fake":
		return fake.New(), nil
	case "3os":
		client := threeos.NewClient(http.DefaultClient, cfg.BaseURL)
		uploader := threeos.NewOSSUploader(http.DefaultClient)
		return threeos.NewProvider(client, uploader), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s (accepted: fake, 3os)", name)
	}
}

func classifyDeployError(err error) reporting.ResultCode {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "auth failed"):
		return reporting.CodeAuthFailed
	case strings.Contains(msg, "list games failed"):
		return reporting.CodeListFailed
	case strings.Contains(msg, "lookup"):
		return reporting.CodeLookupFailed
	case strings.Contains(msg, "upload") || strings.Contains(msg, "get upload policy"):
		return reporting.CodeUploadFailed
	case strings.Contains(msg, "review"):
		return reporting.CodeReviewFailed
	case strings.Contains(msg, "create game") || strings.Contains(msg, "update"):
		return reporting.CodePublishFailed
	default:
		return reporting.CodeDeployFailed
	}
}

func buildDeployConfig(opts DeployOptions) provider.DeployConfig {
	cfg := provider.DeployConfig{
		ServerPath:     opts.ServerPath,
		AppName:        opts.AppName,
		Env:            opts.Env,
		Region:         opts.Region,
		BaseURL:        opts.BaseURL,
		Identifier:     opts.Identifier,
		Password:       opts.Password,
		Mode:           provider.PublishMode(opts.Mode),
		GameURI:        opts.GameURI,
		GameName:       opts.GameName,
		GameDescription: opts.GameDesc,
		GameLogo:       opts.GameLogo,
		PackagePath:    opts.PackagePath,
		SQLPath:        opts.SQLPath,
		Version:        opts.VersionStr,
		ChangeLog:      opts.ChangeLog,
		Page:           opts.Page,
		PageSize:       opts.PageSize,
		ReviewURI:      opts.ReviewURI,
	}

	if opts.ScreenType > 0 {
		cfg.ScreenConfig = &provider.ScreenConfig{
			ScreenType:  opts.ScreenType,
			HalfSupport: opts.HalfSupport,
			HalfRatio:   opts.HalfRatio,
		}
	}

	if opts.BackendDir != "" || opts.SocketDir != "" || opts.FrontendDir != "" {
		cfg.BuildConfig = &provider.BuildConfig{
			Backend:  provider.BuildConfigEntry{WorkDir: opts.BackendDir, Cmd: opts.BackendCmd},
			Frontend: provider.BuildConfigEntry{WorkDir: opts.FrontendDir, Cmd: opts.FrontendCmd},
			Socket:   provider.BuildConfigEntry{WorkDir: opts.SocketDir, Cmd: opts.SocketCmd},
		}
	}

	return cfg
}

func validateProductionConfig(cfg provider.DeployConfig) error {
	if cfg.Identifier == "" {
		return fmt.Errorf("--identifier or GD_IDENTIFIER env is required for production")
	}
	if cfg.Password == "" {
		return fmt.Errorf("--password or GD_PASSWORD env is required for production")
	}

	switch cfg.Mode {
	case provider.PublishModeCreate:
		if cfg.PackagePath == "" {
			return fmt.Errorf("--package-path is required for create mode")
		}
	case provider.PublishModeUpdateVersion:
		if cfg.GameURI == "" {
			return fmt.Errorf("--game-uri is required for update-version mode")
		}
		if cfg.PackagePath == "" {
			return fmt.Errorf("--package-path is required for update-version mode")
		}
	case provider.PublishModeUpdateInfo:
		if cfg.GameURI == "" {
			return fmt.Errorf("--game-uri is required for update-info mode")
		}
	case provider.PublishModeApplyReview:
		if cfg.ReviewURI == "" {
			return fmt.Errorf("--review-uri is required for apply-review mode")
		}
	case provider.PublishModeList:
	default:
		return fmt.Errorf("unsupported mode: %s (accepted: create, update-info, update-version, list, apply-review)", cfg.Mode)
	}
	return nil
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
