---
description: How to install
---

# Installation

## Homebrew (macOS)

For macOS users, you can install using [Homebrew](https://brew.sh/):

```shell-session
$ brew tap "rocketblend/homebrew-tap"
$ brew install rocketblend
```

## Scoop (Windows)

For Windows users, you can install using [Scoop](https://scoop.sh/):

```powershell
scoop bucket add rocketblend "https://github.com/rocketblend/scoop-bucket"
scoop install rocketblend/rocketblend
```

## Pre-compiled binaries

For users who prefer to install pre-compiled binaries, we provide pre-built executables for Windows, macOS, and Linux. These binaries are available on the [Releases](https://github.com/rocketblend/rocketblend/releases) page of our GitHub repository.

### Windows

To install pre-compiled binaries on Windows, follow these steps:

#### Prerequisites

1. [PowerShell](https://docs.microsoft.com/en-us/powershell/scripting/install/installing-powershell) (pre-installed on Windows 10 and later)

#### Installation steps

1. Open a PowerShell prompt.
2. Set the execution policy to allow running local scripts if not already set: `Set-ExecutionPolicy RemoteSigned -Scope CurrentUser`
3. Download and run the `install.ps1` script: `Invoke-WebRequest -Uri "https://raw.githubusercontent.com/rocketblend/rocketblend/master/install.ps1" -OutFile "install.ps1" .\install.ps1`
4. Restart your PowerShell or Command Prompt session for the updated `PATH` to take effect.

```powershell
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/rocketblend/rocketblend/master/install.ps1" -OutFile "install.ps1"
.\install.ps1
```

### Linux and macOS

To install pre-compiled binaries on Linux or macOS, follow these steps:

1. Open a terminal window.
2. Download and run the `install.sh` script:

```shell-session
$ curl -LO "https://raw.githubusercontent.com/rocketblend/rocketblend/master/install.sh"
$ chmod +x install.sh
$ ./install.sh
```

## Source

For users wanting to install directly from source, you can use the `go install` command:

```shell-session
$ go install github.com/rocketblend/rocketblend/cmd/rocketblend@latest
```

This command will download the latest version of the `rocketblend` source code and compile the binary for your platform. Ensure you have Go 1.19 or later installed on your system.



