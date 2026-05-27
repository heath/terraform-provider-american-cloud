terraform {
  required_providers {
    americancloud = {
      source = "registry.terraform.io/americancloud/americancloud"
    }
  }
}

# Configure via environment variables:
#   AC_BEARER_TOKEN  or  AC_CLIENT_ID + AC_CLIENT_SECRET
#   AC_API_URL (optional, defaults to https://api.americancloud.com)
provider "americancloud" {}
