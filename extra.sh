#!/usr/bin/env bash
[[ "$DEBUG" ]] && set -x # Print commands and their arguments as they are executed.

set -e # Exit immediately if a command exits with a non-zero status.

dissasemble(){
    cd "$APP"/build/_output/bin
    go tool objdump main
}

