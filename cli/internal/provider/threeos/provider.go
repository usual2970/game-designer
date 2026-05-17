package threeos

import (
	"context"
	"fmt"

	"github.com/example/game-designer-cli/internal/provider"
)

type ThreeOSProvider struct {
	client   *Client
	uploader *OSSUploader
}

func NewProvider(client *Client, uploader *OSSUploader) *ThreeOSProvider {
	return &ThreeOSProvider{
		client:   client,
		uploader: uploader,
	}
}

func (p *ThreeOSProvider) Name() string { return "3os" }

func (p *ThreeOSProvider) Deploy(ctx context.Context, config provider.DeployConfig) (*provider.DeployResult, error) {
	switch config.Mode {
	case provider.PublishModeList:
		return p.deployList(ctx, config)
	case provider.PublishModeCreate:
		return p.deployCreate(ctx, config)
	case provider.PublishModeUpdateInfo:
		return p.deployUpdateInfo(ctx, config)
	case provider.PublishModeUpdateVersion:
		return p.deployUpdateVersion(ctx, config)
	case provider.PublishModeApplyReview:
		return p.deployApplyReview(ctx, config)
	default:
		return nil, fmt.Errorf("unsupported mode: %s", config.Mode)
	}
}

func (p *ThreeOSProvider) Status(ctx context.Context, config provider.DeployConfig) (*provider.StatusResult, error) {
	return &provider.StatusResult{
		Status:  "unknown",
		Healthy: true,
		URL:     config.BaseURL,
	}, nil
}

func (p *ThreeOSProvider) HealthCheck(ctx context.Context, url string) (*provider.HealthResult, error) {
	return &provider.HealthResult{
		Healthy: true,
		Latency: "0ms",
		URL:     url,
	}, nil
}

func (p *ThreeOSProvider) login(ctx context.Context, config provider.DeployConfig) error {
	_, err := p.client.Login(ctx, config.Identifier, config.Password)
	if err != nil {
		return fmt.Errorf("auth failed: %w", err)
	}
	return nil
}

func (p *ThreeOSProvider) deployList(ctx context.Context, config provider.DeployConfig) (*provider.DeployResult, error) {
	if err := p.login(ctx, config); err != nil {
		return nil, err
	}

	page := config.Page
	if page <= 0 {
		page = 1
	}
	pageSize := config.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	listResp, err := p.client.ListGames(ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("list games failed: %w", err)
	}

	games := make([]provider.GameListEntry, 0, len(listResp.Data))
	for _, g := range listResp.Data {
		games = append(games, provider.GameListEntry{
			URI:         g.URI,
			Name:        g.Name,
			Description: g.Description,
			Logo:        g.Logo,
		})
	}

	return &provider.DeployResult{
		Provider: "3os",
		Mode:     provider.PublishModeList,
		GameList: &provider.GameListResult{
			Page:       listResp.Page,
			PageSize:   listResp.PageSize,
			TotalPages: listResp.TotalPages,
			TotalCount: listResp.TotalCount,
			Games:      games,
		},
	}, nil
}

func (p *ThreeOSProvider) deployCreate(ctx context.Context, config provider.DeployConfig) (*provider.DeployResult, error) {
	if err := p.login(ctx, config); err != nil {
		return nil, err
	}

	assets, err := p.uploadAssets(ctx, config)
	if err != nil {
		return nil, err
	}

	packageURL := assetURLByLabel(assets, "package")
	sqlURL := assetURLByLabel(assets, "sql")

	screenConfig := ScreenConfig{}
	if config.ScreenConfig != nil {
		screenConfig = ScreenConfig{
			ScreenType:  config.ScreenConfig.ScreenType,
			HalfSupport: config.ScreenConfig.HalfSupport,
			HalfRatio:   config.ScreenConfig.HalfRatio,
		}
	}

	buildConfig := map[string]BuildConfigEntry{}
	if config.BuildConfig != nil {
		buildConfig = map[string]BuildConfigEntry{
			"backend":  {WorkDir: config.BuildConfig.Backend.WorkDir, Cmd: config.BuildConfig.Backend.Cmd},
			"frontend": {WorkDir: config.BuildConfig.Frontend.WorkDir, Cmd: config.BuildConfig.Frontend.Cmd},
			"socket":   {WorkDir: config.BuildConfig.Socket.WorkDir, Cmd: config.BuildConfig.Socket.Cmd},
		}
	}

	gameResp, err := p.client.CreateWithVersion(ctx, &GameCreateWithVersionReq{
		Name:        config.GameName,
		Description: config.GameDescription,
		Logo:        config.GameLogo,
		Version: GameVersionCreateReq{
			Version:      config.Version,
			ChangeLog:    config.ChangeLog,
			FileUrl:      packageURL,
			InitSqlUrl:   sqlURL,
			BuildConfig:  buildConfig,
			ScreenConfig: screenConfig,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("create game failed: %w", err)
	}

	result := &provider.DeployResult{
		Provider: "3os",
		Mode:     provider.PublishModeCreate,
		GameURI:  gameResp.URI,
		URL:      gameResp.AccessUrl,
		Version:  config.Version,
		Assets:   assets,
	}

	// Optional review application
	if config.ReviewURI != "" {
		if err := p.client.ApplyReview(ctx, config.ReviewURI); err != nil {
			result.ReviewApplied = false
			result.ReviewURI = config.ReviewURI
			return result, fmt.Errorf("publish succeeded but review application failed: %w", err)
		}
		result.ReviewApplied = true
		result.ReviewURI = config.ReviewURI
	}

	return result, nil
}

func (p *ThreeOSProvider) deployUpdateInfo(ctx context.Context, config provider.DeployConfig) (*provider.DeployResult, error) {
	if err := p.login(ctx, config); err != nil {
		return nil, err
	}

	gameResp, err := p.client.UpdateGame(ctx, config.GameURI, &GameUpdateReq{
		URI:         config.GameURI,
		Description: config.GameDescription,
		Logo:        config.GameLogo,
	})
	if err != nil {
		return nil, fmt.Errorf("update game info failed: %w", err)
	}

	return &provider.DeployResult{
		Provider: "3os",
		Mode:     provider.PublishModeUpdateInfo,
		GameURI:  gameResp.URI,
		URL:      gameResp.AccessUrl,
	}, nil
}

func (p *ThreeOSProvider) deployUpdateVersion(ctx context.Context, config provider.DeployConfig) (*provider.DeployResult, error) {
	if err := p.login(ctx, config); err != nil {
		return nil, err
	}

	gameURI := config.GameURI
	if gameURI == "" {
		return nil, fmt.Errorf("game URI required for update-version")
	}

	assets, err := p.uploadAssets(ctx, config)
	if err != nil {
		return nil, err
	}

	packageURL := assetURLByLabel(assets, "package")
	sqlURL := assetURLByLabel(assets, "sql")

	screenConfig := ScreenConfig{}
	if config.ScreenConfig != nil {
		screenConfig = ScreenConfig{
			ScreenType:  config.ScreenConfig.ScreenType,
			HalfSupport: config.ScreenConfig.HalfSupport,
			HalfRatio:   config.ScreenConfig.HalfRatio,
		}
	}

	buildConfig := map[string]BuildConfigEntry{}
	if config.BuildConfig != nil {
		buildConfig = map[string]BuildConfigEntry{
			"backend":  {WorkDir: config.BuildConfig.Backend.WorkDir, Cmd: config.BuildConfig.Backend.Cmd},
			"frontend": {WorkDir: config.BuildConfig.Frontend.WorkDir, Cmd: config.BuildConfig.Frontend.Cmd},
			"socket":   {WorkDir: config.BuildConfig.Socket.WorkDir, Cmd: config.BuildConfig.Socket.Cmd},
		}
	}

	gameResp, err := p.client.UpdateWithVersion(ctx, &GameUpdateWithVersionReq{
		URI:         gameURI,
		Name:        config.GameName,
		Description: config.GameDescription,
		Logo:        config.GameLogo,
		Version: GameVersionCreateReq{
			Version:      config.Version,
			ChangeLog:    config.ChangeLog,
			FileUrl:      packageURL,
			InitSqlUrl:   sqlURL,
			BuildConfig:  buildConfig,
			ScreenConfig: screenConfig,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("update version failed: %w", err)
	}

	result := &provider.DeployResult{
		Provider: "3os",
		Mode:     provider.PublishModeUpdateVersion,
		GameURI:  gameResp.URI,
		URL:      gameResp.AccessUrl,
		Version:  config.Version,
		Assets:   assets,
	}

	if config.ReviewURI != "" {
		if err := p.client.ApplyReview(ctx, config.ReviewURI); err != nil {
			result.ReviewApplied = false
			result.ReviewURI = config.ReviewURI
			return result, fmt.Errorf("publish succeeded but review application failed: %w", err)
		}
		result.ReviewApplied = true
		result.ReviewURI = config.ReviewURI
	}

	return result, nil
}

func (p *ThreeOSProvider) deployApplyReview(ctx context.Context, config provider.DeployConfig) (*provider.DeployResult, error) {
	if err := p.login(ctx, config); err != nil {
		return nil, err
	}

	if config.ReviewURI == "" {
		return nil, fmt.Errorf("review URI required for apply-review")
	}

	if err := p.client.ApplyReview(ctx, config.ReviewURI); err != nil {
		return &provider.DeployResult{
			Provider:     "3os",
			Mode:         provider.PublishModeApplyReview,
			ReviewURI:    config.ReviewURI,
			ReviewApplied: false,
		}, fmt.Errorf("apply review failed: %w", err)
	}

	return &provider.DeployResult{
		Provider:     "3os",
		Mode:         provider.PublishModeApplyReview,
		ReviewURI:    config.ReviewURI,
		ReviewApplied: true,
	}, nil
}

func (p *ThreeOSProvider) uploadAssets(ctx context.Context, config provider.DeployConfig) ([]provider.UploadedAsset, error) {
	policy, err := p.client.GetPolicyToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get upload policy failed: %w", err)
	}

	var assets []provider.UploadedAsset

	pkgResult, err := p.uploader.UploadFile(ctx, config.PackagePath, policy, "package")
	if err != nil {
		return nil, err
	}
	if pkgResult != nil {
		assets = append(assets, provider.UploadedAsset{
			Label: pkgResult.Label,
			Path:  pkgResult.LocalPath,
			URL:   pkgResult.ObjectURL,
		})
	}

	sqlResult, err := p.uploader.UploadFile(ctx, config.SQLPath, policy, "sql")
	if err != nil {
		return nil, err
	}
	if sqlResult != nil {
		assets = append(assets, provider.UploadedAsset{
			Label: sqlResult.Label,
			Path:  sqlResult.LocalPath,
			URL:   sqlResult.ObjectURL,
		})
	}

	return assets, nil
}

func assetURLByLabel(assets []provider.UploadedAsset, label string) string {
	for _, a := range assets {
		if a.Label == label {
			return a.URL
		}
	}
	return ""
}
