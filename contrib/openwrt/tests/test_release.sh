#!/bin/sh

TEST_DIR=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
REPO_ROOT=$(CDPATH= cd -- "$TEST_DIR/../../.." && pwd)
. "$TEST_DIR/testlib.sh"

test_goreleaser_declares_mips_softfloat_builds() {
    config="$REPO_ROOT/.goreleaser.yaml"
    assert_contains 'id: steamhostsync-mips' "$config" || return 1
    assert_contains 'mips' "$config" || return 1
    assert_contains 'mipsle' "$config" || return 1
    assert_contains 'gomips:' "$config" || return 1
    assert_contains 'softfloat' "$config"
}

test_allowlist_matches_configured_domains() {
    allowlist="$REPO_ROOT/contrib/openwrt/files/etc/steamhostsync/allowed-domains"
    expected="${TMPDIR:-/tmp}/steamhostsync-domains.$$.expected"
    actual="${TMPDIR:-/tmp}/steamhostsync-domains.$$.actual"
    awk '
        /^platforms:/ { platforms = 1; next }
        platforms && /^      - / {
            sub(/^      - /, "")
            print tolower($0)
        }
    ' "$REPO_ROOT/config.yaml" | sort -u >"$expected"
    sort -u "$allowlist" >"$actual"
    result=0
    cmp -s "$expected" "$actual" || result=1
    rm -f "$expected" "$actual"
    [ "$result" -eq 0 ] || fail "OpenWrt allowlist differs from config.yaml"
}

test_router_low_resource_config_is_bounded() {
    config="$REPO_ROOT/contrib/openwrt/files/etc/steamhostsync/router-config.yaml"
    assert_file_exists "$config" || return 1
    assert_contains 'concurrency: 1' "$config" || return 1
    assert_contains 'attempts: 1' "$config"
}

run_test test_goreleaser_declares_mips_softfloat_builds
run_test test_allowlist_matches_configured_domains
run_test test_router_low_resource_config_is_bounded
finish_tests
