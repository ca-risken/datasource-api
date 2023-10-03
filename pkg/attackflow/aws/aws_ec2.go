package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/pkg/attackflow"
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

func newEC2Analyzer(ctx context.Context, arn string, cfg *aws.Config, logger logging.Logger) (attackflow.CloudServiceAnalyzer, error) {
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
	// cache
	cachedResource, cachedMeta, err := getEC2AttackFlowCache(e.resource.CloudId, e.resource.ResourceName)
	if err != nil {
		return nil, err
	}
	if cachedResource != nil && cachedMeta != nil {
		e.resource = cachedResource
		e.metadata = cachedMeta
		resp = attackflow.SetNode(cachedMeta.IsPublic, cachedMeta.PublicIpAddress, cachedResource, resp)
		return resp, nil
	}

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

	e.resource.MetaData, err = attackflow.ParseMetadata(e.metadata)
	if err != nil {
		return nil, err
	}
	resp = attackflow.SetNode(e.metadata.IsPublic, e.metadata.PublicIpAddress, e.resource, resp)

	// cache
	if err := attackflow.SetAttackFlowCache(e.resource.CloudId, e.resource.ResourceName, e.resource); err != nil {
		return nil, err
	}
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
	*datasource.AnalyzeAttackFlowResponse, []attackflow.CloudServiceAnalyzer, error,
) {
	analyzers := []attackflow.CloudServiceAnalyzer{}
	// IAM role
	if e.metadata.IamInstanceProfile != "" {
		instanceProfile := getAWSInfoFromARN(e.metadata.IamInstanceProfile)
		resp.Edges = append(resp.Edges, attackflow.GetEdge(e.resource.ResourceName, instanceProfile.ResourceName, "iam role"))
		iamAnalyzer, err := newIAMAnalyzer(instanceProfile.ResourceName, e.awsConfig, e.logger)
		if err != nil {
			return nil, nil, err
		}
		analyzers = append(analyzers, iamAnalyzer)
	}
	return resp, analyzers, nil
}

func getEC2ARNFromPublicDomain(ctx context.Context, domain, accountID string, cfg *aws.Config) (string, error) {
	if !domainPatternEC2.MatchString(domain) {
		return "", fmt.Errorf("invalid domain: %s", domain)
	}
	// EC2 public domain format: ec2-1-2-3-4.ap-northeast-1.compute.amazonaws.com
	var err error
	region := strings.Split(domain, ".")[1]
	if cfg.Region != region {
		cfg, err = retrieveAWSCredentialWithRegion(ctx, *cfg, region)
		if err != nil {
			return "", err
		}
	}
	client := ec2.NewFromConfig(*cfg)
	// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeInstances.html
	instances, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("dns-name"),
				Values: []string{domain},
			},
		},
	})
	if err != nil {
		return "", err
	}
	if len(instances.Reservations) == 0 || len(instances.Reservations[0].Instances) == 0 {
		return "", fmt.Errorf("instance not found in %s, domain: %s", accountID, domain)
	}
	i := instances.Reservations[0].Instances[0] // only 1 instance (because filter by dns-name)
	return fmt.Sprintf("arn:aws:ec2:%s:%s:instance/%s", cfg.Region, accountID, aws.ToString(i.InstanceId)), nil
}

func getEC2AttackFlowCache(cloudID, resourceName string) (*datasource.Resource, *ec2Metadata, error) {
	resource, err := attackflow.GetAttackFlowCache(cloudID, resourceName)
	if err != nil {
		return nil, nil, err
	}
	if resource == nil {
		return nil, nil, nil
	}
	var meta ec2Metadata
	if err := json.Unmarshal([]byte(resource.MetaData), &meta); err != nil {
		return nil, nil, err
	}
	return resource, &meta, nil
}
