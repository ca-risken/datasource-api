package attackflow

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

func (e *elbAnalyzer) analyzeV2(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, error,
) {
	e.resource.ShortName = getV2LoadBalancerName(e.resource.ResourceName) // update

	// https://docs.aws.amazon.com/elasticloadbalancing/2012-06-01/APIReference/API_DescribeLoadBalancers.html
	lb, err := e.v2client.DescribeLoadBalancers(ctx, &elasticloadbalancingv2.DescribeLoadBalancersInput{
		Names: []string{e.resource.ShortName},
	})
	if err != nil {
		return nil, err
	}
	for _, l := range lb.LoadBalancers {
		e.metadata.Name = aws.ToString(l.LoadBalancerName)
		e.metadata.DNSName = aws.ToString(l.DNSName)
		e.metadata.InternetFacing = l.Scheme == types.LoadBalancerSchemeEnumInternetFacing
		e.metadata.SecurityGroups = l.SecurityGroups
		e.metadata.VpcID = aws.ToString(l.VpcId)
	}

	// https://docs.aws.amazon.com/elasticloadbalancing/2012-06-01/APIReference/API_DescribeLoadBalancerAttributes.html
	attrs, err := e.v2client.DescribeLoadBalancerAttributes(ctx, &elasticloadbalancingv2.DescribeLoadBalancerAttributesInput{
		LoadBalancerArn: aws.String(e.resource.ResourceName),
	})
	if err != nil {
		return nil, err
	}
	for _, a := range attrs.Attributes {
		if aws.ToString(a.Key) == "access_logs.s3.enabled" {
			e.metadata.AccessLogging = aws.ToString(a.Value) == "true"
		}
	}

	// https://docs.aws.amazon.com/elasticloadbalancing/latest/APIReference/API_DescribeTargetGroups.html
	tgs, err := e.v2client.DescribeTargetGroups(ctx, &elasticloadbalancingv2.DescribeTargetGroupsInput{
		LoadBalancerArn: aws.String(e.resource.ResourceName),
	})
	if err != nil {
		return nil, err
	}
	for _, tg := range tgs.TargetGroups {
		// https://docs.aws.amazon.com/elasticloadbalancing/latest/APIReference/API_DescribeTargetHealth.html
		health, err := e.v2client.DescribeTargetHealth(ctx, &elasticloadbalancingv2.DescribeTargetHealthInput{
			TargetGroupArn: tg.TargetGroupArn,
		})
		if err != nil {
			return nil, err
		}
		for _, t := range health.TargetHealthDescriptions {
			if t.Target == nil {
				continue
			}
			e.metadata.TargetGroups = append(e.metadata.TargetGroups, targetGroup{
				TargetType: fmt.Sprint(tg.TargetType),
				TargetID:   aws.ToString(t.Target.Id),
			})
		}
	}

	// public
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

func searchElbDomainV2(ctx context.Context, domain string, cfg *aws.Config) (string, error) {
	client := elasticloadbalancingv2.NewFromConfig(*cfg)
	// https://docs.aws.amazon.com/elasticloadbalancing/2012-06-01/APIReference/API_DescribeLoadBalancers.html
	lb, err := client.DescribeLoadBalancers(ctx, &elasticloadbalancingv2.DescribeLoadBalancersInput{})
	if err != nil {
		return "", err
	}
	for _, l := range lb.LoadBalancers {
		if aws.ToString(l.DNSName) == domain {
			return aws.ToString(l.LoadBalancerArn), nil
		}
	}
	return "", nil
}
