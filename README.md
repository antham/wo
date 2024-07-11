# [![Go Report Card](https://goreportcard.com/badge/github.com/antham/wo)](https://goreportcard.com/report/github.com/antham/wo) [![codecov](https://codecov.io/gh/antham/wo/graph/badge.svg?token=l5zT9434GU)](https://codecov.io/gh/antham/wo) [![GitHub tag](https://img.shields.io/github/tag/antham/wo.svg)]()

Wo is a shell workspace manager inspired by the great https://github.com/jamesob/desk project:
* create a workspace for each of your project
* define environments (e.g., staging, production) for each workspace if needed or use the default one
* execute workspace functions from anywhere, using a specific environment if needed
* supports `bash`,`fish` and `zsh`

### Demonstration

We take as an example a workspace created for a cli project that output the environment variable `SECRET`.

[![asciicast](https://asciinema.org/a/yGEwo4mv3bNcmTM3YC0oD4kfN.svg)](https://asciinema.org/a/yGEwo4mv3bNcmTM3YC0oD4kfN)

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

You need to have several environment variable defined:

| Environment variable | Description                               |
|----------------------|-------------------------------------------|
| SHELL                | the location of the current shell program |
| VISUAL/EDITOR        | the editor to use to edit functions files |

It is advised to create a one letter alias for the run function, as your are going to use it a lot, like so:
``` sh
alias r="wo run"
```

Add the following command to your shell init file according to your shell.

You can customize how the aliases are generated (see below in usage what is the goal of those aliases), the default is to prefix them with `c_`, you can change this behaviour with the `-p` flag on the setup command.

### Bash

`source <(wo setup bash)`

### Fish

`wo setup fish | source`

### Zsh

`source <(wo setup zsh)`

## Usage

### Creating a workspace

First you need to create a workspace for your project, use the `create` command, run:
``` sh
wo create api $PWD/projects/api
```

The first parameter is any name you find convenient to refer to your project and the second parameter is the path of your project on your computer.


Once your workspace is created, an alias is created to jump into the project folder, you need to reload your shell to "activate" it, open a new terminal or in the existing run:

``` sh
exec <name_of_your_shell>
```

The alias in our case will be `c_api`, so the `c_` prefix (you can configure that) followed by the name of your workspace.


### Adding functions

When you have create a workspace, you can then add some functions, run:

``` sh
wo edit api
```

A file will be opened with your default editor, the function you add must fit with the shell you are currently using, if you add one comment line right before the function name it will be taken and used as the description of the function or if the shell is `fish` the description added with the `-d` will be used. 
Here are examples of how to define a function for every shell:

#### Bash
``` bash
# Run a curl request
run_curl() {
  curl $1
}
```

#### Zsh

``` zsh
# Run a curl request
run_curl() {
  curl $1
}
```

#### Fish

``` fish
function run_curl -d "Run a curl request"
  curl $argv[1]
end
```

### Running a function

Once you have defined a function you can call it with the `run` command, run:

``` sh
wo run cli run_curl http://google.fr
```

The first parameter is the workspace to use, the second parameter, the function defined in the workspace, all following parameters are additional parameters you can access in your function with the usual way of accessing function parameters according to your shell.


A function is ran from the folder of your project, so you don't need to do anything to access a command relative to your project, let's say a `npm run` for instance.

A function is ran with a `default` environment if you specified nothing. You can edit an environment with:

``` sh
wo edit env default
```

It will open an editor to let you add your environment variables or things you want to run with each functions.

If you want to create additional environment, run:

``` sh
wo create env prod
```

You can edit the environment with the previous `edit` command.

To use it, you simply provide it to the function to run like so:

``` sh
wo run -e prod cli run_curl http://google.fr
```

You get special environment variables that are defined for every functions:

| Environment variable | Usage                            |
|----------------------|----------------------------------|
| WO_ENV               | the name of the environment used |
| WO_NAME              | the name of the workspace used   |

