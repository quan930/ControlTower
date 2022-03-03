#!/bin/bash
# server.sh
set -x

handle_func() {
  while read -r line; do
          echo "The client sended: $line" >&2
          if [ "$line" = hello ]; then
                  kill 1
                  exit 0
         fi
 done
}
export -f handle_func
socat UNIX-LISTEN:/lifecycle/test.sock system:"bash -c handle_func"