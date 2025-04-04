---
description: A quick guide to getting started with RocketBlend.
---

# Quick Start

## Install RocketBlend <a href="#install-helm" id="install-helm"></a>

Before getting started with RocketBlend, ensure it's installed on your system. You can use tools like `go install`, package managers like Homebrew and Scoop or grab a pre-compiled binary at [the official releases page](https://github.com/rocketblend/rocketblend/releases).

For more details, or for other options, see [the installation guide](installation.md).

## Create a new project <a href="#initialize-a-helm-chart-repository" id="initialize-a-helm-chart-repository"></a>

Once you have everything ready, you can create your first project! Open a terminal in your desired location and run the command `rocketblend new hello-world` .

```shell-session
$ rocketblend new hello-world
⣻ Creating project...
```

{% hint style="info" %}
To see a more detailed output of what's happening, add the `-v` flag.
{% endhint %}

Once complete you'll have a new `.blend` project and a `.rocketblend` folder used to store all dependency data for your project.

```shell-session
$ ls
hello-world.blend  .rocketblend
```

## Install a package <a href="#install-an-example-chart" id="install-an-example-chart"></a>

To install a package, you can run the `rocketblend install` command. RocketBlend uses references to install packages from remote sources such as Github. You can create your own packages or use ones already defined online, but the easiest way to get started is to use the [official library](https://github.com/rocketblend/official-library).

```shell-session
$ rocketblend install blender/4.1.1
⣯ Installing package...
```

In this example, we're changing the project to use a Blender 4.1.1 build.

{% hint style="info" %}
If you have purchased a build or add-on, you can create a custom package definition locally. Then use the `rocketblend insert` command to insert it into your library.
{% endhint %}

Running this command without specifying a package will install all dependencies for a project. This is particularly helpful when sharing projects across multiple machines.

## Open a project <a href="#uninstall-a-release" id="uninstall-a-release"></a>

With the project successfully configured, you can initiate it by executing the `rocketblend run` command. Blender will then launch with all the necessary dependencies defined for the project.

By default, RocketBlend preserves your current add-ons and simply adds the new ones you specify, ensuring that projects not using the tool remain unaffected.

If you prefer a more controlled setup, you can enable strict mode by setting it to true in the `.rocketblend/profile.json` file. With strict mode enabled, only the add-ons listed in that file will be active, and all others—including default add-ons like Cycles—will be disabled.

{% hint style="warning" %}
Any changes to add-ons are temporary and only for that Blender session. They won't be saved to your user preferences, therefore the add-on menu in Blender might show incorrectly. This is done to retain any previously defined add-on preferences.
{% endhint %}
