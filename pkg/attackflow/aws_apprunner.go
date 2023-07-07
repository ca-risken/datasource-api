package attackflow

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apprunner"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

type appRunnerAnalyzer struct {
	resource  *datasource.Resource
	metadata  *appRunnerMetadata
	awsConfig *aws.Config
	client    *apprunner.Client
	logger    logging.Logger
}
type appRunnerMetadata struct {
	Name            string `json:"name"`
	AutoScaling     string `json:"auto_scaling"`
	IamRole         string `json:"iam_role"`
	ComputeResource string `json:"compute_resource"`
	AutoDeploy      bool   `json:"auto_deploy"`
	State           string `json:"state"`
	CodeRepository  string `json:"code_repository"`
	ImageRepository string `json:"image_repository"`
	ServiceUrl      string `json:"service_url"`
	EncryptionKey   string `json:"encryption_key"`
	IsPublic        bool   `json:"is_public"`
	VpcID           string `json:"vpc_id"`
}

func newAppRunnerAnalyzer(ctx context.Context, arn string, cfg *aws.Config, logger logging.Logger) (CloudServiceAnalyzer, error) {
	resource := getAWSInfoFromARN(arn)
	var err error
	if cfg.Region != resource.Region {
		cfg, err = retrieveAWSCredentialWithRegion(ctx, *cfg, resource.Region)
		if err != nil {
			return nil, err
		}
	}
	return &appRunnerAnalyzer{
		resource:  resource,
		metadata:  &appRunnerMetadata{},
		awsConfig: cfg,
		client:    apprunner.NewFromConfig(*cfg),
		logger:    logger,
	}, nil
}

func (a *appRunnerAnalyzer) Analyze(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, error,
) {
	// https://docs.aws.amazon.com/apprunner/latest/api/API_DescribeService.html
	service, err := a.client.DescribeService(ctx, &apprunner.DescribeServiceInput{
		ServiceArn: aws.String(a.resource.ResourceName),
	})
	if err != nil {
		return nil, err
	}
	a.metadata.Name = aws.ToString(service.Service.ServiceName)
	a.resource.ShortName = a.metadata.Name // update short name
	a.metadata.AutoScaling = aws.ToString(service.Service.AutoScalingConfigurationSummary.AutoScalingConfigurationName)
	a.metadata.IamRole = aws.ToString(service.Service.InstanceConfiguration.InstanceRoleArn)
	a.metadata.ComputeResource = a.getCpuMemLabel(ctx,
		aws.ToString(service.Service.InstanceConfiguration.Cpu),
		aws.ToString(service.Service.InstanceConfiguration.Memory),
	)
	a.metadata.AutoDeploy = aws.ToBool(service.Service.SourceConfiguration.AutoDeploymentsEnabled)
	a.metadata.State = fmt.Sprint(service.Service.Status)
	if service.Service.SourceConfiguration.CodeRepository != nil {
		a.metadata.CodeRepository = aws.ToString(service.Service.SourceConfiguration.CodeRepository.RepositoryUrl)
	}
	if service.Service.SourceConfiguration.ImageRepository != nil {
		a.metadata.ImageRepository = aws.ToString(service.Service.SourceConfiguration.ImageRepository.ImageIdentifier)
	}
	a.metadata.ServiceUrl = aws.ToString(service.Service.ServiceUrl)
	if service.Service.EncryptionConfiguration != nil {
		a.metadata.EncryptionKey = fmt.Sprint(service.Service.EncryptionConfiguration.KmsKey)
	}
	if a.metadata.EncryptionKey == "" {
		a.metadata.EncryptionKey = "AWS managed"
	}
	a.metadata.IsPublic = aws.ToBool(&service.Service.NetworkConfiguration.IngressConfiguration.IsPubliclyAccessible)
	a.metadata.VpcID = aws.ToString(service.Service.NetworkConfiguration.EgressConfiguration.VpcConnectorArn)

	a.resource.MetaData, err = parseMetadata(a.metadata)
	if err != nil {
		return nil, err
	}
	resp = setNode(a.metadata.IsPublic, "", a.resource, resp)
	return resp, nil
}

func (a *appRunnerAnalyzer) Next(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, []CloudServiceAnalyzer, error,
) {
	analyzers := []CloudServiceAnalyzer{}
	if a.metadata.IamRole != "" {
		resp.Edges = append(resp.Edges, getEdge(a.resource.ResourceName, a.metadata.IamRole, "iam role"))
		iamAnalyzer, err := newIAMAnalyzer(a.metadata.IamRole, a.awsConfig, a.logger)
		if err != nil {
			return nil, nil, err
		}
		analyzers = append(analyzers, iamAnalyzer)
	}
	return resp, analyzers, nil
}

func (a *appRunnerAnalyzer) getCpuMemLabel(ctx context.Context, cpu, mem string) string {
	cpuLabel := cpu
	memLabel := mem
	cpuInt, err := strconv.Atoi(cpu)
	if err != nil {
		a.logger.Warnf(ctx, "Failed to parse cpu: %s, err: %v", cpuLabel, err)
	} else {
		f := float64(cpuInt) / float64(1000)
		cpuLabel = fmt.Sprintf("%.2f", math.Floor(f*100)/100) // To two decimal places
	}
	memInt, err := strconv.Atoi(mem)
	if err != nil {
		a.logger.Warnf(ctx, "Failed to parse mem: %s, err: %v", memLabel, err)
	} else {
		f := float64(memInt) / float64(1000)
		memLabel = fmt.Sprintf("%.2f", math.Floor(f*100)/100) // To two decimal places
	}
	return fmt.Sprintf("%s vCPU, %s GB", cpuLabel, memLabel)

}
