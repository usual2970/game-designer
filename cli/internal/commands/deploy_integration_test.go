package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/example/game-designer-cli/internal/reporting"
)

func TestVersionCmd(t *testing.T) {
	root := NewRootCmd()
	root.SetArgs([]string{"version"})
	// version uses fmt.Printf which writes to stdout directly
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPreflightCmd_ValidPath(t *testing.T) {
	// Create a valid temporary server directory
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

	// Capture stdout
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

	// Read captured output
	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	if !strings.Contains(output, "SUCCESS") {
		t.Errorf("expected SUCCESS in output, got: %s", output)
	}

	// Find JSON line
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
