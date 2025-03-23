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
  -y, --auto-confirm    overwrite any existing files without requiring confirmation
  -c, --continue        continue rendering from the last rendered frame in the output directory
  -e, --end int         frame to end rendering at, 0 for single frame
  -g, --engine string   override render engine (cycles, eevee, workbench)
  -f, --format string   output format for the rendered frames (default "PNG")
  -h, --help            help for render
  -j, --jump int        number of frames to step forward after each rendered frame (default 1)
  -o, --output string   output path for the rendered frames (default "//output/{{.Revision}}/{{.Name}}-#####")
  -r, --revision int    revision number for the output directory, 0 for auto-increment
  -s, --start int       frame to start rendering from (default 1)
```

### Options inherited from parent commands

```shell-session
  -d, --directory string   working directory for the command (default ".")
  -l, --log-level string   log level (debug, info, warn, error) (default "info")
  -v, --verbose            enable verbose logging
```
