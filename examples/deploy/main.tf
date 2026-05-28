terraform {
  required_providers {
    americancloud = {
      source = "registry.terraform.io/americancloud/americancloud"
    }
  }
}

provider "americancloud" {
  # Credentials from AC_CLIENT_ID and AC_CLIENT_SECRET env vars
}

resource "americancloud_vm" "dokploy" {
  name                = "dokploy-server"
  region              = "us-central-0"
  image               = "ubuntu-24.04-050826"
  vm_package          = "standard-custom"
  cpu                 = 2
  memory_mb           = 4096
  root_disk_gb        = 40
  network             = "e146b008-1519-430c-b3d5-29bff77fc2ea"
  subscription_period = "hourly"
  keypairs            = ["hnmatlock@una.edu"]
}

output "vm_id" {
  value = americancloud_vm.dokploy.id
}

output "vm_status" {
  value = americancloud_vm.dokploy.status
}

output "vm_ip" {
  value = americancloud_vm.dokploy.ip_address
}
