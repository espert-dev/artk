#!/usr/bin/env bash
set -eu -o pipefail

# Run the script at the root of the repo.
cd "$(dirname "${BASH_SOURCE[0]}")"

# Ensure that the artifact output directories exist.
mkdir -p .test

# Remove JUnit reports from previous runs.
find .test -type f -delete

build_module() {
    module_root="$1"
    go build -mod=readonly "./${module_root}/..."
}

test_module() {
    module_root="$1"
    gotestsum --junitfile ".test/${module_root/\//-}.junit.xml" -- \
        -mod=readonly -timeout=1m -failfast -cover -race "./${module_root}/..."
}

for module_root in core tech/*; do
    echo -e "\e[1;34mModule ${module_root}\e[0m"
    build_module "$module_root"
    test_module "$module_root"
    echo
done

echo -e "\e[1;32mOK\e[0m"
exit 0
