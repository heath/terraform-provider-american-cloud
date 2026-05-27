package client

import "context"

// FirewallRule represents a firewall rule.
type FirewallRule struct {
	ID             string `json:"id"`
	Protocol       string `json:"protocol"`
	StartPort      int    `json:"start_port"`
	EndPort        int    `json:"end_port"`
	IPAddressID    string `json:"ip_address_id"`
	IPAddress      string `json:"ip_address"`
	State          string `json:"state"`
	SourceCIDRList string `json:"source_cidr_list"`
	Type           string `json:"type"`
	Created        string `json:"created"`
}

// CreateFirewallRuleRequest is the request body for creating a firewall rule.
type CreateFirewallRuleRequest struct {
	Protocol       string `json:"protocol"`
	StartPort      int    `json:"start_port"`
	EndPort        int    `json:"end_port"`
	SourceCIDRList string `json:"source_cidr_list"`
	Type           string `json:"type,omitempty"`
}

// ListFirewallRules returns all firewall rules for a public IP.
func (c *Client) ListFirewallRules(ctx context.Context, ipID string) ([]FirewallRule, error) {
	var rules []FirewallRule
	err := c.doRequest(ctx, "GET", "/api/v1/networks/firewall/rules/ip/"+ipID, nil, &rules)
	return rules, err
}

// CreateFirewallRule creates a new firewall rule for a public IP.
func (c *Client) CreateFirewallRule(ctx context.Context, ipID string, req *CreateFirewallRuleRequest) (*FirewallRule, error) {
	var rule FirewallRule
	err := c.doRequest(ctx, "POST", "/api/v1/networks/firewall/rules/ip/"+ipID, req, &rule)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// DeleteFirewallRule deletes a firewall rule.
func (c *Client) DeleteFirewallRule(ctx context.Context, ruleID string) error {
	return c.doRequest(ctx, "DELETE", "/api/v1/networks/firewall/rules/"+ruleID, nil, nil)
}

// PortForwardingRule represents a port forwarding rule.
type PortForwardingRule struct {
	ID            string `json:"id"`
	PrivatePort   int    `json:"private_port"`
	PublicPort    int    `json:"public_port"`
	Protocol      string `json:"protocol"`
	VMID          string `json:"vm_id"`
	VMName        string `json:"vm_name"`
	IPAddressID   string `json:"ip_address_id"`
	IPAddress     string `json:"ip_address"`
	State         string `json:"state"`
	Created       string `json:"created"`
}

// CreatePortForwardingRuleRequest is the request body for creating a port forwarding rule.
type CreatePortForwardingRuleRequest struct {
	PrivatePort  int    `json:"private_port"`
	PublicPort   int    `json:"public_port"`
	Protocol     string `json:"protocol"`
	VMID         string `json:"vm_id"`
	OpenFirewall bool   `json:"open_firewall,omitempty"`
}

// ListPortForwardingRules returns all port forwarding rules for a public IP.
func (c *Client) ListPortForwardingRules(ctx context.Context, ipID string) ([]PortForwardingRule, error) {
	var rules []PortForwardingRule
	err := c.doRequest(ctx, "GET", "/api/v1/networks/firewall/port-forwarding/ip/"+ipID, nil, &rules)
	return rules, err
}

// CreatePortForwardingRule creates a new port forwarding rule.
func (c *Client) CreatePortForwardingRule(ctx context.Context, ipID string, req *CreatePortForwardingRuleRequest) (*PortForwardingRule, error) {
	var rule PortForwardingRule
	err := c.doRequest(ctx, "POST", "/api/v1/networks/firewall/port-forwarding/ip/"+ipID, req, &rule)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// DeletePortForwardingRule deletes a port forwarding rule.
func (c *Client) DeletePortForwardingRule(ctx context.Context, ruleID string) error {
	return c.doRequest(ctx, "DELETE", "/api/v1/networks/firewall/port-forwarding/"+ruleID, nil, nil)
}

// EgressRule represents an egress firewall rule.
type EgressRule struct {
	ID             string `json:"id"`
	Protocol       string `json:"protocol"`
	StartPort      int    `json:"start_port"`
	EndPort        int    `json:"end_port"`
	NetworkID      string `json:"network_id"`
	State          string `json:"state"`
	SourceCIDRList string `json:"source_cidr_list"`
	DestCIDRList   string `json:"dest_cidr_list"`
	Created        string `json:"created"`
}

// CreateEgressRuleRequest is the request body for creating an egress rule.
type CreateEgressRuleRequest struct {
	Protocol       string `json:"protocol"`
	StartPort      int    `json:"start_port"`
	EndPort        int    `json:"end_port"`
	SourceCIDRList string `json:"source_cidr_list,omitempty"`
	DestCIDRList   string `json:"dest_cidr_list,omitempty"`
	NetworkID      string `json:"network_id"`
}

// GetEgressRule returns a single egress rule by ID.
func (c *Client) GetEgressRule(ctx context.Context, id string) (*EgressRule, error) {
	var rule EgressRule
	err := c.doRequest(ctx, "GET", "/api/v1/networks/firewall/egress/"+id, nil, &rule)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// CreateEgressRule creates a new egress rule.
func (c *Client) CreateEgressRule(ctx context.Context, req *CreateEgressRuleRequest) (*EgressRule, error) {
	var rule EgressRule
	err := c.doRequest(ctx, "POST", "/api/v1/networks/firewall/egress", req, &rule)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// DeleteEgressRule deletes an egress rule.
func (c *Client) DeleteEgressRule(ctx context.Context, id string) error {
	return c.doRequest(ctx, "DELETE", "/api/v1/networks/firewall/egress/"+id, nil, nil)
}

// LoadBalancerRule represents a load balancer rule.
type LoadBalancerRule struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	PublicPort     int    `json:"public_port"`
	PrivatePort    int    `json:"private_port"`
	Algorithm      string `json:"algorithm"`
	Protocol       string `json:"protocol"`
	State          string `json:"state"`
	IPAddressID    string `json:"ip_address_id"`
	IPAddress      string `json:"ip_address"`
	NetworkID      string `json:"network_id"`
	SourceCIDRList string `json:"source_cidr_list"`
	Created        string `json:"created"`
}

// CreateLoadBalancerRuleRequest is the request body for creating a LB rule.
type CreateLoadBalancerRuleRequest struct {
	Name           string `json:"name"`
	Algorithm      string `json:"algorithm"`
	PublicPort     int    `json:"public_port"`
	PrivatePort    int    `json:"private_port"`
	Protocol       string `json:"protocol,omitempty"`
	SourceCIDRList string `json:"source_cidr_list,omitempty"`
	Description    string `json:"description,omitempty"`
}

// UpdateLoadBalancerRuleRequest is the request body for updating a LB rule.
type UpdateLoadBalancerRuleRequest struct {
	Name           string `json:"name,omitempty"`
	Algorithm      string `json:"algorithm,omitempty"`
	Description    string `json:"description,omitempty"`
}

// ListLoadBalancerRules returns all LB rules for a public IP.
func (c *Client) ListLoadBalancerRules(ctx context.Context, ipID string) ([]LoadBalancerRule, error) {
	var rules []LoadBalancerRule
	err := c.doRequest(ctx, "GET", "/api/v1/networks/firewall/load-balancer/ip/"+ipID, nil, &rules)
	return rules, err
}

// CreateLoadBalancerRule creates a new LB rule.
func (c *Client) CreateLoadBalancerRule(ctx context.Context, ipID string, req *CreateLoadBalancerRuleRequest) (*LoadBalancerRule, error) {
	var rule LoadBalancerRule
	err := c.doRequest(ctx, "POST", "/api/v1/networks/firewall/load-balancer/ip/"+ipID, req, &rule)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// UpdateLoadBalancerRule updates an existing LB rule.
func (c *Client) UpdateLoadBalancerRule(ctx context.Context, ruleID string, req *UpdateLoadBalancerRuleRequest) (*LoadBalancerRule, error) {
	var rule LoadBalancerRule
	err := c.doRequest(ctx, "PUT", "/api/v1/networks/firewall/load-balancer/"+ruleID, req, &rule)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// DeleteLoadBalancerRule deletes a LB rule.
func (c *Client) DeleteLoadBalancerRule(ctx context.Context, ruleID string) error {
	return c.doRequest(ctx, "DELETE", "/api/v1/networks/firewall/load-balancer/"+ruleID, nil, nil)
}
