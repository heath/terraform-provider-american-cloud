package client

import "context"

// VMPackage represents a VM sizing package.
type VMPackage struct {
	ID            string   `json:"id"`
	Label         string   `json:"label"`
	Description   string   `json:"description"`
	MinCPU        int      `json:"minCpu"`
	MaxCPU        int      `json:"maxCpu"`
	MinMemoryMB   int      `json:"minMemoryMb"`
	MaxMemoryMB   int      `json:"maxMemoryMb"`
	MinRootDiskGB int      `json:"minRootDiskGb"`
	MaxRootDiskGB int      `json:"maxRootDiskGb"`
	Tier          string   `json:"tier"`
	Zones         []string `json:"zones"`
}

// ListPackages returns all available VM packages.
func (c *Client) ListPackages(ctx context.Context) ([]VMPackage, error) {
	var resp ListResponse[VMPackage]
	err := c.doRequest(ctx, "GET", "/api/v1/compute/packages", nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetPackage returns a single package by ID.
func (c *Client) GetPackage(ctx context.Context, id string) (*VMPackage, error) {
	var pkg VMPackage
	err := c.doRequest(ctx, "GET", "/api/v1/compute/packages/"+id, nil, &pkg)
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

// GetPackageByLabel returns a single package by label.
func (c *Client) GetPackageByLabel(ctx context.Context, label string) (*VMPackage, error) {
	var pkg VMPackage
	err := c.doRequest(ctx, "GET", "/api/v1/compute/packages/label/"+label, nil, &pkg)
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}
