package attackflow

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

type S3Metadata struct {
	Encryption string
	IsPublic   bool
	Versioning bool
}

func (a *AWSAttackFlowAnalyzer) analyzeS3(ctx context.Context, arn string) (*datasource.AnalyzeAttackFlowResponse, error) {
	// analyze s3 resource
	r, meta, err := a.analyzeS3Resource(ctx, arn)
	if err != nil {
		return nil, err
	}

	a.logger.Infof(ctx, "s3: resource=%v, meta=%v", r, meta)
	return &datasource.AnalyzeAttackFlowResponse{
		Nodes: a.nodes,
		Edges: a.edges,
	}, nil
}

func (a *AWSAttackFlowAnalyzer) analyzeS3Resource(ctx context.Context, arn string) (*datasource.Resource, *S3Metadata, error) {
	r := getAWSInfoFromARN(arn)
	bucketName := aws.String(r.ShortName)

	// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketLocation.html
	location, err := a.s3.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
		Bucket: bucketName,
	})
	if err != nil {
		return nil, nil, err
	}
	regionCode := fmt.Sprint(location.LocationConstraint)
	a.logger.Debugf(ctx, "s3: location=%v", regionCode)
	if regionCode != "" {
		r.Region = regionCode
		if err := a.updateS3ClientWithRegion(ctx, regionCode); err != nil {
			return nil, nil, err
		}
	}

	// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketEncryption.html
	encryption, err := a.s3.GetBucketEncryption(ctx, &s3.GetBucketEncryptionInput{
		Bucket: bucketName,
	})
	if err != nil {
		return nil, nil, err
	}

	// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketPolicyStatus.html
	policyStatus, err := a.s3.GetBucketPolicyStatus(ctx, &s3.GetBucketPolicyStatusInput{
		Bucket: bucketName,
	})
	if err != nil {
		return nil, nil, err
	}

	// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketVersioning.html
	versioning, err := a.s3.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
		Bucket: bucketName,
	})
	if err != nil {
		return nil, nil, err
	}

	sseEncrypt := ""
	for _, rule := range encryption.ServerSideEncryptionConfiguration.Rules {
		sseEncrypt = fmt.Sprint(rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm)
		break
	}

	meta := &S3Metadata{
		Encryption: sseEncrypt,
		IsPublic:   policyStatus.PolicyStatus != nil && policyStatus.PolicyStatus.IsPublic,
		Versioning: versioning.Status == types.BucketVersioningStatusEnabled,
	}
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return nil, nil, err
	}
	r.Layer = LAYER_DATASTORE
	r.MetaData = string(metaJSON)

	// add node
	if meta.IsPublic {
		a.addInternetNode(r.ResourceName, "")
	}
	a.nodes = append(a.nodes, r)
	return r, meta, nil
}

func (a *AWSAttackFlowAnalyzer) updateS3ClientWithRegion(ctx context.Context, region string) error {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithRetryMaxAttempts(RETRY_MAX_ATTEMPT),
	)
	if err != nil {
		return err
	}
	cfg.Credentials = a.awsConfig.Credentials
	a.s3 = s3.NewFromConfig(cfg) // update s3 client
	return nil
}

func getS3ARNFromDomain(domain string) string {
	if !domainPatternS3.MatchString(domain) {
		return ""
	}
	// bucket-name.com.s3.{region}.amazonaws.com -> bucket-name.com
	bucketName := domainPatternS3.ReplaceAll([]byte(domain), []byte{})
	return fmt.Sprintf("arn:aws:s3:::%s", bucketName)
}
