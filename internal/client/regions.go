package client

import "context"

// Region represents a geographic region.
type Region struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	DisplayName string `json:"displayName"`
}

// ListRegions returns all available regions.
func (c *Client) ListRegions(ctx context.Context) ([]Region, error) {
	var regions []Region
	err := c.doRequest(ctx, "GET", "/api/v1/compute/regions", nil, &regions)
	if err != nil {
		return nil, err
	}
	return regions, nil
}

// GetRegion returns a single region by ID.
func (c *Client) GetRegion(ctx context.Context, id string) (*Region, error) {
	var region Region
	err := c.doRequest(ctx, "GET", "/api/v1/compute/regions/"+id, nil, &region)
	if err != nil {
		return nil, err
	}
	return &region, nil
}

// GetRegionByLabel returns a single region by label.
func (c *Client) GetRegionByLabel(ctx context.Context, label string) (*Region, error) {
	var region Region
	err := c.doRequest(ctx, "GET", "/api/v1/compute/regions/label/"+label, nil, &region)
	if err != nil {
		return nil, err
	}
	return &region, nil
}
