#!/usr/bin/env sh

set -xu

export WO_DEBUG=true
export VISUAL=cat
export SHELL=/bin/sh
export APP=sh

create_function() {
echo '
# Hello world function
hello() {
  echo "Hello world !"
}
' > ~/.config/wo/api/functions/functions.sh
}
