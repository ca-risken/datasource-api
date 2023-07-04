package attackflow

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/ca-risken/common/pkg/logging"
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
	AllowedService  []*string `json:"allowed_service"`
	AccessedService []*string `json:"accessed_service"`
}

func newIAMAnalyzer(arn string, cfg *aws.Config, logger logging.Logger) (CloudServiceAnalyzer, error) {
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
	jobID, err := i.generateServiceLastAccessedDetails(ctx, i.resource.ResourceName)
	if err != nil {
		return nil, err
	}
	allowedServices, accessedServices, err := i.analyzeServiceLastAccessedDetails(ctx, jobID)
	if err != nil {
		return nil, err
	}
	i.metadata.AllowedService = allowedServices
	i.metadata.AccessedService = accessedServices

	i.resource.MetaData, err = parseMetadata(i.metadata)
	if err != nil {
		return nil, err
	}
	resp = setNode(false, "", i.resource, resp)
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
	*datasource.AnalyzeAttackFlowResponse, []CloudServiceAnalyzer, error,
) {
	return resp, []CloudServiceAnalyzer{}, nil
}
