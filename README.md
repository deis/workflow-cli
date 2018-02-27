
|![](https://upload.wikimedia.org/wikipedia/commons/thumb/1/17/Warning.svg/156px-Warning.svg.png) | Deis Workflow is no longer maintained.<br />Please [read the announcement](https://deis.com/blog/2017/deis-workflow-final-release/) for more detail. |
|---:|---|
| 09/07/2017 | Deis Workflow [v2.18][] final release before entering maintenance mode |
| 03/01/2018 | End of Workflow maintenance: critical patches no longer merged |
| | [Hephy](https://github.com/teamhephy/workflow) is a fork of Workflow that is actively developed and accepts code contributions. |

# Deis Client

[![Build Status](https://ci.deis.io/buildStatus/icon?job=Deis/workflow-cli/master)](https://ci.deis.io/job/Deis/job/workflow-cli/job/master/)
[![Go Report Card](https://goreportcard.com/badge/github.com/deis/workflow-cli)](https://goreportcard.com/report/github.com/deis/workflow-cli)
[![codebeat badge](https://codebeat.co/badges/05d314a8-ca61-4211-b69e-e7a3033662c8)](https://codebeat.co/projects/github-com-deis-workflow-cli)
[![codecov](https://codecov.io/gh/deis/workflow-cli/branch/master/graph/badge.svg)](https://codecov.io/gh/deis/workflow-cli)

Download Links:

- [64 Bit Linux](https://storage.googleapis.com/workflow-cli-master/deis-latest-linux-amd64)
- [32 Bit Linux](https://storage.googleapis.com/workflow-cli-master/deis-latest-linux-386)
- [64 Bit Mac OS X](https://storage.googleapis.com/workflow-cli-master/deis-latest-darwin-amd64)
- [32 Bit Max OS X](https://storage.googleapis.com/workflow-cli-master/deis-latest-darwin-386)
- [64 Bit Windows](https://storage.googleapis.com/workflow-cli-master/deis-latest-windows-amd64.exe)
- [32 Bit Windows](https://storage.googleapis.com/workflow-cli-master/deis-latest-windows-386.exe)

`deis` is a command line utility used to interact with the [Deis](http://deis.io) open source PaaS.

Please add any [issues](https://github.com/deis/workflow-cli/issues) you find with this software to the [Deis Workflow CLI Project](https://github.com/deis/workflow-cli).

## Installation

### Pre-built Binary

See the appropriate sub-section below for your system to download and install the latest build of this software.

#### 64 Bit Linux

Master:

```console
curl -o deis https://storage.googleapis.com/workflow-cli-master/deis-latest-linux-amd64 && chmod +x deis
```

Latest stable release:

```
curl -o deis https://storage.googleapis.com/workflow-cli-release/deis-stable-linux-amd64 && chmod +x deis
```

#### 32 Bit Linux

Master:

```console
curl -o deis https://storage.googleapis.com/workflow-cli-master/deis-latest-linux-386 && chmod +x deis
```

Latest stable release:

```
curl -o deis https://storage.googleapis.com/workflow-cli-release/deis-stable-linux-386 && chmod +x deis
```

#### 64 Bit Mac OS X

Master:

```console
curl -o deis https://storage.googleapis.com/workflow-cli-master/deis-latest-darwin-amd64 && chmod +x deis
```

Latest stable release:

```
curl -o deis https://storage.googleapis.com/workflow-cli-release/deis-stable-darwin-amd64 && chmod +x deis
```

#### 32 Bit Max OS X

Master:

```console
curl -o deis https://storage.googleapis.com/workflow-cli-master/deis-latest-darwin-386 && chmod +x deis
```

Latest stable release:

```
curl -o deis https://storage.googleapis.com/workflow-cli-release/deis-stable-darwin-386 && chmod +x deis
```

#### 64 Bit Windows

Master:

```console
powershell -NoProfile -ExecutionPolicy Bypass -Command "(new-object net.webclient).DownloadFile('https://storage.googleapis.com/workflow-cli-master/deis-latest-windows-amd64.exe', 'deis.exe')"
```

Latest stable release:

```
powershell -NoProfile -ExecutionPolicy Bypass -Command "(new-object net.webclient).DownloadFile('https://storage.googleapis.com/workflow-cli-release/deis-stable-windows-amd64.exe', 'deis.exe')"
```

#### 32 Bit Windows

Master:

```console
powershell -NoProfile -ExecutionPolicy Bypass -Command "(new-object net.webclient).DownloadFile('https://storage.googleapis.com/workflow-cli-master/deis-latest-windows-386.exe', 'deis.exe')"
```

Latest stable release:

```
powershell -NoProfile -ExecutionPolicy Bypass -Command "(new-object net.webclient).DownloadFile('https://storage.googleapis.com/workflow-cli-release/deis-stable-windows-386.exe', 'deis.exe')"
```


After you execute the appropriate command for your system, you'll have a `deis` binary in the current directory.

Run the following to see the version:

```console
$ ./deis --version
```

You can then move it anywhere in your path:

```console
$ mv deis /usr/local/bin
```

### From Scratch on OS X and Linux

To compile the client from scratch, ensure you have Docker installed and run

    $ make

### From Scratch on Windows

To compile the client from scratch, open PowerShell and execute the following commands in the source directory.

    $ .\make bootstrap
    $ .\make build

`.\make bootstrap` will fetch all required dependencies, while `.\make build` will compile and install
the client in the current directory.

    $ .\deis --version

## Usage

Running `deis help` will give you a up to date list of `deis` commands.
To learn more about a command run `deis help <command>`.

## License

see [LICENSE](https://github.com/deis/workflow-cli/blob/master/LICENSE)

[v2.18]: https://github.com/deis/workflow/releases/tag/v2.18.0
