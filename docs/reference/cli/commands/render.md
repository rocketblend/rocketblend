---
description: Rendering a project
---

# Render

## Synopsis

Render project

```shell-session
$ rocketblend render [flags]
```

## Options

```shell-session
  -f, --format string     set the render format (default "PNG")
  -e, --frame-end int     end frame
  -s, --frame-start int   start frame
  -t, --frame-step int    frame step (default 1)
  -h, --help              help for render
  -o, --output string     set the render path and file name (default "//output/{{.Project}}-#####")
```

### Options inherited from parent commands

```shell-session
  -d, --directory string   working directory for the command (default ".")
  -v, --verbose            enable verbose logging
```
