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
  }
}
```
{% endcode %}
