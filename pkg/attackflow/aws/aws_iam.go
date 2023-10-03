package aws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/pkg/attackflow"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

type iamAnalyzer struct {
	resource  *datasource.Resource
	metadata  *iamMetadata
	awsConfig *aws.Config
	client    *iam.Client
	logger    logging.Logger
}
type iamMetadata struct {
	IamRoleArn      string    `json:"iam_role_arn"`
	AllowedService  []*string `json:"allowed_service"`
	AccessedService []*string `json:"accessed_service"`
}

func newIAMAnalyzer(arn string, cfg *aws.Config, logger logging.Logger) (attackflow.CloudServiceAnalyzer, error) {
	return &iamAnalyzer{
		resource:  getAWSInfoFromARN(arn),
		metadata:  &iamMetadata{},
		awsConfig: cfg,
		client:    iam.NewFromConfig(*cfg),
		logger:    logger,
	}, nil
}

func (i *iamAnalyzer) Analyze(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, error,
) {
	// cache
	cachedResource, cachedMeta, err := getIAMAttackFlowCache(i.resource.CloudId, i.resource.ResourceName)
	if err != nil {
		return nil, err
	}
	if cachedResource != nil && cachedMeta != nil {
		i.resource = cachedResource
		i.metadata = cachedMeta
		resp = attackflow.SetNode(false, "", cachedResource, resp)
		return resp, nil
	}

	i.metadata.IamRoleArn = i.resource.ResourceName
	// instance-profile
	if strings.HasPrefix(i.resource.ResourceName, fmt.Sprintf("arn:aws:iam::%s:instance-profile/", i.resource.CloudId)) {
		role, err := i.client.GetInstanceProfile(ctx, &iam.GetInstanceProfileInput{
			InstanceProfileName: &i.resource.ShortName,
		})
		if err != nil {
			return nil, err
		}
		if len(role.InstanceProfile.Roles) == 0 {
			return nil, errors.New("no IAM role found from instance profile")
		}
		// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html#ec2-instance-profile
		// > An instance profile can contain only one IAM role.
		i.metadata.IamRoleArn = aws.ToString(role.InstanceProfile.Roles[0].Arn) // update
	}

	jobID, err := i.generateServiceLastAccessedDetails(ctx, i.metadata.IamRoleArn)
	if err != nil {
		return nil, err
	}
	allowedServices, accessedServices, err := i.analyzeServiceLastAccessedDetails(ctx, jobID)
	if err != nil {
		return nil, err
	}
	i.metadata.AllowedService = allowedServices
	i.metadata.AccessedService = accessedServices

	i.resource.MetaData, err = attackflow.ParseMetadata(i.metadata)
	if err != nil {
		return nil, err
	}
	resp = attackflow.SetNode(false, "", i.resource, resp)

	// cache
	if err := attackflow.SetAttackFlowCache(i.resource.CloudId, i.resource.ResourceName, i.resource); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *iamAnalyzer) generateServiceLastAccessedDetails(ctx context.Context, arn string) (string, error) {
	// https://docs.aws.amazon.com/IAM/latest/APIReference/API_GenerateServiceLastAccessedDetails.html
	out, err := i.client.GenerateServiceLastAccessedDetails(ctx, &iam.GenerateServiceLastAccessedDetailsInput{
		Arn:         &arn,
		Granularity: types.AccessAdvisorUsageGranularityTypeServiceLevel,
	})
	if err != nil {
		return "", err
	}
	return *out.JobId, nil
}

const MAX_RETRY = 3

func (i *iamAnalyzer) analyzeServiceLastAccessedDetails(ctx context.Context, jobID string) ([]*string, []*string, error) {
	allowedServices := []*string{}
	accessedServices := []*string{}
	for idx := 0; idx < MAX_RETRY; idx++ {
		// https://docs.aws.amazon.com/IAM/latest/APIReference/API_GetServiceLastAccessedDetails.html
		out, err := i.client.GetServiceLastAccessedDetails(ctx, &iam.GetServiceLastAccessedDetailsInput{
			JobId: &jobID,
		})
		if err != nil {
			return nil, nil, err
		}
		if out.JobStatus == types.JobStatusTypeFailed {
			errMsg := fmt.Sprintf("failed to GetServiceLastAccessedDetails, jobID=%s", jobID)
			if out.Error != nil {
				errMsg += fmt.Sprintf(" error_code=%s, message=%s", *out.Error.Code, *out.Error.Message)
			}
			return nil, nil, errors.New(errMsg)
		}
		if out.JobStatus == types.JobStatusTypeInProgress {
			time.Sleep(time.Second)
			continue
		}
		if out.JobStatus == types.JobStatusTypeCompleted {
			for _, accessed := range out.ServicesLastAccessed {
				allowedServices = append(allowedServices, accessed.ServiceName)
				if accessed.LastAuthenticated != nil {
					accessedServices = append(accessedServices, accessed.ServiceName)
				}
			}
			break
		}
		// unknown status
		return nil, nil, fmt.Errorf("unknown Job Status for GetServiceLastAccessedDetails: jobID=%s, status=%s", jobID, out.JobStatus)
	}
	return allowedServices, accessedServices, nil
}

func (i *iamAnalyzer) Next(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, []attackflow.CloudServiceAnalyzer, error,
) {
	return resp, []attackflow.CloudServiceAnalyzer{}, nil
}

func getIAMAttackFlowCache(cloudID, resourceName string) (*datasource.Resource, *iamMetadata, error) {
	resource, err := attackflow.GetAttackFlowCache(cloudID, resourceName)
	if err != nil {
		return nil, nil, err
	}
	if resource == nil {
		return nil, nil, nil
	}
	var meta iamMetadata
	if err := json.Unmarshal([]byte(resource.MetaData), &meta); err != nil {
		return nil, nil, err
	}
	return resource, &meta, nil
}
