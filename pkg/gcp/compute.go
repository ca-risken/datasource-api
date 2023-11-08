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
		if !c.onSameVPCNetwork(fw) {
			continue
		}
		if fw.TargetTags == nil || len(fw.TargetTags) == 0 {
			// No target tags means it applies to all instances in the network
			c.addFirewall(fw)
			continue
		}
		for _, targetTag := range instance.Tags.Items {
			for _, fwTargetTag := range fw.TargetTags {
				if targetTag == fwTargetTag {
					// Match the specific network tag
					c.addFirewall(fw)
				}
			}
		}
	}
	return c, nil
}

func (c *Compute) onSameVPCNetwork(fw *compute.Firewall) bool {
	for _, networkInterface := range c.Instance.NetworkInterfaces {
		if networkInterface.Network == fw.Network {
			return true
		}
	}
	return false
}

func (c *Compute) addFirewall(fw *compute.Firewall) {
	c.Firewalls = append(c.Firewalls, fw)
	if !c.IsPublic {
		c.IsPublic = c.isPublicVM(fw)
	}
}

func (c *Compute) isPublicVM(fw *compute.Firewall) bool {
	if c.Instance == nil {
		return false // no instance
	}
	if c.Instance.Status != "RUNNING" {
		return false // not running
	}
	if !c.hasPublicIP() {
		return false // no public IP
	}

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

func (c *Compute) hasPublicIP() bool {
	if c.Instance.NetworkInterfaces == nil || len(c.Instance.NetworkInterfaces) == 0 {
		return false // no network interfaces
	}
	for _, ni := range c.Instance.NetworkInterfaces {
		if ni.AccessConfigs == nil || len(ni.AccessConfigs) == 0 {
			continue
		}
		for _, ac := range ni.AccessConfigs {
			if ac.NatIP != "" {
				return true
			}
		}
	}
	return false
}
