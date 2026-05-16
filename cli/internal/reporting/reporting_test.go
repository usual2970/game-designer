package reporting

import (
	"encoding/json"
	"testing"
)

func TestSuccessResult(t *testing.T) {
	r := SuccessResult("deployed", map[string]string{"url": "https://example.com"})
	if !r.Success {
		t.Error("expected success=true")
	}
	if r.Code != CodeSuccess {
		t.Errorf("expected code=SUCCESS, got %s", r.Code)
	}
	if r.Message != "deployed" {
		t.Errorf("expected message=deployed, got %s", r.Message)
	}
}

func TestFailResult(t *testing.T) {
	r := FailResult(CodePreflightFailed, "checks failed", nil)
	if r.Success {
		t.Error("expected success=false")
	}
	if r.Code != CodePreflightFailed {
		t.Errorf("expected code=PREFLIGHT_FAILED, got %s", r.Code)
	}
}

func TestResult_ToJSON(t *testing.T) {
	r := SuccessResult("ok", "details")
	jsonStr := r.ToJSON()

	var parsed Result
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !parsed.Success {
		t.Error("expected success=true in JSON output")
	}
}
