#!/usr/bin/env zsh

set -xu

export WO_DEBUG=true
export VISUAL=cat
export SHELL=/bin/zsh
export APP=zsh

create_function() {
echo '
# Hello world function
hello() {
  echo "Hello world !"
}
' > ~/.config/wo/workspaces/api/functions/functions.zsh
}
