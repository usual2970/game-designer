package reporting

import "encoding/json"

type ResultCode string

const (
	CodeSuccess          ResultCode = "SUCCESS"
	CodePreflightFailed  ResultCode = "PREFLIGHT_FAILED"
	CodeDeployFailed     ResultCode = "DEPLOY_FAILED"
	CodeHealthCheckFail  ResultCode = "HEALTH_CHECK_FAILED"
	CodeInternalError    ResultCode = "INTERNAL_ERROR"
	CodeConfigError      ResultCode = "CONFIG_ERROR"
	CodeAuthFailed       ResultCode = "AUTH_FAILED"
	CodeListFailed       ResultCode = "LIST_FAILED"
	CodeLookupFailed     ResultCode = "LOOKUP_FAILED"
	CodeUploadFailed     ResultCode = "UPLOAD_FAILED"
	CodePublishFailed    ResultCode = "PUBLISH_FAILED"
	CodeReviewFailed     ResultCode = "REVIEW_FAILED"
	CodePartialSuccess   ResultCode = "PARTIAL_SUCCESS"
)

type Result struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Code    ResultCode  `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

func SuccessResult(message string, details interface{}) Result {
	return Result{Success: true, Message: message, Code: CodeSuccess, Details: details}
}

func FailResult(code ResultCode, message string, details interface{}) Result {
	return Result{Success: false, Message: message, Code: code, Details: details}
}

func (r Result) ToJSON() string {
	b, _ := json.Marshal(r)
	return string(b)
}
