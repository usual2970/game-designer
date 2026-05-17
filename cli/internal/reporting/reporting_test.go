package reporting

import (
	"encoding/json"
	"strings"
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

func TestAllResultCodes(t *testing.T) {
	codes := []ResultCode{
		CodeSuccess, CodePreflightFailed, CodeDeployFailed,
		CodeHealthCheckFail, CodeInternalError, CodeConfigError,
		CodeAuthFailed, CodeListFailed, CodeLookupFailed,
		CodeUploadFailed, CodePublishFailed, CodeReviewFailed,
		CodePartialSuccess,
	}
	for _, code := range codes {
		r := FailResult(code, "test", nil)
		if r.Code != code {
			t.Errorf("expected code=%s, got %s", code, r.Code)
		}
	}
}

func TestResultJSON_NoSensitiveData(t *testing.T) {
	r := SuccessResult("deployed", map[string]string{
		"gameUri": "game123",
		"url":     "https://game.example.com",
	})
	jsonStr := r.ToJSON()
	if strings.Contains(jsonStr, "password") {
		t.Error("JSON should not contain 'password'")
	}
	if strings.Contains(jsonStr, "token") {
		t.Error("JSON should not contain 'token'")
	}
}

func TestFailResultJSON_Parsable(t *testing.T) {
	r := FailResult(CodeAuthFailed, "auth failed: invalid credentials", map[string]string{
		"endpoint": "/common/v1/auth/login",
	})
	var parsed Result
	if err := json.Unmarshal([]byte(r.ToJSON()), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if parsed.Success {
		t.Error("expected success=false")
	}
	if parsed.Code != "AUTH_FAILED" {
		t.Errorf("expected AUTH_FAILED, got %s", parsed.Code)
	}
}
