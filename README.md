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

### With the install script

To install the binary in `/usr/local/bin`, run:

``` sh
curl -sSf https://raw.githubusercontent.com/antham/wo/main/installer.sh | sudo sh
```

To select another path, run:

``` sh
curl -sSf https://raw.githubusercontent.com/antham/wo/main/installer.sh | sh -s -- -o "<install_path>"
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

You can set the theme with the `-t` flag, it could be either `dark` or `light`, the default is the `light` theme.

### Bash

`source <(wo setup bash)`

### Fish

`wo setup fish | source`

### Zsh

`source <(wo setup zsh)`

## Usage

### Creating a workspace

To create a workspace for your project, use the `create` command, run:
``` sh
wo create cli $PWD/projects/cli
```

Once your workspace is created, an alias is created to jump into the project folder, you need to reload your shell to "activate" it, open a new terminal or in the existing run:

``` sh
exec <name_of_your_shell>
```

The alias in our case will be `c_cli`, so the `c_` prefix (you can configure that) followed by the name of your workspace.


### Adding functions to a workspace

To add some functions, run:

``` sh
wo edit cli
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

To run a function into a workspace, call the `run` command:

``` sh
wo run cli run_curl http://google.fr
```

The first parameter is the workspace to use, the second parameter, the function defined in the workspace, all following parameters are additional parameters you can access in your function with the usual way of accessing function parameters according to your shell.


A function is ran from the folder of your project, so you don't need to do anything to access a command relative to your project, let's say a `npm run` for instance.

### Running a function in an environment

All functions are ran in a `default` environment if you specified nothing, you can edit this environment with:

``` sh
wo env edit cli default
```

It will open an editor to let you add your environment variables or things you want to run with each functions.

If you want to create additional environment, run:

``` sh
wo env create cli prod
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

### Changing the path of an existing workspace

Run:

``` sh
wo config set cli $PWD/project/cli2
```

### Committing the workspaces

You can commit and push the folder containing all workspaces on a repository, it is located at:

``` sh
wo global get config-dir
```

A default `.gitignore` is provided to exclude all environment variables. At the moment the process of committing and pushing the workspaces is manual. 

When you restore a backup from git run `wo fix` to restore the default environment as the folder containing all the environments are not committed.

### To go further

Check the help of the command line

``` sh
wo help
```
