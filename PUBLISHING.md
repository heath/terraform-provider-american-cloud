# Publishing to the Terraform Registry

This guide walks through publishing the American Cloud Terraform provider to the public Terraform Registry.

## Prerequisites

- GitHub account with admin access to this repo
- [GoReleaser](https://goreleaser.com/install/) (only needed for local releases)
- GPG installed (`brew install gnupg` on macOS)

## Step 1: Rename the GitHub repo

The Terraform Registry requires repos to follow the naming convention `terraform-provider-{NAME}`.

1. Go to this repository's **Settings** in Github and select the **General** tab.
2. Rename the repository to `terraform-provider-americancloud`

This sets your registry namespace to `{YOUR_ORG}/americancloud`.

Update your local remote after renaming:

```bash
git remote set-url origin git@github.com:{YOUR_ORG}/terraform-provider-americancloud.git
```

## Step 2: Generate a GPG signing key

The Terraform Registry requires RSA or DSA keys (not ECC, which is GPG's default).

```bash
gpg --full-generate-key
```

When prompted:

- **Key type:** RSA and RSA
- **Key size:** 4096
- **Expiration:** your choice (0 for no expiry)
- **Name/email:** your identity

### Export the public key

You'll paste this into the Terraform Registry:

```bash
gpg --armor --export "your-email@example.com"
```

### Export the private key

You'll add this as a GitHub Actions secret:

```bash
gpg --armor --export-secret-keys "your-email@example.com"
```

## Step 3: Add the GPG public key to Terraform Registry

1. Sign in at [registry.terraform.io](https://registry.terraform.io) with your GitHub account
2. Go to **User Settings > Signing Keys**
3. Paste the ASCII-armored **public** key from Step 2

## Step 4: Add GitHub Actions secrets

In your repo, go to **Settings > Secrets and variables > Actions** and add two secrets:

| Secret | Value |
|--------|-------|
| `GPG_PRIVATE_KEY` | Output of `gpg --armor --export-secret-keys` |
| `PASSPHRASE` | The passphrase you set for the GPG key |

## Step 5: Register the provider on the Terraform Registry

1. Go to [registry.terraform.io/publish/provider](https://registry.terraform.io/publish/provider)
2. Select your repository (`terraform-provider-americancloud`)
3. The registry automatically creates a webhook on your repo subscribed to release events

If the webhook is missing later, go to the provider's settings page on the registry and click **Resync**.

## Step 6: Create your first release

Commit the release automation files (already in the repo):

```bash
git add .goreleaser.yml terraform-registry-manifest.json .github/
git commit -m "Add release automation for Terraform registry"
git push origin main
```

Tag and push to trigger the release:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The GitHub Actions workflow will:

1. Check out the code with full tag history
2. Set up Go from `go.mod`
3. Import the GPG key from secrets
4. Run GoReleaser, which builds binaries for all platforms, signs the checksums, and creates a GitHub release

The registry webhook picks up the release and publishes it automatically.

## After publishing

Users can install the provider with:

```hcl
terraform {
  required_providers {
    americancloud = {
      source  = "{YOUR_ORG}/americancloud"
      version = "~> 0.1"
    }
  }
}
```

Once the registry version is live, you can remove the local `dev_overrides` from `~/.terraformrc`.

## Release files reference

These files were added to the repo as part of this setup:

| File | Purpose |
|------|---------|
| `.goreleaser.yml` | Builds binaries for linux/darwin/windows/freebsd across amd64/arm64/arm/386, signs checksums with GPG, includes the registry manifest |
| `terraform-registry-manifest.json` | Declares protocol version `6.0` (plugin framework) |
| `.github/workflows/release.yml` | GitHub Actions workflow triggered by `v*` tags |

## Subsequent releases

For future versions, just tag and push:

```bash
git tag v0.2.0
git push origin v0.2.0
```

Tags must follow [Semantic Versioning](https://semver.org/) (e.g., `v1.2.3`). Prerelease versions like `v1.0.0-beta.1` are supported. Never modify or replace an already-released version.

## Releasing locally (without GitHub Actions)

If you prefer to release from your machine instead of CI:

```bash
# Cache GPG passphrase
gpg --armor --detach-sign

# Set GitHub token
export GITHUB_TOKEN="your-personal-access-token"  # needs public_repo scope

# Tag
git tag v0.1.0

# Build, sign, upload
goreleaser release --clean
```

## Troubleshooting

- **"GPG key type not supported"** — The registry only accepts RSA and DSA keys, not ECC. Regenerate with `--full-generate-key` and select RSA.
- **Webhook not firing** — Go to the provider settings on the registry and click Resync. Remove any stale webhooks for `registry.terraform.io` from your GitHub repo settings.
- **GoReleaser passphrase error** — GoReleaser doesn't support interactive passphrase prompts. Cache it first with `gpg --armor --detach-sign` or use a key without a passphrase.

## References

- [Publish providers to the Terraform registry](https://developer.hashicorp.com/terraform/registry/providers/publishing)
- [Release and publish tutorial](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-release-publish)
- [terraform-provider-scaffolding-framework](https://github.com/hashicorp/terraform-provider-scaffolding-framework)
- [GoReleaser documentation](https://goreleaser.com)
