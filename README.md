### [Discussions](https://github.com/rocketblend/rocketblend/discussions) │ [Documentation](https://docs.rocketblend.io) │ [Latest Release](https://github.com/rocketblend/rocketblend/releases/latest)

# RocketBlend

[![Github tag](https://badgen.net/github/tag/rocketblend/rocketblend)](https://github.com/rocketblend/rocketblend/tags)
[![Go Doc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/rocketblend/rocketblend)
[![Go Report Card](https://goreportcard.com/badge/github.com/rocketblend/rocketblend)](https://goreportcard.com/report/github.com/rocketblend/rocketblend)
[![GitHub](https://img.shields.io/github/license/rocketblend/rocketblend)](https://github.com/rocketblend/rocketblend/blob/master/LICENSE)

![Hero image of RocketBlend CLI](docs/assets/rocketblend-about.svg)

> RocketBlend is a CLI tool that streamlines the process of managing builds and addons for [Blender](https://www.blender.org/) projects.

## Getting Started

See [Quick Start](https://docs.rocketblend.io/getting-started/quick-start) in our documentation.

## Installation

For more detailed installation instructions, please refer to the [full documentation](https://docs.rocketblend.io/getting-started/installation).

### Scoop (Windows)

Windows users can install `rocketblend` using [Scoop](https://scoop.sh/).

```powershell
scoop bucket add rocketblend "https://github.com/rocketblend/scoop-bucket"
scoop install "rocketblend/rocketblend"
```

### Homebrew (macOS)

MacOS users can install `rocketblend` using [Homebrew](https://brew.sh/).

```bash
brew tap "rocketblend/homebrew-tap"
brew install rocketblend
```

### Pre-compiled binaries

To install pre-compiled binaries for `rocketblend`, you can either manually download the latest release or use the provided scripts to automate the process:

1. **Manual installation**: Download the latest release from the [releases page](https://github.com/rocketblend/rocketblend/releases) and extract the binary to a directory included in your `PATH`.
2. **Automated installation**: Run the appropriate script for your platform to download and install `rocketblend`:
   - Windows: [install.ps1](https://raw.githubusercontent.com/rocketblend/rocketblend/master/install.ps1)
   - macOS/Linux: [install.sh](https://raw.githubusercontent.com/rocketblend/rocketblend/master/install.sh)

### Source

To install `rocketblend` directly from source code, you can use the `go install` command:

```bash
go install github.com/rocketblend/rocketblend/cmd/rocketblend@latest
```

This command will download the latest version of the `rocketblend` source code and compile the binary for your platform. Ensure you have Go 1.19 or later installed on your system.

## See Also

- [RocketBlend Launcher](https://github.com/rocketblend/rocketblend-launcher) - Replacement launcher for Blender that utilises RocketBlend.
- [RocketBlend Collector](https://github.com/rocketblend/rocketblend-collector) - CLI tool for generating build collections from offical blender releases.
- [RocketBlend Companion](https://github.com/rocketblend/rocketblend-companion) - Blender addon to aid with working with RocketBlend. **NOTE: WIP**
- [Official Library](https://github.com/rocketblend/official-library) - Collection of builds and addons.

## Roadmap
- CI/CD pipeline for releases.
- Companion blender addon.
- GUI launcher project.
- Searchable build and addon website similar to [hub.docker.com](https://hub.docker.com/) or [pkg.go.dev](pkg.go.dev).

## Inspired By

- [Blender Launcher](https://github.com/DotBow/Blender-Launcher)
- [Scoop](https://scoop.sh/)