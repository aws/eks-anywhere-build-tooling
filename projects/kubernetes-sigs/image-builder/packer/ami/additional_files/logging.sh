#!/bin/bash

log::error() {
  local message="${1}"
  timestamp=$(date --iso-8601=seconds)
  echo "!!! [${timestamp}] ${1}" >&2
}

# Print a status line. Formatted to show up in a stream of output.
log::info() {
  timestamp=$(date --iso-8601=seconds)
  echo "+++ [${timestamp}] ${1}"
}
