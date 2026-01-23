---
description: Customising RocketBlend
---

# Configuration

### Default Config

{% code title="config.json" fullWidth="false" %}
```json
{
  "platform": "<windows/macos/linux>",
  "defaultbuild": "github.com/rocketblend/official-library/packages/v0/builds/blender/4.2.2",
  "installationspath": "./installations",
  "packagespath": "./packages",
  "loglevel": "info",
  "aliases": {
    "github.com/rocketblend/official-library/packages/v0/addons": "addons",
    "github.com/rocketblend/official-library/packages/v0/builds": "builds",
    "github.com/rocketblend/official-library/packages/v0/builds/blender": "blender"
  },
  "trusteddomains": [
    "download.blender.org"
  ],
  "untrusteddomainpolicy": "warn"
}
```
{% endcode %}

### Download Security

RocketBlend includes download security to help protect against downloading files from untrusted domains.

**Trusted Domains** (`trusteddomains`): List of domains from which files can be safely downloaded. Default: `download.blender.org`

**Untrusted Domain Policy** (`untrusteddomainpolicy`): How to handle downloads from domains not in the trusted list:
- `block` - Prevents downloads (most secure)
- `warn` - Logs warning but allows download (default)
- `allow` - Permits downloads without warnings

#### Examples

```bash
# Set trusted domains (replaces entire list)
rocketblend config trusteddomains --set '["download.blender.org", "github.com"]'

# Block downloads from untrusted domains
rocketblend config untrusteddomainpolicy --set "block"

# Allow all downloads
rocketblend config untrusteddomainpolicy --set "allow"
```
