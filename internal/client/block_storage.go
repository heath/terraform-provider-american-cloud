package client

import "context"

// Volume represents a block storage volume.
type Volume struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	SizeGB  int    `json:"size_gb"`
	Status  string `json:"status"`
	VMID    string `json:"vm_id"`
	Region  string `json:"region"`
	Created string `json:"created"`
}

// CreateVolumeRequest is the request body for creating a volume.
type CreateVolumeRequest struct {
	Name   string `json:"name"`
	SizeGB int    `json:"size_gb"`
	Region string `json:"region"`
}

// ResizeVolumeRequest is the request body for resizing a volume.
type ResizeVolumeRequest struct {
	SizeGB int `json:"size_gb"`
}

// AttachVolumeRequest is the request body for attaching a volume to a VM.
type AttachVolumeRequest struct {
	VMID string `json:"vm_id"`
}

// ListVolumes returns all block storage volumes.
func (c *Client) ListVolumes(ctx context.Context) ([]Volume, error) {
	var volumes []Volume
	err := c.doRequest(ctx, "GET", "/api/v1/storage/block-storage", nil, &volumes)
	return volumes, err
}

// GetVolume returns a single volume by ID.
func (c *Client) GetVolume(ctx context.Context, id string) (*Volume, error) {
	var volume Volume
	err := c.doRequest(ctx, "GET", "/api/v1/storage/block-storage/"+id, nil, &volume)
	if err != nil {
		return nil, err
	}
	return &volume, nil
}

// CreateVolume creates a new block storage volume.
func (c *Client) CreateVolume(ctx context.Context, req *CreateVolumeRequest) (*Volume, error) {
	var volume Volume
	err := c.doRequest(ctx, "POST", "/api/v1/storage/block-storage", req, &volume)
	if err != nil {
		return nil, err
	}
	return &volume, nil
}

// DeleteVolume deletes a block storage volume.
func (c *Client) DeleteVolume(ctx context.Context, id string) error {
	return c.doRequest(ctx, "DELETE", "/api/v1/storage/block-storage/"+id, nil, nil)
}

// ResizeVolume resizes a block storage volume (grow only).
func (c *Client) ResizeVolume(ctx context.Context, id string, sizeGB int) error {
	req := &ResizeVolumeRequest{SizeGB: sizeGB}
	return c.doRequest(ctx, "PUT", "/api/v1/storage/block-storage/"+id+"/resize", req, nil)
}

// AttachVolume attaches a volume to a VM.
func (c *Client) AttachVolume(ctx context.Context, id, vmID string) error {
	req := &AttachVolumeRequest{VMID: vmID}
	return c.doRequest(ctx, "POST", "/api/v1/storage/block-storage/"+id+"/attach", req, nil)
}

// DetachVolume detaches a volume from its VM.
func (c *Client) DetachVolume(ctx context.Context, id string) error {
	return c.doRequest(ctx, "POST", "/api/v1/storage/block-storage/"+id+"/detach", nil, nil)
}
