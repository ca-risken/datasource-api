package attackflow

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

const (
	// region
	REGION_US_EAST_1 = "us-east-1"

	// service
	SERVICE_CLOUDFRONT   = "cloudfront"
	SERVICE_S3           = "s3"
	SERVICE_LAMBDA       = "lambda"
	SERVICE_SQS          = "sqs"
	SERVICE_SNS          = "sns"
	SERVICE_EVENT_BRIDGE = "events"
	SERVICE_IAM          = "iam"
	SERVICE_API_GATEWAY  = "apigateway"
	SERVICE_EC2          = "ec2"
	SERVICE_ELB          = "elasticloadbalancing"

	RETRY_MAX_ATTEMPT = 10
)

var (
	domainPatternCloudFront = regexp.MustCompile(`\.cloudfront\.net$`)
	domainPatternS3         = regexp.MustCompile(`\.s3\..*\.amazonaws\.com$`)  // https://docs.aws.amazon.com/ja_jp/AmazonS3/latest/userguide/VirtualHosting.html
	domainPatternLambdaURL  = regexp.MustCompile(`\.lambda-url\..*\.on\.aws$`) // https://docs.aws.amazon.com/lambda/latest/dg/lambda-urls.html

	supportedAWSServices = map[string]bool{
		SERVICE_CLOUDFRONT:  true,
		SERVICE_S3:          true,
		SERVICE_LAMBDA:      true,
		SERVICE_API_GATEWAY: true,
		SERVICE_EC2:         true,
		SERVICE_ELB:         true,
		// TODO support below services
		// "app-runner":    true,
	}
)

type AWS struct {
	cloudID        string
	region         string
	initialService string
	logger         logging.Logger

	// aws client
	role      *model.AWSRelDataSource
	awsConfig *aws.Config
}

func NewAWS(
	ctx context.Context,
	req *datasource.AnalyzeAttackFlowRequest,
	awsrepo db.AWSRepoInterface,
	logger logging.Logger,
) (CSP, error) {

	r := getAWSInfoFromARN(req.ResourceName)
	role, err := awsrepo.GetAWSRelDataSourceByAccountID(ctx, req.ProjectId, req.CloudId)
	if err != nil {
		return nil, err
	}
	csp := &AWS{
		cloudID:        req.CloudId,
		region:         r.Region,
		initialService: r.Service,
		role:           role,
		logger:         logger,
	}
	if err := csp.retrieveAWSCredential(ctx); err != nil {
		return nil, err
	}
	return csp, nil
}

func getAWSInfoFromARN(arn string) *datasource.Resource {
	// arn:aws:iam::123456789012:user/MyUser -> Service: iam, Region: global, ShortName: MyUser
	splitArn := strings.Split(arn, ":")
	if len(splitArn) < 5 {
		return nil
	}

	// shortName
	shortName := strings.Join(splitArn[5:], "/")
	if strings.Contains(shortName, "/") {
		splitName := strings.Split(shortName, "/")
		shortName = splitName[len(splitName)-1]
	}

	// region
	region := splitArn[3]
	if region == "" {
		region = REGION_GLOBAL
	}
	return &datasource.Resource{
		ResourceName: arn,
		ShortName:    shortName,
		CloudType:    splitArn[1],
		CloudId:      splitArn[4],
		Service:      splitArn[2],
		Region:       region,
		Layer:        getLayerFromAWSService(splitArn[2]),
	}
}

func (a *AWS) retrieveAWSCredential(ctx context.Context) error {
	region := a.region
	if region == REGION_GLOBAL {
		region = REGION_US_EAST_1
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithRetryMaxAttempts(RETRY_MAX_ATTEMPT),
	)
	if err != nil {
		return err
	}
	stsClient := sts.NewFromConfig(cfg)
	provider := stscreds.NewAssumeRoleProvider(stsClient, a.role.AssumeRoleArn,
		func(p *stscreds.AssumeRoleOptions) {
			p.RoleSessionName = "RISKEN"
			p.ExternalID = &a.role.ExternalID
		},
	)
	cfg.Credentials = aws.NewCredentialsCache(provider)
	_, err = cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return err
	}
	a.awsConfig = &cfg
	return nil
}

func retrieveAWSCredentialWithRegion(ctx context.Context, awsConfig aws.Config, region string) (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithRetryMaxAttempts(RETRY_MAX_ATTEMPT),
	)
	if err != nil {
		return nil, err
	}
	cfg.Credentials = awsConfig.Credentials
	return &cfg, nil
}

func (a *AWS) GetInitialServiceAnalyzer(ctx context.Context, req *datasource.AnalyzeAttackFlowRequest) (
	CloudServiceAnalyzer, error,
) {
	var err error
	var serviceAnalyzer CloudServiceAnalyzer

	// check supported initial service
	switch a.initialService {
	case SERVICE_CLOUDFRONT:
		serviceAnalyzer, err = newCloudFrontAnalyzer(req.ResourceName, a.awsConfig, a.logger)
	case SERVICE_S3:
		serviceAnalyzer, err = newS3Analyzer(req.ResourceName, a.awsConfig, a.logger)
	case SERVICE_LAMBDA:
		serviceAnalyzer, err = newLambdaAnalyzer(ctx, req.ResourceName, a.awsConfig, a.logger)
	case SERVICE_API_GATEWAY:
		serviceAnalyzer, err = newAPIGatewayAnalyzer(ctx, req.ResourceName, a.awsConfig, a.logger)
	case SERVICE_EC2:
		serviceAnalyzer, err = newEC2Analyzer(ctx, req.ResourceName, a.awsConfig, a.logger)
	case SERVICE_ELB:
		serviceAnalyzer, err = newELBAnalyzer(ctx, req.ResourceName, a.awsConfig, a.logger)
	default:
		return nil, fmt.Errorf("not supported service: %s", a.initialService)
	}
	if err != nil {
		return nil, err
	}
	return serviceAnalyzer, nil
}

func isSupportedAWSService(serviceName string) bool {
	return supportedAWSServices[serviceName]
}

func findAWSServiceFromDomain(domain string) string {
	switch {
	case domainPatternCloudFront.MatchString(domain):
		return SERVICE_CLOUDFRONT
	case domainPatternS3.MatchString(domain):
		return SERVICE_S3
	case domainPatternLambdaURL.MatchString(domain):
		return SERVICE_LAMBDA
	default:
		return ""
	}
}

func getLayerFromAWSService(service string) string {
	switch service {
	case SERVICE_CLOUDFRONT:
		return LAYER_CDN
	case SERVICE_API_GATEWAY:
		return LAYER_GATEWAY
	case SERVICE_LAMBDA:
		return LAYER_COMPUTE
	case SERVICE_S3, SERVICE_SQS, SERVICE_SNS, SERVICE_EVENT_BRIDGE:
		return LAYER_DATASTORE
	case SERVICE_IAM:
		return LAYER_LATERAL_MOVEMENT
	default:
		return ""
	}
}

func isPublicSecurityGroup(ctx context.Context, ec2Client *ec2.Client, groupID string) (bool, error) {
	// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeSecurityGroups.html
	groups, err := ec2Client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{groupID},
	})
	if err != nil {
		return false, err
	}
	for _, sgs := range groups.SecurityGroups {
		for _, ipPermissions := range sgs.IpPermissions {
			for _, ipRange := range ipPermissions.IpRanges {
				if aws.ToString(ipRange.CidrIp) == "0.0.0.0/0" {
					return true, nil
				}
			}
			for _, ipv6Range := range ipPermissions.Ipv6Ranges {
				if aws.ToString(ipv6Range.CidrIpv6) == "::/0" {
					return true, nil
				}
			}
		}
	}
	return false, nil
}
