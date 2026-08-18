#!/bin/sh

TEST_COUNT=0
FAIL_COUNT=0

fail() {
    printf '    %s\n' "$*" >&2
    return 1
}

assert_eq() {
    expected=$1
    actual=$2
    message=${3:-expected values to match}
    [ "$expected" = "$actual" ] || fail "$message"
}

assert_file_exists() {
    [ -f "$1" ] || fail "expected file to exist: $1"
}

assert_file_absent() {
    [ ! -e "$1" ] || fail "expected path to be absent: $1"
}

assert_files_equal() {
    cmp -s "$1" "$2" || fail "files differ: $1 $2"
}

assert_contains() {
    needle=$1
    file=$2
    grep -F -e "$needle" "$file" >/dev/null 2>&1 || fail "missing expected text in $file"
}

assert_line_count() {
    expected=$1
    pattern=$2
    file=$3
    actual=$(grep -cF -e "$pattern" "$file" 2>/dev/null || true)
    assert_eq "$expected" "$actual" "unexpected line count in $file"
}

run_test() {
    name=$1
    TEST_COUNT=$((TEST_COUNT + 1))
    if ("$name"); then
        printf 'ok %d - %s\n' "$TEST_COUNT" "$name"
    else
        FAIL_COUNT=$((FAIL_COUNT + 1))
        printf 'not ok %d - %s\n' "$TEST_COUNT" "$name"
    fi
}

finish_tests() {
    if [ "$FAIL_COUNT" -ne 0 ]; then
        printf '%d of %d tests failed\n' "$FAIL_COUNT" "$TEST_COUNT" >&2
        exit 1
    fi
    printf '%d tests passed\n' "$TEST_COUNT"
}

new_sandbox() {
    TEST_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/steamhostsync-test.XXXXXX") || exit 1
    export TEST_ROOT
    mkdir -p "$TEST_ROOT/bin" "$TEST_ROOT/etc/crontabs" "$TEST_ROOT/var/lib/steamhostsync" "$TEST_ROOT/var/lock"
}

install_mock() {
    source=$1
    name=$2
    cp "$source" "$TEST_ROOT/bin/$name" || return 1
    chmod +x "$TEST_ROOT/bin/$name"
}
