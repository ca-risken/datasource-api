package activity

import (
	"github.com/ca-risken/common/pkg/logging"
	awsClient "github.com/ca-risken/datasource-api/proto/aws"
)

type ActivityService struct {
	awsClient        awsClient.AWSServiceClient
	cloudTrailClient cloudTrailAPI
	configClient     configServiceAPI
	logger           logging.Logger
}

func NewActivityService(a awsClient.AWSServiceClient, region string, l logging.Logger) *ActivityService {
	cfg := newConfigServiceClient(region, l)
	ct := newCloudTrailClient(region, l)
	return &ActivityService{
		awsClient:        a,
		cloudTrailClient: ct,
		configClient:     cfg,
		logger:           l,
	}
}
