package code

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/code"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

func convertDataSource(data *model.CodeDataSource) *code.CodeDataSource {
	if data == nil {
		return &code.CodeDataSource{}
	}
	return &code.CodeDataSource{
		CodeDataSourceId: data.CodeDataSourceID,
		Name:             data.Name,
		Description:      data.Description,
		MaxScore:         data.MaxScore,
		CreatedAt:        data.CreatedAt.Unix(),
		UpdatedAt:        data.UpdatedAt.Unix(),
	}
}

func (c *CodeService) ListDataSource(ctx context.Context, req *code.ListDataSourceRequest) (*code.ListDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := c.repository.ListCodeDataSource(ctx, req.CodeDataSourceId, req.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &code.ListDataSourceResponse{}, nil
		}
		return nil, err
	}
	data := code.ListDataSourceResponse{}
	for _, d := range *list {
		data.CodeDataSource = append(data.CodeDataSource, convertDataSource(&d))
	}
	return &data, nil
}

const maskData = "xxxxxxxxxx"

func convertGitHubSetting(gitHubSetting *model.CodeGitHubSetting, gitleaksSetting *model.CodeGitleaksSetting, dependencySetting *model.CodeDependencySetting, maskKey bool) *code.GitHubSetting {
	var convertedGithubSetting code.GitHubSetting
	if gitHubSetting == nil {
		return &convertedGithubSetting
	}
	convertedGithubSetting = code.GitHubSetting{
		GithubSettingId:     gitHubSetting.CodeGitHubSettingID,
		Name:                gitHubSetting.Name,
		ProjectId:           gitHubSetting.ProjectID,
		Type:                getType(gitHubSetting.Type),
		BaseUrl:             gitHubSetting.BaseURL,
		TargetResource:      gitHubSetting.TargetResource,
		GithubUser:          gitHubSetting.GitHubUser,
		PersonalAccessToken: gitHubSetting.PersonalAccessToken,
		CreatedAt:           gitHubSetting.CreatedAt.Unix(),
		UpdatedAt:           gitHubSetting.UpdatedAt.Unix(),
	}
	if convertedGithubSetting.PersonalAccessToken != "" && maskKey {
		convertedGithubSetting.PersonalAccessToken = maskData // Masking sensitive data.
	}
	if gitleaksSetting != nil {
		convertedGithubSetting.GitleaksSetting = convertGitleaksSetting(gitleaksSetting)
	}
	if dependencySetting != nil {
		convertedGithubSetting.DependencySetting = convertDependencySetting(dependencySetting)
	}
	return &convertedGithubSetting
}
func convertGitleaksSetting(data *model.CodeGitleaksSetting) *code.GitleaksSetting {
	var gitleaksSetting code.GitleaksSetting
	if data == nil {
		return &gitleaksSetting
	}
	gitleaksSetting = code.GitleaksSetting{
		GithubSettingId:   data.CodeGitHubSettingID,
		CodeDataSourceId:  data.CodeDataSourceID,
		ProjectId:         data.ProjectID,
		RepositoryPattern: data.RepositoryPattern,
		ScanPublic:        data.ScanPublic,
		ScanInternal:      data.ScanInternal,
		ScanPrivate:       data.ScanPrivate,
		Status:            getStatus(data.Status),
		StatusDetail:      data.StatusDetail,
		CreatedAt:         data.CreatedAt.Unix(),
		UpdatedAt:         data.UpdatedAt.Unix(),
	}
	if !zero.IsZeroVal(data.ScanAt) {
		gitleaksSetting.ScanAt = data.ScanAt.Unix()
	}
	return &gitleaksSetting
}

func convertDependencySetting(data *model.CodeDependencySetting) *code.DependencySetting {
	var dependencySetting code.DependencySetting
	if data == nil {
		return &dependencySetting
	}
	dependencySetting = code.DependencySetting{
		GithubSettingId:  data.CodeGitHubSettingID,
		CodeDataSourceId: data.CodeDataSourceID,
		ProjectId:        data.ProjectID,
		Status:           getStatus(data.Status),
		StatusDetail:     data.StatusDetail,
		CreatedAt:        data.CreatedAt.Unix(),
		UpdatedAt:        data.UpdatedAt.Unix(),
	}
	if !zero.IsZeroVal(data.ScanAt) {
		dependencySetting.ScanAt = data.ScanAt.Unix()
	}
	return &dependencySetting
}

func (c *CodeService) ListGitHubSetting(ctx context.Context, req *code.ListGitHubSettingRequest) (*code.ListGitHubSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	gitHubSettings, err := c.repository.ListGitHubSetting(ctx, req.ProjectId, req.GithubSettingId)
	if err != nil {
		return nil, err
	}
	data := code.ListGitHubSettingResponse{}
	if len(*gitHubSettings) == 0 {
		return &data, nil
	}
	gitleaksSettings, err := c.repository.ListGitleaksSetting(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	dependencySettings, err := c.repository.ListDependencySetting(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	mapGitleaksSetting := map[uint32]model.CodeGitleaksSetting{}
	for _, gitleaksSetting := range *gitleaksSettings {
		mapGitleaksSetting[gitleaksSetting.CodeGitHubSettingID] = gitleaksSetting
	}
	mapDependencySetting := map[uint32]model.CodeDependencySetting{}
	for _, dependencySetting := range *dependencySettings {
		mapDependencySetting[dependencySetting.CodeGitHubSettingID] = dependencySetting
	}
	for _, gitHubSetting := range *gitHubSettings {
		var gitleaks *model.CodeGitleaksSetting
		var dependency *model.CodeDependencySetting
		valGitleaks, ok := mapGitleaksSetting[gitHubSetting.CodeGitHubSettingID]
		if ok {
			gitleaks = &valGitleaks
		}
		valDependency, ok := mapDependencySetting[gitHubSetting.CodeGitHubSettingID]
		if ok {
			dependency = &valDependency
		}
		data.GithubSetting = append(data.GithubSetting, convertGitHubSetting(&gitHubSetting, gitleaks, dependency, true))
	}
	return &data, nil
}

func (c *CodeService) GetGitHubSetting(ctx context.Context, req *code.GetGitHubSettingRequest) (*code.GetGitHubSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	githubSetting, err := c.repository.GetGitHubSetting(ctx, req.ProjectId, req.GithubSettingId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &code.GetGitHubSettingResponse{}, nil
		}
		return nil, err
	}
	gitleaksSetting, err := c.repository.GetGitleaksSetting(ctx, githubSetting.ProjectID, githubSetting.CodeGitHubSettingID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	dependencySetting, err := c.repository.GetDependencySetting(ctx, githubSetting.ProjectID, githubSetting.CodeGitHubSettingID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &code.GetGitHubSettingResponse{GithubSetting: convertGitHubSetting(githubSetting, gitleaksSetting, dependencySetting, false)}, nil
}

func (c *CodeService) PutGitHubSetting(ctx context.Context, req *code.PutGitHubSettingRequest) (*code.PutGitHubSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if req.GithubSetting.PersonalAccessToken != "" && req.GithubSetting.PersonalAccessToken != maskData {
		encrypted, err := encryptWithBase64(&c.cipherBlock, req.GithubSetting.PersonalAccessToken)
		if err != nil {
			c.logger.Errorf(ctx, "Failed to encrypt PAT: err=%+v", err)
			return nil, err
		}
		req.GithubSetting.PersonalAccessToken = encrypted
	} else {
		req.GithubSetting.PersonalAccessToken = "" // for not update token.
	}
	registeredGitHubSetting, err := c.repository.UpsertGitHubSetting(ctx, req.GithubSetting)
	if err != nil {
		return nil, err
	}
	return &code.PutGitHubSettingResponse{GithubSetting: convertGitHubSetting(registeredGitHubSetting, nil, nil, true)}, nil
}

func (c *CodeService) DeleteGitHubSetting(ctx context.Context, req *code.DeleteGitHubSettingRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := c.repository.DeleteGitleaksCache(ctx, req.GithubSettingId); err != nil {
		return nil, err
	}
	if err := c.repository.DeleteGitleaksSetting(ctx, req.ProjectId, req.GithubSettingId); err != nil {
		return nil, err
	}
	if err := c.repository.DeleteDependencySetting(ctx, req.ProjectId, req.GithubSettingId); err != nil {
		return nil, err
	}
	organizations, err := c.repository.ListGitHubEnterpriseOrg(ctx, req.ProjectId, req.GithubSettingId)
	if err != nil {
		return nil, err
	}
	for _, org := range *organizations {
		err = c.repository.DeleteGitHubEnterpriseOrg(ctx, org.ProjectID, org.CodeGitHubSettingID, org.Organization)
		if err != nil {
			return nil, err
		}
	}
	if err := c.repository.DeleteGitHubSetting(ctx, req.ProjectId, req.GithubSettingId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (c *CodeService) PutGitleaksSetting(ctx context.Context, req *code.PutGitleaksSettingRequest) (*code.PutGitleaksSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	registerd, err := c.repository.UpsertGitleaksSetting(ctx, req.GitleaksSetting)
	if err != nil {
		return nil, err
	}
	return &code.PutGitleaksSettingResponse{GitleaksSetting: convertGitleaksSetting(registerd)}, nil
}

func (c *CodeService) DeleteGitleaksSetting(ctx context.Context, req *code.DeleteGitleaksSettingRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := c.repository.DeleteGitleaksCache(ctx, req.GithubSettingId); err != nil {
		return nil, err
	}
	if err := c.repository.DeleteGitleaksSetting(ctx, req.ProjectId, req.GithubSettingId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertGitleaksCache(data *model.CodeGitleaksCache) *code.GitleaksCache {
	var converted code.GitleaksCache
	if data == nil {
		return &converted
	}
	converted = code.GitleaksCache{
		GithubSettingId:    data.CodeGitHubSettingID,
		RepositoryFullName: data.RepositoryFullName,
		CreatedAt:          data.CreatedAt.Unix(),
		UpdatedAt:          data.UpdatedAt.Unix(),
	}
	if !data.ScanAt.IsZero() {
		converted.ScanAt = data.ScanAt.Unix()
	}
	return &converted
}

func (c *CodeService) GetGitleaksCache(ctx context.Context, req *code.GetGitleaksCacheRequest) (*code.GetGitleaksCacheResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := c.repository.GetGitleaksCache(ctx, req.ProjectId, req.GithubSettingId, req.RepositoryFullName, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &code.GetGitleaksCacheResponse{}, nil
		}
		return nil, err
	}
	return &code.GetGitleaksCacheResponse{GitleaksCache: convertGitleaksCache(data)}, nil
}

func (c *CodeService) PutGitleaksCache(ctx context.Context, req *code.PutGitleaksCacheRequest) (*code.PutGitleaksCacheResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	// gitleaks setting data in project must be exists
	if _, err := c.repository.GetGitleaksSetting(ctx, req.ProjectId, req.GitleaksCache.GithubSettingId); err != nil {
		return nil, err
	}
	data, err := c.repository.UpsertGitleaksCache(ctx, req.ProjectId, req.GitleaksCache)
	if err != nil {
		return nil, err
	}
	return &code.PutGitleaksCacheResponse{GitleaksCache: convertGitleaksCache(data)}, nil
}

func (c *CodeService) PutDependencySetting(ctx context.Context, req *code.PutDependencySettingRequest) (*code.PutDependencySettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	registerd, err := c.repository.UpsertDependencySetting(ctx, req.DependencySetting)
	if err != nil {
		return nil, err
	}
	return &code.PutDependencySettingResponse{DependencySetting: convertDependencySetting(registerd)}, nil
}

func (c *CodeService) DeleteDependencySetting(ctx context.Context, req *code.DeleteDependencySettingRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := c.repository.DeleteDependencySetting(ctx, req.ProjectId, req.GithubSettingId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func getType(s string) code.Type {
	typeKey := strings.ToUpper(s)
	if _, ok := code.Type_value[typeKey]; !ok {
		return code.Type_UNKNOWN_TYPE
	}
	switch typeKey {
	case code.Type_ENTERPRISE.String():
		return code.Type_ENTERPRISE
	case code.Type_ORGANIZATION.String():
		return code.Type_ORGANIZATION
	case code.Type_USER.String():
		return code.Type_USER
	default:
		return code.Type_UNKNOWN_TYPE
	}
}

func getStatus(s string) code.Status {
	statusKey := strings.ToUpper(s)
	if _, ok := code.Status_value[statusKey]; !ok {
		return code.Status_UNKNOWN
	}
	switch statusKey {
	case code.Status_OK.String():
		return code.Status_OK
	case code.Status_CONFIGURED.String():
		return code.Status_CONFIGURED
	case code.Status_IN_PROGRESS.String():
		return code.Status_IN_PROGRESS
	case code.Status_ERROR.String():
		return code.Status_ERROR
	default:
		return code.Status_UNKNOWN
	}
}

func convertGitHubEnterpriseOrg(data *model.CodeGitHubEnterpriseOrg) *code.GitHubEnterpriseOrg {
	if data == nil {
		return &code.GitHubEnterpriseOrg{}
	}
	return &code.GitHubEnterpriseOrg{
		GithubSettingId: data.CodeGitHubSettingID,
		Organization:    data.Organization,
		ProjectId:       data.ProjectID,
		CreatedAt:       data.CreatedAt.Unix(),
		UpdatedAt:       data.CreatedAt.Unix(),
	}
}

func (c *CodeService) ListGitHubEnterpriseOrg(ctx context.Context, req *code.ListGitHubEnterpriseOrgRequest) (*code.ListGitHubEnterpriseOrgResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := c.repository.ListGitHubEnterpriseOrg(ctx, req.ProjectId, req.GithubSettingId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &code.ListGitHubEnterpriseOrgResponse{}, nil
		}
		return nil, err
	}
	data := code.ListGitHubEnterpriseOrgResponse{}
	for _, d := range *list {
		data.GithubEnterpriseOrg = append(data.GithubEnterpriseOrg, convertGitHubEnterpriseOrg(&d))
	}
	return &data, nil
}

func (c *CodeService) PutGitHubEnterpriseOrg(ctx context.Context, req *code.PutGitHubEnterpriseOrgRequest) (*code.PutGitHubEnterpriseOrgResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	registered, err := c.repository.UpsertGitHubEnterpriseOrg(ctx, req.GithubEnterpriseOrg)
	if err != nil {
		return nil, err
	}
	return &code.PutGitHubEnterpriseOrgResponse{GithubEnterpriseOrg: convertGitHubEnterpriseOrg(registered)}, nil
}

func (c *CodeService) DeleteGitHubEnterpriseOrg(ctx context.Context, req *code.DeleteGitHubEnterpriseOrgRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := c.repository.DeleteGitHubEnterpriseOrg(ctx, req.ProjectId, req.GithubSettingId, req.Organization)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (c *CodeService) InvokeScanGitleaks(ctx context.Context, req *code.InvokeScanGitleaksRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := c.repository.GetGitleaksSetting(ctx, req.ProjectId, req.GithubSettingId)
	if err != nil {
		return nil, err
	}
	resp, err := c.sqs.Send(ctx, c.codeGitleaksQueueURL, &message.CodeQueueMessage{
		GitHubSettingID: data.CodeGitHubSettingID,
		ProjectID:       data.ProjectID,
		ScanOnly:        req.ScanOnly,
		FullScan:        req.FullScan,
	})
	if err != nil {
		return nil, err
	}
	if _, err = c.repository.UpsertGitleaksSetting(ctx, &code.GitleaksSettingForUpsert{
		GithubSettingId:   data.CodeGitHubSettingID,
		CodeDataSourceId:  data.CodeDataSourceID,
		ProjectId:         data.ProjectID,
		RepositoryPattern: data.RepositoryPattern,
		ScanPublic:        data.ScanPublic,
		ScanInternal:      data.ScanInternal,
		ScanPrivate:       data.ScanPrivate,
		Status:            code.Status_IN_PROGRESS,
		StatusDetail:      fmt.Sprintf("Start scan at %+v", time.Now().Format(time.RFC3339)),
		ScanAt:            data.ScanAt.Unix(),
	}); err != nil {
		return nil, err
	}
	c.logger.Infof(ctx, "Invoke scanned, messageId: %v", resp.MessageId)
	return &empty.Empty{}, nil
}

func (c *CodeService) InvokeScanDependency(ctx context.Context, req *code.InvokeScanDependencyRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := c.repository.GetDependencySetting(ctx, req.ProjectId, req.GithubSettingId)
	if err != nil {
		return nil, err
	}
	resp, err := c.sqs.Send(ctx, c.codeDependencyQueueURL, &message.CodeQueueMessage{
		GitHubSettingID: data.CodeGitHubSettingID,
		ProjectID:       data.ProjectID,
		ScanOnly:        req.ScanOnly,
	})
	if err != nil {
		return nil, err
	}
	if _, err = c.repository.UpsertDependencySetting(ctx, &code.DependencySettingForUpsert{
		GithubSettingId:  data.CodeGitHubSettingID,
		CodeDataSourceId: data.CodeDataSourceID,
		ProjectId:        data.ProjectID,
		Status:           code.Status_IN_PROGRESS,
		StatusDetail:     fmt.Sprintf("Start scan at %+v", time.Now().Format(time.RFC3339)),
		ScanAt:           data.ScanAt.Unix(),
	}); err != nil {
		return nil, err
	}
	c.logger.Infof(ctx, "Invoke scanned, messageId: %v", resp.MessageId)
	return &empty.Empty{}, nil
}

func (c *CodeService) InvokeScanAll(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	listGitleaks, err := c.repository.ListGitleaksSetting(ctx, 0)
	if err != nil {
		return nil, err
	}
	for _, g := range *listGitleaks {
		if zero.IsZeroVal(g.ProjectID) || zero.IsZeroVal(g.CodeDataSourceID) {
			continue
		}
		if resp, err := c.projectClient.IsActive(ctx, &project.IsActiveRequest{ProjectId: g.ProjectID}); err != nil {
			c.logger.Errorf(ctx, "Failed to project.IsActive API, err=%+v", err)
			return nil, err
		} else if !resp.Active {
			c.logger.Infof(ctx, "Skip deactive project, project_id=%d", g.ProjectID)
			continue
		}
		if _, err := c.InvokeScanGitleaks(ctx, &code.InvokeScanGitleaksRequest{
			GithubSettingId: g.CodeGitHubSettingID,
			ProjectId:       g.ProjectID,
			ScanOnly:        true,
		}); err != nil {
			c.logger.Errorf(ctx, "InvokeScanGitleaks error occured: code_github_setting_id=%d, err=%+v", g.CodeGitHubSettingID, err)
			return nil, err
		}
	}
	listDependency, err := c.repository.ListDependencySetting(ctx, 0)
	if err != nil {
		return nil, err
	}
	for _, g := range *listDependency {
		if zero.IsZeroVal(g.ProjectID) || zero.IsZeroVal(g.CodeDataSourceID) {
			continue
		}
		if resp, err := c.projectClient.IsActive(ctx, &project.IsActiveRequest{ProjectId: g.ProjectID}); err != nil {
			c.logger.Errorf(ctx, "Failed to project.IsActive API, err=%+v", err)
			return nil, err
		} else if !resp.Active {
			c.logger.Infof(ctx, "Skip deactive project, project_id=%d", g.ProjectID)
			continue
		}
		if _, err := c.InvokeScanDependency(ctx, &code.InvokeScanDependencyRequest{
			GithubSettingId: g.CodeGitHubSettingID,
			ProjectId:       g.ProjectID,
			ScanOnly:        true,
		}); err != nil {
			c.logger.Errorf(ctx, "InvokeScanDependency error occured: code_github_setting_id=%d, err=%+v", g.CodeGitHubSettingID, err)
			return nil, err
		}
	}
	return &empty.Empty{}, nil
}
