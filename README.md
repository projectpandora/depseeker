[![License](https://img.shields.io/badge/license-MIT-_red.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/projectpandora/depseeker)](https://goreportcard.com/report/github.com/projectpandora/depseeker)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/projectpandora/depseeker/issues)
[![GitHub Release](https://img.shields.io/github/release/projectpandora/depseeker)](https://github.com/projectpandora/depseeker/releases)
[![Docker Images](https://img.shields.io/docker/pulls/projectpandora/depseeker.svg)](https://hub.docker.com/r/projectpandora/depseeker)

depseeker is a fast and multi-purpose toolkit for finding npm dependencies in web applications, it is designed to maintain the result reliability with increased threads.

# Resources

- [Resources](#resources)
- [Features](#features)
- [Installation Instructions](#installation-instructions)
  - [From Binary](#from-binary)
  - [From Source](#from-source)
  - [From Github](#from-github)
- [Usage](#usage)
- [Thanks](#thanks)

# Features

- Check for whether dependencies are public or private.
- Simple and modular code base making it easy to contribute.
- Fast and fully configurable flags for many usecases.

# Installation Instructions

### From Binary

The installation is easy. You can download the pre-built binaries for your platform from the [Releases](https://github.com/projectpandora/depseeker/releases/) page. Extract them using tar, move it to your `$PATH` and you're ready to go.

```sh
Download latest binary from https://github.com/projectpandora/depseeker/releases

▶ tar -xvf depseeker-linux-amd64.tar
▶ mv depseeker-linux-amd64 /usr/local/bin/depseeker
▶ depseeker -h
```

### From Source

depseeker requires **GO 1.14+** to install successfully. Run the following command to get the repo -

```sh
▶ GO111MODULE=on go get -v github.com/projectpandora/depseeker/cmd/depseeker
```

### From Github

```sh
▶ git clone https://github.com/projectpandora/depseeker.git; cd depseeker/cmd/depseeker; go build; mv depseeker /usr/local/bin/; depseeker -version
```

# Usage

```sh
depseeker -h
```

This will display help for the tool.

### Running depseeker with stdin

This will run the tool against all the urls in `urls.txt`.

```
cat urls.txt | depseeker
```

### Running depseeker with file input

This will run the tool against all the hosts and subdomains in urls.txt.

```sh
depseeker -l urls.txt -silent
```

### Running depseeker with subfinder and httpx

```sh
subfinder -d hackerone.com -silent | httpx -silent | depseeker
```

# Thanks

depseeker is inspired by [projectdiscovery](https://projectdiscovery.io) works :heart:
