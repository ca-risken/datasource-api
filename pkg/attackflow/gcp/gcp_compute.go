package gcp

import (
	"context"
	"encoding/json"

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
	State           string   `json:"state"`
	Zone            string   `json:"zone"`
	InstanceID      uint64   `json:"instance_id"`
	MachineType     string   `json:"machine_type"`
	IsPublic        bool     `json:"is_public"`
	VPCNetwork      []string `json:"vpc_network"`
	ServiceAccount  []string `json:"service_account"`
	PublicDnsName   []string `json:"public_dns_name"`
	PublicIpAddress []string `json:"public_ip_address"`
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
	// cache
	cachedResource, cachedMeta, err := getGcpComputeCache(c.resource.CloudId, c.resource.ResourceName)
	if err != nil {
		return nil, err
	}
	if cachedResource != nil {
		c.resource = cachedResource
		resp = attackflow.SetNode(cachedMeta.IsPublic, "", cachedResource, resp)
		return resp, nil
	}

	compute, err := c.client.DescribeInstance(ctx, c.resource.CloudId, c.resource.Region, c.resource.ShortName)
	if err != nil {
		return nil, err
	}
	c.metadata.State = compute.Instance.Status
	c.metadata.Zone = getShortNameFromURL(compute.Instance.Zone) // https://www.googleapis.com/compute/v1/projects/project-id/zones/northamerica-northeast1-a -> northamerica-northeast1-a
	c.metadata.InstanceID = compute.Instance.Id
	c.metadata.IsPublic = compute.IsPublic
	c.metadata.MachineType = getShortNameFromURL(compute.Instance.MachineType) // https://www.googleapis.com/compute/v1/projects/project-id/zones/zone-name/machineTypes/e2-micro -> e2-micro
	if compute.Instance.NetworkInterfaces != nil {
		for _, ni := range compute.Instance.NetworkInterfaces {
			c.metadata.VPCNetwork = append(c.metadata.VPCNetwork, getShortNameFromURL(ni.Network)) // https://www.googleapis.com/compute/v1/projects/project-id/global/networks/vpc-name -> vpc-name
		}
	}
	if compute.Instance.ServiceAccounts != nil {
		for _, sa := range compute.Instance.ServiceAccounts {
			c.metadata.ServiceAccount = append(c.metadata.ServiceAccount, sa.Email)
		}
	}
	if compute.Instance.NetworkInterfaces != nil {
		for _, ni := range compute.Instance.NetworkInterfaces {
			if ni.AccessConfigs != nil {
				for _, ac := range ni.AccessConfigs {
					c.metadata.PublicDnsName = append(c.metadata.PublicDnsName, ac.PublicPtrDomainName)
					c.metadata.PublicIpAddress = append(c.metadata.PublicIpAddress, ac.NatIP)
				}
			}
		}
	}

	c.resource.MetaData, err = attackflow.ParseMetadata(c.metadata)
	if err != nil {
		return nil, err
	}
	resp = attackflow.SetNode(c.metadata.IsPublic, "", c.resource, resp)

	// cache
	if err := attackflow.SetAttackFlowCache(c.resource.CloudId, c.resource.ResourceName, c.resource); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *computeAnalyzer) Next(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, []attackflow.CloudServiceAnalyzer, error,
) {
	return resp, nil, nil
}

func getGcpComputeCache(cloudID, resourceName string) (*datasource.Resource, *computeMetadata, error) {
	resource, err := attackflow.GetAttackFlowCache(cloudID, resourceName)
	if err != nil {
		return nil, nil, err
	}
	if resource == nil {
		return nil, nil, nil
	}
	var meta computeMetadata
	if err := json.Unmarshal([]byte(resource.MetaData), &meta); err != nil {
		return nil, nil, err
	}
	return resource, &meta, nil
}
