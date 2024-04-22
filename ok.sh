#!/usr/bin/env bash
set -eu -o pipefail

section() {
    echo -e "\e[1;34m$*\e[0m"
}

# Run the script at the root of the repo.
cd "$(dirname "${BASH_SOURCE[0]}")"

# Remove JUnit reports from previous runs.
rm -f junit.xml

modules=(
    "./..."
    "./tech/grpcerror/..."
)

section Building...
go build -mod=readonly "${modules[@]}"
echo

section Testing...
gotestsum --junitfile "junit.xml" --\
    -mod=readonly -timeout=1m -failfast -cover -race "${modules[@]}"
echo

section Vetting...
linter_config="$(realpath .golangci.yaml)"
for module in "${modules[@]}"; do
    # Some linters such as musttag fail unless the module starts at the CWD.
    pushd "$(dirname "${module}")"
    golangci-lint run --config="${linter_config}" ./...
    popd
done
echo

echo -e "\e[1;32mOK\e[0m"
exit 0
