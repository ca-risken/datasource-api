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

func convertGitleaks(githubSetting *model.CodeGitHubSetting, gitleaksSetting *model.CodeGitleaksSetting, maskKey bool) *code.Gitleaks {
	var gitleaks code.Gitleaks
	if githubSetting == nil && gitleaksSetting == nil {
		return &gitleaks
	}
	gitleaks = code.Gitleaks{
		GitleaksId:          githubSetting.CodeGitHubSettingID,
		CodeDataSourceId:    gitleaksSetting.CodeDataSourceID,
		Name:                githubSetting.Name,
		ProjectId:           githubSetting.ProjectID,
		Type:                getType(githubSetting.Type),
		BaseUrl:             githubSetting.BaseURL,
		TargetResource:      githubSetting.TargetResource,
		RepositoryPattern:   gitleaksSetting.RepositoryPattern,
		GithubUser:          githubSetting.GitHubUser,
		PersonalAccessToken: githubSetting.PersonalAccessToken,
		ScanPublic:          gitleaksSetting.ScanPublic,
		ScanInternal:        gitleaksSetting.ScanInternal,
		ScanPrivate:         gitleaksSetting.ScanPrivate,
		Status:              getStatus(gitleaksSetting.Status),
		StatusDetail:        gitleaksSetting.StatusDetail,
		CreatedAt:           gitleaksSetting.CreatedAt.Unix(),
		UpdatedAt:           gitleaksSetting.UpdatedAt.Unix(),
	}
	if gitleaks.PersonalAccessToken != "" && maskKey {
		gitleaks.PersonalAccessToken = maskData // Masking sensitive data.
	}
	if !zero.IsZeroVal(gitleaksSetting.ScanAt) {
		gitleaks.ScanAt = gitleaksSetting.ScanAt.Unix()
	}
	return &gitleaks
}

func (c *CodeService) ListGitleaks(ctx context.Context, req *code.ListGitleaksRequest) (*code.ListGitleaksResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := c.repository.ListGitHubSetting(ctx, req.ProjectId, req.GitleaksId)
	if err != nil {
		return nil, err
	}
	data := code.ListGitleaksResponse{}
	for _, d := range *list {
		gitleaksSetting, err := c.repository.GetGitleaksSetting(ctx, d.ProjectID, d.CodeGitHubSettingID)
		if err != nil {
			return nil, err
		}
		data.Gitleaks = append(data.Gitleaks, convertGitleaks(&d, gitleaksSetting, true))
	}
	return &data, nil
}

func (c *CodeService) GetGitleaks(ctx context.Context, req *code.GetGitleaksRequest) (*code.GetGitleaksResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	githubSetting, err := c.repository.GetGitHubSetting(ctx, req.ProjectId, req.GitleaksId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &code.GetGitleaksResponse{}, nil
		}
		return nil, err
	}
	gitleaksSetting, err := c.repository.GetGitleaksSetting(ctx, githubSetting.ProjectID, githubSetting.CodeGitHubSettingID)
	if err != nil {
		return nil, err
	}
	return &code.GetGitleaksResponse{Gitleaks: convertGitleaks(githubSetting, gitleaksSetting, false)}, nil
}

func (c *CodeService) PutGitleaks(ctx context.Context, req *code.PutGitleaksRequest) (*code.PutGitleaksResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if req.Gitleaks.PersonalAccessToken != "" && req.Gitleaks.PersonalAccessToken != maskData {
		encrypted, err := encryptWithBase64(&c.cipherBlock, req.Gitleaks.PersonalAccessToken)
		if err != nil {
			c.logger.Errorf(ctx, "Failed to encrypt PAT: err=%+v", err)
			return nil, err
		}
		req.Gitleaks.PersonalAccessToken = encrypted
	} else {
		req.Gitleaks.PersonalAccessToken = "" // for not update token.
	}
	registeredGitHubSetting, err := c.repository.UpsertGitHubSetting(ctx, req.Gitleaks)
	if err != nil {
		return nil, err
	}
	req.Gitleaks.GitleaksId = registeredGitHubSetting.CodeGitHubSettingID
	registeredGitleaksSetting, err := c.repository.UpsertGitleaksSetting(ctx, req.Gitleaks)
	if err != nil {
		return nil, err
	}
	return &code.PutGitleaksResponse{Gitleaks: convertGitleaks(registeredGitHubSetting, registeredGitleaksSetting, true)}, nil
}

func (c *CodeService) DeleteGitleaks(ctx context.Context, req *code.DeleteGitleaksRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := c.repository.DeleteGitHubSetting(ctx, req.ProjectId, req.GitleaksId)
	if err != nil {
		return nil, err
	}
	err = c.repository.DeleteGitleaksSetting(ctx, req.ProjectId, req.GitleaksId)
	if err != nil {
		return nil, err
	}
	organizations, err := c.repository.ListGitHubEnterpriseOrg(ctx, req.ProjectId, req.GitleaksId)
	if err != nil {
		return nil, err
	}
	for _, org := range *organizations {
		err = c.repository.DeleteGitHubEnterpriseOrg(ctx, org.ProjectID, org.CodeGitHubSettingID, org.Organization)
		if err != nil {
			return nil, err
		}
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

func convertEnterpriseOrg(data *model.CodeGitHubEnterpriseOrg) *code.EnterpriseOrg {
	if data == nil {
		return &code.EnterpriseOrg{}
	}
	return &code.EnterpriseOrg{
		GitleaksId: data.CodeGitHubSettingID,
		Login:      data.Organization,
		ProjectId:  data.ProjectID,
		CreatedAt:  data.CreatedAt.Unix(),
		UpdatedAt:  data.CreatedAt.Unix(),
	}
}

func (c *CodeService) ListEnterpriseOrg(ctx context.Context, req *code.ListEnterpriseOrgRequest) (*code.ListEnterpriseOrgResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := c.repository.ListGitHubEnterpriseOrg(ctx, req.ProjectId, req.GitleaksId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &code.ListEnterpriseOrgResponse{}, nil
		}
		return nil, err
	}
	data := code.ListEnterpriseOrgResponse{}
	for _, d := range *list {
		data.EnterpriseOrg = append(data.EnterpriseOrg, convertEnterpriseOrg(&d))
	}
	return &data, nil
}

func (c *CodeService) PutEnterpriseOrg(ctx context.Context, req *code.PutEnterpriseOrgRequest) (*code.PutEnterpriseOrgResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	registered, err := c.repository.UpsertGitHubEnterpriseOrg(ctx, req.EnterpriseOrg)
	if err != nil {
		return nil, err
	}
	return &code.PutEnterpriseOrgResponse{EnterpriseOrg: convertEnterpriseOrg(registered)}, nil
}

func (c *CodeService) DeleteEnterpriseOrg(ctx context.Context, req *code.DeleteEnterpriseOrgRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := c.repository.DeleteGitHubEnterpriseOrg(ctx, req.ProjectId, req.GitleaksId, req.Login)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (c *CodeService) InvokeScanGitleaks(ctx context.Context, req *code.InvokeScanGitleaksRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := c.repository.GetGitleaksSetting(ctx, req.ProjectId, req.GitleaksId)
	if err != nil {
		return nil, err
	}
	resp, err := c.sqs.Send(ctx, c.sqs.CodeGitleaksQueueURL, &message.GitleaksQueueMessage{
		GitleaksID: data.CodeGitHubSettingID,
		ProjectID:  data.ProjectID,
		ScanOnly:   req.ScanOnly,
	})
	if err != nil {
		return nil, err
	}
	if _, err = c.repository.UpsertGitleaksSetting(ctx, &code.GitleaksForUpsert{
		GitleaksId:        data.CodeGitHubSettingID,
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

func (c *CodeService) InvokeScanAllGitleaks(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	list, err := c.repository.ListGitleaksSetting(ctx, 0)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &empty.Empty{}, nil
		}
		return nil, err
	}
	for _, g := range *list {
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
			GitleaksId: g.CodeGitHubSettingID,
			ProjectId:  g.ProjectID,
			ScanOnly:   true,
		}); err != nil {
			c.logger.Errorf(ctx, "InvokeScanGitleaks error occured: code_github_setting_id=%d, err=%+v", g.CodeGitHubSettingID, err)
			return nil, err
		}
	}
	return &empty.Empty{}, nil
}
