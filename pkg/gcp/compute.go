package gcp

import (
	"context"
	"fmt"

	"google.golang.org/api/compute/v1"
)

type Compute struct {
	Instance  *compute.Instance   `json:"instance"`
	Firewalls []*compute.Firewall `json:"firewalls"`
	IsPublic  bool                `json:"is_public"`
}

func (g *GcpClient) DescribeInstance(ctx context.Context, projectID, zone, instanceName string) (*Compute, error) {
	instance, err := g.compute.Instances.Get(projectID, zone, instanceName).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}
	fwList, err := g.compute.Firewalls.List(projectID).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list firewalls: %w", err)
	}

	c := &Compute{
		Instance: instance,
		IsPublic: false,
	}
	for _, fw := range fwList.Items {
		for _, targetTag := range instance.Tags.Items {
			for _, fwTargetTag := range fw.TargetTags {
				if targetTag == fwTargetTag {
					c.Firewalls = append(c.Firewalls, fw)
					if !c.IsPublic {
						c.IsPublic = isPublicFirewallRule(fw)
					}
				}
			}
		}
	}
	return c, nil
}

func isPublicFirewallRule(fw *compute.Firewall) bool {
	if fw.Allowed == nil && len(fw.Allowed) == 0 {
		return false
	}
	for _, srcRange := range fw.SourceRanges {
		if srcRange == "0.0.0.0/0" || srcRange == "::/0" {
			return true
		}
	}
	return false
}
