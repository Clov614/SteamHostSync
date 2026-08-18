#!/bin/sh

TEST_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
REPO_ROOT=$(CDPATH= cd -- "$TEST_DIR/../../.." && pwd)
. "$TEST_DIR/testlib.sh"

UPDATE_SCRIPT="$REPO_ROOT/contrib/openwrt/files/usr/libexec/steamhostsync-update"

setup_update() {
    new_sandbox
    install_mock "$TEST_DIR/mocks/uclient-fetch" uclient-fetch
    install_mock "$TEST_DIR/mocks/dnsmasq-service" dnsmasq-service
    install_mock "$TEST_DIR/mocks/local-generator" local-generator

    export STEAMHOSTSYNC_STATE_DIR="$TEST_ROOT/var/lib/steamhostsync"
    export STEAMHOSTSYNC_LOCK_DIR="$TEST_ROOT/var/lock/steamhostsync.lock"
    export STEAMHOSTSYNC_ALLOWLIST="$TEST_DIR/fixtures/allowed-domains"
    export STEAMHOSTSYNC_FETCH_BIN="$TEST_ROOT/bin/uclient-fetch"
    export STEAMHOSTSYNC_DNSMASQ_BIN="$TEST_ROOT/bin/dnsmasq-service"
    export STEAMHOSTSYNC_PRIMARY_URL="https://example.test/Hosts"
    export STEAMHOSTSYNC_FALLBACK_URL=""
    export STEAMHOSTSYNC_MAX_BYTES=1048576
    export MOCK_DNSMASQ_COUNT="$TEST_ROOT/dnsmasq.count"
    export MOCK_GENERATOR_ARGS="$TEST_ROOT/generator.args"
    export MOCK_GENERATOR_FILE="$TEST_DIR/fixtures/valid-new"
    cp "$TEST_DIR/fixtures/valid-old" "$STEAMHOSTSYNC_STATE_DIR/Hosts"
}

test_valid_remote_file_is_published_atomically() {
    setup_update
    export STEAMHOSTSYNC_SOURCE=remote
    export MOCK_FETCH_FILE="$TEST_DIR/fixtures/valid-new"

    sh "$UPDATE_SCRIPT" || return 1

    assert_files_equal "$TEST_DIR/fixtures/valid-new" "$STEAMHOSTSYNC_STATE_DIR/Hosts" || return 1
    assert_files_equal "$TEST_DIR/fixtures/valid-old" "$STEAMHOSTSYNC_STATE_DIR/Hosts.previous" || return 1
    assert_eq 1 "$(cat "$MOCK_DNSMASQ_COUNT")" || return 1
    assert_file_absent "$TEST_ROOT/etc/hosts"
}

test_invalid_downloads_keep_last_known_good() {
    for fixture in invalid-html invalid-ip invalid-hostname private-ip; do
        setup_update
        export STEAMHOSTSYNC_SOURCE=remote
        export MOCK_FETCH_FILE="$TEST_DIR/fixtures/$fixture"

        if sh "$UPDATE_SCRIPT"; then
            fail "expected $fixture to be rejected"
            return 1
        fi
        assert_files_equal "$TEST_DIR/fixtures/valid-old" "$STEAMHOSTSYNC_STATE_DIR/Hosts" || return 1
        assert_file_absent "$MOCK_DNSMASQ_COUNT" || return 1
    done
}

test_unchanged_content_does_not_reload_dnsmasq() {
    setup_update
    export STEAMHOSTSYNC_SOURCE=remote
    export MOCK_FETCH_FILE="$TEST_DIR/fixtures/valid-old"

    sh "$UPDATE_SCRIPT" || return 1
    assert_file_absent "$MOCK_DNSMASQ_COUNT"
}

test_reload_failure_restores_previous_file() {
    setup_update
    export STEAMHOSTSYNC_SOURCE=remote
    export MOCK_FETCH_FILE="$TEST_DIR/fixtures/valid-new"
    export MOCK_DNSMASQ_FAIL_ONCE="$TEST_ROOT/fail-once"
    touch "$MOCK_DNSMASQ_FAIL_ONCE"

    if sh "$UPDATE_SCRIPT"; then
        fail "expected update to report reload failure"
        return 1
    fi

    assert_files_equal "$TEST_DIR/fixtures/valid-old" "$STEAMHOSTSYNC_STATE_DIR/Hosts" || return 1
    assert_eq 2 "$(cat "$MOCK_DNSMASQ_COUNT")"
}

test_existing_live_lock_prevents_second_update() {
    setup_update
    export STEAMHOSTSYNC_SOURCE=remote
    export MOCK_FETCH_FILE="$TEST_DIR/fixtures/valid-new"
    mkdir "$STEAMHOSTSYNC_LOCK_DIR"
    printf '%s\n' "$$" >"$STEAMHOSTSYNC_LOCK_DIR/pid"

    sh "$UPDATE_SCRIPT"
    status=$?

    assert_eq 75 "$status" || return 1
    assert_files_equal "$TEST_DIR/fixtures/valid-old" "$STEAMHOSTSYNC_STATE_DIR/Hosts"
}

test_successful_update_releases_lock() {
    setup_update
    export STEAMHOSTSYNC_SOURCE=remote
    export MOCK_FETCH_FILE="$TEST_DIR/fixtures/valid-new"

    sh "$UPDATE_SCRIPT" || return 1
    assert_file_absent "$STEAMHOSTSYNC_LOCK_DIR" || return 1
    sh "$UPDATE_SCRIPT" || return 1
    assert_file_absent "$STEAMHOSTSYNC_LOCK_DIR"
}

test_oversized_download_keeps_last_known_good() {
    setup_update
    export STEAMHOSTSYNC_SOURCE=remote
    export MOCK_FETCH_FILE="$TEST_DIR/fixtures/valid-new"
    export STEAMHOSTSYNC_MAX_BYTES=32

    if sh "$UPDATE_SCRIPT"; then
        fail "expected oversized file to be rejected"
        return 1
    fi
    assert_files_equal "$TEST_DIR/fixtures/valid-old" "$STEAMHOSTSYNC_STATE_DIR/Hosts" || return 1
    assert_file_absent "$MOCK_DNSMASQ_COUNT"
}

test_remote_sources_must_reach_consensus() {
    setup_update
    export STEAMHOSTSYNC_SOURCE=remote
    export STEAMHOSTSYNC_REQUIRE_CONSENSUS=1
    export STEAMHOSTSYNC_PRIMARY_URL=https://primary.example.test/Hosts
    export STEAMHOSTSYNC_FALLBACK_URL=https://fallback.example.test/Hosts
    export MOCK_FETCH_PRIMARY_FILE="$TEST_DIR/fixtures/valid-new"
    export MOCK_FETCH_FALLBACK_FILE="$TEST_DIR/fixtures/valid-old"

    if sh "$UPDATE_SCRIPT"; then
        fail "expected mismatched mirrors to be rejected"
        return 1
    fi
    assert_files_equal "$TEST_DIR/fixtures/valid-old" "$STEAMHOSTSYNC_STATE_DIR/Hosts" || return 1
    assert_file_absent "$MOCK_DNSMASQ_COUNT"
}

test_local_source_runs_generator_and_publishes_output() {
    setup_update
    export STEAMHOSTSYNC_SOURCE=local
    export STEAMHOSTSYNC_LOCAL_BINARY="$TEST_ROOT/bin/local-generator"
    export STEAMHOSTSYNC_LOCAL_CONFIG="$TEST_ROOT/router-config.yaml"
    : >"$STEAMHOSTSYNC_LOCAL_CONFIG"

    sh "$UPDATE_SCRIPT" || return 1

    assert_files_equal "$TEST_DIR/fixtures/valid-new" "$STEAMHOSTSYNC_STATE_DIR/Hosts" || return 1
    assert_contains '-readme' "$MOCK_GENERATOR_ARGS" || return 1
    assert_contains '-out' "$MOCK_GENERATOR_ARGS"
}

run_test test_valid_remote_file_is_published_atomically
run_test test_invalid_downloads_keep_last_known_good
run_test test_unchanged_content_does_not_reload_dnsmasq
run_test test_reload_failure_restores_previous_file
run_test test_existing_live_lock_prevents_second_update
run_test test_successful_update_releases_lock
run_test test_oversized_download_keeps_last_known_good
run_test test_remote_sources_must_reach_consensus
run_test test_local_source_runs_generator_and_publishes_output
finish_tests
