#!/usr/bin/env fish

set -g -x fish_trace on

set -g -x SHELL /usr/bin/fish
set -g -x WO_DEBUG true
set -g -x VISUAL cat
set -g -x APP fish


function create_function
    echo '
function hello -d "Hello world function"
  echo "Hello world !"
end
' >~/.config/wo/api/functions/functions.fish
end
