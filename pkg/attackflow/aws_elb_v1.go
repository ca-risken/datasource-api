package attackflow

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

func (e *elbAnalyzer) analyzeV1(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, error,
) {
	// DescribeLoadBalancers
	lb, err := e.v1client.DescribeLoadBalancers(ctx, &elasticloadbalancing.DescribeLoadBalancersInput{
		LoadBalancerNames: []string{e.resource.ShortName},
	})
	if err != nil {
		return nil, err
	}

	instances := []string{}
	for _, l := range lb.LoadBalancerDescriptions {
		e.metadata.Name = aws.ToString(l.LoadBalancerName)
		e.metadata.DNSName = aws.ToString(l.DNSName)
		e.metadata.InternetFacing = aws.ToString(l.Scheme) == "internet-facing"
		e.metadata.SecurityGroups = l.SecurityGroups
		e.metadata.VpcID = aws.ToString(l.VPCId)
		for _, i := range l.Instances {
			instances = append(instances, aws.ToString(i.InstanceId))
		}
	}

	// TargetGroups
	for _, i := range instances {
		e.metadata.TargetGroups = append(e.metadata.TargetGroups, targetGroup{
			TargetID:   i,
			TargetType: "instance",
		})
	}

	// Attributes
	attrs, err := e.v1client.DescribeLoadBalancerAttributes(ctx, &elasticloadbalancing.DescribeLoadBalancerAttributesInput{
		LoadBalancerName: aws.String(e.resource.ShortName),
	})
	if err != nil {
		return nil, err
	}
	if attrs.LoadBalancerAttributes != nil && attrs.LoadBalancerAttributes.AccessLog != nil {
		e.metadata.AccessLogging = attrs.LoadBalancerAttributes.AccessLog.Enabled
	}

	// Public
	if err := e.setPublic(ctx); err != nil {
		return nil, err
	}

	e.resource.MetaData, err = parseMetadata(e.metadata)
	if err != nil {
		return nil, err
	}
	resp = setNode(e.metadata.IsPublic, "", e.resource, resp)
	return resp, nil
}
