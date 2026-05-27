package client

import "context"

// PublicIP represents a public IP address.
type PublicIP struct {
	ID                    string `json:"id"`
	IPAddress             string `json:"ip_address"`
	NetworkID             string `json:"network_id"`
	VMID                  string `json:"vm_id"`
	VMName                string `json:"vm_name"`
	Status                string `json:"status"`
	Region                string `json:"region"`
	IsSourceNat           bool   `json:"is_source_nat"`
	IsStaticNat           bool   `json:"is_static_nat"`
	HasFirewallRules      bool   `json:"has_firewall_rules"`
	VPCID                 string `json:"vpc_id"`
	VPCName               string `json:"vpc_name"`
	AssociatedNetworkID   string `json:"associated_network_id"`
	AssociatedNetworkName string `json:"associated_network_name"`
}

// ReservePublicIPRequest is the request body for reserving a public IP.
type ReservePublicIPRequest struct {
	NetworkID string `json:"networkId,omitempty"`
	VPCID     string `json:"vpcId,omitempty"`
	Region    string `json:"region,omitempty"`
}

// EnableStaticNatRequest is the request body for enabling static NAT.
type EnableStaticNatRequest struct {
	VirtualMachineID string `json:"virtualMachineId"`
}

// ListPublicIPs returns all public IPs.
func (c *Client) ListPublicIPs(ctx context.Context) ([]PublicIP, error) {
	var ips []PublicIP
	err := c.doRequest(ctx, "GET", "/api/v1/networks/public-ip", nil, &ips)
	return ips, err
}

// GetPublicIP returns a single public IP by ID.
func (c *Client) GetPublicIP(ctx context.Context, id string) (*PublicIP, error) {
	var ip PublicIP
	err := c.doRequest(ctx, "GET", "/api/v1/networks/public-ip/"+id, nil, &ip)
	if err != nil {
		return nil, err
	}
	return &ip, nil
}

// ReservePublicIP reserves a new public IP.
func (c *Client) ReservePublicIP(ctx context.Context, req *ReservePublicIPRequest) (*PublicIP, error) {
	var ip PublicIP
	err := c.doRequest(ctx, "POST", "/api/v1/networks/public-ip", req, &ip)
	if err != nil {
		return nil, err
	}
	return &ip, nil
}

// ReleasePublicIP releases a public IP.
func (c *Client) ReleasePublicIP(ctx context.Context, id string) error {
	return c.doRequest(ctx, "DELETE", "/api/v1/networks/public-ip/"+id, nil, nil)
}

// EnableStaticNat enables static NAT from a public IP to a VM.
func (c *Client) EnableStaticNat(ctx context.Context, ipID, vmID string) error {
	req := &EnableStaticNatRequest{VirtualMachineID: vmID}
	return c.doRequest(ctx, "PUT", "/api/v1/networks/public-ip/"+ipID+"/enable-static-nat", req, nil)
}

// DisableStaticNat disables static NAT on a public IP.
func (c *Client) DisableStaticNat(ctx context.Context, ipID string) error {
	return c.doRequest(ctx, "PUT", "/api/v1/networks/public-ip/"+ipID+"/disable-static-nat", nil, nil)
}
