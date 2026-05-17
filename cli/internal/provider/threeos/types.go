package threeos

import (
	"encoding/json"
	"fmt"
)

// APIResponse mirrors the backend's pkg/resp/resp.go envelope.
type APIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func (r APIResponse) IsSuccess() bool {
	return r.Code == 0
}

// AuthLoginReq mirrors domain.AuthLoginReq.
type AuthLoginReq struct {
	Identifier string `json:"identifier"`
	Type       string `json:"type"`
	Data       string `json:"data"`
}

// AuthLoginResp mirrors domain.AuthLoginResp.
type AuthLoginResp struct {
	AccessToken string `json:"accessToken"`
	ExpiresAt   int64  `json:"expiresIn"`
}

// FilePolicyTokenResp mirrors domain.FilePolicyTokenResp.
type FilePolicyTokenResp struct {
	Policy           string `json:"policy"`
	SecurityToken    string `json:"security_token"`
	SignatureVersion string `json:"x_oss_signature_version"`
	Credential       string `json:"x_oss_credential"`
	Date             string `json:"x_oss_date"`
	Signature        string `json:"signature"`
	Host             string `json:"host"`
	Dir              string `json:"dir"`
}

// ScreenConfig mirrors domain.ScreenConfig.
type ScreenConfig struct {
	ScreenType  int    `json:"screenType"`
	HalfSupport int    `json:"halfSupport"`
	HalfRatio   string `json:"halfRatio"`
}

// BuildConfigEntry mirrors domain.GameVersionBuildConfig.
type BuildConfigEntry struct {
	WorkDir string `json:"workDir"`
	Cmd     string `json:"cmd"`
}

// GameVersionCreateReq mirrors domain.GameVersionCreateReq (subset used by CLI).
type GameVersionCreateReq struct {
	Version      string                      `json:"version"`
	ChangeLog    string                      `json:"changeLog"`
	FileUrl      string                      `json:"fileUrl"`
	InitSqlUrl   string                      `json:"initSqlUrl"`
	BuildConfig  map[string]BuildConfigEntry `json:"buildConfig"`
	ScreenConfig ScreenConfig                `json:"screenConfig"`
}

// GameCreateWithVersionReq mirrors domain.GameCreateWithVersionReq.
type GameCreateWithVersionReq struct {
	Name          string               `json:"name"`
	Description   string               `json:"description"`
	Logo          string               `json:"logo"`
	ResourceItems []ResourceInfoItem   `json:"resourceItems,omitempty"`
	Version       GameVersionCreateReq `json:"version"`
}

// ResourceInfoItem mirrors domain.ResourceInfoItem.
type ResourceInfoItem struct {
	URI         string `json:"uri,omitempty"`
	Img         string `json:"img"`
	Description string `json:"description"`
}

// GameUpdateReq mirrors domain.GameUpdateReq (subset for CLI).
type GameUpdateReq struct {
	URI           string             `json:"uri"`
	Description   string             `json:"description"`
	Logo          string             `json:"logo"`
	ResourceItems []ResourceInfoItem `json:"resourceItems,omitempty"`
}

// GameUpdateWithVersionReq mirrors domain.GameUpdateWithVersionReq.
type GameUpdateWithVersionReq struct {
	URI           string               `json:"uri"`
	Name          string               `json:"name"`
	Description   string               `json:"description"`
	Logo          string               `json:"logo"`
	ResourceItems []ResourceInfoItem   `json:"resourceItems,omitempty"`
	Version       GameVersionCreateReq `json:"version"`
}

// GameInfoResp mirrors domain.GameInfoResp (fields the CLI uses).
type GameInfoResp struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	AccessUrl   string `json:"accessUrl"`
	SocketUrl   string `json:"socketUrl"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
	ScreenType  int8   `json:"screenType"`
	Listed      int    `json:"listed"`
}

// GameListResp mirrors domain.constant.PaginatedResponse[GameInfoResp].
type GameListResp struct {
	Page       int            `json:"page"`
	PageSize   int            `json:"pageSize"`
	TotalPages int            `json:"totalPages"`
	TotalCount int64          `json:"totalCount"`
	Data       []GameInfoResp `json:"data"`
}

// GameReviewApplyReq mirrors domain.GameReviewApplyReq.
type GameReviewApplyReq struct {
	URI string `json:"uri"`
}

// ClientError is a typed error for API client failures.
type ClientError struct {
	Endpoint string
	Message  string
	Code     int
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("%s: %s (code=%d)", e.Endpoint, e.Message, e.Code)
}
