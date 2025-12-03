package code

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/datasource-api/pkg/github"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/code"
	"github.com/golang/protobuf/ptypes/empty"
	ghub "github.com/google/go-github/v44/github"
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

func convertGitHubSetting(
	gitHubSetting *model.CodeGitHubSetting,
	gitleaksSetting *model.CodeGitleaksSetting,
	dependencySetting *model.CodeDependencySetting,
	codeScanSetting *model.CodeCodeScanSetting,
	maskKey bool,
) *code.GitHubSetting {
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
	if codeScanSetting != nil {
		convertedGithubSetting.CodeScanSetting = convertCodeScanSetting(codeScanSetting)
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
		GithubSettingId:   data.CodeGitHubSettingID,
		CodeDataSourceId:  data.CodeDataSourceID,
		ProjectId:         data.ProjectID,
		RepositoryPattern: data.RepositoryPattern,
		Status:            getStatus(data.Status),
		StatusDetail:      data.StatusDetail,
		CreatedAt:         data.CreatedAt.Unix(),
		UpdatedAt:         data.UpdatedAt.Unix(),
	}
	if !zero.IsZeroVal(data.ScanAt) {
		dependencySetting.ScanAt = data.ScanAt.Unix()
	}
	return &dependencySetting
}

func convertCodeScanSetting(data *model.CodeCodeScanSetting) *code.CodeScanSetting {
	var converted code.CodeScanSetting
	if data == nil {
		return &converted
	}
	converted = code.CodeScanSetting{
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
		converted.ScanAt = data.ScanAt.Unix()
	}
	return &converted
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
	codeScanSettings, err := c.repository.ListCodeScanSetting(ctx, req.ProjectId)
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
	mapCodeScanSetting := map[uint32]model.CodeCodeScanSetting{}
	for _, codeScanSetting := range *codeScanSettings {
		mapCodeScanSetting[codeScanSetting.CodeGitHubSettingID] = codeScanSetting
	}

	for _, gitHubSetting := range *gitHubSettings {
		var gitleaks *model.CodeGitleaksSetting
		var dependency *model.CodeDependencySetting
		var codescan *model.CodeCodeScanSetting
		valGitleaks, ok := mapGitleaksSetting[gitHubSetting.CodeGitHubSettingID]
		if ok {
			gitleaks = &valGitleaks
		}
		valDependency, ok := mapDependencySetting[gitHubSetting.CodeGitHubSettingID]
		if ok {
			dependency = &valDependency
		}
		valCodeScan, ok := mapCodeScanSetting[gitHubSetting.CodeGitHubSettingID]
		if ok {
			codescan = &valCodeScan
		}
		data.GithubSetting = append(data.GithubSetting, convertGitHubSetting(&gitHubSetting, gitleaks, dependency, codescan, true))
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
	codeScanSetting, err := c.repository.GetCodeScanSetting(ctx, githubSetting.ProjectID, githubSetting.CodeGitHubSettingID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &code.GetGitHubSettingResponse{GithubSetting: convertGitHubSetting(githubSetting, gitleaksSetting, dependencySetting, codeScanSetting, false)}, nil
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
	return &code.PutGitHubSettingResponse{GithubSetting: convertGitHubSetting(registeredGitHubSetting, nil, nil, nil, true)}, nil
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
	if err := c.repository.DeleteCodeScanSetting(ctx, req.ProjectId, req.GithubSettingId); err != nil {
		return nil, err
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
	registered, err := c.repository.UpsertGitleaksSetting(ctx, req.GitleaksSetting)
	if err != nil {
		return nil, err
	}
	if !registered.ErrorNotifiedAt.IsZero() && registered.Status != code.Status_ERROR.String() {
		if err := c.repository.UpdateCodeGitleaksErrorNotifiedAt(ctx, gorm.Expr("NULL"), registered.CodeGitHubSettingID, registered.ProjectID); err != nil {
			return nil, err
		}
	}
	return &code.PutGitleaksSettingResponse{GitleaksSetting: convertGitleaksSetting(registered)}, nil
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

func (c *CodeService) ListGitleaksCache(ctx context.Context, req *code.ListGitleaksCacheRequest) (*code.ListGitleaksCacheResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := c.repository.ListGitleaksCache(ctx, req.ProjectId, req.GithubSettingId)
	if err != nil {
		return nil, err
	}
	gitleaksCache := code.ListGitleaksCacheResponse{}
	for _, d := range *data {
		gitleaksCache.GitleaksCache = append(gitleaksCache.GitleaksCache, convertGitleaksCache(&d))
	}

	return &gitleaksCache, nil
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
	registered, err := c.repository.UpsertDependencySetting(ctx, req.DependencySetting)
	if err != nil {
		return nil, err
	}
	if !registered.ErrorNotifiedAt.IsZero() && registered.Status != code.Status_ERROR.String() {
		if err := c.repository.UpdateCodeDependencyErrorNotifiedAt(ctx, gorm.Expr("NULL"), registered.CodeGitHubSettingID, registered.ProjectID); err != nil {
			return nil, err
		}
	}
	return &code.PutDependencySettingResponse{DependencySetting: convertDependencySetting(registered)}, nil
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

func (c *CodeService) PutCodeScanSetting(ctx context.Context, req *code.PutCodeScanSettingRequest) (*code.PutCodeScanSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	registered, err := c.repository.UpsertCodeScanSetting(ctx, req.CodeScanSetting)
	if err != nil {
		return nil, err
	}
	if !registered.ErrorNotifiedAt.IsZero() && registered.Status != code.Status_ERROR.String() {
		if err := c.repository.UpdateCodeCodeScanErrorNotifiedAt(ctx, gorm.Expr("NULL"), registered.CodeGitHubSettingID, registered.ProjectID); err != nil {
			return nil, err
		}
	}
	return &code.PutCodeScanSettingResponse{CodeScanSetting: convertCodeScanSetting(registered)}, nil
}

func (c *CodeService) DeleteCodeScanSetting(ctx context.Context, req *code.DeleteCodeScanSettingRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	// Delete all associated repositories (bulk delete)
	if err := c.repository.DeleteCodeScanRepository(ctx, req.ProjectId, req.GithubSettingId); err != nil {
		return nil, err
	}
	// Delete CodeScanSetting
	err := c.repository.DeleteCodeScanSetting(ctx, req.ProjectId, req.GithubSettingId)
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
		RepositoryName:  req.RepositoryName,
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
		RepositoryName:  req.RepositoryName,
	})
	if err != nil {
		return nil, err
	}
	if _, err = c.repository.UpsertDependencySetting(ctx, &code.DependencySettingForUpsert{
		GithubSettingId:   data.CodeGitHubSettingID,
		CodeDataSourceId:  data.CodeDataSourceID,
		ProjectId:         data.ProjectID,
		RepositoryPattern: data.RepositoryPattern,
		Status:            code.Status_IN_PROGRESS,
		StatusDetail:      fmt.Sprintf("Start scan at %+v", time.Now().Format(time.RFC3339)),
		ScanAt:            data.ScanAt.Unix(),
	}); err != nil {
		return nil, err
	}
	c.logger.Infof(ctx, "Invoke scanned, messageId: %v", resp.MessageId)
	return &empty.Empty{}, nil
}

func (c *CodeService) InvokeScanCodeScan(ctx context.Context, req *code.InvokeScanCodeScanRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := c.repository.GetCodeScanSetting(ctx, req.ProjectId, req.GithubSettingId)
	if err != nil {
		return nil, err
	}

	// Get list of repositories filtered by CodeScanSetting (RepositoryPattern, ScanPublic/Internal/Private, etc.)
	repos, err := c.listCodescanTargetRepository(ctx, req.ProjectId, req.GithubSettingId)
	if err != nil {
		errMsg := strings.ToLower(err.Error())
		// Check if error is authentication error (bad credentials, unauthorized)
		if strings.Contains(errMsg, "bad credentials") || strings.Contains(errMsg, "unauthorized") || strings.Contains(errMsg, "401") {
			c.logger.Errorf(ctx, "GitHub API authentication error when listing repositories: project_id=%d, github_setting_id=%d, err=%+v (PAT may be expired or invalid)", req.ProjectId, req.GithubSettingId, err)
			// Update status to ERROR
			if _, updateErr := c.repository.UpsertCodeScanSetting(ctx, &code.CodeScanSettingForUpsert{
				GithubSettingId:   data.CodeGitHubSettingID,
				CodeDataSourceId:  data.CodeDataSourceID,
				ProjectId:         data.ProjectID,
				RepositoryPattern: data.RepositoryPattern,
				ScanPublic:        data.ScanPublic,
				ScanInternal:      data.ScanInternal,
				ScanPrivate:       data.ScanPrivate,
				Status:            code.Status_ERROR,
				StatusDetail:      fmt.Sprintf("Authentication error: %v", err),
				ScanAt:            data.ScanAt.Unix(),
			}); updateErr != nil {
				c.logger.Errorf(ctx, "Failed to update status to ERROR: project_id=%d, github_setting_id=%d, err=%+v", req.ProjectId, req.GithubSettingId, updateErr)
				return nil, fmt.Errorf("failed to update status: %w", updateErr)
			}
			// Return nil to allow InvokeScanAll to continue with other settings
			return &empty.Empty{}, nil
		}
		// For other errors, return error as before
		return nil, err
	}

	if len(repos) == 0 {
		c.logger.Warnf(ctx, "No repositories found for scanning: project_id=%d, github_setting_id=%d (this may be expected if filter criteria exclude all repositories)",
			req.ProjectId, req.GithubSettingId)
		return &empty.Empty{}, nil
	}

	var messageIDs []string
	for _, repo := range repos {
		if repo.FullName == nil {
			c.logger.Errorf(ctx, "Repository with nil FullName found: project_id=%d, github_setting_id=%d, repo_id=%v, succeeded=%d before failure",
				req.ProjectId, req.GithubSettingId, repo.ID, len(messageIDs))
			return nil, fmt.Errorf("repository with nil FullName found (repo_id=%v)", repo.ID)
		}
		resp, err := c.sqs.Send(ctx, c.codeCodeScanQueueURL, &message.CodeQueueMessage{
			GitHubSettingID: data.CodeGitHubSettingID,
			ProjectID:       data.ProjectID,
			ScanOnly:        req.ScanOnly,
			RepositoryName:  *repo.FullName,
		})
		if err != nil {
			c.logger.Errorf(ctx, "Failed to send message for repository %s: project_id=%d, github_setting_id=%d, succeeded=%d before failure, err=%+v",
				*repo.FullName, req.ProjectId, req.GithubSettingId, len(messageIDs), err)
			return nil, fmt.Errorf("failed to send message for repository %s", *repo.FullName)
		}
		if resp.MessageId != nil {
			messageIDs = append(messageIDs, *resp.MessageId)
		}
	}

	// Update status only if all messages were sent successfully
	statusDetail := fmt.Sprintf("Start scan at %+v, attempted=%d, succeeded=%d", time.Now().Format(time.RFC3339), len(repos), len(messageIDs))

	if _, err = c.repository.UpsertCodeScanSetting(ctx, &code.CodeScanSettingForUpsert{
		GithubSettingId:   data.CodeGitHubSettingID,
		CodeDataSourceId:  data.CodeDataSourceID,
		ProjectId:         data.ProjectID,
		RepositoryPattern: data.RepositoryPattern,
		ScanPublic:        data.ScanPublic,
		ScanInternal:      data.ScanInternal,
		ScanPrivate:       data.ScanPrivate,
		Status:            code.Status_IN_PROGRESS,
		StatusDetail:      statusDetail,
		ScanAt:            data.ScanAt.Unix(),
	}); err != nil {
		return nil, err
	}

	c.logger.Infof(ctx, "Invoke scanned: project_id=%d, github_setting_id=%d, attempted=%d, succeeded=%d, messageIds: %v", req.ProjectId, req.GithubSettingId, len(repos), len(messageIDs), messageIDs)
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
	listCodeScan, err := c.repository.ListCodeScanSetting(ctx, 0)
	if err != nil {
		return nil, err
	}
	for _, codescan := range *listCodeScan {
		if codescan.ProjectID == 0 || codescan.CodeDataSourceID == 0 {
			continue
		}
		if resp, err := c.projectClient.IsActive(ctx, &project.IsActiveRequest{ProjectId: codescan.ProjectID}); err != nil {
			c.logger.Errorf(ctx, "Failed to project.IsActive API, err=%+v", err)
			return nil, err
		} else if !resp.Active {
			c.logger.Infof(ctx, "Skip deactive project, project_id=%d", codescan.ProjectID)
			continue
		}
		if _, err := c.InvokeScanCodeScan(ctx, &code.InvokeScanCodeScanRequest{
			GithubSettingId: codescan.CodeGitHubSettingID,
			ProjectId:       codescan.ProjectID,
			ScanOnly:        true,
		}); err != nil {
			errMsg := strings.ToLower(err.Error())
			// Check if error is from GetCodeScanSetting (database error - should return error)
			if strings.Contains(errMsg, "code scan setting not found") || strings.Contains(errMsg, "invalid db") || strings.Contains(errMsg, "record not found") {
				c.logger.Errorf(ctx, "InvokeScanCodeScan database error: project_id=%d, code_github_setting_id=%d, err=%+v", codescan.ProjectID, codescan.CodeGitHubSettingID, err)
				return nil, err
			}
			c.logger.Errorf(ctx, "InvokeScanCodeScan error occured: project_id=%d, code_github_setting_id=%d, err=%+v (skipping this setting)", codescan.ProjectID, codescan.CodeGitHubSettingID, err)
			continue
		}
	}
	return &empty.Empty{}, nil
}

func (c *CodeService) PutCodeScanRepository(ctx context.Context, req *code.PutCodeScanRepositoryRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	// Upsert repository status
	_, err := c.repository.UpsertCodeScanRepository(ctx, req.ProjectId, req.CodeScanRepository)
	if err != nil {
		return nil, err
	}
	c.logger.Infof(ctx, "PutCodeScanRepository: project_id=%d, github_setting_id=%d, repository=%s, status=%s",
		req.ProjectId, req.CodeScanRepository.GithubSettingId, req.CodeScanRepository.RepositoryFullName, req.CodeScanRepository.Status.String())
	c.logger.Debugf(ctx, "PutCodeScanRepository: status_detail=%s", req.CodeScanRepository.StatusDetail)
	return &empty.Empty{}, nil
}

// ListCodescanTargetRepository lists repositories filtered by CodeScanSetting (RPC handler)
func (c *CodeService) ListCodescanTargetRepository(ctx context.Context, req *code.ListCodescanTargetRepositoryRequest) (*code.ListCodescanTargetRepositoryResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	repos, err := c.listCodescanTargetRepository(ctx, req.ProjectId, req.GithubSettingId)
	if err != nil {
		return nil, err
	}
	// Convert GitHub repositories to response format
	repositoryList := code.ListCodescanTargetRepositoryResponse{}
	for _, repo := range repos {
		repositoryList.Repository = append(repositoryList.Repository, convertGitHubRepository(repo))
	}
	return &repositoryList, nil
}

// listCodescanTargetRepository lists repositories filtered by CodeScanSetting (internal function)
func (c *CodeService) listCodescanTargetRepository(ctx context.Context, projectID, githubSettingID uint32) ([]*ghub.Repository, error) {
	// Get GitHub setting to use for API call
	githubSetting, err := c.repository.GetGitHubSetting(ctx, projectID, githubSettingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*ghub.Repository{}, nil
		}
		return nil, err
	}
	// Get CodeScanSetting from database to use saved filter options
	codeScanSetting, err := c.repository.GetCodeScanSetting(ctx, projectID, githubSettingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.logger.Errorf(ctx, "CodeScanSetting not found: project_id=%d, github_setting_id=%d", projectID, githubSettingID)
			return nil, fmt.Errorf("code scan setting not found: project_id=%d, github_setting_id=%d", projectID, githubSettingID)
		}
		return nil, err
	}

	// Decrypt PAT before using it for GitHub API call
	decryptedPAT := ""
	if githubSetting.PersonalAccessToken != "" {
		decrypted, err := decryptWithBase64(&c.cipherBlock, githubSetting.PersonalAccessToken)
		if err != nil {
			c.logger.Errorf(ctx, "Failed to decrypt PAT: err=%+v", err)
			return nil, fmt.Errorf("failed to decrypt PAT: %w", err)
		}
		decryptedPAT = decrypted
		if decryptedPAT == "" {
			c.logger.Errorf(ctx, "Failed to decrypt PAT: decrypted PAT is empty, project_id=%d, github_setting_id=%d", projectID, githubSettingID)
			return nil, fmt.Errorf("decrypted PAT is empty")
		}
	} else {
		c.logger.Warnf(ctx, "PersonalAccessToken is empty for github_setting_id=%d", githubSettingID)
	}

	// Convert model to proto for GitHub API call
	protoGitHubSetting := convertGitHubSetting(githubSetting, nil, nil, nil, false)
	// Override with decrypted PAT
	protoGitHubSetting.PersonalAccessToken = decryptedPAT

	// Call GitHub API to list repositories
	repos, err := c.githubClient.ListRepository(ctx, protoGitHubSetting, "")
	if err != nil {
		return nil, err
	}

	// Apply size-based filtering
	repos = github.FilterExcludedRepositories(repos, c.limitRepositorySizeKb)

	// Apply filter options based on CodeScanSetting
	if codeScanSetting != nil {
		filterOpts := &github.FilterOptions{
			RepositoryPattern: codeScanSetting.RepositoryPattern,
			ScanPublic:        codeScanSetting.ScanPublic,
			ScanInternal:      codeScanSetting.ScanInternal,
			ScanPrivate:       codeScanSetting.ScanPrivate,
		}
		repos = github.ApplyFilters(repos, filterOpts)
	}

	return repos, nil
}

func convertGitHubRepository(repo *ghub.Repository) *code.GitHubRepository {
	if repo == nil {
		return &code.GitHubRepository{}
	}
	converted := code.GitHubRepository{
		Private: false,
	}
	if repo.ID != nil {
		converted.Id = *repo.ID
	}
	if repo.FullName != nil {
		converted.FullName = *repo.FullName
	}
	if repo.CloneURL != nil {
		converted.CloneUrl = *repo.CloneURL
	}
	if repo.Private != nil {
		converted.Private = *repo.Private
	}
	if repo.DefaultBranch != nil {
		converted.DefaultBranch = *repo.DefaultBranch
	}
	return &converted
}
