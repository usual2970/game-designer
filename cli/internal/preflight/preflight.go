package preflight

import (
	"fmt"
	"os"
	"os/exec"
)

type CheckResult struct {
	Name    string `json:"name"`
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

func RunChecks(serverPath string) []CheckResult {
	var results []CheckResult

	results = append(results, checkServerPath(serverPath))
	results = append(results, checkGoBuild(serverPath))

	return results
}

func checkServerPath(path string) CheckResult {
	if path == "" {
		return CheckResult{
			Name:    "server-path",
			Passed:  false,
			Message: "server path is empty",
		}
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return CheckResult{
			Name:    "server-path",
			Passed:  false,
			Message: fmt.Sprintf("server path does not exist: %s", path),
		}
	}
	return CheckResult{
		Name:    "server-path",
		Passed:  true,
		Message: "server path exists",
	}
}

func checkGoBuild(serverPath string) CheckResult {
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = serverPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return CheckResult{
			Name:    "go-build",
			Passed:  false,
			Message: fmt.Sprintf("build failed: %s", string(output)),
		}
	}
	return CheckResult{
		Name:    "go-build",
		Passed:  true,
		Message: "server builds successfully",
	}
}

func AllPassed(results []CheckResult) bool {
	for _, r := range results {
		if !r.Passed {
			return false
		}
	}
	return true
}
