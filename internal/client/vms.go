package client

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// VM represents a virtual machine.
type VM struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Status         string   `json:"status"`
	Region         string   `json:"region"`
	Image          string   `json:"image"`
	VMPackage      string   `json:"vm_package"`
	CPU            int      `json:"cpu"`
	MemoryMB       int      `json:"memory_mb"`
	RootDiskGB     int      `json:"root_disk_gb"`
	IPAddress      string   `json:"ip_address"`
	NetworkID      string   `json:"network_id"`
	Tags           []string `json:"tags"`
	Created        string   `json:"created"`
	Password       string   `json:"password,omitempty"`
	Hostname       string   `json:"hostname"`
	KeypairNames   []string `json:"keypairs"`
	NetworkAccess  string   `json:"network_access"`
}

// VMSpecs specifies custom CPU/memory/disk for VM creation.
type VMSpecs struct {
	VCPU       int `json:"vcpu"`
	MemoryMB   int `json:"memory_mb"`
	RootDiskGB int `json:"root_disk_gb"`
}

// CreateVMRequest is the request body for creating a VM.
type CreateVMRequest struct {
	Name              string   `json:"name"`
	Region            string   `json:"region"`
	VMPackage         string   `json:"vm_package,omitempty"`
	VMSpecs           *VMSpecs `json:"vm_specs,omitempty"`
	Image             string   `json:"image"`
	Network           string   `json:"network,omitempty"`
	SubscriptionPeriod string  `json:"subscription_period,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	Keypairs          []string `json:"keypairs,omitempty"`
	Userdata          string   `json:"userdata,omitempty"`
	NetworkAccess     string   `json:"network_access,omitempty"`
}

// ListVMs returns all virtual machines.
func (c *Client) ListVMs(ctx context.Context) ([]VM, error) {
	var vms []VM
	err := c.doRequest(ctx, "GET", "/api/v1/compute/vms", nil, &vms)
	if err != nil {
		return nil, err
	}
	return vms, nil
}

// GetVM returns a single VM by ID.
func (c *Client) GetVM(ctx context.Context, id string) (*VM, error) {
	var vm VM
	err := c.doRequest(ctx, "GET", "/api/v1/compute/vms/"+id, nil, &vm)
	if err != nil {
		return nil, err
	}
	return &vm, nil
}

// CreateVM creates a new virtual machine.
func (c *Client) CreateVM(ctx context.Context, req *CreateVMRequest) (*VM, error) {
	var vm VM
	err := c.doRequest(ctx, "POST", "/api/v1/compute/vms", req, &vm)
	if err != nil {
		return nil, err
	}
	return &vm, nil
}

// DeleteVM deletes a virtual machine by ID.
func (c *Client) DeleteVM(ctx context.Context, id string) error {
	return c.doRequest(ctx, "DELETE", "/api/v1/compute/vms/"+id, nil, nil)
}

// PowerVM controls the power state of a VM (start, stop, reboot).
func (c *Client) PowerVM(ctx context.Context, id, action string, force bool) error {
	query := url.Values{}
	query.Set("action", action)
	if force {
		query.Set("force", "true")
	}
	return c.doRequestWithQuery(ctx, "PUT", "/api/v1/compute/vms/"+id+"/power", query, nil, nil)
}

// ScaleVM scales a VM's CPU and/or memory.
func (c *Client) ScaleVM(ctx context.Context, id string, cpu, memoryMB int) error {
	query := url.Values{}
	if cpu > 0 {
		query.Set("cpu", fmt.Sprintf("%d", cpu))
	}
	if memoryMB > 0 {
		query.Set("memory_mb", fmt.Sprintf("%d", memoryMB))
	}
	return c.doRequestWithQuery(ctx, "PUT", "/api/v1/compute/vms/"+id+"/scale", query, nil, nil)
}

// ResizeVMDisk resizes the root disk of a VM.
func (c *Client) ResizeVMDisk(ctx context.Context, id string, sizeGB int, reboot bool) error {
	body := map[string]interface{}{
		"size_gb": sizeGB,
		"reboot":  reboot,
	}
	return c.doRequest(ctx, "PUT", "/api/v1/compute/vms/"+id+"/disk-resize", body, nil)
}

// UpdateVMHostname updates the hostname of a VM.
func (c *Client) UpdateVMHostname(ctx context.Context, id, hostname string) error {
	query := url.Values{}
	query.Set("hostname", hostname)
	return c.doRequestWithQuery(ctx, "PUT", "/api/v1/compute/vms/"+id+"/hostname", query, nil, nil)
}

// WaitForVMReady polls until the VM reaches "Running" status and has an IP
// address assigned, or the context is cancelled.
func (c *Client) WaitForVMReady(ctx context.Context, id string, pollInterval time.Duration) (*VM, error) {
	for {
		vm, err := c.GetVM(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("polling VM %s: %w", id, err)
		}

		switch vm.Status {
		case "Running", "running", "STARTED", "Started":
			// Don't return until the IP is assigned — the API may report
			// Running before networking is fully configured.
			if vm.IPAddress != "" {
				return vm, nil
			}
		case "Error", "error", "Failed", "failed", "ERROR", "FAILED":
			return vm, fmt.Errorf("VM %s reached error state: %s", id, vm.Status)
		}

		select {
		case <-ctx.Done():
			return vm, fmt.Errorf("timed out waiting for VM %s to be ready (current status: %s, ip: %s)", id, vm.Status, vm.IPAddress)
		case <-time.After(pollInterval):
			// continue polling
		}
	}
}

// WaitForVMDeleted polls until the VM returns 404, reaches a terminal
// deletion status, or the context is cancelled.
func (c *Client) WaitForVMDeleted(ctx context.Context, id string, pollInterval time.Duration) error {
	for {
		vm, err := c.GetVM(ctx, id)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return fmt.Errorf("polling VM %s for deletion: %w", id, err)
		}

		// Treat terminal states as successful deletion — some APIs keep
		// the resource visible with a final status before removing it.
		switch vm.Status {
		case "Destroyed", "destroyed", "DESTROYED",
			"Deleted", "deleted", "DELETED",
			"Expunged", "expunged", "EXPUNGED":
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for VM %s to be deleted (current status: %s)", id, vm.Status)
		case <-time.After(pollInterval):
			// continue polling
		}
	}
}
