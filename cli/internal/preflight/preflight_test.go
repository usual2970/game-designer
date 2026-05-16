package preflight

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckServerPath_Exists(t *testing.T) {
	tmpDir := t.TempDir()
	result := checkServerPath(tmpDir)
	if !result.Passed {
		t.Errorf("expected pass for existing path, got: %s", result.Message)
	}
}

func TestCheckServerPath_NotExists(t *testing.T) {
	result := checkServerPath("/nonexistent/path/that/does/not/exist")
	if result.Passed {
		t.Error("expected fail for nonexistent path")
	}
}

func TestCheckServerPath_Empty(t *testing.T) {
	result := checkServerPath("")
	if result.Passed {
		t.Error("expected fail for empty path")
	}
}

func TestCheckGoBuild_ValidServer(t *testing.T) {
	tmpDir := t.TempDir()
	modContent := `module example.com/test
go 1.24
`
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(modContent), 0644)
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\nfunc main() {}\n"), 0644)

	result := checkGoBuild(tmpDir)
	if !result.Passed {
		t.Errorf("expected pass for valid server, got: %s", result.Message)
	}
}

func TestRunChecks_AllPassed(t *testing.T) {
	tmpDir := t.TempDir()
	modContent := `module example.com/test
go 1.24
`
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(modContent), 0644)
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\nfunc main() {}\n"), 0644)

	results := RunChecks(tmpDir)
	if !AllPassed(results) {
		t.Error("expected all checks to pass")
	}
}

func TestAllPassed_WithFailure(t *testing.T) {
	results := []CheckResult{
		{Name: "a", Passed: true},
		{Name: "b", Passed: false},
	}
	if AllPassed(results) {
		t.Error("expected false when one check fails")
	}
}
