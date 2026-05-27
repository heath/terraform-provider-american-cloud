package client

import "context"

// IsolatedNetwork represents an isolated network.
type IsolatedNetwork struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Region      string `json:"region"`
	Netmask     string `json:"netmask"`
	Gateway     string `json:"gateway"`
	Status      string `json:"status"`
	Created     string `json:"created"`
}

// DetailedIsolatedNetwork includes IP address info.
type DetailedIsolatedNetwork struct {
	IsolatedNetwork
	IPAddresses []PublicIP `json:"ipAddresses"`
}

// CreateIsolatedNetworkRequest is the request body for creating an isolated network.
type CreateIsolatedNetworkRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Region      string `json:"region"`
	Netmask     string `json:"netmask,omitempty"`
	Gateway     string `json:"gateway,omitempty"`
}

// UpdateIsolatedNetworkRequest is the request body for updating an isolated network.
type UpdateIsolatedNetworkRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// ListIsolatedNetworks returns all isolated networks.
func (c *Client) ListIsolatedNetworks(ctx context.Context) ([]IsolatedNetwork, error) {
	var networks []IsolatedNetwork
	err := c.doRequest(ctx, "GET", "/api/v1/networks/isolated", nil, &networks)
	return networks, err
}

// GetIsolatedNetwork returns a single isolated network by ID.
func (c *Client) GetIsolatedNetwork(ctx context.Context, id string) (*DetailedIsolatedNetwork, error) {
	var network DetailedIsolatedNetwork
	err := c.doRequest(ctx, "GET", "/api/v1/networks/isolated/"+id, nil, &network)
	if err != nil {
		return nil, err
	}
	return &network, nil
}

// CreateIsolatedNetwork creates a new isolated network.
func (c *Client) CreateIsolatedNetwork(ctx context.Context, req *CreateIsolatedNetworkRequest) (*IsolatedNetwork, error) {
	var network IsolatedNetwork
	err := c.doRequest(ctx, "POST", "/api/v1/networks/isolated", req, &network)
	if err != nil {
		return nil, err
	}
	return &network, nil
}

// UpdateIsolatedNetwork updates an existing isolated network.
func (c *Client) UpdateIsolatedNetwork(ctx context.Context, id string, req *UpdateIsolatedNetworkRequest) (*IsolatedNetwork, error) {
	var network IsolatedNetwork
	err := c.doRequest(ctx, "PUT", "/api/v1/networks/isolated/"+id, req, &network)
	if err != nil {
		return nil, err
	}
	return &network, nil
}

// DeleteIsolatedNetwork deletes an isolated network.
func (c *Client) DeleteIsolatedNetwork(ctx context.Context, id string) error {
	return c.doRequest(ctx, "DELETE", "/api/v1/networks/isolated/"+id, nil, nil)
}

// VPC represents a Virtual Private Cloud.
type VPC struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CIDR        string `json:"cidr"`
	Status      string `json:"status"`
	Region      string `json:"region"`
	Created     string `json:"created"`
}

// DetailedVPC includes tiers and public IPs.
type DetailedVPC struct {
	VPC
	Tiers     []VPCTier  `json:"tiers"`
	PublicIPs []PublicIP `json:"public_ips"`
}

// CreateVPCRequest is the request body for creating a VPC.
type CreateVPCRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Region      string `json:"region"`
	CIDR        string `json:"cidr"`
}

// UpdateVPCRequest is the request body for updating a VPC.
type UpdateVPCRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// ListVPCs returns all VPCs.
func (c *Client) ListVPCs(ctx context.Context) ([]VPC, error) {
	var vpcs []VPC
	err := c.doRequest(ctx, "GET", "/api/v1/networks/vpc", nil, &vpcs)
	return vpcs, err
}

// GetVPC returns a single VPC by ID.
func (c *Client) GetVPC(ctx context.Context, id string) (*DetailedVPC, error) {
	var vpc DetailedVPC
	err := c.doRequest(ctx, "GET", "/api/v1/networks/vpc/"+id, nil, &vpc)
	if err != nil {
		return nil, err
	}
	return &vpc, nil
}

// CreateVPC creates a new VPC.
func (c *Client) CreateVPC(ctx context.Context, req *CreateVPCRequest) (*VPC, error) {
	var vpc VPC
	err := c.doRequest(ctx, "POST", "/api/v1/networks/vpc", req, &vpc)
	if err != nil {
		return nil, err
	}
	return &vpc, nil
}

// UpdateVPC updates a VPC.
func (c *Client) UpdateVPC(ctx context.Context, id string, req *UpdateVPCRequest) (*VPC, error) {
	var vpc VPC
	err := c.doRequest(ctx, "PUT", "/api/v1/networks/vpc/"+id, req, &vpc)
	if err != nil {
		return nil, err
	}
	return &vpc, nil
}

// DeleteVPC deletes a VPC.
func (c *Client) DeleteVPC(ctx context.Context, id string) error {
	return c.doRequest(ctx, "DELETE", "/api/v1/networks/vpc/"+id, nil, nil)
}

// VPCTier represents a tier (subnet) within a VPC.
type VPCTier struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	VPCID       string `json:"vpc_id"`
	Gateway     string `json:"gateway"`
	Netmask     string `json:"netmask"`
	ACLID       string `json:"acl_id"`
	Status      string `json:"status"`
	Created     string `json:"created"`
}

// CreateVPCTierRequest is the request body for creating a VPC tier.
type CreateVPCTierRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	VPCID       string `json:"vpc_id"`
	Gateway     string `json:"gateway"`
	Netmask     string `json:"netmask"`
	ACLID       string `json:"acl_id,omitempty"`
}

// CreateVPCTier creates a new tier within a VPC.
func (c *Client) CreateVPCTier(ctx context.Context, req *CreateVPCTierRequest) (*VPCTier, error) {
	var tier VPCTier
	err := c.doRequest(ctx, "POST", "/api/v1/networks/vpc/tier", req, &tier)
	if err != nil {
		return nil, err
	}
	return &tier, nil
}
