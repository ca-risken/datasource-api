package gcp

import (
	"context"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/pkg/attackflow"
	"github.com/ca-risken/datasource-api/pkg/gcp"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

type computeAnalyzer struct {
	resource *datasource.Resource
	metadata *computeMetadata
	client   gcp.GcpServiceClient
	logger   logging.Logger
}
type computeMetadata struct {
	ServiceAccount  string `json:"service_account"`
	State           string `json:"state"`
	Zone            string `json:"zone"`
	InstanceID      uint64 `json:"instance_id"`
	IsPublic        bool   `json:"is_public"`
	PublicDnsName   string `json:"public_dns_name"`
	PublicIpAddress string `json:"public_ip_address"`
	MachineType     string `json:"machine_type"`
}

func newComputeAnalyzer(ctx context.Context, r *datasource.Resource, client gcp.GcpServiceClient, logger logging.Logger) (
	attackflow.CloudServiceAnalyzer, error,
) {
	return &computeAnalyzer{
		resource: r,
		metadata: &computeMetadata{},
		client:   client,
		logger:   logger,
	}, nil
}

func (c *computeAnalyzer) Analyze(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, error,
) {
	compute, err := c.client.DescribeInstance(ctx, c.resource.CloudId, c.resource.Region, c.resource.ShortName)
	if err != nil {
		return nil, err
	}
	c.metadata.State = compute.Instance.Status
	c.metadata.Zone = compute.Instance.Zone
	c.metadata.InstanceID = compute.Instance.Id
	c.metadata.IsPublic = compute.IsPublic
	c.metadata.MachineType = compute.Instance.MachineType
	if compute.Instance.ServiceAccounts != nil &&
		len(compute.Instance.ServiceAccounts) > 0 {
		c.metadata.ServiceAccount = compute.Instance.ServiceAccounts[0].Email
	}
	if compute.Instance.NetworkInterfaces != nil &&
		len(compute.Instance.NetworkInterfaces) > 0 &&
		compute.Instance.NetworkInterfaces[0].AccessConfigs != nil &&
		len(compute.Instance.NetworkInterfaces[0].AccessConfigs) > 0 {
		c.metadata.PublicDnsName = compute.Instance.NetworkInterfaces[0].AccessConfigs[0].PublicPtrDomainName
		c.metadata.PublicIpAddress = compute.Instance.NetworkInterfaces[0].AccessConfigs[0].NatIP
	}

	c.resource.MetaData, err = attackflow.ParseMetadata(c.metadata)
	if err != nil {
		return nil, err
	}
	resp = attackflow.SetNode(c.metadata.IsPublic, c.metadata.PublicIpAddress, c.resource, resp)
	return resp, nil
}

func (c *computeAnalyzer) Next(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, []attackflow.CloudServiceAnalyzer, error,
) {
	return resp, nil, nil
}
