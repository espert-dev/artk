#!/usr/bin/env bash
set -eu -o pipefail

build_module() {
    go build "$1/..."
}

test_module() {
    gotestsum --junitfile junit.xml -- \
        -timeout=1m -failfast -cover -race "$1/..."
}

for module in ./core ./tech/*; do
    echo -e "\e[1;34mModule $module\e[0m"
    build_module "$module"
    test_module "$module"
    echo
done

echo -e "\e[1;32mOK\e[0m"
exit 0
