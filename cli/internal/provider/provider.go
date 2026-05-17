package provider

import "context"

type PublishMode string

const (
	PublishModeCreate        PublishMode = "create"
	PublishModeUpdateInfo    PublishMode = "update-info"
	PublishModeUpdateVersion PublishMode = "update-version"
	PublishModeList          PublishMode = "list"
	PublishModeApplyReview   PublishMode = "apply-review"
)

type ScreenConfig struct {
	ScreenType  int    `json:"screenType"`
	HalfSupport int    `json:"halfSupport"`
	HalfRatio   string `json:"halfRatio"`
}

type BuildConfigEntry struct {
	WorkDir string `json:"workDir"`
	Cmd     string `json:"cmd"`
}

type BuildConfig struct {
	Backend  BuildConfigEntry `json:"backend"`
	Frontend BuildConfigEntry `json:"frontend"`
	Socket   BuildConfigEntry `json:"socket"`
}

type ResourceInfoItem struct {
	URI         string `json:"uri,omitempty"`
	Img         string `json:"img"`
	Description string `json:"description"`
}

type DeployConfig struct {
	ServerPath string `json:"serverPath"`
	Env        string `json:"env"`
	Region     string `json:"region"`
	AppName    string `json:"appName"`

	// Production auth
	BaseURL    string `json:"baseUrl,omitempty"`
	Identifier string `json:"-"`
	Password   string `json:"-"`

	// Publish mode
	Mode PublishMode `json:"mode,omitempty"`

	// Game metadata (create / update-info)
	GameName        string `json:"gameName,omitempty"`
	GameDescription string `json:"gameDescription,omitempty"`
	GameLogo        string `json:"gameLogo,omitempty"`
	GameURI         string `json:"gameUri,omitempty"`

	// Version (create / update-version)
	Version   string `json:"version,omitempty"`
	ChangeLog string `json:"changeLog,omitempty"`

	// Local asset paths
	PackagePath string `json:"packagePath,omitempty"`
	SQLPath     string `json:"sqlPath,omitempty"`

	// Screen and build config
	ScreenConfig *ScreenConfig      `json:"screenConfig,omitempty"`
	BuildConfig  *BuildConfig       `json:"buildConfig,omitempty"`

	// Resource items (update-info, update-version)
	ResourceItems []ResourceInfoItem `json:"resourceItems,omitempty"`

	// List mode
	Page     int `json:"page,omitempty"`
	PageSize int `json:"pageSize,omitempty"`

	// Review
	ReviewURI string `json:"reviewUri,omitempty"`
}

type UploadedAsset struct {
	Label string `json:"label"`
	Path  string `json:"path"`
	URL   string `json:"url"`
}

type GameListEntry struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
}

type GameListResult struct {
	Page       int             `json:"page"`
	PageSize   int             `json:"pageSize"`
	TotalPages int             `json:"totalPages"`
	TotalCount int64           `json:"totalCount"`
	Games      []GameListEntry `json:"games"`
}

type DeployResult struct {
	URL      string `json:"url"`
	Version  string `json:"version"`
	AppName  string `json:"appName"`
	Provider string `json:"provider"`

	// Production-specific
	Mode         PublishMode     `json:"mode,omitempty"`
	GameURI      string          `json:"gameUri,omitempty"`
	ReviewURI    string          `json:"reviewUri,omitempty"`
	ReviewApplied bool           `json:"reviewApplied,omitempty"`
	Assets       []UploadedAsset `json:"assets,omitempty"`
	GameList     *GameListResult `json:"gameList,omitempty"`
}

type StatusResult struct {
	Status   string `json:"status"`
	Healthy  bool   `json:"healthy"`
	URL      string `json:"url"`
	Version  string `json:"version"`
}

type HealthResult struct {
	Healthy bool   `json:"healthy"`
	Latency string `json:"latency"`
	URL     string `json:"url"`
}

type Provider interface {
	Name() string
	Deploy(ctx context.Context, config DeployConfig) (*DeployResult, error)
	Status(ctx context.Context, config DeployConfig) (*StatusResult, error)
	HealthCheck(ctx context.Context, url string) (*HealthResult, error)
}
