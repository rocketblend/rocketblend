---
description: An improved command-line interface for Blender.
---

# Rocketblend

## Synopsis

A dependency manager for Blender projects.

Common actions for RocketBlend:

* rocketblend new: create a new project
* rocketblend describe: gives information for a given package
* rocketblend install: add a dependency to a project
* rocketblend run: runs a project
* rocketblend render: renders out a project

Documentation is available at https://docs.rocketblend.io/

### Usage

```shell-session
$ rocketblend [commands] [flags]
```

### Options

```shell-session
  -d, --directory string   working directory for the command (default ".")
  -h, --help               help for rocketblend
  -l, --log-level string   log level (debug, info, warn, error) (default "info")
  -v, --verbose            enable verbose logging
      --version            version for rocketblend
```
