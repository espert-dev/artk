#!/bin/bash
set -eu -o pipefail

# Print error message and exit with status 1.
die() {
    echo -e "\e[1;31m$*\e[0m" >&2
    exit 1
}

warn() {
    echo -e "\e[1;33m$*\e[0m" >&2
}

# Detect module directories.
declare -a module_directories
detect_module_directories() {
    for gomod in $(find -type f -name go.mod | sort); do
        module_directories+=("$(dirname "${gomod}" | sed -E 's#^\.?\/?##g')")
    done
}

# Extract release tags from new commit messages.
release_tags() {
    git log -1 --pretty=%B\
    | grep '^Release:'\
    | sed 's/^Release://g'\
    | xargs -n1\
    | sort -V\
    | uniq
}

# Tags must belong to one of the modules.
validate_tags() {
    for release_tag in $(release_tags); do
        path="$(echo "${release_tag}" | sed -E 's/\/?v[0-9]+(\.[0-9]+)+//g')"

        valid=0
        for module_directory in "${module_directories[@]}"; do
            if [[ "${path}" == "${module_directory}" ]]; then
                valid=1
                break
            fi
        done

        if [[ $valid -eq 0 ]]; then
            die "Invalid tag '${release_tag}'"
        fi

        echo "The tag '${release_tag}' will be created on main."
    done
}

# Validate that the tags do not exist already.
# This should run in branches other than main.
assert_all_tags_are_new() {
    for release_tag in $(release_tags); do
        if git rev-parse "${release_tag}" >/dev/null 2>&1; then
            die "The tag '${release_tag}' already exists!"
        fi
    done
}

# Create missing tags locally, but not on origin.
# This should only be run on main.
#
# We tolerate existing tags so that we can re-run a release job if it fails
# midway.
create_tags() {
    for release_tag in $(release_tags); do
        if git rev-parse "$release_tag" >/dev/null 2>/dev/null; then
            warn "Tag '$release_tag' already exists! Skipping."
            continue
        fi

        git tag "$release_tag"
        echo "Created tag ${release_tag}"
    done
}

if [[ $# -ne 1 ]]; then
    die "Usage: $(basename $0) (verify|create)"
fi

cmd="$1"
case "${cmd}" in
    verify)
        echo "Verifying tags..."
        detect_module_directories
        validate_tags
        assert_all_tags_are_new
        exit 0
        ;;
    create)
        echo "Creating tags..."
        detect_module_directories
        validate_tags
        create_tags
        exit 0
        ;;
    *)
        die "Unknown command '${cmd}'."
        ;;
esac

exit 0
