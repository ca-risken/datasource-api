package attackflow

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

type ec2Analyzer struct {
	resource  *datasource.Resource
	metadata  *ec2Metadata
	awsConfig *aws.Config
	client    *ec2.Client
	logger    logging.Logger
}
type ec2Metadata struct {
	IamInstanceProfile string `json:"iam_instance_profile"`
	State              string `json:"state"`
	VpcID              string `json:"vpc_id"`
	SubnetID           string `json:"subnet_id"`
	IsPublic           bool   `json:"is_public"`
	PublicDnsName      string `json:"public_dns_name"`
	PublicIpAddress    string `json:"public_ip_address"`
	InstanceType       string `json:"instance_type"`
	Platform           string `json:"platform"`
	DefaultEncrypt     bool   `json:"default_encrypt"`
}

func newEC2Analyzer(ctx context.Context, arn string, cfg *aws.Config, logger logging.Logger) (CloudServiceAnalyzer, error) {
	resource := getAWSInfoFromARN(arn)
	var err error
	if cfg.Region != resource.Region {
		cfg, err = retrieveAWSCredentialWithRegion(ctx, *cfg, resource.Region)
		if err != nil {
			return nil, err
		}
	}
	return &ec2Analyzer{
		resource:  resource,
		metadata:  &ec2Metadata{},
		awsConfig: cfg,
		client:    ec2.NewFromConfig(*cfg),
		logger:    logger,
	}, nil
}

func (e *ec2Analyzer) Analyze(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, error,
) {
	// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeInstances.html
	instances, err := e.client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{e.resource.ShortName},
	})
	if err != nil {
		return nil, err
	}
	if len(instances.Reservations) == 0 || len(instances.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("instance not found: %s", e.resource.ResourceName)
	}
	instance := instances.Reservations[0].Instances[0]
	if instance.IamInstanceProfile != nil {
		e.metadata.IamInstanceProfile = aws.ToString(instance.IamInstanceProfile.Arn)
	}
	e.metadata.State = fmt.Sprint(instance.State.Name)
	e.metadata.VpcID = aws.ToString(instance.VpcId)
	e.metadata.SubnetID = aws.ToString(instance.SubnetId)
	e.metadata.PublicDnsName = aws.ToString(instance.PublicDnsName)
	e.metadata.PublicIpAddress = aws.ToString(instance.PublicIpAddress)
	e.metadata.InstanceType = fmt.Sprint(instance.InstanceType)
	e.metadata.Platform = fmt.Sprint(instance.Platform)
	hasPublicSG, err := e.hasPublicSecurityGroups(ctx, instance)
	if err != nil {
		return nil, err
	}
	if hasPublicSG && e.metadata.PublicIpAddress != "" && e.metadata.State == "running" {
		e.metadata.IsPublic = true
	}

	// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_GetEbsEncryptionByDefault.html
	encryption, err := e.client.GetEbsEncryptionByDefault(ctx, &ec2.GetEbsEncryptionByDefaultInput{})
	if err != nil {
		return nil, err
	}
	e.metadata.DefaultEncrypt = aws.ToBool(encryption.EbsEncryptionByDefault)

	e.resource.MetaData, err = parseMetadata(e.metadata)
	if err != nil {
		return nil, err
	}
	resp = setNode(e.metadata.IsPublic, e.metadata.PublicIpAddress, e.resource, resp)
	return resp, nil
}

func (e *ec2Analyzer) hasPublicSecurityGroups(ctx context.Context, instance types.Instance) (bool, error) {
	for _, sg := range instance.SecurityGroups {
		isPublic, err := isPublicSecurityGroup(ctx, e.client, aws.ToString(sg.GroupId))
		if err != nil {
			return false, err
		}
		if isPublic {
			return true, nil
		}
	}
	return false, nil
}

func (e *ec2Analyzer) Next(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, []CloudServiceAnalyzer, error,
) {
	analyzers := []CloudServiceAnalyzer{}
	// IAM role
	if e.metadata.IamInstanceProfile != "" {
		instanceProfile := getAWSInfoFromARN(e.metadata.IamInstanceProfile)
		resp.Edges = append(resp.Edges, getEdge(e.resource.ResourceName, instanceProfile.ResourceName, "iam role"))
		iamAnalyzer, err := newIAMAnalyzer(instanceProfile.ResourceName, e.awsConfig, e.logger)
		if err != nil {
			return nil, nil, err
		}
		analyzers = append(analyzers, iamAnalyzer)
	}
	return resp, analyzers, nil
}
