package attackflow

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
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
	SERVICE_CLOUDFRONT = "cloudfront"
	SERVICE_S3         = "s3"

	RETRY_MAX_ATTEMPT = 10
)

var (
	domainPatternCloudFront = regexp.MustCompile(`\.cloudfront\.net$`)
	domainPatternS3         = regexp.MustCompile(`\.s3\..*\.amazonaws\.com$`) // https://docs.aws.amazon.com/ja_jp/AmazonS3/latest/userguide/VirtualHosting.html

	supportedAWSServices = map[string]bool{
		SERVICE_CLOUDFRONT: true,
		SERVICE_S3:         true,
		// TODO support below services
		// "alb":        true,
		// "ec2":        true,
		// "lambda":        true,
		// "api-gateway":   true,
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
	splitName := strings.Split(splitArn[5], "/")
	region := splitArn[3]
	if region == "" {
		region = REGION_GLOBAL
	}
	return &datasource.Resource{
		ResourceName: arn,
		ShortName:    splitName[len(splitName)-1],
		CloudType:    splitArn[1],
		Service:      splitArn[2],
		Region:       region,
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

func (a *AWS) GetInitialServiceAnalyzer(ctx context.Context, req *datasource.AnalyzeAttackFlowRequest) (
	CloudServiceAnalyzer, error,
) {
	var err error
	var serviceAnalyzer CloudServiceAnalyzer

	// check supported initial service
	switch a.initialService {
	case SERVICE_CLOUDFRONT:
		serviceAnalyzer, err = newCloudFrontAnalyzer(req.ResourceName, a.awsConfig, a.logger)
		if err != nil {
			return nil, err
		}
	case SERVICE_S3:
		serviceAnalyzer, err = newS3Analyzer(req.ResourceName, a.awsConfig, a.logger)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("not supported service: %s", a.initialService)
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
	default:
		return ""
	}
}
