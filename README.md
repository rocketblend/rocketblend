### [Discussions](https://github.com/rocketblend/rocketblend/discussions) │ [Documentation](https://docs.rocketblend.io) │ [Latest Release](https://github.com/rocketblend/rocketblend/releases/latest)

# RocketBlend

[![Github tag](https://badgen.net/github/tag/rocketblend/rocketblend)](https://github.com/rocketblend/rocketblend/tags)
[![Go Doc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/rocketblend/rocketblend)
[![Go Report Card](https://goreportcard.com/badge/github.com/rocketblend/rocketblend)](https://goreportcard.com/report/github.com/rocketblend/rocketblend)
[![GitHub](https://img.shields.io/github/license/rocketblend/rocketblend)](https://github.com/rocketblend/rocketblend/blob/master/LICENSE)

RocketBlend is an open-source tool that offers build and addon management for the 3D graphics software, [Blender](https://www.blender.org/).

## About

RocketBlend consists of two parts: a CLI tool and a Launcher.

### CLI

The command line tool is the primary application for managing your local installations of addons and Blender versions.

### Launcher

The Launcher is an alternative launcher for .blend files that utilizes a `rocketfile.json` file to determine the correct version of Blender and the addons to run.

**Note: The Launcher is optional, but it allows you to open .blend files as you normally would while still maintaining addon and build management.**

## Getting Started

See [Quick Start](https://docs.rocketblend.io/getting-started/quick-start) in our documentation.

## See Also

- [Official Library](https://github.com/rocketblend/official-library) - Collection of builds and addons.
- [RocketBlend Collector](https://github.com/rocketblend/rocketblend-collector) - CLI tool for generating build collections from offical blender releases.
- [RocketBlend Companion](https://github.com/rocketblend/rocketblend-companion) - Blender addon to aid with working with RocketBlend. **NOTE: WIP**

## Roadmap
- CI/CD pipeline for releases.
- Companion blender addon.
- GUI project.
- Searchable build and addon website similar to [hub.docker.com](https://hub.docker.com/) or [pkg.go.dev](pkg.go.dev).

## Acknowledgments

- Inspired by [Blender Launcher](https://github.com/DotBow/Blender-Launcher)