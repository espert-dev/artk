#!/usr/bin/env bash
set -eu -o pipefail

find \( -name 'go.sum' -o -name 'go.work.sum' \) -exec cat '{}' \;\
| sort\
| uniq\
> go.sum

exit 0
