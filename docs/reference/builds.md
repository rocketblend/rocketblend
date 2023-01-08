# Builds

## Example

Below is an example `build.json` file for a Blender 3.4.0 build with the Rocketblend 0.1.0 package:

{% code title="build.json" %}
```json
{
   "reference": "github.com/rocketblend/official-library/builds/rocketblend/stable/3.4.0",
   "args":"",
   "source":[
      {
         "platform":"linux",
         "executable":"blender-3.4.0-linux-x64.tar/blender",
         "url":"https://download.blender.org/release/Blender3.4/blender-3.4.0-linux-x64.tar.xz"
      },
      {
         "platform":"macos/apple",
         "executable":"blender-3.4.0-macos-arm64/blender.dmg",
         "url":"https://download.blender.org/release/Blender3.4/blender-3.4.0-macos-arm64.dmg"
      },
      {
         "platform":"macos/intel",
         "executable":"blender-3.4.0-macos-x64/blender.dmg",
         "url":"https://download.blender.org/release/Blender3.4/blender-3.4.0-macos-x64.dmg"
      },
      {
         "platform":"windows",
         "executable":"blender-3.4.0-windows-x64/blender.exe",
         "url":"https://download.blender.org/release/Blender3.4/blender-3.4.0-windows-x64.zip"
      }
   ],
   "packages":[
      "github.com/rocketblend/official-library/packages/rocketblend/0.1.0"
   ]
}
```
{% endcode %}
