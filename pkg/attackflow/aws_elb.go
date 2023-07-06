package attackflow

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

type elbAnalyzer struct {
	resource  *datasource.Resource
	metadata  *elbMetadata
	awsConfig *aws.Config
	v1client  *elasticloadbalancing.Client
	v2client  *elasticloadbalancingv2.Client
	ec2client *ec2.Client
	logger    logging.Logger
}
type elbMetadata struct {
	Name           string        `json:"name"`
	DNSName        string        `json:"dns_name"`
	InternetFacing bool          `json:"internet_facing"`
	SecurityGroups []string      `json:"security_groups"`
	VpcID          string        `json:"vpc_id"`
	AccessLogging  bool          `json:"access_logging"`
	TargetGroups   []targetGroup `json:"target_groups"`
	IsPublic       bool          `json:"is_public"`
}

type targetGroup struct {
	// The ID of the target.
	// If the target type of the target group is instance, specify an instance ID.
	// If the target type is ip, specify an IP address.
	// If the target type is lambda, specify the ARN of the Lambda function.
	// If the target type is alb, specify the ARN of the Application Load Balancer target.
	TargetID   string `json:"target_id"`
	TargetType string `json:"target_type"`
}

func newELBAnalyzer(ctx context.Context, arn string, cfg *aws.Config, logger logging.Logger) (CloudServiceAnalyzer, error) {
	resource := getAWSInfoFromARN(arn)
	var err error
	if cfg.Region != resource.Region {
		cfg, err = retrieveAWSCredentialWithRegion(ctx, *cfg, resource.Region)
		if err != nil {
			return nil, err
		}
	}
	return &elbAnalyzer{
		resource:  resource,
		metadata:  &elbMetadata{},
		awsConfig: cfg,
		v1client:  elasticloadbalancing.NewFromConfig(*cfg),   // Classic Load Balancer
		v2client:  elasticloadbalancingv2.NewFromConfig(*cfg), // Application Load Balancer etc.
		ec2client: ec2.NewFromConfig(*cfg),
		logger:    logger,
	}, nil
}

func (e *elbAnalyzer) Analyze(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, error,
) {
	if isV2LoadBalancer(e.resource.ResourceName) {
		return e.analyzeV2(ctx, resp)
	} else {
		return e.analyzeV1(ctx, resp)
	}
}

func (e *elbAnalyzer) Next(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, []CloudServiceAnalyzer, error,
) {
	analyzers := []CloudServiceAnalyzer{}
	for _, target := range e.metadata.TargetGroups {
		switch target.TargetType {
		case "instance":
			ec2Arn := fmt.Sprintf("arn:aws:ec2:%s:%s:instance/%s", e.resource.Region, e.resource.CloudId, target.TargetID)
			ec2Analyzer, err := newEC2Analyzer(ctx, ec2Arn, e.awsConfig, e.logger)
			if err != nil {
				return nil, nil, err
			}
			analyzers = append(analyzers, ec2Analyzer)
			resp.Edges = append(resp.Edges, getEdge(e.resource.ResourceName, ec2Arn, "target"))
		case "lambda":
			lambdaAnalyzer, err := newLambdaAnalyzer(ctx, target.TargetID, e.awsConfig, e.logger)
			if err != nil {
				return nil, nil, err
			}
			analyzers = append(analyzers, lambdaAnalyzer)
			resp.Edges = append(resp.Edges, getEdge(e.resource.ResourceName, target.TargetID, "target"))
		case "alb":
			albAnalyzer, err := newELBAnalyzer(ctx, target.TargetID, e.awsConfig, e.logger)
			if err != nil {
				return nil, nil, err
			}
			analyzers = append(analyzers, albAnalyzer)
			resp.Edges = append(resp.Edges, getEdge(e.resource.ResourceName, target.TargetID, "target"))
		case "ip":
			resp.Edges = append(resp.Edges, getEdge(e.resource.ResourceName, target.TargetID, "target"))
			resp.Nodes = append(resp.Nodes, getInternalServiceNode(target.TargetID, e.awsConfig.Region))
		}
	}
	return resp, analyzers, nil
}

// Check if the ARN is a v2 format
// v1 ARN: arn:aws:elasticloadbalancing:{region}:{account-id}:loadbalancer/{name} (RISKEN original format because No document & cannot see arn)
// v2 ARN: arn:aws:elasticloadbalancing:{region}:{account-id}:loadbalancer/{type}/{name}/{id}
func isV2LoadBalancer(arn string) bool {
	splitArn := strings.Split(arn, ":")
	if len(splitArn) < 5 {
		return false
	}
	lbPart := strings.Join(splitArn[5:], "")
	return strings.HasPrefix(lbPart, "loadbalancer/") && len(strings.Split(lbPart, "/")) >= 4
}

func getV2LoadBalancerName(arn string) string {
	splitArn := strings.Split(arn, ":")
	if len(splitArn) < 5 {
		return ""
	}
	lbPart := strings.Split(strings.Join(splitArn[5:], ""), "/")
	if len(lbPart) < 3 {
		return ""
	}
	return lbPart[2]
}

func (e *elbAnalyzer) setPublic(ctx context.Context) error {
	// public or not
	for _, sg := range e.metadata.SecurityGroups {
		isPublic, err := isPublicSecurityGroup(ctx, e.ec2client, sg)
		if err != nil {
			return err
		}
		if isPublic {
			e.metadata.IsPublic = true
			break
		}
	}
	return nil
}
