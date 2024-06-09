#!/usr/bin/env bash

set -xu

export WO_DEBUG=true
export VISUAL=cat
export SHELL=/bin/bash
export APP=bash

shopt -s expand_aliases

create_function() {
echo '
# Hello world function
hello() {
  echo "Hello world !"
}
' > ~/.config/wo/api/functions/functions.bash
}
