# Deis Client

[![Build Status](https://travis-ci.org/deis/workflow-cli.svg?branch=master)](https://travis-ci.org/deis/workflow-cli)
[![Go Report Card](http://goreportcard.com/badge/deis/workflow-cli)](http://goreportcard.com/report/deis/workflow-cli)
[![codebeat badge](https://codebeat.co/badges/05d314a8-ca61-4211-b69e-e7a3033662c8)](https://codebeat.co/projects/github-com-deis-workflow-cli)

Download Links:

- [64 Bit Linux](https://storage.googleapis.com/workflow-cli/deis-latest-linux-amd64)
- [32 Bit Linux](https://storage.googleapis.com/workflow-cli/deis-latest-linux-386)
- [64 Bit Mac OS X](https://storage.googleapis.com/workflow-cli/deis-latest-darwin-amd64)
- [32 Bit Max OS X](https://storage.googleapis.com/workflow-cli/deis-latest-darwin-386)

(Note: Windows builds are not yet supported. [#26](https://github.com/deis/workflow-cli/issues/26) currently tracks the work to support them).

`deis` is a command line utility used to interact with the [Deis](http://deis.io) open source PaaS.

Please add any [issues](https://github.com/deis/workflow-cli/issues) you find with this software to the [Deis Workflow CLI Project](https://github.com/deis/workflow-cli).

## Installation

### Pre-built Binary

Run the appropriate command for your system to download and install a `deis` binary:

- 64 Bit Linux: `curl https://storage.googleapis.com/workflow-cli/deis-latest-linux-amd64 > deis && chmod +x deis`
- 32 Bit Linux: `curl https://storage.googleapis.com/workflow-cli/deis-latest-linux-386 > deis && chmod +x deis`
- 64 Bit Mac OS X: `curl https://storage.googleapis.com/workflow-cli/deis-latest-darwin-amd64 > deis && chmod +x deis`
- 32 Bit Max OS X: `curl https://storage.googleapis.com/workflow-cli/deis-latest-darwin-386 > deis && chmod +x deis`

(Note: Windows builds are not yet supported. [#26](https://github.com/deis/workflow-cli/issues/26) currently tracks the work to support them).

### From Scratch

To compile the client from scratch, ensure you have Docker installed and run

	$ make bootstrap
	$ make build

`make bootstrap` will fetch all required dependencies, while `make build` will compile and install
the client in the current directory.

	$ ./deis --version

## Usage

Running `deis help` will give you a up to date list of `deis` commands.
To learn more about a command run `deis help <command>`.

## Windows Support

`deis` has experimental support for Windows. To build deis for Windows, you need to install
[go](https://golang.org/) and [glide](https://github.com/Masterminds/glide). Then run the `make.bat` script.

## License

see [LICENSE](https://github.com/deis/workflow-cli/blob/master/LICENSE)
