#!/usr/bin/env bash
set -eu -o pipefail

# Delete old generated files.
find -name '*\.gen\.go' -o -name '*_string.go' -delete

# Generate across all modules.
while IFS= read -r -d '' go_mod
do
    module="$(dirname "$go_mod")"
    echo "Regenerating $module ..."
    go generate "${module}/..."
done < <(find . -name 'go.mod' -print0)
