# American Cloud Terraform Provider

A Terraform provider for managing infrastructure on [American Cloud](https://americancloud.com).

## Authentication

The provider supports two authentication methods: **JWT bearer token** or **API key** (client ID + client secret).

Set credentials via environment variables:

```bash
# API Key auth
export AC_CLIENT_ID="your-client-id"
export AC_CLIENT_SECRET="your-client-secret"

# Or bearer token auth
export AC_BEARER_TOKEN="your-token"

# Optional: override the API base URL
export AC_API_URL="https://api.americancloud.com"
```

Or configure them directly in the provider block:

```hcl
provider "americancloud" {
  client_id     = "your-client-id"
  client_secret = "your-client-secret"
}
```

## Resources

| Resource | Description |
|---|---|
| `americancloud_vm` | Virtual machines with custom CPU, memory, and disk |
| `americancloud_isolated_network` | Isolated networks |
| `americancloud_vpc` | Virtual private clouds |
| `americancloud_vpc_tier` | VPC network tiers |
| `americancloud_public_ip` | Public IP addresses with optional static NAT |
| `americancloud_firewall_rule` | Ingress firewall rules |
| `americancloud_port_forwarding_rule` | Port forwarding rules |
| `americancloud_egress_rule` | Egress firewall rules |
| `americancloud_load_balancer_rule` | Load balancer rules |
| `americancloud_block_storage` | Block storage volumes |
| `americancloud_ssh_key` | SSH key pairs |

## Data Sources

| Data Source | Description |
|---|---|
| `americancloud_images` | Available OS images |
| `americancloud_regions` | Available regions |
| `americancloud_packages` | Available VM packages |

## Example

```hcl
terraform {
  required_providers {
    americancloud = {
      source = "registry.terraform.io/americancloud/americancloud"
    }
  }
}

provider "americancloud" {}

data "americancloud_images" "ubuntu" {
  os      = "Ubuntu"
  version = "24.04"
}

resource "americancloud_vm" "web" {
  name                = "web-server"
  region              = "us-central-0"
  image               = data.americancloud_images.ubuntu.images[0].label
  vm_package          = "standard-custom"
  cpu                 = 2
  memory_mb           = 4096
  root_disk_gb        = 50
  network             = "your-network-uuid"
  subscription_period = "hourly"
  keypairs            = ["my-ssh-key"]
}
```

## Development

Build and install the provider locally:

```bash
make build
make install
```

Create a `.terraformrc` file to use the local build:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/americancloud/americancloud" = "/path/to/go/bin"
  }
  direct {}
}
```

Then run Terraform with `TF_CLI_CONFIG_FILE=.terraformrc`.

## License

Apache 2.0. See [LICENSE](LICENSE).
