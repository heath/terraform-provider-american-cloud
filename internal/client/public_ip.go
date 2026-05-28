package client

import "context"

// PublicIP represents a public IP address.
type PublicIP struct {
	ID                    string `json:"id"`
	IPAddress             string `json:"ipAddress"`
	NetworkID             string `json:"networkId"`
	VMID                  string `json:"vmId"`
	VMName                string `json:"vmName"`
	Status                string `json:"status"`
	Region                string `json:"region"`
	IsSourceNat           bool   `json:"isSourceNat"`
	IsStaticNat           bool   `json:"isStaticNat"`
	HasFirewallRules      bool   `json:"hasFirewallRules"`
	VPCID                 string `json:"vpcId"`
	VPCName               string `json:"vpcName"`
	AssociatedNetworkID   string `json:"associatedNetworkId"`
	AssociatedNetworkName string `json:"associatedNetworkName"`
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
	var resp ListResponse[PublicIP]
	err := c.doRequest(ctx, "GET", "/api/v1/networks/public-ip", nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
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
