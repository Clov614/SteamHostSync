#!/bin/sh
set -u

TEST_DIR=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
failed=0

for test_file in "$TEST_DIR"/test_*.sh; do
    printf '\n# %s\n' "$(basename "$test_file")"
    if ! sh "$test_file"; then
        failed=1
    fi
done

exit "$failed"
