# [![Go Report Card](https://goreportcard.com/badge/github.com/antham/wo)](https://goreportcard.com/report/github.com/antham/wo) [![codecov](https://codecov.io/gh/antham/wo/graph/badge.svg?token=l5zT9434GU)](https://codecov.io/gh/antham/wo) [![GitHub tag](https://img.shields.io/github/tag/antham/wo.svg)]()

Wo is a shell workspace manager inspired by the great https://github.com/jamesob/desk project.
To sum up, you can define helpers functions to run related to a project and define environment variables separated per environment (staging, prod, ...).

## Install

### With Go

If you have `go` installed you can run:

``` sh
go install github.com/antham/wo@latest
```

### With Archlinux

You will find it on aur, run with `yay`:

``` sh
yay -S wo-bin
```

### With Homebrew

Run:

``` sh
brew tap antham/homebrew-wo
brew install wo
```

### Other systems

You can find `deb`, `rpm` and `apk` packages on the release page : https://github.com/antham/wo/releases

### Binaries

You can find binaries for `linux` and `darwin` for `arm64` and `amd64` on the release page: https://github.com/antham/wo/releases


## Setup

Add the following command to your shell init file according to your shell.

You can customize how the aliases are generated, the default is to prefix them with `c_`, you can change this with the `-p` flag.

### Bash

`source <(wo setup bash)`

### Fish

`wo setup fish | source`

### Zsh

`source <(wo setup zsh)`
