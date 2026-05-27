package client

import (
	"context"
	"net/url"
)

// Image represents an OS image available for VM provisioning.
type Image struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
	OS          string `json:"os"`
	Version     string `json:"version"`
}

// ListImages returns all available images, optionally filtered by OS and version.
func (c *Client) ListImages(ctx context.Context, os, version string) ([]Image, error) {
	query := url.Values{}
	if os != "" {
		query.Set("os", os)
	}
	if version != "" {
		query.Set("version", version)
	}

	var images []Image
	err := c.doRequestWithQuery(ctx, "GET", "/api/v1/compute/images", query, nil, &images)
	if err != nil {
		return nil, err
	}
	return images, nil
}

// GetImage returns a single image by ID.
func (c *Client) GetImage(ctx context.Context, id string) (*Image, error) {
	var image Image
	err := c.doRequest(ctx, "GET", "/api/v1/compute/images/"+id, nil, &image)
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// GetImageByLabel returns a single image by label.
func (c *Client) GetImageByLabel(ctx context.Context, label string) (*Image, error) {
	var image Image
	err := c.doRequest(ctx, "GET", "/api/v1/compute/images/label/"+label, nil, &image)
	if err != nil {
		return nil, err
	}
	return &image, nil
}
