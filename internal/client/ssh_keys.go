package client

import "context"

// SSHKey represents an SSH key.
type SSHKey struct {
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
}

// CreateSSHKeyRequest is the request body for creating an SSH key.
type CreateSSHKeyRequest struct {
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
}

// CreateSSHKeyResponse is the response from creating an SSH key.
type CreateSSHKeyResponse struct {
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
	PrivateKey  string `json:"private_key,omitempty"`
}

// ListSSHKeys returns all SSH keys.
func (c *Client) ListSSHKeys(ctx context.Context) ([]SSHKey, error) {
	var keys []SSHKey
	err := c.doRequest(ctx, "GET", "/api/v1/compute/ssh-keys", nil, &keys)
	return keys, err
}

// CreateSSHKey creates a new SSH key.
func (c *Client) CreateSSHKey(ctx context.Context, req *CreateSSHKeyRequest) (*CreateSSHKeyResponse, error) {
	var resp CreateSSHKeyResponse
	err := c.doRequest(ctx, "POST", "/api/v1/compute/ssh-keys", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteSSHKey deletes an SSH key by name.
func (c *Client) DeleteSSHKey(ctx context.Context, name string) error {
	return c.doRequest(ctx, "DELETE", "/api/v1/compute/ssh-keys/"+name, nil, nil)
}
