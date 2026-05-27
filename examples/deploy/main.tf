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

resource "americancloud_vm" "test" {
  name               = "tf-test-vm"
  region             = "us-central-0"
  image              = "ubuntu-24.04-050826"
  vm_package         = "standard-custom"
  cpu                = 1
  memory_mb          = 1024
  root_disk_gb       = 25
  network            = "e146b008-1519-430c-b3d5-29bff77fc2ea"
  subscription_period = "hourly"
  keypairs           = ["hnmatlock@una.edu"]
}

output "vm_id" {
  value = americancloud_vm.test.id
}

output "vm_status" {
  value = americancloud_vm.test.status
}
