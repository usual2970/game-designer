package commands

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/example/game-designer-cli/internal/reporting"
)

func TestVersionCmd(t *testing.T) {
	root := NewRootCmd()
	root.SetArgs([]string{"version"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPreflightCmd_ValidPath(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module example.com/test\ngo 1.24\n"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\nfunc main() {}\n"), 0644)

	root := NewRootCmd()
	root.SetArgs([]string{"preflight", "--server-path", tmpDir})

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeployCmd_FakeProvider(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module example.com/test\ngo 1.24\n"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\nfunc main() {}\n"), 0644)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	root := NewRootCmd()
	root.SetArgs([]string{
		"deploy",
		"--server-path", tmpDir,
		"--app-name", "test-game",
		"--provider", "fake",
	})

	err := root.Execute()
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	if !strings.Contains(output, "SUCCESS") {
		t.Errorf("expected SUCCESS in output, got: %s", output)
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	var lastLine string
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "{") {
			lastLine = lines[i]
			break
		}
	}

	if lastLine == "" {
		t.Fatalf("no JSON line found in output: %s", output)
	}

	var result reporting.Result
	if err := json.Unmarshal([]byte(lastLine), &result); err != nil {
		t.Fatalf("invalid JSON output: %v, line: %s", err, lastLine)
	}
	if !result.Success {
		t.Error("expected success=true in deploy result")
	}
}

// runSubprocessTest runs the current test binary as a subprocess with the
// given env flag, capturing combined output. This is needed because the
// deploy command calls os.Exit(1) on failure, which kills the test process.
func runSubprocessTest(t *testing.T, envKey string) string {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
	cmd.Env = append(os.Environ(), envKey+"=1")
	output, _ := cmd.CombinedOutput()
	return string(output)
}

func TestDeployCmd_UnsupportedProvider(t *testing.T) {
	const envKey = "GD_TEST_UNSUPPORTED_PROVIDER"
	if os.Getenv(envKey) == "1" {
		tmpDir := t.TempDir()
		os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module example.com/test\ngo 1.24\n"), 0644)
		os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\nfunc main() {}\n"), 0644)

		root := NewRootCmd()
		root.SetArgs([]string{
			"deploy",
			"--server-path", tmpDir,
			"--provider", "unknown",
		})
		root.Execute()
		return
	}

	output := runSubprocessTest(t, envKey)
	if !strings.Contains(output, "CONFIG_ERROR") {
		t.Errorf("expected CONFIG_ERROR for unsupported provider, got: %s", output)
	}
}

func TestDeployCmd_ProductionMissingCredentials(t *testing.T) {
	const envKey = "GD_TEST_MISSING_CREDS"
	if os.Getenv(envKey) == "1" {
		tmpDir := t.TempDir()
		os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module example.com/test\ngo 1.24\n"), 0644)
		os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\nfunc main() {}\n"), 0644)

		root := NewRootCmd()
		root.SetArgs([]string{
			"deploy",
			"--server-path", tmpDir,
			"--provider", "3os",
			"--mode", "list",
		})
		root.Execute()
		return
	}

	output := runSubprocessTest(t, envKey)
	if !strings.Contains(output, "CONFIG_ERROR") {
		t.Errorf("expected CONFIG_ERROR for missing credentials, got: %s", output)
	}
}

func TestDeployCmd_ProductionCreateMissingPackage(t *testing.T) {
	const envKey = "GD_TEST_CREATE_NO_PKG"
	if os.Getenv(envKey) == "1" {
		tmpDir := t.TempDir()
		os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module example.com/test\ngo 1.24\n"), 0644)
		os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\nfunc main() {}\n"), 0644)

		root := NewRootCmd()
		root.SetArgs([]string{
			"deploy",
			"--server-path", tmpDir,
			"--provider", "3os",
			"--mode", "create",
			"--identifier", "test@example.com",
			"--password", "testpass",
		})
		root.Execute()
		return
	}

	output := runSubprocessTest(t, envKey)
	if !strings.Contains(output, "CONFIG_ERROR") {
		t.Errorf("expected CONFIG_ERROR for missing package-path, got: %s", output)
	}
	if !strings.Contains(output, "package-path") {
		t.Errorf("expected error to mention package-path, got: %s", output)
	}
}

func TestDeployCmd_ProductionUpdateVersionMissingGameURI(t *testing.T) {
	const envKey = "GD_TEST_UPDATE_NO_URI"
	if os.Getenv(envKey) == "1" {
		tmpDir := t.TempDir()
		os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module example.com/test\ngo 1.24\n"), 0644)
		os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\nfunc main() {}\n"), 0644)

		root := NewRootCmd()
		root.SetArgs([]string{
			"deploy",
			"--server-path", tmpDir,
			"--provider", "3os",
			"--mode", "update-version",
			"--identifier", "test@example.com",
			"--password", "testpass",
			"--package-path", "/tmp/game.zip",
		})
		root.Execute()
		return
	}

	output := runSubprocessTest(t, envKey)
	if !strings.Contains(output, "CONFIG_ERROR") {
		t.Errorf("expected CONFIG_ERROR for missing game-uri, got: %s", output)
	}
	if !strings.Contains(output, "game-uri") {
		t.Errorf("expected error to mention game-uri, got: %s", output)
	}
}

func TestDeployCmd_ProductionInvalidMode(t *testing.T) {
	const envKey = "GD_TEST_INVALID_MODE"
	if os.Getenv(envKey) == "1" {
		tmpDir := t.TempDir()
		os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module example.com/test\ngo 1.24\n"), 0644)
		os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\nfunc main() {}\n"), 0644)

		root := NewRootCmd()
		root.SetArgs([]string{
			"deploy",
			"--server-path", tmpDir,
			"--provider", "3os",
			"--mode", "invalid",
			"--identifier", "test@example.com",
			"--password", "testpass",
		})
		root.Execute()
		return
	}

	output := runSubprocessTest(t, envKey)
	if !strings.Contains(output, "CONFIG_ERROR") {
		t.Errorf("expected CONFIG_ERROR for invalid mode, got: %s", output)
	}
	if !strings.Contains(output, "unsupported mode") {
		t.Errorf("expected error to mention unsupported mode, got: %s", output)
	}
}

func TestBuildDeployConfig_ScreenConfig(t *testing.T) {
	opts := DeployOptions{
		ServerPath:  ".",
		AppName:     "test",
		ScreenType:  1,
		HalfSupport: 2,
		HalfRatio:   "0.75",
	}
	cfg := buildDeployConfig(opts)
	if cfg.ScreenConfig == nil {
		t.Fatal("expected ScreenConfig to be set")
	}
	if cfg.ScreenConfig.ScreenType != 1 {
		t.Errorf("expected ScreenType=1, got %d", cfg.ScreenConfig.ScreenType)
	}
}

func TestBuildDeployConfig_BuildConfig(t *testing.T) {
	opts := DeployOptions{
		ServerPath:  ".",
		AppName:     "test",
		BackendDir:  "admin",
		BackendCmd:  "./server",
		SocketDir:   "logic",
		SocketCmd:   "./server -type logic",
	}
	cfg := buildDeployConfig(opts)
	if cfg.BuildConfig == nil {
		t.Fatal("expected BuildConfig to be set")
	}
	if cfg.BuildConfig.Backend.WorkDir != "admin" {
		t.Errorf("expected Backend.WorkDir=admin, got %s", cfg.BuildConfig.Backend.WorkDir)
	}
	if cfg.BuildConfig.Frontend.WorkDir != "" {
		t.Errorf("expected Frontend.WorkDir empty, got %s", cfg.BuildConfig.Frontend.WorkDir)
	}
}

func TestBuildDeployConfig_NoScreenOrBuildWhenEmpty(t *testing.T) {
	opts := DeployOptions{
		ServerPath: ".",
		AppName:    "test",
	}
	cfg := buildDeployConfig(opts)
	if cfg.ScreenConfig != nil {
		t.Error("expected ScreenConfig to be nil when no screen flags set")
	}
	if cfg.BuildConfig != nil {
		t.Error("expected BuildConfig to be nil when no build flags set")
	}
}
