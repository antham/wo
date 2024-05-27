#!/usr/bin/env bash

export SHELL=/bin/bash
export VISUAL=cat

mkdir -p ~/api ~/front ~/db

# Create workspaces
wo create api ~/api
wo create front ~/front
wo create db ~/db

# List workspaces
echo "Workspaces
---
* api
* db
* front" > /tmp/expected-workspaces

wo list 2> /tmp/actual-workspaces

diff /tmp/expected-workspaces /tmp/actual-workspaces

# Show functions in a workspace

echo '
# Hello world function
test() {
  echo "Hello world !"
}
' > ~/.config/wo/functions/api.bash

echo "Workspace api
---
Functions

* test : Hello world function
---
Envs

* default" > /tmp/expected-api-workspace

wo show api 2> /tmp/actual-api-workspace

diff /tmp/expected-api-workspace /tmp/actual-api-workspace

# Run a function in a workspace

test "$(wo r api test)" = "Hello world !"
